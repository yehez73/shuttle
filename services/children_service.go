package services

import (
	// "strings"
	// "time"
	// "database/sql"
	"encoding/json"
	"log"
	// "errors"
	"database/sql"
	// "shuttle/errors"
	// "encoding/json"
	// "fmt"
	"shuttle/models/dto"
	"shuttle/models/entity"

	// "shuttle/models/entity"
	"shuttle/repositories"

	// "github.com/google/uuid"
	// "github.com/jmoiron/sqlx"
)

type ChildernServiceInterface interface {
	GetAllChilderns(id string) ([]dto.StudentResponseDTO, int, error)
	GetSpecChildern(id string) (dto.StudentResponseDTO, error)
	UpdateChildern(id string, req dto.StudentRequestByParentDTO, username string) error
}

type ChildernService struct {
	ChildernRepository repositories.ChildernRepositoryInterface
}

func NewChildernService(childernRepository repositories.ChildernRepositoryInterface) ChildernServiceInterface {
	return &ChildernService{
		ChildernRepository: childernRepository,
	}
}

func (service *ChildernService) GetAllChilderns(id string) ([]dto.StudentResponseDTO, int, error) {
    childerns, err := service.ChildernRepository.FetchAllChilderns(id)
    if err != nil {
        log.Println("Error fetching students from repository:", err)
        return nil, 0, err
    }

    var childernsDTO []dto.StudentResponseDTO

    for _, childern := range childerns {
        childernsDTO = append(childernsDTO, dto.StudentResponseDTO{
            UUID:          childern.UUID.String(),
            FirstName:     childern.FirstName,
            LastName:      childern.LastName,
            Grade:         childern.Grade,
            Gender:        childern.Gender,
            SchoolUUID:    childern.SchoolUUID.String(),
			SchoolName: childern.SchoolName,
        })
    }

    total := len(childernsDTO)

    return childernsDTO, total, nil
}


func (service *ChildernService) GetSpecChildern(id string) (dto.StudentResponseDTO, error) {
	childern, err := service.ChildernRepository.FetchSpecChildern(id)
	if err != nil {
		return dto.StudentResponseDTO{}, err
	}
	Address := ""
		if childern.StudentAddress.Valid {
			Address = childern.StudentAddress.String
		}
	
		// Memeriksa PickupPoint
		var pickupPoint string
			if childern.StudentPickupPoint.Valid {
				pickupPoint = childern.StudentPickupPoint.String // Ini sudah berupa string
		}

	studentDTO := dto.StudentResponseDTO{
		UUID:       childern.UUID.String(),
		FirstName:  childern.FirstName,
		LastName:   childern.LastName,
		Gender:     childern.Gender,
		StudentAddress:     Address,
		PickupPoint: pickupPoint,
		Grade:      childern.Grade,
		SchoolUUID: childern.SchoolUUID.String(),
	}

	return studentDTO, nil
}

func (service *ChildernService) UpdateChildern(id string, req dto.StudentRequestByParentDTO, username string) error {
	log.Println("Service UpdateChildern started for ID:", id)

	// Serialize pickup point ke JSON
	var pickupPointJSON []byte
	var err error
	if req.StudentPickupPoint != nil {
		pickupPointJSON, err = json.Marshal(req.StudentPickupPoint)
		if err != nil {
			log.Println("Failed to serialize StudentPickupPoint:", err)
			return err
		}
	}

	// Buat entity Student dari DTO
	student := entity.Student{
		FirstName:         req.StudentFirstName,
		LastName:          req.StudentLastName,
		Gender:            string(req.StudentGender),
		StudentAddress:    sql.NullString{String: req.StudentAddress, Valid: req.StudentAddress != ""},
		StudentPickupPoint: sql.NullString{String: string(pickupPointJSON), Valid: req.StudentPickupPoint != nil},
		UpdatedBy:         sql.NullString{String: username, Valid: username != ""},
	}

	log.Println("Prepared Student entity for update:", student)

	// Panggil repository untuk update data
	err = service.ChildernRepository.UpdateChildern(student, id)
	if err != nil {
		log.Println("Failed to update student in repository:", err)
		return err
	}

	log.Println("Service UpdateChildern completed successfully for ID:", id)
	return nil
}
