// Package middleware contains all the middlewares used for this project
package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/AadumKhor/bitespeed-backend-task/src/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

// ValidatePhoneNumber is the middleware to verify phone number
func ValidatePhoneNumber() gin.HandlerFunc {
	/*
		Implementing this since there is a validator for email in
		Gin (check models/identify.go) but not for mobile number
	*/
	return func(ctx *gin.Context) {
		// create a new trace for the entire flow
		// this trace will be passed to the next handler as well
		trace := xid.New().String()

		// extract the entire body to a variable as a []byte
		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.GetIdentifyErrorMessage(models.ErrMessageInvalidRequest, trace))
			ctx.Abort()
			return
		}

		// unmarshal the incoming request to our model
		var identifyRequest models.IdentifyRequest
		if err := json.Unmarshal(body, &identifyRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, models.GetIdentifyErrorMessage(models.ErrMessageInvalidRequest, trace))
			ctx.Abort()
			return
		}

		// if phone number is not empty, check if it's valid with a regex
		if identifyRequest.PhoneNumber != 0 {
			phoneNumberStr := strings.Trim(strings.ReplaceAll(fmt.Sprintf("%v", identifyRequest.PhoneNumber), ".", ""), "0")
			ok, _ := regexp.MatchString(models.RegexPhone, phoneNumberStr)
			if !ok {
				ctx.JSON(http.StatusBadRequest, models.GetIdentifyErrorMessage(models.ErrMessageInvalidPhoneNumber, trace))
				ctx.Abort()
				return
			}
		}

		// reset the same request body for next handler
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		// setting the same trace
		ctx.Set(models.TraceIDKey, trace)
		ctx.Next()
	}
}
