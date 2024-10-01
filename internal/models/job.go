package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	URL               string             `bson:"url" json:"url"`
	Title             string             `bson:"title" json:"title"`
	Description       string             `bson:"description" json:"description"`
	Company           string             `bson:"company" json:"company"`
	Location          string             `bson:"location" json:"location"`
	EmploymentType    string             `bson:"employmentType" json:"employmentType"`
	PostingDate       time.Time          `bson:"postingDate" json:"postingDate"`
	ExpirationDate    time.Time          `bson:"expirationDate" json:"expirationDate"`
	IsActive          bool               `bson:"isActive" json:"isActive"`
	JobCategories     []string           `bson:"jobCategories" json:"jobCategories"`
	MustSkills        []string           `bson:"mustSkills" json:"mustSkills"`
	OptionalSkills    []string           `bson:"optionalSkills" json:"optionalSkills"`
	Salary            string             `bson:"salary" json:"salary"`
	YearsOfExperience int                `bson:"yearsOfExperience" json:"yearsOfExperience"`
	EducationLevel    string             `bson:"educationLevel" json:"educationLevel"`
	Benefits          []string           `bson:"benefits" json:"benefits"`
	CompanySize       int                `bson:"companySize" json:"companySize"`
	WorkCulture       string             `bson:"workCulture" json:"workCulture"`
	Remote            bool               `bson:"remote" json:"remote"`
	Languages         []string           `bson:"languages" json:"languages"`
}
