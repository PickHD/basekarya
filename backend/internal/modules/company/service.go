package company

import (
	"basekarya-backend/pkg/constants"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service interface {
	GetProfile(ctx context.Context) (*CompanyProfileResponse, error)
	UpdateProfile(ctx context.Context, req *UpdateCompanyProfileRequest, file *multipart.FileHeader) error
}

type service struct {
	repo    Repository
	cache   CacheProvider
	storage StorageProvider
}

func NewService(repo Repository, cache CacheProvider, storage StorageProvider) Service {
	return &service{repo, cache, storage}
}

func (s *service) GetProfile(ctx context.Context) (*CompanyProfileResponse, error) {
	cacheKey := fmt.Sprintf(constants.COMPANY_PROFILE_CACHE_KEY, 1)

	cacheData, err := s.cache.Get(ctx, cacheKey)
	if err == redis.Nil {
		data, err := s.repo.FindByID(ctx, 1)
		if err != nil {
			return nil, err
		}

		parsedData, err := json.Marshal(&CompanyProfileResponse{
			ID:          data.ID,
			Name:        data.Name,
			Address:     data.Address,
			Email:       data.Email,
			PhoneNumber: data.PhoneNumber,
			Website:     data.Website,
			TaxNumber:   data.TaxNumber,
			LogoURL:     data.LogoURL,
		})
		if err != nil {
			return nil, err
		}

		err = s.cache.Set(ctx, cacheKey, parsedData, 24*time.Hour)
		if err != nil {
			return nil, err
		}

		return &CompanyProfileResponse{
			ID:          data.ID,
			Name:        data.Name,
			Address:     data.Address,
			Email:       data.Email,
			PhoneNumber: data.PhoneNumber,
			Website:     data.Website,
			TaxNumber:   data.TaxNumber,
			LogoURL:     data.LogoURL,
		}, nil
	} else if err != nil {
		return nil, err
	}

	var resp CompanyProfileResponse
	err = json.Unmarshal([]byte(cacheData), &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *service) UpdateProfile(ctx context.Context, req *UpdateCompanyProfileRequest, file *multipart.FileHeader) error {
	data, err := s.repo.FindByID(ctx, 1)
	if err != nil {
		return err
	}

	company, err := s.buildCompanyProfileData(ctx, data, req, file)
	if err != nil {
		return err
	}

	err = s.repo.Update(ctx, company)
	if err != nil {
		return err
	}

	err = s.cache.Del(ctx, fmt.Sprintf(constants.COMPANY_PROFILE_CACHE_KEY, 1))
	if err != nil {
		return err
	}

	return nil
}

func (s *service) buildCompanyProfileData(ctx context.Context, curr *Company, update *UpdateCompanyProfileRequest, file *multipart.FileHeader) (*Company, error) {
	if update.Name != "" {
		curr.Name = update.Name
	}

	if update.Address != "" {
		curr.Address = update.Address
	}

	if update.Email != "" {
		curr.Email = update.Email
	}

	if update.PhoneNumber != "" {
		curr.PhoneNumber = update.PhoneNumber
	}

	if update.Website != "" {
		curr.Website = update.Website
	}

	if update.TaxNumber != "" {
		curr.TaxNumber = update.TaxNumber
	}

	if file != nil {
		fileName := fmt.Sprintf("companies/%d/logo-%d.jpg", curr.ID, time.Now().Unix())
		fileURL, err := s.storage.UploadFileMultipart(ctx, file, fileName)
		if err != nil {
			return nil, err
		}

		curr.LogoURL = fileURL
	}

	return curr, nil
}
