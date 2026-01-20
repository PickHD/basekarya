package reimburstment

type ReimbursementFilter struct {
	UserID uint
	Status string
	Limit  int
	Offset int
}