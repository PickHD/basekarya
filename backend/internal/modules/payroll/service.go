package payroll

import (
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/logger"
	"basekarya-backend/pkg/response"
	"basekarya-backend/pkg/utils"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/signintech/gopdf"
)

type Service interface {
	GenerateAll(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
	GetList(ctx context.Context, filter *PayrollFilter) ([]PayrollListResponse, *response.Meta, error)
	GetDetail(ctx context.Context, id uint) (*PayrollDetailResponse, error)
	GeneratePayslipPDF(ctx context.Context, id uint) (*gopdf.GoPdf, *Payroll, error)
	MarkAsPaid(ctx context.Context, id uint) error
	BlastPayslipEmail(ctx context.Context, id uint) error
}

type service struct {
	repo               Repository
	user               UserProvider
	reimbursement      ReimbursementProvider
	attendance         AttendanceProvider
	company            CompanyProvider
	notification       NotificationProvider
	transactionManager infrastructure.TransactionManager
	client             *http.Client
	email              EmailProvider
	loan               LoanProvider
	overtime           OvertimeProvider
}

func NewService(repo Repository,
	user UserProvider,
	reimbursement ReimbursementProvider,
	attendance AttendanceProvider,
	company CompanyProvider,
	notification NotificationProvider,
	transactionManager infrastructure.TransactionManager,
	client *http.Client,
	email EmailProvider,
	loan LoanProvider,
	overtime OvertimeProvider) Service {
	return &service{repo, user, reimbursement, attendance, company, notification, transactionManager, client, email, loan, overtime}
}

func (s *service) GenerateAll(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	employees, err := s.user.FindAllEmployeeActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all employee active: %w", err)
	}

	employeeIds := make([]uint, len(employees))
	for i, emp := range employees {
		employeeIds[i] = emp.ID
	}

	existingPayrollMap, err := s.repo.GetExistingEmployeeID(req.Month, req.Year)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch existing employee id: %w", err)
	}

	attendanceMap, err := s.attendance.GetBulkLateDuration(ctx, req.Month, req.Year)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bulk late duration: %w", err)
	}

	reimburseMap, err := s.reimbursement.GetBulkApprovedAmount(ctx, req.Month, req.Year)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bulk approved amount: %w", err)
	}

	loanMap, err := s.loan.GetBulkActiveLoansByEmployeeIds(ctx, employeeIds)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bulk active loans by employee ids: %w", err)
	}

	overtimeMap, err := s.overtime.GetBulkActiveOvertimesByEmployeeIds(ctx, req.Month, req.Year, employeeIds)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bulk overtime amounts: %w", err)
	}

	successCount := 0
	periodDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.Local)

	var payrollsToInsert []Payroll

	for _, emp := range employees {
		// if already exist on this year & month, skip
		if existingPayrollMap[emp.ID] {
			continue
		}

		// take data with O(1) lookup
		baseSalary := emp.BaseSalary
		totalLateMinutes := attendanceMap[emp.ID]
		reimburseAmount := reimburseMap[emp.UserID]

		// calculate loan
		loanData := loanMap[emp.ID]
		loanAmount := loanData.InstallmentAmount
		loanData.RemainingAmount -= loanAmount

		// calculate overtime nominal
		// hourly wage = 1/173 * base salary
		hourlyWage := baseSalary / 173
		totalOvertimeMinutes := overtimeMap[emp.ID]
		overtimeAmount := 0.0

		if totalOvertimeMinutes > 0 {
			overtimeHours := float64(totalOvertimeMinutes) / 60.0

			// simplified rule: 1.5x for first hour, 2.0x for the rest
			if overtimeHours <= 1.0 {
				overtimeAmount = overtimeHours * 1.5 * hourlyWage
			} else {
				// first hour is 1.5x, remaining hours are 2.0x
				overtimeAmount = (1.0 * 1.5 * hourlyWage) + ((overtimeHours - 1.0) * 2.0 * hourlyWage)
			}
		}

		// calculate net salary
		latePenaltyAmount := float64(totalLateMinutes * constants.PenaltyPerMinuteLate)
		totalAllowance := baseSalary + reimburseAmount + overtimeAmount
		totalDeduction := latePenaltyAmount + loanAmount
		netSalary := totalAllowance - totalDeduction

		// construct object
		payroll := Payroll{
			EmployeeID:     emp.ID,
			PeriodDate:     periodDate,
			BaseSalary:     baseSalary,
			TotalAllowance: totalAllowance,
			TotalDeduction: totalDeduction,
			NetSalary:      netSalary,
			Status:         constants.PayrollStatusDraft,
			Details:        []PayrollDetail{},
		}

		payroll.Details = append(payroll.Details, PayrollDetail{
			Title:  "Base Salary",
			Type:   constants.DetailTypeAllowance,
			Amount: baseSalary,
		})

		// check if reimburse amount not zero
		if reimburseAmount > 0 {
			payroll.Details = append(payroll.Details, PayrollDetail{
				Title:  "Reimbursement",
				Type:   constants.DetailTypeAllowance,
				Amount: reimburseAmount,
			})
		}

		// check if overtime amount not zero
		if overtimeAmount > 0 {
			payroll.Details = append(payroll.Details, PayrollDetail{
				Title:  fmt.Sprintf("Uang Lembur (%d jam %d menit)", totalOvertimeMinutes/60, totalOvertimeMinutes%60),
				Type:   constants.DetailTypeAllowance,
				Amount: overtimeAmount,
			})
		}

		// check if late penalty amount not zero
		if latePenaltyAmount > 0 {
			payroll.Details = append(payroll.Details, PayrollDetail{
				Title:  fmt.Sprintf("Potongan Terlambat (%d menit)", totalLateMinutes),
				Type:   constants.DetailTypeDeduction,
				Amount: latePenaltyAmount,
			})
		}

		// check if loan amount not zero
		if loanAmount > 0 {
			payroll.Details = append(payroll.Details, PayrollDetail{
				Title:  "Potongan Kasbon",
				Type:   constants.DetailTypeDeduction,
				Amount: loanAmount,
			})
		}

		// insert to slice & update success count
		payrollsToInsert = append(payrollsToInsert, payroll)
		successCount++
	}

	// check if payrollsToInsert empty, return 0
	if len(payrollsToInsert) == 0 {
		return nil, nil
	}

	// bulk insert payrolls
	if err := s.repo.CreateBulk(&payrollsToInsert); err != nil {
		logger.Errorf("Failed create bulk payrolls %w", err)

		successCount = 0
		return nil, err
	}

	return &GenerateResponse{
		SuccessCount: successCount,
		Year:         req.Year,
		Month:        req.Month,
	}, nil
}

func (s *service) GetList(ctx context.Context, filter *PayrollFilter) ([]PayrollListResponse, *response.Meta, error) {
	data, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, nil, err
	}

	if len(data) == 0 {
		return []PayrollListResponse{}, nil, nil
	}

	var responses []PayrollListResponse
	for _, p := range data {
		empName := "Unknown"
		empNIK := "-"
		if p.Employee != nil {
			empName = p.Employee.FullName
			empNIK = p.Employee.NIK
		}

		responses = append(responses, PayrollListResponse{
			ID:           p.ID,
			EmployeeName: empName,
			EmployeeNIK:  empNIK,
			PeriodDate:   p.PeriodDate.Format(constants.DefaultTimeFormat),
			NetSalary:    p.NetSalary,
			Status:       string(p.Status),
			CreatedAt:    p.CreatedAt,
		})
	}

	meta := response.NewMetaOffset(filter.Page, filter.Limit, total)
	return responses, meta, nil
}

func (s *service) GetDetail(ctx context.Context, id uint) (*PayrollDetailResponse, error) {
	payroll, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if payroll.Employee == nil {
		return nil, errors.New("employee not found")
	}

	emp := payroll.Employee

	details := make([]Detail, len(payroll.Details))
	for _, detail := range payroll.Details {
		details = append(details, Detail{
			ID:        detail.ID,
			PayrollID: detail.PayrollID,
			Title:     detail.Title,
			Type:      detail.Type,
			Amount:    detail.Amount,
		})
	}

	payrollDetail := PayrollDetailResponse{
		ID:                        payroll.ID,
		EmployeeID:                emp.ID,
		EmployeeName:              emp.FullName,
		EmployeeNIK:               emp.NIK,
		EmployeeBankNumber:        emp.BankAccountNumber,
		EmployeeBankName:          emp.BankName,
		EmployeeBankAccountHolder: emp.BankAccountHolder,
		PeriodDate:                payroll.PeriodDate.Format(constants.DefaultTimeFormat),
		BaseSalary:                payroll.BaseSalary,
		TotalAllowance:            payroll.TotalAllowance,
		TotalDeduction:            payroll.TotalDeduction,
		NetSalary:                 payroll.NetSalary,
		Status:                    string(payroll.Status),
		CreatedAt:                 payroll.CreatedAt,
		Details:                   details,
	}

	return &payrollDetail, nil
}

func (s *service) GeneratePayslipPDF(ctx context.Context, id uint) (*gopdf.GoPdf, *Payroll, error) {
	payroll, err := s.repo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}

	company, err := s.company.FindByID(ctx, 1)
	if err != nil {
		return nil, nil, err
	}

	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	// Pastikan warna teks default hitam
	pdf.SetTextColor(0, 0, 0)

	// --- SETUP FONT ---
	err = pdf.AddTTFFont("Roboto", "assets/fonts/Roboto-Regular.ttf")
	if err != nil {
		return nil, nil, fmt.Errorf("failed load font regular: %w", err)
	}
	err = pdf.AddTTFFont("Roboto-Bold", "assets/fonts/Roboto-Bold.ttf")
	if err != nil {
		return nil, nil, fmt.Errorf("failed load font bold: %w", err)
	}

	marginLeft := 30.0
	marginRight := 565.0
	contentWidth := marginRight - marginLeft // 535.0 pt

	startY := 30.0
	currentY := startY
	logoRendered := false

	// --- SECTION: HEADER ---
	if company.LogoURL != "" {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, company.LogoURL, nil)

		if err == nil {
			resp, err := s.client.Do(req)
			if err == nil && resp.StatusCode == http.StatusOK {
				defer resp.Body.Close()

				imgHolder, err := gopdf.ImageHolderByReader(resp.Body)
				if err == nil {
					_ = pdf.ImageByHolder(imgHolder, marginLeft, startY, &gopdf.Rect{W: 70, H: 70})

					textStartX := marginLeft + 90.0

					pdf.SetXY(textStartX, startY+15)
					_ = pdf.SetFont("Roboto-Bold", "", 20)
					_ = pdf.Cell(nil, company.Name)

					pdf.SetXY(textStartX, startY+35)
					_ = pdf.SetFont("Roboto", "", 11)
					_ = pdf.Cell(nil, company.Address)

					pdf.SetXY(textStartX, startY+50)
					_ = pdf.Cell(nil, fmt.Sprintf("Telp: %s | Email: %s", company.PhoneNumber, company.Email))

					currentY = startY + 90
					logoRendered = true
				}
			}
		}
	}

	if !logoRendered {
		pdf.SetXY(marginLeft, startY)
		_ = pdf.SetFont("Roboto-Bold", "", 20)
		_ = pdf.CellWithOption(&gopdf.Rect{W: contentWidth, H: 25}, company.Name, gopdf.CellOption{Align: gopdf.Center})

		pdf.SetXY(marginLeft, startY+25)
		_ = pdf.SetFont("Roboto", "", 11)
		_ = pdf.CellWithOption(&gopdf.Rect{W: contentWidth, H: 15}, company.Address, gopdf.CellOption{Align: gopdf.Center})

		pdf.SetXY(marginLeft, startY+40)
		_ = pdf.CellWithOption(&gopdf.Rect{W: contentWidth, H: 15}, fmt.Sprintf("Telp: %s | Email: %s", company.PhoneNumber, company.Email), gopdf.CellOption{Align: gopdf.Center})

		currentY = startY + 70
	}

	pdf.SetLineWidth(1)
	pdf.Line(marginLeft, currentY, marginRight, currentY)
	currentY += 20

	// --- SECTION: EMPLOYEE INFO ---
	printInfo := func(x float64, y float64, label, value string) {
		pdf.SetXY(x, y)
		_ = pdf.SetFont("Roboto-Bold", "", 11)
		_ = pdf.Cell(nil, label)

		pdf.SetXY(x+70, y)
		_ = pdf.SetFont("Roboto", "", 11)
		_ = pdf.Cell(nil, ": "+value)
	}

	periodStr := payroll.PeriodDate.Format("January 2006")

	printInfo(marginLeft, currentY, "Name", payroll.Employee.FullName)
	printInfo(320, currentY, "Period", periodStr)
	currentY += 20

	printInfo(marginLeft, currentY, "NIK", payroll.Employee.NIK)
	printInfo(320, currentY, "Status", string(payroll.Status))
	currentY += 30

	// --- SECTION: PAYROLL TABLE ---
	halfWidth := contentWidth / 2
	col2X := marginLeft + halfWidth

	// Header Tabel
	pdf.SetFillColor(240, 240, 240)
	pdf.RectFromUpperLeftWithStyle(marginLeft, currentY, halfWidth, 25, "F")
	pdf.RectFromUpperLeftWithStyle(col2X, currentY, halfWidth, 25, "F")

	pdf.SetTextColor(0, 0, 0)

	pdf.SetXY(marginLeft, currentY)
	_ = pdf.SetFont("Roboto-Bold", "", 11)
	_ = pdf.CellWithOption(&gopdf.Rect{W: halfWidth, H: 25}, "  EARNINGS (PENDAPATAN)", gopdf.CellOption{Border: gopdf.AllBorders, Align: gopdf.Middle | gopdf.Left})

	pdf.SetXY(col2X, currentY)
	_ = pdf.CellWithOption(&gopdf.Rect{W: halfWidth, H: 25}, "  DEDUCTIONS (POTONGAN)", gopdf.CellOption{Border: gopdf.AllBorders, Align: gopdf.Middle | gopdf.Left})
	currentY += 25

	var earnings, deductions []PayrollDetail
	for _, d := range payroll.Details {
		if d.Type == constants.DetailTypeAllowance {
			earnings = append(earnings, d)
		} else {
			deductions = append(deductions, d)
		}
	}

	maxRows := len(earnings)
	if len(deductions) > maxRows {
		maxRows = len(deductions)
	}

	_ = pdf.SetFont("Roboto", "", 10)
	formatCurrency := func(amount float64) string {
		return fmt.Sprintf("Rp %s", utils.FormatNumber(amount))
	}

	rowH := 20.0
	labelW := 150.0
	valueW := halfWidth - labelW

	// Data Rows
	for i := 0; i < maxRows; i++ {
		// (Earnings)
		pdf.SetXY(marginLeft, currentY)
		if i < len(earnings) {
			_ = pdf.CellWithOption(&gopdf.Rect{W: labelW, H: rowH}, "   "+earnings[i].Title, gopdf.CellOption{Border: gopdf.Left, Align: gopdf.Middle | gopdf.Left})
			pdf.SetXY(marginLeft+labelW, currentY)
			_ = pdf.CellWithOption(&gopdf.Rect{W: valueW, H: rowH}, formatCurrency(earnings[i].Amount)+"   ", gopdf.CellOption{Border: gopdf.Right, Align: gopdf.Middle | gopdf.Right})
		} else {
			_ = pdf.CellWithOption(&gopdf.Rect{W: labelW, H: rowH}, "", gopdf.CellOption{Border: gopdf.Left})
			pdf.SetXY(marginLeft+labelW, currentY)
			_ = pdf.CellWithOption(&gopdf.Rect{W: valueW, H: rowH}, "", gopdf.CellOption{Border: gopdf.Right})
		}

		// (Deductions)
		pdf.SetXY(col2X, currentY)
		if i < len(deductions) {
			_ = pdf.CellWithOption(&gopdf.Rect{W: labelW, H: rowH}, "   "+deductions[i].Title, gopdf.CellOption{Border: gopdf.Left, Align: gopdf.Middle | gopdf.Left})
			pdf.SetXY(col2X+labelW, currentY)
			_ = pdf.CellWithOption(&gopdf.Rect{W: valueW, H: rowH}, formatCurrency(deductions[i].Amount)+"   ", gopdf.CellOption{Border: gopdf.Right, Align: gopdf.Middle | gopdf.Right})
		} else {
			_ = pdf.CellWithOption(&gopdf.Rect{W: labelW, H: rowH}, "", gopdf.CellOption{Border: gopdf.Left})
			pdf.SetXY(col2X+labelW, currentY)
			_ = pdf.CellWithOption(&gopdf.Rect{W: valueW, H: rowH}, "", gopdf.CellOption{Border: gopdf.Right})
		}
		currentY += rowH
	}

	// Total Row
	_ = pdf.SetFont("Roboto-Bold", "", 11)

	pdf.Line(marginLeft, currentY, marginRight, currentY)

	pdf.SetXY(marginLeft, currentY)
	_ = pdf.CellWithOption(&gopdf.Rect{W: labelW, H: 25}, "   Total Earnings", gopdf.CellOption{Border: gopdf.Left | gopdf.Bottom, Align: gopdf.Middle | gopdf.Left})
	pdf.SetXY(marginLeft+labelW, currentY)
	_ = pdf.CellWithOption(&gopdf.Rect{W: valueW, H: 25}, formatCurrency(payroll.TotalAllowance)+"   ", gopdf.CellOption{Border: gopdf.Right | gopdf.Bottom, Align: gopdf.Middle | gopdf.Right})

	pdf.SetXY(col2X, currentY)
	_ = pdf.CellWithOption(&gopdf.Rect{W: labelW, H: 25}, "   Total Deductions", gopdf.CellOption{Border: gopdf.Left | gopdf.Bottom, Align: gopdf.Middle | gopdf.Left})
	pdf.SetXY(col2X+labelW, currentY)
	_ = pdf.CellWithOption(&gopdf.Rect{W: valueW, H: 25}, formatCurrency(payroll.TotalDeduction)+"   ", gopdf.CellOption{Border: gopdf.Right | gopdf.Bottom, Align: gopdf.Middle | gopdf.Right})

	currentY += 40

	// --- SECTION: TAKE HOME PAY ---
	thpLabelW := 350.0
	thpValueW := contentWidth - thpLabelW

	pdf.SetFillColor(220, 230, 241)
	pdf.RectFromUpperLeftWithStyle(marginLeft, currentY, contentWidth, 30, "F")

	pdf.SetTextColor(0, 0, 0)

	_ = pdf.SetFont("Roboto-Bold", "", 14)
	pdf.SetXY(marginLeft, currentY)
	_ = pdf.CellWithOption(&gopdf.Rect{W: thpLabelW, H: 30}, "   TAKE HOME PAY ", gopdf.CellOption{Border: gopdf.Left | gopdf.Top | gopdf.Bottom, Align: gopdf.Middle | gopdf.Left})

	pdf.SetXY(marginLeft+thpLabelW, currentY)
	_ = pdf.CellWithOption(&gopdf.Rect{W: thpValueW, H: 30}, formatCurrency(payroll.NetSalary)+"   ", gopdf.CellOption{Border: gopdf.Right | gopdf.Top | gopdf.Bottom, Align: gopdf.Middle | gopdf.Right})

	currentY += 80

	// --- SECTION: SIGNATURE ---
	signatureX := marginRight - 150.0
	_ = pdf.SetFont("Roboto", "", 11)
	pdf.SetXY(signatureX, currentY)
	_ = pdf.CellWithOption(&gopdf.Rect{W: 150, H: 15}, "Authorized Signature,", gopdf.CellOption{Align: gopdf.Center})

	currentY += 70
	_ = pdf.SetFont("Roboto-Bold", "", 11)
	pdf.SetXY(signatureX, currentY)
	_ = pdf.CellWithOption(&gopdf.Rect{W: 150, H: 15}, "( HR Manager )", gopdf.CellOption{Border: gopdf.Top, Align: gopdf.Center})

	return pdf, payroll, nil
}

func (s *service) MarkAsPaid(ctx context.Context, id uint) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		payroll, err := s.repo.FindByID(id)
		if err != nil {
			return err
		}

		if payroll.Status == constants.PayrollStatusPaid {
			return nil
		}

		if err := s.repo.UpdateStatus(ctx, id, constants.PayrollStatusPaid); err != nil {
			return err
		}

		var deductedLoan float64
		for _, detail := range payroll.Details {
			if detail.Title == "Potongan Kasbon" && detail.Type == constants.DetailTypeDeduction {
				deductedLoan += detail.Amount
			}
		}

		if deductedLoan > 0 {
			loanMap, err := s.loan.GetBulkActiveLoansByEmployeeIds(ctx, []uint{payroll.EmployeeID})
			if err != nil {
				return fmt.Errorf("failed to fetch active loan: %w", err)
			}

			if activeLoan, exists := loanMap[payroll.EmployeeID]; exists {
				activeLoan.RemainingAmount -= deductedLoan
				if activeLoan.RemainingAmount <= 0 {
					activeLoan.RemainingAmount = 0
					activeLoan.Status = constants.LoanStatusPaidOff
				}

				if err := s.loan.Update(ctx, &activeLoan); err != nil {
					return fmt.Errorf("failed to update loan status: %w", err)
				}
			}
		}

		periodMonth := int(payroll.PeriodDate.Month())
		periodYear := payroll.PeriodDate.Year()
		if err := s.overtime.UpdateBulkStatusByEmployeeId(ctx, payroll.EmployeeID, periodMonth, periodYear, constants.OvertimeStatusPaid); err != nil {
			return fmt.Errorf("failed to update overtimes status to paid: %w", err)
		}

		go func() {
			_ = s.notification.SendNotification(
				payroll.Employee.UserID,
				string(constants.NotificationTypePayrollPaid),
				"Payroll Sudah dibayarkan",
				fmt.Sprintf("Payroll %s sudah dibayarkan.", payroll.PeriodDate.Format(constants.PayrollTimeFormat)),
				id,
			)
		}()

		return nil
	})
}

func (s *service) BlastPayslipEmail(ctx context.Context, id uint) error {
	pdfBytes, payroll, err := s.generatePayslipPDFBytes(ctx, id)
	if err != nil {
		return fmt.Errorf("failed generate pdf: %w", err)
	}

	if payroll.Status != constants.PayrollStatusPaid {
		return fmt.Errorf("payroll status must paid")
	}

	if payroll.Employee.Email == "" {
		return fmt.Errorf("email required, make sure to update first")
	}

	periodStr := payroll.PeriodDate.Format(constants.PayrollTimeFormat)
	subject := fmt.Sprintf("Payslip: %s - %s", periodStr, payroll.Employee.FullName)
	fileName := fmt.Sprintf("Payslip_%s_%s.pdf", strings.ReplaceAll(payroll.Employee.FullName, " ", "-"), payroll.PeriodDate.Format("Jan2006"))

	htmlBody := fmt.Sprintf(`
		<h3>Hello %s,</h3>
		<p>Terlampir adalah slip gaji Anda untuk periode <strong>%s</strong>.</p>
		<p>Harap jaga kerahasiaan dokumen ini. Jika ada pertanyaan, silakan hubungi tim HR.</p>
		<br>
		<p>Salam,</p>
		<p><strong>HR Manager</strong></p>
	`, payroll.Employee.FullName, periodStr)

	err = s.email.SendWithAttachment(
		payroll.Employee.Email,
		subject,
		htmlBody,
		fileName,
		pdfBytes,
	)
	if err != nil {
		return fmt.Errorf("failed to send email %s: %w", payroll.Employee.Email, err)
	}

	return nil
}

func (s *service) generatePayslipPDFBytes(ctx context.Context, id uint) ([]byte, *Payroll, error) {
	pdf, payroll, err := s.GeneratePayslipPDF(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	pdfBytes := pdf.GetBytesPdf()

	return pdfBytes, payroll, nil
}
