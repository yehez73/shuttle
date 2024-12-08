package services

import (
	"database/sql"
	"time"

	"shuttle/models/dto"
	"shuttle/models/entity"
	"shuttle/repositories"

	"github.com/google/uuid"
)

type SchoolServiceInterface interface {
	GetAllSchools(page, limit int, sortField, sortDirection string) ([]dto.SchoolResponseDTO, int, error)
	GetSpecSchool(uuid string) (dto.SchoolResponseDTO, error)
	AddSchool(req dto.SchoolRequestDTO, username string) error
	UpdateSchool(id string, req dto.SchoolRequestDTO, username string) error
	DeleteSchool(id string, username string) error
}

type SchoolService struct {
	schoolRepository repositories.SchoolRepositoryInterface
}

func NewSchoolService(schoolRepository repositories.SchoolRepositoryInterface) SchoolService {
	return SchoolService{
		schoolRepository: schoolRepository,
	}
}

func (service *SchoolService) GetAllSchools(page, limit int, sortField, sortDirection string) ([]dto.SchoolResponseDTO, int, error) {
	offset := (page - 1) * limit

	schools, admin, err := service.schoolRepository.FetchAllSchools(offset, limit, sortField, sortDirection)
	if err != nil {
		return nil, 0, err
	}

	total, err := service.schoolRepository.CountSchools()
	if err != nil {
		return nil, 0, err
	}

	var schoolsDTO []dto.SchoolResponseDTO
	for _, school := range schools {

		var adminFullName string

		if admin[school.UUID.String()].SchoolUUID == uuid.Nil {
			adminFullName = "N/A"
		} else if admin[school.UUID.String()].SchoolUUID != uuid.Nil {
			adminFullName = admin[school.UUID.String()].FirstName + " " + admin[school.UUID.String()].LastName
		}

		schoolsDTO = append(schoolsDTO, dto.SchoolResponseDTO{
			UUID:        school.UUID.String(),
			Name:        school.Name,
			AdminName:   adminFullName,
			Address:     school.Address,
			Contact:     school.Contact,
			Email:       school.Email,
		})
	}

	return schoolsDTO, total, nil
}

func (service *SchoolService) GetSpecSchool(id string) (dto.SchoolResponseDTO, error) {
	school, admin, err := service.schoolRepository.FetchSpecSchool(id)
	if err != nil {
		return dto.SchoolResponseDTO{}, err
	}

	var userUUID, userFullName string
	if admin.SchoolUUID == uuid.Nil {
		userUUID = "N/A"
		userFullName = "N/A"
	} else {
		userUUID = admin.UserUUID.String()
		userFullName = admin.FirstName + " " + admin.LastName
	}

	schoolDTO := dto.SchoolResponseDTO{
		UUID:        school.UUID.String(),
		Name:        school.Name,
		AdminUUID:   userUUID,
		AdminName:   userFullName,
		Address:     school.Address,
		Contact:     school.Contact,
		Email:       school.Email,
		Description: school.Description,
		CreatedAt:   safeTimeFormat(school.CreatedAt),
		CreatedBy:   safeStringFormat(school.CreatedBy),
		UpdatedAt:   safeTimeFormat(school.UpdatedAt),
		UpdatedBy:   safeStringFormat(school.UpdatedBy),
	}

	return schoolDTO, nil
}

func (service *SchoolService) AddSchool(req dto.SchoolRequestDTO, username string) error {
	school := entity.School{
		ID:          time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
		UUID:        uuid.New(),
		Name:        req.Name,
		Address:     req.Address,
		Contact:     req.Contact,
		Email:       req.Email,
		Description: req.Description,
		CreatedBy:   toNullString(username),
	}

	if err := service.schoolRepository.SaveSchool(school); err != nil {
		return err
	}

	return nil
}

func (service *SchoolService) UpdateSchool(id string, req dto.SchoolRequestDTO, username string) error {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	school := entity.School{
		UUID:        parsedUUID,
		Name:        req.Name,
		Address:     req.Address,
		Contact:     req.Contact,
		Email:       req.Email,
		Description: req.Description,
		UpdatedAt:   toNullTime(time.Now()),
		UpdatedBy:   toNullString(username),
	}

	if err := service.schoolRepository.UpdateSchool(school); err != nil {
		return err
	}

	return nil
}

func (service *SchoolService) DeleteSchool(id string, username string) error {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	school := entity.School{
		UUID:      parsedUUID,
		DeletedAt: toNullTime(time.Now()),
		DeletedBy: toNullString(username),
	}

	if err := service.schoolRepository.DeleteSchool(school); err != nil {
		return err
	}

	return nil
}

func safeStringFormat(s sql.NullString) string {
	if !s.Valid || s.String == "" {
		return "N/A"
	}
	return s.String
}

func safeTimeFormat(t sql.NullTime) string {
	if !t.Valid {
		return "N/A"
	}
	return t.Time.Format(time.RFC3339)
}

func toNullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{String: value, Valid: false}
	}
	return sql.NullString{String: value, Valid: true}
}

func toNullTime(value time.Time) sql.NullTime {
	return sql.NullTime{Time: value, Valid: !value.IsZero()}
}
