package repositories

import (
	// "encoding/json"
	// "errors"
	"log"
	"shuttle/models/entity"

	"github.com/jmoiron/sqlx"
)

type ChildernRepositoryInterface interface {
	FetchAllChilderns(id string) ([]entity.Student, error)
	FetchSpecChildern(id string) (entity.Student, error)
	UpdateChildern(student entity.Student, studentUUID string) error
}

type childernRepository struct {
	DB *sqlx.DB
}

func NewChildernRepository(DB *sqlx.DB) ChildernRepositoryInterface {
	return &childernRepository{
		DB: DB,
	}
}

func (repositories *childernRepository) FetchAllChilderns(id string) ([]entity.Student, error) {
	var childerns []entity.Student

	query := `
        SELECT 
        s.student_uuid,
        s.student_first_name,
        s.student_last_name,
        s.student_gender,
        s.student_grade,
        s.parent_uuid,
        s.school_uuid
    FROM students s
    WHERE s.parent_uuid = $1
    `

	rows, err := repositories.DB.Queryx(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var childern entity.Student

		if err := rows.Scan(
			&childern.UUID,
			&childern.FirstName,
			&childern.LastName,
			&childern.Gender,
			&childern.Grade,
			&childern.ParentUUID,
			&childern.SchoolUUID,
		); err != nil {
			return nil, err
		}

		childerns = append(childerns, childern)
	}

	return childerns, nil
}

func (repositories *childernRepository) FetchSpecChildern(id string) (entity.Student, error) {
	var childern entity.Student

	query := `
		SELECT 
            s.student_uuid,
            s.student_first_name,
            s.student_last_name,
            s.student_gender,
			s.student_address,
			s.student_pickup_point,
            s.student_grade,
            s.parent_uuid,
            s.school_uuid,
            sc.school_name
        FROM students s
        JOIN schools sc ON s.school_uuid = sc.school_uuid
        WHERE s.student_uuid = $1
	`

	err := repositories.DB.QueryRowx(query, id).Scan(
		&childern.UUID,
		&childern.FirstName,
		&childern.LastName,
		&childern.Gender,
		&childern.StudentAddress,
		&childern.StudentPickupPoint,
		&childern.Grade,
		&childern.ParentUUID,
		&childern.SchoolUUID,
		&childern.SchoolName,
	)
	log.Println("SchoolUUID:", childern.SchoolUUID)

	if err != nil {
		return entity.Student{}, err
	}

	return childern, nil
}

// 

func (repo *childernRepository) UpdateChildern(student entity.Student, studentUUID string) error {
	log.Println("Repository UpdateChildern started for ID:", studentUUID)

	query := `
		UPDATE students
		SET 
			student_first_name = $1, 
			student_last_name = $2, 
			student_gender = $3, 
			student_address = $4, 
			student_pickup_point = $5, 
			updated_at = NOW(), 
			updated_by = $6
		WHERE student_uuid = $7
	`

	_, err := repo.DB.Exec(query,
		student.FirstName,
		student.LastName,
		student.Gender,
		student.StudentAddress,
		student.StudentPickupPoint,
		student.UpdatedBy,
		studentUUID,
	)

	if err != nil {
		log.Println("Failed to execute update query:", err)
		return err
	}

	log.Println("Repository UpdateChildern completed successfully for ID:", studentUUID)
	return nil
}
