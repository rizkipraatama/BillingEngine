package service

import (
	"BillingEngine/internal/domain"
	"BillingEngine/internal/repository"
	"fmt"
	"time"
)

type LoanService interface {
	CreateLoan(id domain.LoanID) (*domain.Loan, error)
	GetOutstanding(id domain.LoanID) (int, error)
	IsDelinquent(id domain.LoanID, currentDate time.Time) (bool, error)
	MakePayment(id domain.LoanID, amount int, paymentDate time.Time) error
	GetBillingInfo(id domain.LoanID, currentDate time.Time) ([]string, error)
}

type loanService struct {
	repo repository.LoanRepository
}

func NewLoanService(repo repository.LoanRepository) LoanService {
	return &loanService{repo: repo}
}

func (s *loanService) CreateLoan(id domain.LoanID) (*domain.Loan, error) {
	startDate := time.Now().AddDate(0, 0, 7)
	schedule := make([]domain.Payment, domain.TotalWeeks)
	for i := 0; i < domain.TotalWeeks; i++ {
		dueDate := startDate.AddDate(0, 0, 7*i)
		schedule[i] = domain.Payment{
			Week:    i + 1,
			Amount:  domain.WeeklyPayment,
			Paid:    false,
			DueDate: dueDate,
		}
	}

	loan := &domain.Loan{
		ID:           id,
		Schedule:     schedule,
		Outstanding:  domain.LoanAmount + (domain.WeeklyPayment * domain.TotalWeeks) - domain.LoanAmount,
		LastPaidWeek: 0,
		StartDate:    startDate,
	}

	err := s.repo.Save(loan)
	if err != nil {
		return nil, fmt.Errorf("failed to create loan: %w", err)
	}

	return loan, nil
}

func (s *loanService) GetOutstanding(id domain.LoanID) (int, error) {
	loan, err := s.repo.FindByID(id)
	if err != nil {
		return 0, fmt.Errorf("failed to get outstanding: %w", err)
	}
	return loan.Outstanding, nil
}

func (s *loanService) IsDelinquent(id domain.LoanID, currentDate time.Time) (bool, error) {
	loan, err := s.repo.FindByID(id)
	if err != nil {
		return false, fmt.Errorf("failed to check delinquency: %w", err)
	}

	return domain.IsDelinquent(loan, currentDate), nil
}

func (s *loanService) MakePayment(id domain.LoanID, amount int, paymentDate time.Time) error {
	loan, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to make payment: %w", err)
	}

	remainingAmount := amount
	var paymentMade bool

	for i := loan.LastPaidWeek; i < len(loan.Schedule); i++ {
		if !loan.Schedule[i].Paid {
			if remainingAmount >= domain.WeeklyPayment {
				// Full payment for this week
				loan.Schedule[i].Paid = true
				loan.Outstanding -= domain.WeeklyPayment
				remainingAmount -= domain.WeeklyPayment
				loan.LastPaidWeek = i + 1
				paymentMade = true
			} else if remainingAmount > 0 {
				// Partial payment for this week
				loan.Schedule[i].Paid = true
				loan.Outstanding -= remainingAmount
				loan.LastPaidWeek = i + 1
				paymentMade = true
				remainingAmount = 0
			}

			if remainingAmount == 0 {
				break
			}
		}
	}

	if !paymentMade {
		return fmt.Errorf("no eligible payments found")
	}

	if remainingAmount > 0 {
		// If there's still remaining amount, apply it to the next unpaid week
		for i := loan.LastPaidWeek; i < len(loan.Schedule); i++ {
			if !loan.Schedule[i].Paid {
				loan.Schedule[i].Paid = true
				loan.Outstanding -= remainingAmount
				loan.LastPaidWeek = i + 1
				break
			}
		}
	}

	err = s.repo.Save(loan)
	if err != nil {
		return fmt.Errorf("failed to save loan after payment: %w", err)
	}

	return nil
}

func (s *loanService) GetBillingInfo(id domain.LoanID, currentDate time.Time) ([]string, error) {
	loan, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing info: %w", err)
	}

	return loan.GetBillingInfo(currentDate), nil
}
