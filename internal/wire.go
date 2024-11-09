//go:build wireinject
// +build wireinject

package internal

import (
	"BillingEngine/internal/repository"
	"BillingEngine/internal/service"
	"github.com/google/wire"
)

func InitializeLoanService() (service.LoanService, error) {
	wire.Build(
		repository.NewInMemoryLoanRepository,
		service.NewLoanService,
	)
	return nil, nil
}
