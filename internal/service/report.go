package service

import (
	"fmt"

	"zl0y-billing/internal/repository"
)

const ReportCost = 500 // 5.00 in cents

type ReportService struct {
	reportRepo *repository.ReportRepository
	userRepo   *repository.UserRepository
}

func NewReportService(reportRepo *repository.ReportRepository, userRepo *repository.UserRepository) *ReportService {
	return &ReportService{
		reportRepo: reportRepo,
		userRepo:   userRepo,
	}
}

func (s *ReportService) PurchaseReport(userID int, reportID string) error {
	// Get the report
	report, err := s.reportRepo.GetReportByID(reportID)
	if err != nil {
		return fmt.Errorf("report not found")
	}

	// Check if report belongs to user
	if report.UserID == nil || *report.UserID != userID {
		return fmt.Errorf("report not found")
	}

	// Check if already purchased
	if report.IsPurchased {
		return fmt.Errorf("report already purchased")
	}

	// Get user to check balance
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Check if user has sufficient balance
	if user.Balance < ReportCost {
		return fmt.Errorf("insufficient balance")
	}

	// Simulate transaction by performing both operations
	// In a real scenario, this would be wrapped in a distributed transaction
	// or use a saga pattern for cross-database consistency

	// 1. Deduct balance from user
	if err := s.userRepo.DeductBalance(userID, ReportCost); err != nil {
		return fmt.Errorf("failed to deduct balance: %w", err)
	}

	// 2. Mark report as purchased
	if err := s.reportRepo.MarkReportAsPurchased(reportID); err != nil {
		// In a real scenario, you would need to compensate by adding balance back
		// or implement proper distributed transaction handling
		return fmt.Errorf("failed to mark report as purchased: %w", err)
	}

	return nil
}
