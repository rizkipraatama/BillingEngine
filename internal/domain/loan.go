package domain

import (
	"fmt"
	"time"
)

type LoanID int

type Payment struct {
	Week    int
	Amount  int
	Paid    bool
	DueDate time.Time
}

type Loan struct {
	ID           LoanID
	Schedule     []Payment
	Outstanding  int
	LastPaidWeek int
	StartDate    time.Time
}

const (
	TotalWeeks       = 50
	LoanAmount       = 5_000_000
	WeeklyPayment    = 110_000
	DelinquencyWeeks = 2
)

func IsDelinquent(loan *Loan, currentDate time.Time) bool {
	consecutiveMissedPayments := 0
	for _, payment := range loan.Schedule {
		if payment.DueDate.After(currentDate) {
			break
		}
		if !payment.Paid {
			consecutiveMissedPayments++
			if consecutiveMissedPayments > DelinquencyWeeks {
				return true
			}
		} else {
			consecutiveMissedPayments = 0
		}
	}
	return false
}

func (l *Loan) GetBillingInfo(currentDate time.Time) []string {
	var info []string
	info = append(info, fmt.Sprintf("Loan ID: %d", l.ID))
	info = append(info, fmt.Sprintf("Total Loan Amount: %d", LoanAmount))
	info = append(info, fmt.Sprintf("Outstanding Amount: %d", l.Outstanding))
	info = append(info, fmt.Sprintf("Delinquent: %v", IsDelinquent(l, currentDate)))
	info = append(info, "\nPayment Schedule:")
	info = append(info, "Week | Due Date | Amount | Status")
	info = append(info, "-------------------------------------")

	for _, payment := range l.Schedule {
		status := "Unpaid"
		if payment.Paid {
			status = "Paid"
		} else if payment.DueDate.Before(currentDate) {
			status = "Overdue"
		}
		info = append(info, fmt.Sprintf("%4d | %s | %7d | %s",
			payment.Week,
			payment.DueDate.Format("2006-01-02"),
			payment.Amount,
			status))
	}

	return info
}
