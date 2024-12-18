package services

import (
	"database/sql"
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

	students, parent, err := service.studentRepository.FetchAllStudentsWithParents(offset, limit, sortField, sortDirection, schoolUUIDStr)
	if err != nil {
		return nil, 0, err
	}

	total, err := service.studentRepository.CountAllStudentsWithParents(schoolUUIDStr)
	if err != nil {
		return nil, 0, err
	}

	var studentsWithParents []dto.SchoolStudentParentResponseDTO

	for _, student := range students {
		var parentName string
		if student.ParentUUID == uuid.Nil {
			parentName = "N/A"
		} else {
			parentName = parent.FirstName + " " + parent.LastName
		}

		studentsWithParents = append(studentsWithParents, dto.SchoolStudentParentResponseDTO{
			StudentUUID:      student.StudentUUID.String(),
			ParentName:       parentName,
			StudentFirstName: student.StudentFirstName,
			StudentLastName:  student.StudentLastName,
			StudentGender:    dto.Gender(student.StudentGender),
			StudentGrade:     student.StudentGrade,
			Address:          parent.Address,
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
	if student.ParentUUID == uuid.Nil {
		parentName = "N/A"
	} else {
		parentName = parent.FirstName + " " + parent.LastName
	}

	return dto.SchoolStudentParentResponseDTO{
		StudentUUID:      student.StudentUUID.String(),
		ParentUUID:       student.ParentUUID.String(),
		ParentName:       parentName,
		ParentPhone:      parent.Phone,
		StudentFirstName: student.StudentFirstName,
		StudentLastName:  student.StudentLastName,
		StudentGender:    dto.Gender(student.StudentGender),
		StudentGrade:     student.StudentGrade,
		Address:          parent.Address,
		CreatedAt:        safeTimeFormat(student.CreatedAt),
		CreatedBy:        safeStringFormat(student.CreatedBy),
		UpdatedAt:        safeTimeFormat(student.UpdatedAt),
		UpdatedBy:        safeStringFormat(student.UpdatedBy),
	}, nil
}

func (service *StudentService) AddSchoolStudentWithParents(student dto.SchoolStudentParentRequestDTO, schoolUUID string, username string) error {
	var parentID uuid.UUID

	parentExists, err := service.userRepository.CheckEmailExist("", student.Parent.Email)
	if err != nil {
		return err
	}

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
		parentID, err = service.userRepository.FetchUUIDByEmail(student.Parent.Email)
		if err != nil {
			return transactionError
		}
	}

	newStudent := &entity.Student{
		StudentID:        time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
		StudentUUID:      uuid.New(),
		ParentUUID:       parentID,
		SchoolUUID:       *parseSafeUUID(schoolUUID),
		StudentFirstName: student.Student.StudentFirstName,
		StudentLastName:  student.Student.StudentLastName,
		StudentGender:    string(student.Student.StudentGender),
		StudentGrade:     student.Student.StudentGrade,
		CreatedBy:        sql.NullString{String: username, Valid: true},
	}

	err = service.studentRepository.SaveStudent(*newStudent)
	if err != nil {
		transactionError = err
		return transactionError
	}

	return nil
}

func (service *StudentService) UpdateSchoolStudentWithParents(id string, student dto.SchoolStudentParentRequestDTO, schoolUUID, username string) error {

	println("first line")

	studentUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	println("second line")

	_, err = service.userRepository.CheckEmailExist(id, student.Parent.Email)
	if err != nil {
		return err
	}

	println("third line")
	tx, err := service.userRepository.BeginTransaction()
	if err != nil {
		return err
	}

	println("fourth line")

	var transactionError error
	defer func() {
		if transactionError != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	println("fifth line")
	println("email: ", student.Parent.Email)

	var parentData entity.ParentDetails
	_, parentData, err = service.studentRepository.FetchSpecStudentWithParents(studentUUID, schoolUUID)
	if err != nil {
		transactionError = err
		return transactionError
	}

	println("sixth line")

	studentEntity := &entity.Student{
		StudentUUID:      studentUUID,
		ParentUUID:       parentData.UserUUID,
		SchoolUUID:       *parseSafeUUID(schoolUUID),
		StudentFirstName: student.Student.StudentFirstName,
		StudentLastName:  student.Student.StudentLastName,
		StudentGender:    string(student.Student.StudentGender),
		StudentGrade:     student.Student.StudentGrade,
		UpdatedBy:        sql.NullString{String: username, Valid: true},
	}

	
	println("sampe sini")

	err = service.studentRepository.UpdateStudent(*studentEntity)
	if err != nil {
		transactionError = err
		return transactionError
	}

	err = service.userService.UpdateUser(parentData.UserUUID.String(), student.Parent, username, nil)
	if err != nil {
		transactionError = err
		return transactionError
	}

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