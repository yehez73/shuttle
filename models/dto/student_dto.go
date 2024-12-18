package dto

import ()

type StudentRequestDTO struct {
	StudentFirstName string `json:"student_first_name" validate:"required"`
	StudentLastName  string `json:"student_last_name" validate:"required"`
	StudentGender    Gender `json:"student_gender" validate:"required"`
	StudentGrade     string `json:"student_grade" validate:"required"`
}

type SchoolStudentParentRequestDTO struct {
	Student StudentRequestDTO `json:"student" validate:"required"`
	Parent  UserRequestsDTO   `json:"parent" validate:"required"`
}

type SchoolStudentParentResponseDTO struct {
	StudentUUID      string `json:"student_uuid"`
	ParentUUID       string `json:"parent_uuid,omitempty"`
	ParentName       string `json:"parent_name"`
	ParentPhone      string `json:"parent_phone"`
	StudentFirstName string `json:"student_first_name"`
	StudentLastName  string `json:"student_last_name"`
	StudentGender    Gender `json:"student_gender"`
	StudentGrade     string `json:"student_grade"`
	Address          string `json:"student_address"`
	CreatedAt        string `json:"created_at,omitempty"`
	CreatedBy        string `json:"created_by,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
	UpdatedBy        string `json:"updated_by,omitempty"`
}
