package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"shuttle/errors"
	"shuttle/models/dto"
	"shuttle/models/entity"
	"shuttle/repositories"
	"time"

	"github.com/google/uuid"
)

type StudentServiceInterface interface {
	GetAllStudentsWithParents(page int, limit int, sortField string, sortDirection string, schoolUUIDStr string) ([]dto.SchoolStudentParentResponseDTO, int, error)
	GetSpecStudentWithParents(id, schoolUUIDStr string) (dto.SchoolStudentParentResponseDTO, error)
	AddSchoolStudentWithParents(student dto.SchoolStudentParentRequestDTO, schoolUUID string, username string) error
	UpdateSchoolStudentWithParents(id string, student dto.SchoolStudentParentRequestDTO, schoolUUID, username string) error
	DeleteSchoolStudentWithParentsIfNeccessary(id, schoolUUID, username string) error
}

type StudentService struct {
	userService       UserServiceInterface
	studentRepository repositories.StudentRepositoryInterface
	userRepository    repositories.UserRepositoryInterface
}

func NewStudentService(studentRepository repositories.StudentRepositoryInterface, userService UserServiceInterface, userRepository repositories.UserRepositoryInterface) StudentService {
	return StudentService{
		userService:       userService,
		studentRepository: studentRepository,
		userRepository:    userRepository,
	}
}

func (service *StudentService) GetAllStudentsWithParents(page int, limit int, sortField string, sortDirection string, schoolUUIDStr string) ([]dto.SchoolStudentParentResponseDTO, int, error) {
	offset := (page - 1) * limit

	students, parents, err := service.studentRepository.FetchAllStudentsWithParents(offset, limit, sortField, sortDirection, schoolUUIDStr)
	if err != nil {
		return nil, 0, err
	}

	total, err := service.studentRepository.CountAllStudentsWithParents(schoolUUIDStr)
	if err != nil {
		return nil, 0, err
	}

	var studentsWithParents []dto.SchoolStudentParentResponseDTO

	for i, student := range students {
		var parentName string
		if !student.ParentUUID.Valid || student.ParentUUID.String == "" {
			parentName = "N/A"
		} else {
			parentName = parents[i].FirstName + " " + parents[i].LastName
		}

		studentsWithParents = append(studentsWithParents, dto.SchoolStudentParentResponseDTO{
			StudentUUID:      student.UUID.String(),
			ParentName:       parentName,
			StudentFirstName: student.FirstName,
			StudentLastName:  student.LastName,
			StudentGender:    dto.Gender(student.Gender),
			StudentGrade:     student.Grade,
			Address:          parents[i].Address,
			CreatedAt:        safeTimeFormat(student.CreatedAt),
		})
	}

	return studentsWithParents, total, nil
}

func (service *StudentService) GetSpecStudentWithParents(id, schoolUUIDStr string) (dto.SchoolStudentParentResponseDTO, error) {
    studentUUID, err := uuid.Parse(id)
    if err != nil {
        return dto.SchoolStudentParentResponseDTO{}, err
    }

    student, parent, err := service.studentRepository.FetchSpecStudentWithParents(studentUUID, schoolUUIDStr)
    if err != nil {
        return dto.SchoolStudentParentResponseDTO{}, err
    }

    var parentName string
    if !student.ParentUUID.Valid || student.ParentUUID.String == "" {
        parentName = "N/A"
    } else {
        parentName = parent.FirstName + " " + parent.LastName
    }

    Address := ""
    if student.StudentAddress.Valid {
        Address = student.StudentAddress.String
    }

    var pickupPoint string
    if student.StudentPickupPoint.Valid {
        pickupPoint = student.StudentPickupPoint.String
    }

    return dto.SchoolStudentParentResponseDTO{
        StudentUUID:      student.UUID.String(),
        ParentUUID:       student.ParentUUID.String,
        ParentName:       parentName,
		ParentFirstName: parent.FirstName,
		ParentlastName: parent.LastName,
        ParentPhone:      parent.Phone,
		ParentEmail: student.UserEmail,
        ParentUsername:   student.UserUsername,  // Kirimkan username orang tua
        StudentFirstName: student.FirstName,
        StudentLastName:  student.LastName,
        StudentGender:    dto.Gender(student.Gender),
        StudentGrade:     student.Grade,
        Address:          Address,
        PickupPoint:      pickupPoint,
        CreatedAt:        safeTimeFormat(student.CreatedAt),
        CreatedBy:        safeStringFormat(student.CreatedBy),
        UpdatedAt:        safeTimeFormat(student.UpdatedAt),
        UpdatedBy:        safeStringFormat(student.UpdatedBy),
    }, nil
}

func (service *StudentService) AddSchoolStudentWithParents(student dto.SchoolStudentParentRequestDTO, schoolUUID string, username string) error {
	var parentID uuid.UUID

	// Periksa apakah email parent sudah ada di database
	parentExists, err := service.userRepository.CheckEmailExist("", student.Parent.Email)
	if err != nil {
		return err
	}

	// Mulai transaksi
	tx, err := service.userRepository.BeginTransaction()
	if err != nil {
		return err
	}

	var transactionError error
	defer func() {
		if transactionError != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if !parentExists {
		// Jika parent belum ada, tambahkan data parent baru
		newParent := &dto.UserRequestsDTO{
			Username:  student.Parent.Username,
			FirstName: student.Parent.FirstName,
			LastName:  student.Parent.LastName,
			Gender:    student.Parent.Gender,
			Email:     student.Parent.Email,
			Password:  student.Parent.Password,
			Role:      dto.Role(entity.Parent),
			RoleCode:  "P",
			Phone:     student.Parent.Phone,
			Address:   student.Parent.Address,
		}

		parentID, err = service.userService.AddUser(*newParent, username)
		if err != nil {
			transactionError = err
			return transactionError
		}
	} else {
		// Jika parent sudah ada, ambil UUID-nya
		parentID, err = service.userRepository.FetchUUIDByEmail(student.Parent.Email)
		if err != nil {
			transactionError = err
			return transactionError
		}
	}

	// Proses pickup point menjadi JSON jika perlu
	var pickupPointJSON []byte
	if student.Student.StudentPickupPoint != nil {
		pickupPointJSON, err = json.Marshal(student.Student.StudentPickupPoint)
		if err != nil {
			transactionError = err
			return transactionError
		}
	}

	// Buat data siswa baru
	newStudent := &entity.Student{
		ID:        time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
		UUID:      uuid.New(),
		ParentUUID: sql.NullString{String: parentID.String(), Valid: true},
		SchoolUUID:       *parseSafeUUID(schoolUUID),
		FirstName: student.Student.StudentFirstName,
		LastName:  student.Student.StudentLastName,
		Gender:    string(student.Student.StudentGender),
		Grade:     student.Student.StudentGrade,
		StudentAddress:   sql.NullString{String: student.Student.StudentAddress, Valid: true},
		StudentPickupPoint: sql.NullString{String: string(pickupPointJSON), Valid: true},
		CreatedBy:        sql.NullString{String: username, Valid: true},
	}

	// Simpan data siswa
	err = service.studentRepository.SaveStudent(*newStudent)
	if err != nil {
		transactionError = err
		return transactionError
	}

	return nil
}

func (service *StudentService) UpdateSchoolStudentWithParents(id string, student dto.StudentRequestDTO, schoolUUID, username string) error {
	studentUUID, err := uuid.Parse(id)
	if err != nil {
		log.Println("ERROR: Failed to parse student UUID:", err)
		return err
	}
	log.Println("INFO: Parsed student UUID:", studentUUID)

	// Proses pickup point menjadi JSON
	pickupPointJSON, err := json.Marshal(student.StudentPickupPoint)
	if err != nil {
		log.Println("ERROR: Failed to marshal pickup point:", err)
		return err
	}
	log.Println("INFO: Pickup point successfully marshaled:", string(pickupPointJSON))

	// Update data siswa
	studentEntity := entity.Student{
		UUID:             studentUUID,
		SchoolUUID:       *parseSafeUUID(schoolUUID),
		FirstName:        student.StudentFirstName,
		LastName:         student.StudentLastName,
		Gender:           string(student.StudentGender),
		Grade:            student.StudentGrade,
		StudentAddress:   sql.NullString{String: student.StudentAddress, Valid: true},
		StudentPickupPoint: sql.NullString{String: string(pickupPointJSON), Valid: true},
		UpdatedBy:        sql.NullString{String: username, Valid: true},
	}
	log.Println("INFO: Preparing student entity for update:", studentEntity)

	err = service.studentRepository.UpdateStudent(studentEntity)
	if err != nil {
		log.Println("ERROR: Failed to update student data:", err)
		return err
	}
	log.Println("INFO: Student data updated successfully")

	return nil
}

func (service *StudentService) DeleteSchoolStudentWithParentsIfNeccessary(id, schoolUUID, username string) error {
	studentUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	_, _,err = service.studentRepository.FetchSpecStudentWithParents(studentUUID, schoolUUID)
	if err != nil {
		return errors.New("student not found", 404)	
	}

	return service.studentRepository.DeleteStudentWithParents(studentUUID, schoolUUID, username)
}
