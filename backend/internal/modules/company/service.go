package company

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"
)

type Service interface {
	GetProfile(ctx context.Context) (*CompanyProfileResponse, error)
	UpdateProfile(ctx context.Context, req *UpdateCompanyProfileRequest, file *multipart.FileHeader) error
}

type service struct {
	repo    Repository
	storage StorageProvider
}

func NewService(repo Repository, storage StorageProvider) Service {
	return &service{repo, storage}
}

func (s *service) GetProfile(ctx context.Context) (*CompanyProfileResponse, error) {
	data, err := s.repo.FindByID(ctx, 1)
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

	return s.repo.Update(ctx, company)
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
