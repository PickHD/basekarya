package department

import (
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service interface {
	GetAll(ctx context.Context) ([]LookupResponse, error)
	GetByID(ctx context.Context, id uint) (*LookupResponse, error)
	Create(ctx context.Context, req *CreateDepartmentRequest) (*LookupResponse, error)
	Update(ctx context.Context, id uint, req *UpdateDepartmentRequest) (*LookupResponse, error)
	Delete(ctx context.Context, id uint) error
}

type service struct {
	repo  Repository
	cache CacheProvider
}

func NewService(repo Repository, cache CacheProvider) Service {
	return &service{repo, cache}
}

func (s *service) GetAll(ctx context.Context) ([]LookupResponse, error) {
	cacheData, err := s.cache.Get(context.Background(), constants.DEPARTMEN_CACHE_KEY)
	if err == redis.Nil {
		var results []LookupResponse
		data, err := s.repo.FindAll(ctx)
		if err != nil {
			return nil, err
		}

		for _, d := range data {
			results = append(results, LookupResponse{ID: d.ID, Name: d.Name})
		}

		parsedData, err := json.Marshal(&results)
		if err != nil {
			return nil, err
		}

		err = s.cache.Set(context.Background(), constants.DEPARTMEN_CACHE_KEY, parsedData, 24*time.Hour)
		if err != nil {
			return nil, err
		}

		return results, nil
	}

	var results []LookupResponse
	err = json.Unmarshal([]byte(cacheData), &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*LookupResponse, error) {
	dept, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &LookupResponse{ID: dept.ID, Name: dept.Name}, nil
}

func (s *service) Create(ctx context.Context, req *CreateDepartmentRequest) (*LookupResponse, error) {
	exists, err := s.repo.ExistsByName(ctx, req.Name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("department name already exists")
	}

	dept := &Department{
		Name:      req.Name,
		CompanyID: utils.GetCompanyIDFromCtx(ctx),
	}
	if err := s.repo.Create(ctx, dept); err != nil {
		return nil, err
	}

	_ = s.cache.Del(context.Background(), constants.DEPARTMEN_CACHE_KEY)

	return &LookupResponse{ID: dept.ID, Name: dept.Name}, nil
}

func (s *service) Update(ctx context.Context, id uint, req *UpdateDepartmentRequest) (*LookupResponse, error) {
	dept, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("department not found")
		}
		return nil, err
	}

	exists, err := s.repo.ExistsByName(ctx, req.Name, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("department name already exists")
	}

	dept.Name = req.Name
	if err := s.repo.Update(ctx, dept); err != nil {
		return nil, err
	}

	_ = s.cache.Del(context.Background(), constants.DEPARTMEN_CACHE_KEY)

	return &LookupResponse{ID: dept.ID, Name: dept.Name}, nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("department not found")
		}
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	_ = s.cache.Del(context.Background(), constants.DEPARTMEN_CACHE_KEY)

	return nil
}
