package payment_model

type InquiryModel struct {
	DoesExists bool `json:"doesExists,omitempty"`
	IsPaid     bool `json:"isPaid,omitempty"`
}
