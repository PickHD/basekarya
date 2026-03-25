package master

import (
	"basekarya-backend/pkg/constants"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service interface {
	GetAllDepartments() ([]LookupResponse, error)
	GetAllShifts() ([]LookupResponse, error)
	GetAllLeaveTypes() ([]LookupLeaveTypeResponse, error)
}

type service struct {
	repo  Repository
	cache CacheProvider
}

func NewService(repo Repository, cache CacheProvider) Service {
	return &service{repo, cache}
}

func (s *service) GetAllDepartments() ([]LookupResponse, error) {
	cacheData, err := s.cache.Get(context.Background(), constants.DEPARTMEN_CACHE_KEY)
	if err == redis.Nil {
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

func (s *service) GetAllShifts() ([]LookupResponse, error) {
	cacheData, err := s.cache.Get(context.Background(), constants.SHIFT_CACHE_KEY)
	if err == redis.Nil {
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

		parsedData, err := json.Marshal(&results)
		if err != nil {
			return nil, err
		}

		err = s.cache.Set(context.Background(), constants.SHIFT_CACHE_KEY, parsedData, 24*time.Hour)
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

func (s *service) GetAllLeaveTypes() ([]LookupLeaveTypeResponse, error) {
	cacheData, err := s.cache.Get(context.Background(), constants.LEAVE_TYPE_CACHE_KEY)
	if err == redis.Nil {
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

		parsedData, err := json.Marshal(&results)
		if err != nil {
			return nil, err
		}

		err = s.cache.Set(context.Background(), constants.LEAVE_TYPE_CACHE_KEY, parsedData, 24*time.Hour)
		if err != nil {
			return nil, err
		}

		return results, nil
	}

	var results []LookupLeaveTypeResponse
	err = json.Unmarshal([]byte(cacheData), &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
