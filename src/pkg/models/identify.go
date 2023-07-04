// Package models contains all models used in this project
package models

import "github.com/gin-gonic/gin"

type (
	// IdentifyRequest is the incoming request model
	IdentifyRequest struct {
		Email       string `json:"email" binding:"omitempty"`
		PhoneNumber int    `json:"phoneNumber" binding:"omitempty"`
	}

	// IdentifyResponse encapsulates the response from this API
	IdentifyResponse struct {
		PrimaryContactID    int
		Emails              []string
		PhoneNumbers        []string
		SecondaryContactIDs []int
	}
)

// GetIdentifyErrorMessage is utility function to return a standard error response
func GetIdentifyErrorMessage(message string, trace string) gin.H {
	return gin.H{
		"error":    message,
		"trace_id": trace,
	}
}

// GetIdentifySuccessResponse is a utility function to return a standard success response
func GetIdentifySuccessResponse(response IdentifyResponse) gin.H {
	return gin.H{
		"contact": map[string]any{
			"primaryContactId":   response.PrimaryContactID,
			"emails":              response.Emails,
			"phoneNumbers":        response.PhoneNumbers,
			"secondaryContactIds": response.SecondaryContactIDs,
		},
	}
}
