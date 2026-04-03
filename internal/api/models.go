package api

import (
	"time"
)

type Campaign struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"` // draft, scheduled, sending, sent, stopped
	Type      string    `json:"type"`   // one-time, drip
	CreatedAt time.Time `json:"createdAt"`
}

type Contact struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone,omitempty"`
	Description string    `json:"description,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Status      string    `json:"status"` // active, unsubscribed, bounced
	CreatedAt   time.Time `json:"createdAt"`
}

type Stats struct {
	TotalSent           int     `json:"totalSent"`
	TotalDelivered      int     `json:"totalDelivered"`
	TotalOpens          int     `json:"totalOpens"`
	TotalClicks         int     `json:"totalClicks"`
	CTR                 float64 `json:"ctr"`
	TotalFailed         int     `json:"totalFailed"`
	TotalSpamComplaints int     `json:"totalSpamComplaints"`
}

type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role,omitempty"`
}

type BrandConfig struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Tone         string `json:"tone"`
	CustomPrompt string `json:"customPrompt"`
	IsDefault    bool   `json:"isDefault"`
}

type ReferralInfo struct {
	Account *struct {
		ID           string `json:"id"`
		ReferralCode string `json:"referralCode"`
		Count        struct {
			Clicks      int `json:"clicks"`
			Referrals   int `json:"referrals"`
			Commissions int `json:"commissions"`
		} `json:"_count"`
	} `json:"account"`
	Balances struct {
		Pending  int `json:"pending"`
		Approved int `json:"approved"`
		Paid     int `json:"paid"`
		Voided   int `json:"voided"`
	} `json:"balances"`
}

type AffiliateAnalytics struct {
	Summary struct {
		TotalClicks      int     `json:"totalClicks"`
		Signups          int     `json:"signups"`
		PaidConversions  int     `json:"paidConversions"`
		ConversionRate   float64 `json:"conversionRate"`
		TotalEarnings    int     `json:"totalEarnings"`
		PendingEarnings  int     `json:"pendingEarnings"`
	} `json:"summary"`
	RecentClicks []ClickRecord    `json:"recentClicks"`
	Referrals    []ReferralRecord `json:"referrals"`
}

type ClickRecord struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	IP        string    `json:"ip"`
	UA        string    `json:"ua"`
	Converted bool      `json:"converted"`
}

type ReferralRecord struct {
	ID          string    `json:"id"`
	TeamName    string    `json:"teamName"`
	Status      string    `json:"status"` // active, trial
	CreatedAt   time.Time `json:"createdAt"`
	FirstPaidAt *time.Time `json:"firstPaidAt,omitempty"`
}

type PayoutRecord struct {
	ID        string    `json:"id"`
	Amount    int       `json:"amount"` // in cents
	Status    string    `json:"status"` // paid, pending, failed
	Method    string    `json:"method"` // paypal
	Reference string    `json:"reference,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
