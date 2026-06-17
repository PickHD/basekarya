package constants

var PlanModules = map[string][]string{
	"free":  {"attendance", "leave"},
	"basic": {"attendance", "leave", "overtime", "loan", "reimbursement", "payroll", "contract", "finance", "bpjs", "tax"},
	"pro":   {"attendance", "leave", "overtime", "loan", "reimbursement", "payroll", "contract", "finance", "recruitment", "onboarding", "asset", "bpjs", "tax"},
}

var ModulePermissionGroups = map[string][]string{
	"attendance":    {"Attendance"},
	"leave":         {"Leave"},
	"overtime":      {"Overtime"},
	"loan":          {"Loan"},
	"reimbursement": {"Reimbursement"},
	"payroll":       {"Payroll"},
	"contract":      {"Contract"},
	"finance":       {"Finance"},
	"recruitment":   {"Recruitment"},
	"onboarding":    {"Onboarding"},
	"asset":         {"Asset"},
	"bpjs":          {"BPJS"},
	"tax":           {"Tax"},
}

var AlwaysAvailableGroups = []string{
	"Permission",
	"Role",
	"Master",
	"Employee",
	"Company",
	"Announcement",
}
