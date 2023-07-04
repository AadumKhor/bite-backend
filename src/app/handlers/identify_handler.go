// Package handlers for all handlers in this project
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/AadumKhor/bitespeed-backend-task/src/pkg/models"
	"github.com/AadumKhor/bitespeed-backend-task/src/pkg/utils"
	"github.com/gin-gonic/gin"
)

// IdentifyHandler is the struct that calls the Handle function and contains relevant data
type IdentifyHandler struct {
	Store IdentifyHandlerStore
}

// IdentifyHandlerStore is the store that defines all DB methods required by this handler
type IdentifyHandlerStore interface {
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CheckPhoneNumberExists(ctx context.Context, phoneNumber string) (bool, error)
	AddNewContact(ctx context.Context, contact *models.Contact) (*models.Contact, error)
	GetPrimaryContact(ctx context.Context, email string, phone string) ([]models.Contact, error)
	UpdateExistingContact(ctx context.Context, contact *models.Contact) error
	GetExistingContactsLinkedID(ctx context.Context, linkedID string) ([]models.Contact, error)
}

// Utitlity function to send error response with code and response body
func handleError(ctx *gin.Context, status int, response gin.H) {
	ctx.JSON(status, response)
	ctx.Abort()
}

// Utility function to send success response with code and response body
func handleSuccess(ctx *gin.Context, response gin.H) {
	ctx.JSON(http.StatusOK, response)
}

// Checks if list does not contain the same item already & if either item is not empty string
func canAddItemToList(item string, existingItem string) bool {
	return item != existingItem && item != ""
}

// Handle handles the `identify` route
func (handler IdentifyHandler) Handle(ctx *gin.Context) {
	trace := ctx.GetString(models.TraceIDKey)
	backgroundContext := context.Background()

	// parse request
	var identifyRequest models.IdentifyRequest
	err := json.NewDecoder(ctx.Request.Body).Decode(&identifyRequest)
	if err != nil {
		handleError(ctx, http.StatusBadRequest, models.GetIdentifyErrorMessage(models.ErrMessageInvalidRequest, trace))
		return
	}

	// form the right phone number if its empty
	var phone string
	if identifyRequest.PhoneNumber == 0 {
		phone = ""
	} else {
		phone = fmt.Sprintf("%d", identifyRequest.PhoneNumber)
	}

	email := identifyRequest.Email
	if email == "" && phone == "" {
		handleError(ctx, http.StatusBadRequest, models.GetIdentifyErrorMessage(models.ErrMessageInvalidRequest, trace))
		return
	}

	// check if email already exists
	var phoneExists, emailExists bool
	if email == "" {
		emailExists = false
	} else {
		emailExists, err = handler.Store.CheckEmailExists(backgroundContext, email)
		if err != nil {
			utils.DefaultLogger.Error(err.Error())
			log.Default().Printf("error while checking if email exists for email: %s", email)
			handleError(ctx, http.StatusInternalServerError, models.GetIdentifyErrorMessage(models.ErrInternalServerError, trace))
			return
		}
	}

	// check if phone number already exists
	if phone == "" {
		phoneExists = false
	} else {
		phoneExists, err = handler.Store.CheckPhoneNumberExists(backgroundContext, phone)
		if err != nil {
			utils.DefaultLogger.Error(err.Error())
			log.Default().Printf("error while checking if phone number exists for number: %s", phone)
			handleError(ctx, http.StatusInternalServerError, models.GetIdentifyErrorMessage(models.ErrInternalServerError, trace))
			return
		}
	}

	// both don't exist, create a new record
	if !emailExists && !phoneExists {
		// add row to contact table
		contact, err := handler.createNewContact(backgroundContext, &identifyRequest)
		if err != nil {
			utils.DefaultLogger.Error(err.Error())
			log.Default().Printf("error while adding new contact")
			handleError(ctx, http.StatusInternalServerError, models.GetIdentifyErrorMessage(models.ErrInternalServerError, trace))
		}

		// form response with data from DB
		response := models.IdentifyResponse{
			PrimaryContactID:    int(contact.ID),
			Emails:              []string{*contact.Email},
			PhoneNumbers:        []string{*contact.PhoneNumber},
			SecondaryContactIDs: []int{},
		}

		handleSuccess(ctx, models.GetIdentifySuccessResponse(response))
		return
	}

	// either one or both exist
	primaryContacts, err := handler.Store.GetPrimaryContact(ctx, email, phone)
	if err != nil {
		utils.DefaultLogger.Error(err.Error())
		log.Default().Printf("error while fetching primary contact details for email: %s & phone: %s", email, phone)
		handleError(ctx, http.StatusInternalServerError, models.GetIdentifyErrorMessage(models.ErrInternalServerError, trace))
	}

	if len(primaryContacts) == 0 {
		panic(errors.New("logical error detected"))
	}

	// list is ordered by `created_at` hence first element will remain primary
	// rest will be updated and linked to this element
	primaryContact := primaryContacts[0]
	if len(primaryContacts) > 1 {
		// link contacts
		err := handler.linkContacts(primaryContacts, primaryContact, ctx)

		// return internal server error
		if err != nil {
			handleError(ctx, http.StatusInternalServerError, models.GetIdentifyErrorMessage(models.ErrInternalServerError, trace))
			return
		}
	}

	// get contact list linked to the primary
	linkedID := fmt.Sprintf("%d", primaryContact.ID)
	contacts, err := handler.Store.GetExistingContactsLinkedID(ctx, linkedID)
	emails, phoneNumbers := []string{*primaryContact.Email}, []string{*primaryContact.PhoneNumber}
	secondaryContactIDs := []int{}
	for _, contact := range contacts {
		if canAddItemToList(*contact.Email, *primaryContact.Email) {
			emails = append(emails, *contact.Email)
		}
		if canAddItemToList(*contact.PhoneNumber, *primaryContact.PhoneNumber) {
			phoneNumbers = append(phoneNumbers, *contact.PhoneNumber)
		}
		secondaryContactIDs = append(secondaryContactIDs, int(contact.ID))
	}

	// check if current data contains new information
	newEmail := !utils.StringExistsInList(email, emails)
	newPhone := !utils.StringExistsInList(phone, phoneNumbers)

	var secondaryContact *models.Contact
	if newEmail || newPhone {
		fmt.Println("Creating secondary contact")
		// if new information exist, create secondary contact
		tempSecondaryContact := models.Contact{
			PhoneNumber:    &phone,
			Email:          &email,
			LinkedID:       &primaryContact.ID,
			LinkPrecedence: models.Secondary,
		}
		secondaryContact, err = handler.Store.AddNewContact(ctx, &tempSecondaryContact)
		secondaryContactIDs = append(secondaryContactIDs, int(secondaryContact.ID))
	}

	// add current information for response
	if newEmail {
		emails = append(emails, *secondaryContact.Email)
	}
	if newPhone {
		phoneNumbers = append(phoneNumbers, *secondaryContact.PhoneNumber)
	}

	// form response from above list
	response := models.IdentifyResponse{
		PrimaryContactID:    int(primaryContact.ID),
		Emails:              emails,
		PhoneNumbers:        phoneNumbers,
		SecondaryContactIDs: secondaryContactIDs,
	}

	handleSuccess(ctx, models.GetIdentifySuccessResponse(response))
	return
}

func (handler IdentifyHandler) linkContacts(primaryContacts []models.Contact, primaryContact models.Contact, ctx *gin.Context) error {
	var err error
	for i := range primaryContacts {
		if i == 0 {
			continue
		}

		tempContact := primaryContacts[i]

		tempContact.LinkedID = &primaryContact.ID
		tempContact.LinkPrecedence = models.Secondary

		err = handler.Store.UpdateExistingContact(ctx, &tempContact)
		if err != nil {
			utils.DefaultLogger.Error(err.Error())
			log.Default().Printf("error while linking contact : %+v", tempContact)
			break
		}
	}
	return err
}

func (handler IdentifyHandler) createNewContact(ctx context.Context, identifyRequest *models.IdentifyRequest) (*models.Contact, error) {
	// create a new contact model
	var phone string
	// guard condition for since phone is by default an `int`
	// so default value will be 0 in the DB which is wrong
	if identifyRequest.PhoneNumber != 0 {
		phone = fmt.Sprintf("%d", identifyRequest.PhoneNumber)
	}
	contact := models.Contact{
		PhoneNumber:    &phone,
		Email:          &identifyRequest.Email,
		LinkPrecedence: models.Primary,
	}

	// add it to db
	// result is also of type *models.Contact
	result, err := handler.Store.AddNewContact(ctx, &contact)
	if err != nil {
		return nil, err
	}

	return result, err
}
