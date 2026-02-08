package master

type LookupResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type LookupLeaveTypeResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	DefaultQuota int    `json:"default_quota"`
	IsDeducted   bool   `json:"is_deducted"`
}
