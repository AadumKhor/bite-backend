package models

// Route constants
const (
	// IdentifyRoute constant for `identify` route
	IdentifyRoute = "/identify"
)

// Error message constants
const (
	ErrMessageInvalidRequest     = "Invalid request body"
	ErrMessageInvalidPhoneNumber = "Invalid phone number"
)

// Miscellaneous
const (
	RegexPhone = "^[1-9][0-9]{9}$"
	TraceIDKey = "trace_id"
)
