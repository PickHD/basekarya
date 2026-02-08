package master

type Service interface {
	GetAllDepartments() ([]LookupResponse, error)
	GetAllShifts() ([]LookupResponse, error)
	GetAllLeaveTypes() ([]LookupLeaveTypeResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetAllDepartments() ([]LookupResponse, error) {
	var results []LookupResponse
	data, err := s.repo.FindAllDepartments()
	if err != nil {
		return nil, err
	}

	for _, d := range data {
		result := LookupResponse{
			ID:   d.ID,
			Name: d.Name,
		}

		results = append(results, result)
	}

	return results, nil
}

func (s *service) GetAllShifts() ([]LookupResponse, error) {
	var results []LookupResponse
	data, err := s.repo.FindAllShifts()
	if err != nil {
		return nil, err
	}

	for _, d := range data {
		result := LookupResponse{
			ID:   d.ID,
			Name: d.Name,
		}

		results = append(results, result)
	}

	return results, nil
}

func (s *service) GetAllLeaveTypes() ([]LookupLeaveTypeResponse, error) {
	var results []LookupLeaveTypeResponse
	data, err := s.repo.FindAllLeaveTypes()
	if err != nil {
		return nil, err
	}

	for _, d := range data {
		result := LookupLeaveTypeResponse{
			ID:           d.ID,
			Name:         d.Name,
			DefaultQuota: d.DefaultQuota,
			IsDeducted:   d.IsDeducted,
		}

		results = append(results, result)
	}

	return results, nil
}
