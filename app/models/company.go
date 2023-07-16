package models

import (
	"errors"

	"github.com/7leven7/xm-company-crud/app/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	maxNameLength        = 15
	maxDescriptionLength = 3000
)

type Company struct {
	gorm.Model
	ID          string `json:"id" gorm:"type:uuid;primaryKey;not null"`
	Name        string `json:"name" gorm:"size:15;not null;unique"`
	Description string `json:"description,omitempty" gorm:"size:3000"`
	Employees   int    `json:"employees" gorm:"not null"`
	Registered  bool   `json:"registered" gorm:"not null"`
	Type        string `json:"type" gorm:"not null;check:type IN ('Corporations', 'NonProfit', 'Cooperative', 'Sole Proprietorship')"`
}

func (c *Company) Validate() error {
	if _, err := uuid.Parse(c.ID); err != nil {
		return errors.New("ID is required and must be a valid UUID")
	}

	if len(c.Name) > maxNameLength {
		return errors.New("name must be maximum 15 characters")
	}

	if len(c.Description) > maxDescriptionLength {
		return errors.New("description must be maximum 3000 characters")
	}

	if c.Employees <= 0 {
		return errors.New("employees must be greater than zero")
	}

	return nil
}

func (c *Company) SaveCompany() (*Company, error) {
	c.ID = uuid.New().String()
	if err := c.Validate(); err != nil {
		return nil, err
	}

	err := database.DB.Create(&c).Error

	if err != nil {
		return &Company{}, err
	}

	return c, nil
}
func (c *Company) UpdateCompany(updatedCompany *Company) (*Company, error) {
	if updatedCompany.Name != "" {
		c.Name = updatedCompany.Name
	}

	if updatedCompany.Description != "" {
		c.Description = updatedCompany.Description
	}

	if updatedCompany.Employees != 0 {
		c.Employees = updatedCompany.Employees
	}

	c.Registered = updatedCompany.Registered

	if updatedCompany.Type != "" {
		c.Type = updatedCompany.Type
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	err := database.DB.Save(&c).Error
	if err != nil {
		return &Company{}, err
	}

	return c, nil
}

func (c *Company) DeleteCompany() error {
	err := database.DB.Delete(&c).Error
	return err
}

func GetCompanyByID(id string) (*Company, error) {
	var company Company
	err := database.DB.Where("id = ?", id).First(&company).Error

	if err != nil {
		return nil, err
	}

	return &company, nil
}
