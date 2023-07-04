package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

// LinkPrecedence is used as an enum for DB values
type LinkPrecedence string

const (
	// Secondary is one of the value for the enum LinkPrecedence
	Secondary LinkPrecedence = "secondary"

	// Primary is one of the values for the enum LinkPrecedence
	Primary LinkPrecedence = "primary"
)

// Contact struct represents the table contact
type Contact struct {
	bun.BaseModel `bun:"contact,alias:c"`
	ID             int64          `bun:"id,autoincrement,notnull"`
	PhoneNumber    *string        `bun:"phone_number,nullzero"`
	Email          *string        `bun:"email,nullzero"`
	LinkedID       *int64         `bun:"linked_id,nullzero"`
	LinkPrecedence LinkPrecedence `bun:"link_precedence,type:link_precedence"`
	CreatedAt      time.Time      `bun:"created_at,default:current_timestamp()"`
	UpdatedAt      time.Time      `bun:"updated_at,default:current_timestamp()"`
	DeletedAt      *time.Time     `bun:"deleted_at,nullzero"`
}

// BeforeUpdate is a helper function tha replaces the need for a trigger on the DB
func (c *Contact) BeforeUpdate(_ *bun.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}

// AfterDelete is a helper function that replaces the need for a trigger on the DB
func (c *Contact) AfterDelete(db *bun.DB) error {
	now := time.Now()
	c.DeletedAt = &now
	_, err := db.NewUpdate().Model(c).
		Set("deleted_at = ?", c.DeletedAt).
		Where("id = ?", c.ID).
		Exec(context.Background())
	return err
}
