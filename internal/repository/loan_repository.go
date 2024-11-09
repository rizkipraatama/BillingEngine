package repository

import (
	"BillingEngine/internal/domain"
	"errors"
)

var ErrLoanNotFound = errors.New("loan not found")

type LoanRepository interface {
	Save(loan *domain.Loan) error
	FindByID(id domain.LoanID) (*domain.Loan, error)
}
