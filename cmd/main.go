package main

import (
	"BillingEngine/internal"
	"BillingEngine/internal/domain"
	"fmt"
	"log"
	"time"
)

func main() {
	loanService, err := internal.InitializeLoanService()
	if err != nil {
		log.Fatalf("Failed to initialize loan service: %v", err)
	}

	loanID := domain.LoanID(100)
	loan, err := loanService.CreateLoan(loanID)
	if err != nil {
		log.Fatalf("Error creating loan: %v\n", err)
	}
	fmt.Printf("Loan created with ID: %d\n", loan.ID)

	billingInfo, err := loanService.GetBillingInfo(loanID, time.Now())
	if err != nil {
		log.Fatalf("Error getting billing info: %v\n", err)
	}

	fmt.Println("\nBilling Information:")
	for _, line := range billingInfo {
		fmt.Println(line)
	}

	outstanding, err := loanService.GetOutstanding(loanID)
	if err != nil {
		log.Fatalf("Error getting outstanding: %v\n", err)
	}
	fmt.Printf("Initial outstanding: %d\n", outstanding)

	delinquent, err := loanService.IsDelinquent(loanID, time.Now())
	if err != nil {
		log.Fatalf("Error checking delinquency: %v\n", err)
	}
	fmt.Printf("Is delinquent: %v\n", delinquent)

	for i := 0; i < 5; i++ {
		err := loanService.MakePayment(loanID, domain.WeeklyPayment, time.Now().AddDate(0, 0, 7*i))
		if err != nil {
			log.Fatalf("Error making payment: %v\n", err)
		}
		fmt.Printf("Payment %d made successfully\n", i+1)
	}

	// Get outstanding amount after payments
	outstanding, err = loanService.GetOutstanding(loanID)
	if err != nil {
		log.Fatalf("Error getting outstanding: %v\n", err)
	}
	fmt.Printf("Outstanding after 5 payments: %d\n", outstanding)

	// Check delinquency status after payments
	delinquent, err = loanService.IsDelinquent(loanID, time.Now().AddDate(0, 0, 35))
	if err != nil {
		log.Fatalf("Error checking delinquency: %v\n", err)
	}
	fmt.Printf("Is delinquent after 5 weeks: %v\n", delinquent)

	billingInfoAfterPayment, err := loanService.GetBillingInfo(loanID, time.Now())
	if err != nil {
		log.Fatalf("Error getting billing info: %v\n", err)
	}

	fmt.Println("\nBilling Information After Payment:")
	for _, line := range billingInfoAfterPayment {
		fmt.Println(line)
	}
	// Simulate missing 3 payments
	futureData := time.Now().AddDate(0, 0, 56)
	delinquent, err = loanService.IsDelinquent(loanID, futureData)
	if err != nil {
		log.Fatalf("Error checking delinquency: %v\n", err)
	}
	fmt.Printf("Is delinquent at %s after missing 2 payments: %v\n", futureData, delinquent)
}
