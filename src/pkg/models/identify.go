// Package models contains all models used in this project
package models

import "github.com/gin-gonic/gin"

type (
	// IdentifyRequest is the incoming request model
	IdentifyRequest struct {
		Email       string `json:"email" binding:"omitempty,email"`
		PhoneNumber int    `json:"phoneNumber" binding:"omitempty"`
	}

	// IdentifyResponse is the response we sent from server
	IdentifyResponse struct {
	}
)

// GetIdentifyErrorMessage is utility function to return a standard error response
func GetIdentifyErrorMessage(message string, trace string) gin.H {
	return gin.H{
		"error":    message,
		"trace_id": trace,
	}
}
