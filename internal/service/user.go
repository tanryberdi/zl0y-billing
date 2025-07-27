package service

import (
	"fmt"

	"zl0y-billing/internal/models"
	"zl0y-billing/internal/repository"
)

type UserService struct {
	userRepo   *repository.UserRepository
	reportRepo *repository.ReportRepository
}

func NewUserService(userRepo *repository.UserRepository, reportRepo *repository.ReportRepository) *UserService {
	return &UserService{
		userRepo:   userRepo,
		reportRepo: reportRepo,
	}
}

func (s *UserService) LinkAnonymousReport(clientGeneratedID string, userID int) (int, error) {
	// Verify if the user exists
	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return 0, fmt.Errorf("user not found: %w", err)
	}

	// Link the anonymous report to the user
	count, err := s.reportRepo.LinkAnonymousReport(clientGeneratedID, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to link anonymous report: %w", err)
	}

	return count, nil
}

func (s *UserService) GetUserReports(userID, limit, offset int) (*models.ReportsResponse, error) {
	// Verify if the user exists
	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Set default pagination values
	if limit <= 0 || limit > 100 {
		limit = 20 // Default limit
	}

	if offset < 0 {
		offset = 0 // Default offset
	}

	// Get user reports
	reports, total, err := s.reportRepo.GetReportsByUserID(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user reports: %w", err)
	}

	return &models.ReportsResponse{
		Reports: reports,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	}, nil
}
