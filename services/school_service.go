package services

import (
	"database/sql"
	"strings"
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
	DeleteSchool(id, username, adminUUID string) error
}

type SchoolService struct {
	schoolRepository repositories.SchoolRepositoryInterface
	userRepository   repositories.UserRepositoryInterface
}

func NewSchoolService(schoolRepository repositories.SchoolRepositoryInterface, userRepository repositories.UserRepositoryInterface) SchoolService {
	return SchoolService{
		schoolRepository: schoolRepository,
		userRepository:   userRepository,
	}
}

func (service *SchoolService) GetAllSchools(page, limit int, sortField, sortDirection string) ([]dto.SchoolResponseDTO, int, error) {
	offset := (page - 1) * limit

	// Fetch data schools dan admin
	schools, adminMap, err := service.schoolRepository.FetchAllSchools(offset, limit, sortField, sortDirection)
	if err != nil {
		return nil, 0, err
	}

	// Hitung total schools
	total, err := service.schoolRepository.CountSchools()
	if err != nil {
		return nil, 0, err
	}

	// Convert schools dan admin menjadi DTO
	var schoolsDTO []dto.SchoolResponseDTO
	for _, school := range schools {
		admins := adminMap[school.UUID.String()] // Ambil slice admin berdasarkan UUID sekolah

		// Buat string adminFullName (gabungkan nama admin)
		var adminFullName string
		if len(admins) == 0 {
			adminFullName = "N/A"
		} else {
			var names []string
			for _, admin := range admins {
				names = append(names, admin.FirstName+" "+admin.LastName)
			}
			adminFullName = strings.Join(names, ", ") // Gabungkan nama admin dengan koma
		}

		// Masukkan data ke DTO
		schoolsDTO = append(schoolsDTO, dto.SchoolResponseDTO{
			UUID:      school.UUID.String(),
			Name:      school.Name,
			AdminName: adminFullName,
			Address:   school.Address,
			Contact:   school.Contact,
			Email:     school.Email,
		})
	}

	return schoolsDTO, total, nil
}

func (service *SchoolService) GetSpecSchool(id string) (dto.SchoolResponseDTO, error) {
	school, admins, err := service.schoolRepository.FetchSpecSchool(id)
	if err != nil {
		return dto.SchoolResponseDTO{}, err
	}

	var adminNames, adminUUIDs []string
	for _, admin := range admins {
		userFullName := "N/A"
		if admin.UserUUID != uuid.Nil {
			userFullName = admin.FirstName + " " + admin.LastName
		}

		adminUUID := "N/A"
		if admin.UserUUID != uuid.Nil {
			adminUUID = admin.UserUUID.String()
		}

		adminUUIDs = append(adminUUIDs, adminUUID)
		adminNames = append(adminNames, userFullName)
	}

	adminUUIDsStr := strings.Join(adminUUIDs, ", ")
	adminNamesStr := strings.Join(adminNames, ", ")

	schoolDTO := dto.SchoolResponseDTO{
		UUID:        school.UUID.String(),
		Name:        school.Name,
		AdminUUID:   adminUUIDsStr,
		AdminName:   adminNamesStr,
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

func (service *SchoolService) DeleteSchool(id, username, adminUUID string) error {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	// Handle multiple admin UUIDs deletion
	if adminUUID != "N/A" && adminUUID != "" {
		uuidList := strings.Split(adminUUID, ", ")

		tx, err := service.userRepository.BeginTransaction()
		if err != nil {
			return err
		}
		var transactionErr error
		defer func() {
			if transactionErr != nil {
				tx.Rollback()
			} else {
				transactionErr = tx.Commit()
			}
		}()

		for _, uuids := range uuidList {
			parsedAdminUUID, err := uuid.Parse(uuids)
			if err != nil {
				continue
			}

			if err := service.userRepository.DeleteSchoolAdmin(tx, parsedAdminUUID, username); err != nil {
				return err
			}
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
	} else {
		tx, err := service.userRepository.BeginTransaction()
		if err != nil {
			return err
		}
		var transactionErr error
		defer func() {
			if transactionErr != nil {
				tx.Rollback()
			} else {
				transactionErr = tx.Commit()
			}
		}()

		// Delete the school
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
