package fakers

import (
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/puskipus/e-commerce/app/models"
	"gorm.io/gorm"
)

func UserFaker(db *gorm.DB) *models.User {
	return &models.User{
		ID:            uuid.New().String(),
		FirstName:     faker.FirstName(),
		LastName:      faker.LastName(),
		Email:         faker.Email(),
		Password:      "123",
		RememberToken: "",
		CreatedAt:     time.Time{},
		UpdateAt:      time.Time{},
		DeletedAt:     gorm.DeletedAt{},
	}
}
