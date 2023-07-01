// Package handlers for all handlers in this project
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AadumKhor/bitespeed-backend-task/src/pkg/models"
	"github.com/gin-gonic/gin"
)

// HandleIdentify handles the `identify` route
func HandleIdentify(ctx *gin.Context) {
	trace := ctx.GetString(models.TraceIDKey)

	var identifyRequest models.IdentifyRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&identifyRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, models.GetIdentifyErrorMessage(models.ErrMessageInvalidRequest, trace))
		ctx.Abort()
		return
	}

	fmt.Printf("Getting request: %+v", identifyRequest)
}
