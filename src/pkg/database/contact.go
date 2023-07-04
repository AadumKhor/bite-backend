package database

import (
	"context"

	"github.com/AadumKhor/bitespeed-backend-task/src/pkg/models"
)

const (
	queryWithOnlyEmail         = "email = ?"
	queryWithOnlyPhone         = "phone_number = ?"
	queryWithBothEmailAndPhone = "(email = ? OR phone_number = ?)"
)

// CheckEmailExists is a function to check if a row exists with given email
func (p PGStore) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	db := p.replica

	exists, err := db.NewSelect().
		Model((*models.Contact)(nil)).
		Where("email = ?", email).
		Where("deleted_at IS NULL").
		Exists(ctx)

	return exists, err
}

// CheckPhoneNumberExists is a function to check if a row exists with given number
func (p PGStore) CheckPhoneNumberExists(ctx context.Context, phone string) (bool, error) {
	db := p.replica

	exists, err := db.NewSelect().
		Model((*models.Contact)(nil)).
		Where("phone_number = ?", phone).
		Where("deleted_at IS NULL").
		Exists(ctx)

	return exists, err
}

// AddNewContact adds a new contact row with given data and returns the row as well
func (p PGStore) AddNewContact(ctx context.Context, contact *models.Contact) (*models.Contact, error) {
	db := p.master

	err := db.NewInsert().Model(contact).Returning("*").Scan(ctx)
	if err != nil {
		return nil, err
	}

	return contact, err
}

// GetPrimaryContact is
func (p PGStore) GetPrimaryContact(ctx context.Context, email string, phone string) ([]models.Contact, error) {
	db := p.replica
	contacts := []models.Contact{}

	query := db.NewSelect().Model(&contacts).
		Where("deleted_at IS NULL").
		Where("link_precedence = ?", models.Primary)

	if email != "" && phone != "0" {
		query = query.Where(queryWithBothEmailAndPhone, email, phone)
	} else if email != "" {
		query = query.Where(queryWithOnlyEmail, email)
	} else {
		query = query.Where(queryWithOnlyPhone, phone)
	}

	err := query.Order("created_at ASC").Scan(ctx)
	if err != nil {
		return []models.Contact{}, nil
	}

	return contacts, err
}

// UpdateExistingContact is
func (p PGStore) UpdateExistingContact(ctx context.Context, contact *models.Contact) error {
	db := p.master

	_, err := db.NewUpdate().Model(contact).Where("id = ?", contact.ID).Exec(ctx)

	return err
}

// GetExistingContactsLinkedID gets existing contact list with given linked ID
func (p PGStore) GetExistingContactsLinkedID(ctx context.Context, linkedID string) ([]models.Contact, error) {
	db := p.replica

	contacts := []models.Contact{}
	err := db.NewSelect().
		Model(&contacts).
		Where("linked_id = ? AND deleted_at IS NULL", linkedID).
		Order("created_at ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return contacts, nil
}
