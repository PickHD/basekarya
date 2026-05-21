package constants

var PlanModules = map[string][]string{
	"free":  {"attendance", "leave"},
	"basic": {"attendance", "leave", "overtime", "loan", "reimbursement", "payroll", "contract", "finance"},
	"pro":   {"attendance", "leave", "overtime", "loan", "reimbursement", "payroll", "contract", "finance", "recruitment", "onboarding"},
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
}

var AlwaysAvailableGroups = []string{
	"Permission",
	"Role",
	"Master",
	"Employee",
	"Company",
	"Announcement",
}
