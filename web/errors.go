package web

// UserError represents an error to be returned in a user-visible
// manner
type UserError struct {
	Err string `json:"error"`
}

// Code returns the HTTP status code this error should return
func (*UserError) Code() int {
	return 400
}

// Error implements the error interface
func (ue *UserError) Error() string {
	return ue.Err
}
