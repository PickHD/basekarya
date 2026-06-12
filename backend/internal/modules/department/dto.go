package department

type LookupResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CreateDepartmentRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type UpdateDepartmentRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}
