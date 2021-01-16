package model

// Payment is model for Pembayaran
type Payment struct {
	ID          int64         `json:"id" db:"id"`
	Name        string        `json:"name" db:"name"`
	Date        string        `json:"date" db:"date"`
	Amount      int           `json:"amount" db:"amount"`
	Method      PaymentMethod `json:"method" db:"method"`
	MethodDesc  string        `json:"method_desc" db:"-"`
	Account     AccountType   `json:"account" db:"account"`
	AccountDesc string        `json:"account_desc" db:"-"`
	AcceptBy    string        `json:"accept_by" db:"accept_by"`
}
