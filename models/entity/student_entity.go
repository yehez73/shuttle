package entity

import (
	"database/sql"

	"github.com/google/uuid"
)

type Student struct {
	StudentID        int64          `db:"student_id"`
	StudentUUID      uuid.UUID      `db:"student_uuid"`
	ParentUUID       uuid.UUID      `db:"parent_uuid"`
	SchoolUUID       uuid.UUID      `db:"school_uuid"`
	StudentFirstName string         `db:"student_first_name"`
	StudentLastName  string         `db:"student_last_name"`
	StudentGender    string         `db:"student_gender"`
	StudentGrade     string         `db:"student_grade"`
	CreatedAt        sql.NullTime   `db:"created_at"`
	CreatedBy        sql.NullString `db:"created_by"`
	UpdatedAt        sql.NullTime   `db:"updated_at"`
	UpdatedBy        sql.NullString `db:"updated_by"`
	DeletedAt        sql.NullTime   `db:"deleted_at"`
	DeletedBy        sql.NullString `db:"deleted_by"`
}
