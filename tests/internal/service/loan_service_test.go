package service_test

import (
	"BillingEngine/internal/domain"
	"BillingEngine/internal/repository"
	"BillingEngine/internal/service"
	"testing"
)

// MockLoanRepository is a mock implementation of the LoanRepository interface
type MockLoanRepository struct {
	loans map[domain.LoanID]*domain.Loan
}

func NewMockLoanRepository() *MockLoanRepository {
	return &MockLoanRepository{
		loans: make(map[domain.LoanID]*domain.Loan),
	}
}

func (m *MockLoanRepository) Save(loan *domain.Loan) error {
	m.loans[loan.ID] = loan
	return nil
}

func (m *MockLoanRepository) FindByID(id domain.LoanID) (*domain.Loan, error) {
	loan, ok := m.loans[id]
	if !ok {
		return nil, repository.ErrLoanNotFound
	}
	return loan, nil
}

func TestCreateLoan(t *testing.T) {
	repo := NewMockLoanRepository()
	loanService := service.NewLoanService(repo)

	loanID := domain.LoanID(1)
	loan, err := loanService.CreateLoan(loanID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if loan.ID != loanID {
		t.Errorf("Expected loan ID %d, got %d", loanID, loan.ID)
	}

	if loan.Outstanding != (domain.LoanAmount + (domain.WeeklyPayment * domain.TotalWeeks) - domain.LoanAmount) {
		t.Errorf("Expected outstanding amount %d, got %d", domain.LoanAmount, loan.Outstanding)
	}

	if len(loan.Schedule) != domain.TotalWeeks {
		t.Errorf("Expected %d payments in schedule, got %d", domain.TotalWeeks, len(loan.Schedule))
	}
}

func TestGetOutstanding(t *testing.T) {
	repo := NewMockLoanRepository()
	loanService := service.NewLoanService(repo)

	loanID := domain.LoanID(1)
	_, err := loanService.CreateLoan(loanID)
	if err != nil {
		t.Fatalf("Failed to create loan: %v", err)
	}

	outstanding, err := loanService.GetOutstanding(loanID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if outstanding != domain.LoanAmount+(domain.WeeklyPayment*domain.TotalWeeks)-domain.LoanAmount {
		t.Errorf("Expected outstanding amount %d, got %d", domain.LoanAmount, outstanding)
	}
}

func TestIsDelinquent(t *testing.T) {
	repo := NewMockLoanRepository()
	loanService := service.NewLoanService(repo)

	loanID := domain.LoanID(1)
	loan, err := loanService.CreateLoan(loanID)
	if err != nil {
		t.Fatalf("Failed to create loan: %v", err)
	}

	delinquent, err := loanService.IsDelinquent(loanID, loan.StartDate.AddDate(0, 0, 7))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if delinquent {
		t.Errorf("Expected loan to not be delinquent")
	}

	delinquent, err = loanService.IsDelinquent(loanID, loan.StartDate.AddDate(0, 0, 21))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !delinquent {
		t.Errorf("Expected loan to be delinquent")
	}
}

func TestMakePayment(t *testing.T) {
	repo := NewMockLoanRepository()
	loanService := service.NewLoanService(repo)

	loanID := domain.LoanID(1)
	loan, err := loanService.CreateLoan(loanID)
	if err != nil {
		t.Fatalf("Failed to create loan: %v", err)
	}

	// Test full payment
	err = loanService.MakePayment(loanID, domain.WeeklyPayment, loan.StartDate.AddDate(0, 0, 7))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	outstanding, _ := loanService.GetOutstanding(loanID)
	expectedOutstanding := (domain.LoanAmount + (domain.WeeklyPayment * domain.TotalWeeks) - domain.LoanAmount) - domain.WeeklyPayment
	if outstanding != expectedOutstanding {
		t.Errorf("Expected outstanding amount %d, got %d", expectedOutstanding, outstanding)
	}

	// Test partial payment
	err = loanService.MakePayment(loanID, domain.WeeklyPayment/2, loan.StartDate.AddDate(0, 0, 14))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	outstanding, _ = loanService.GetOutstanding(loanID)
	expectedOutstanding = (domain.LoanAmount + (domain.WeeklyPayment * domain.TotalWeeks) - domain.LoanAmount) - domain.WeeklyPayment - domain.WeeklyPayment/2
	if outstanding != expectedOutstanding {
		t.Errorf("Expected outstanding amount %d, got %d", expectedOutstanding, outstanding)
	}

	// Test overpayment
	err = loanService.MakePayment(loanID, domain.WeeklyPayment*2, loan.StartDate.AddDate(0, 0, 21))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	outstanding, _ = loanService.GetOutstanding(loanID)
	expectedOutstanding = (domain.LoanAmount + (domain.WeeklyPayment * domain.TotalWeeks) - domain.LoanAmount) - domain.WeeklyPayment - domain.WeeklyPayment/2 - domain.WeeklyPayment*2
	if outstanding != expectedOutstanding {
		t.Errorf("Expected outstanding amount %d, got %d", expectedOutstanding, outstanding)
	}
}
