package repository

import (
	"BillingEngine/internal/domain"
	"fmt"
	"sync"
)

type inMemoryLoanRepository struct {
	loans map[domain.LoanID]*domain.Loan
	mutex sync.RWMutex
}

func NewInMemoryLoanRepository() LoanRepository {
	return &inMemoryLoanRepository{
		loans: make(map[domain.LoanID]*domain.Loan),
	}
}

func (r *inMemoryLoanRepository) Save(loan *domain.Loan) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.loans[loan.ID] = loan
	return nil
}

func (r *inMemoryLoanRepository) FindByID(id domain.LoanID) (*domain.Loan, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	loan, ok := r.loans[id]
	if !ok {
		return nil, fmt.Errorf("loan not found")
	}
	return loan, nil
}
