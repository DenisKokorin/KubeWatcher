package models

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	UUID               uuid.UUID `json:"uuid,omitempty"`
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	Position           string    `json:"position"`
	Team               string    `json:"team"`
	Role               string    `json:"role"`
	HireDate           time.Time `json:"hire_date"`
	ProfileDocumentURL *string   `json:"profile_document_url,omitempty"`
	HashedPassword     string    `json:"hashed_password,omitempty"`
}

type EmployeeQuery struct {
	Search    string `json:"search" query:"search"`
	Team      string `json:"team" query:"team"`
	Position  string `json:"position" query:"position"`
	Role      string `json:"role" query:"role"`
	SortBy    string `json:"sortBy" query:"sortBy"`
	SortOrder string `json:"sortOrder" query:"sortOrder"`
	Page      int    `json:"page" query:"page"`
	Size      int    `json:"size" query:"size"`
}

type EmployeeList struct {
	Employees []*Employee `json:"employees"`
	Total     int         `json:"total"`
	Page      int         `json:"page"`
	Size      int         `json:"size"`
}

type Review struct {
	ID                  uuid.UUID `json:"id,omitempty"`
	EmployeeID          uuid.UUID `json:"employee_id,omitempty"`
	ReviewerID          uuid.UUID `json:"reviewer_id,omitempty"`
	Period              string    `json:"period"`
	Rating              int       `json:"rating"`
	Comments            string    `json:"comments"`
	Goals               []string  `json:"goals"`
	Strengths           []string  `json:"strengths"`
	AreasForImprovement []string  `json:"areas_for_improvement"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type ReviewRequest struct {
	EmployeeID          string   `json:"employee_id" binding:"required"`
	ReviewerID          string   `json:"reviewer_id" binding:"required"`
	Period              string   `json:"period" binding:"required"`
	Rating              int      `json:"rating" binding:"required,min=1,max=5"`
	Comments            string   `json:"comments"`
	Goals               []string `json:"goals"`
	Strengths           []string `json:"strengths"`
	AreasForImprovement []string `json:"areas_for_improvement"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type EmployeeWithReviews struct {
	Employee *Employee `json:"employee"`
	Reviews  []*Review `json:"reviews"`
}

type TeamStats struct {
	Team          string  `json:"team"`
	EmployeeCount int     `json:"employee_count"`
	AverageRating float64 `json:"average_rating"`
	ReviewsCount  int     `json:"reviews_count"`
}
