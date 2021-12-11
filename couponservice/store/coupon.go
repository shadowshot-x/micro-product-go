package store

type Coupon struct {
	Code        string `json:code`
	Name        string `json:name`
	Description string `json:description`
	VendorName  string `json:vendor`
}
