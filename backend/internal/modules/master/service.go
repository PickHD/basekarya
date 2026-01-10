package master

type Service interface {
	GetAllDepartments() ([]LookupResponse, error)
	GetAllShifts() ([]LookupResponse, error)
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
