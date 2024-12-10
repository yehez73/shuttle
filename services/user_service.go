package services

import (
	"database/sql"
	"fmt"
	"strconv"

	"shuttle/databases"
	"shuttle/errors"
	"shuttle/logger"
	"shuttle/models/dto"
	"shuttle/models/entity"
	"shuttle/repositories"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserServiceInterface interface {
	GetAllSuperAdmin(page int, limit int, sortField string, sortDirection string) ([]dto.UserResponseDTO, int, error)
	GetSpecSuperAdmin(uuid string) (dto.UserResponseDTO, error)
	GetAllSchoolAdmin(page int, limit int, sortField string, sortDirection string) ([]dto.UserResponseDTO, int, error)
	GetSpecSchoolAdmin(uuid string) (dto.UserResponseDTO, error)

	GetAllDriverFromAllSchools(page int, limit int, sortField string, sortDirection string) ([]dto.UserResponseDTO, error)
	GetAllDriverForPermittedSchool(page int, limit int, sortField string, sortDirection string, schoolUUID string) ([]dto.UserResponseDTO, int, error)
	GetSpecDriverFromAllSchools(uuid string) (dto.UserResponseDTO, error)

	AddUser(user entity.User, user_name string) (uuid.UUID, error)
	UpdateUser(id string, user dto.UserRequestsDTO, user_name string, file []byte) error

	DeleteSuperAdmin(id string, user_name string) error
	DeleteSchoolAdmin(id string, user_name string) error
	DeleteDriver(id string, user_name string) error

	GetSpecUser(id string) (entity.User, error)
	GetSpecUserWithDetails(id string) (entity.User, error)

	CheckPermittedSchoolAccess(userUUID string) (string, error)
}

type UserService struct {
	userRepository repositories.UserRepositoryInterface
}

func NewUserService(userRepository repositories.UserRepositoryInterface) UserService {
	return UserService{
		userRepository: userRepository,
	}
}

func (service *UserService) GetAllSuperAdmin(page int, limit int, sortField, sortDirection string) ([]dto.UserResponseDTO, int, error) {
	offset := (page - 1) * limit

	users, err := service.userRepository.FetchAllSuperAdmins(offset, limit, sortField, sortDirection)
	if err != nil {
		return nil, 0, err
	}

	total, err := service.userRepository.CountSuperAdmin()
	if err != nil {
		return nil, 0, err
	}

	var usersDTO []dto.UserResponseDTO

	for _, user := range users {
		userDTO := dto.UserResponseDTO{
			UUID:       user.UUID.String(),
			Username:   user.Username,
			Email:      user.Email,
			Status:     user.Status,
			LastActive: safeTimeFormat(user.LastActive),
			CreatedAt:  safeTimeFormat(user.CreatedAt),
		}
		if details, ok := user.Details.(entity.SuperAdminDetails); ok {
			userDTO.Details = dto.SuperAdminDetailsResponseDTO{
				FirstName: details.FirstName,
				LastName:  details.LastName,
				Gender:    dto.Gender(details.Gender),
				Phone:     details.Phone,
			}
		}
		usersDTO = append(usersDTO, userDTO)
	}

	return usersDTO, total, nil
}

func (service *UserService) GetAllSchoolAdmin(page int, limit int, sortField, sortDirection string) ([]dto.UserResponseDTO, int, error) {
	offset := (page - 1) * limit

	users, school, err := service.userRepository.FetchAllSchoolAdmins(offset, limit, sortField, sortDirection)
	if err != nil {
		return nil, 0, err
	}

	total, err := service.userRepository.CountSchoolAdmin()
	if err != nil {
		return nil, 0, err
	}

	var usersDTO []dto.UserResponseDTO

	for _, user := range users {
		userDTO := dto.UserResponseDTO{
			UUID:       user.UUID.String(),
			Username:   user.Username,
			Email:      user.Email,
			Status:     user.Status,
			LastActive: safeTimeFormat(user.LastActive),
			CreatedAt:  safeTimeFormat(user.CreatedAt),
		}
		if details, ok := user.Details.(entity.SchoolAdminDetails); ok {
			userDTO.Details = dto.SchoolAdminDetailsResponseDTO{
				SchoolName: school.Name,
				Picture:    details.Picture,
				FirstName:  details.FirstName,
				LastName:   details.LastName,
				Gender:     dto.Gender(details.Gender),
				Phone:      details.Phone,
			}
		}
		usersDTO = append(usersDTO, userDTO)
	}

	return usersDTO, total, nil
}

func (service *UserService) GetAllDriverFromAllSchools(page int, limit int, sortField string, sortDirection string) ([]dto.UserResponseDTO, int, error) {
	offset := (page - 1) * limit

	users, school, vehicle, err := service.userRepository.FetchAllDrivers(offset, limit, sortField, sortDirection)
	if err != nil {
		return nil, 0, err
	}

	total, err := service.userRepository.CountAllDriver()
	if err != nil {
		return nil, 0, err
	}

	var usersDTO []dto.UserResponseDTO

	for _, user := range users {
		userDTO := dto.UserResponseDTO{
			UUID:       user.UUID.String(),
			Username:   user.Username,
			Email:      user.Email,
			Status:     user.Status,
			LastActive: safeTimeFormat(user.LastActive),
			CreatedAt:  safeTimeFormat(user.CreatedAt),
		}
		if details, ok := user.Details.(entity.DriverDetails); ok {
			var vehicleDetails string
			if vehicle.VehicleNumber == "N/A" || vehicle.UUID == uuid.Nil {
				vehicleDetails = "N/A"
			} else {
				vehicleDetails = fmt.Sprintf("%s (%s)", vehicle.VehicleNumber, vehicle.VehicleName)
			}

			userDTO.Details = dto.DriverDetailsResponseDTO{
				SchoolName:    school.Name,
				VehicleNumber: vehicleDetails,
				Picture:       details.Picture,
				FirstName:     details.FirstName,
				LastName:      details.LastName,
				Gender:        dto.Gender(details.Gender),
				Phone:         details.Phone,
				Address:       details.Address,
				LicenseNumber: details.LicenseNumber,
			}
		}
		usersDTO = append(usersDTO, userDTO)
	}

	return usersDTO, total, nil
}

func (service *UserService) GetAllDriverForPermittedSchool(page int, limit int, sortField string, sortDirection string, schoolUUID string) ([]dto.UserResponseDTO, int, error) {
	offset := (page - 1) * limit

	users, school, vehicle, err := service.userRepository.FetchAllDriversForPermittedSchool(offset, limit, sortField, sortDirection, schoolUUID)
	if err != nil {
		return nil, 0, err
	}

	total, err := service.userRepository.CountSchoolAdmin()
	if err != nil {
		return nil, 0, err
	}

	var usersDTO []dto.UserResponseDTO

	for _, user := range users {
		userDTO := dto.UserResponseDTO{
			UUID:       user.UUID.String(),
			Username:   user.Username,
			Email:      user.Email,
			Status:     user.Status,
			LastActive: safeTimeFormat(user.LastActive),
			CreatedAt:  safeTimeFormat(user.CreatedAt),
		}
		if details, ok := user.Details.(entity.DriverDetails); ok {
			userDTO.Details = dto.DriverDetailsResponseDTO{
				SchoolName:    school.Name,
				VehicleNumber: vehicle.VehicleNumber,
				Picture:       details.Picture,
				FirstName:     details.FirstName,
				LastName:      details.LastName,
				Gender:        dto.Gender(details.Gender),
				Phone:         details.Phone,
				Address:       details.Address,
				LicenseNumber: details.LicenseNumber,
			}
		}
		usersDTO = append(usersDTO, userDTO)
	}

	return usersDTO, total, nil
}

func (service *UserService) GetSpecSuperAdmin(uuid string) (dto.UserResponseDTO, error) {
	user, err := service.userRepository.FetchSpecSuperAdmin(uuid)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}

	userDTO := dto.UserResponseDTO{
		UUID:       user.UUID.String(),
		Username:   user.Username,
		Email:      user.Email,
		Status:     user.Status,
		LastActive: safeTimeFormat(user.LastActive),
		CreatedAt:  safeTimeFormat(user.CreatedAt),
		CreatedBy:  safeStringFormat(user.CreatedBy),
		UpdatedAt:  safeTimeFormat(user.UpdatedAt),
		UpdatedBy:  safeStringFormat(user.UpdatedBy),
	}

	if details, ok := user.Details.(entity.SuperAdminDetails); ok {
		var Picture string
		if details.Picture == "" {
			Picture = "N/A"
		} else {
			Picture = details.Picture
		}
		userDTO.Details = dto.SuperAdminDetailsResponseDTO{
			Picture:   Picture,
			FirstName: details.FirstName,
			LastName:  details.LastName,
			Gender:    dto.Gender(details.Gender),
			Phone:     details.Phone,
			Address:   details.Address,
		}
	}

	return userDTO, nil
}

func (service *UserService) GetSpecSchoolAdmin(uuid string) (dto.UserResponseDTO, error) {
	user, school, err := service.userRepository.FetchSpecSchoolAdmin(uuid)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}

	userDTO := dto.UserResponseDTO{
		UUID:       user.UUID.String(),
		Username:   user.Username,
		Email:      user.Email,
		Status:     user.Status,
		LastActive: safeTimeFormat(user.LastActive),
		CreatedAt:  safeTimeFormat(user.CreatedAt),
		CreatedBy:  safeStringFormat(user.CreatedBy),
		UpdatedAt:  safeTimeFormat(user.UpdatedAt),
		UpdatedBy:  safeStringFormat(user.UpdatedBy),
	}

	if details, ok := user.Details.(entity.SchoolAdminDetails); ok {
		userDTO.Details = dto.SchoolAdminDetailsResponseDTO{
			SchoolUUID: details.SchoolUUID.String(),
			SchoolName: school.Name,
			Picture:    details.Picture,
			FirstName:  details.FirstName,
			LastName:   details.LastName,
			Gender:     dto.Gender(details.Gender),
			Phone:      details.Phone,
			Address:    details.Address,
		}
	}

	return userDTO, nil
}

func (service *UserService) GetSpecDriverFromAllSchools(id string) (dto.UserResponseDTO, error) {
	user, school, vehicle, err := service.userRepository.FetchSpecDriverFromAllSchools(id)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}

	userDTO := dto.UserResponseDTO{
		UUID:       user.UUID.String(),
		Username:   user.Username,
		Email:      user.Email,
		Status:     user.Status,
		LastActive: safeTimeFormat(user.LastActive),
		CreatedAt:  safeTimeFormat(user.CreatedAt),
		CreatedBy:  safeStringFormat(user.CreatedBy),
		UpdatedAt:  safeTimeFormat(user.UpdatedAt),
		UpdatedBy:  safeStringFormat(user.UpdatedBy),
	}

	if details, ok := user.Details.(entity.DriverDetails); ok {
		var vehicleDetails, schoolUUID, vehicleUUID string
		if vehicle.VehicleNumber == "N/A" || vehicle.UUID == uuid.Nil {
			vehicleDetails = "N/A"
		} else {
			vehicleDetails = fmt.Sprintf("%s (%s)", vehicle.VehicleNumber, vehicle.VehicleName)
		}

		if school.UUID == uuid.Nil {
			schoolUUID = "N/A"
		} else {
			schoolUUID = school.UUID.String()
		}

		if vehicle.UUID == uuid.Nil {
			vehicleUUID = "N/A"
		} else {
			vehicleUUID = vehicle.UUID.String()
		}

		userDTO.Details = dto.DriverDetailsResponseDTO{
			SchoolUUID:    schoolUUID,
			SchoolName:    school.Name,
			VehicleUUID:   vehicleUUID,
			VehicleNumber: vehicleDetails,
			Picture:       details.Picture,
			FirstName:     details.FirstName,
			LastName:      details.LastName,
			Gender:        dto.Gender(details.Gender),
			Phone:         details.Phone,
			Address:       details.Address,
			LicenseNumber: details.LicenseNumber,
		}
	}

	return userDTO, nil
}

func (s *UserService) AddUser(req dto.UserRequestsDTO, user_name string) (uuid.UUID, error) {
	exists, err := s.userRepository.CheckEmailExist("", req.Email)
	if err != nil {
		return uuid.Nil, err
	}
	if exists {
		return uuid.Nil, errors.New("email already exists", 409)
	}

	exists, err = s.userRepository.CheckUsernameExist("", req.Username)
	if err != nil {
		return uuid.Nil, err
	}
	if exists {
		return uuid.Nil, errors.New("username already exists", 409)
	}

	if req.Password != "" {
		hashedPassword, err := hashPassword(req.Password)
		if err != nil {
			return uuid.Nil, err
		}
		req.Password = hashedPassword
	}

	tx, err := s.userRepository.BeginTransaction()
	if err != nil {
		return uuid.Nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	var transactionErr error
	defer func() {
		if transactionErr != nil {
			tx.Rollback()
		} else {
			transactionErr = tx.Commit()
		}
	}()

	userEntity := entity.User{
		ID:        time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
		UUID:      uuid.New(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		Role:      entity.Role(req.Role),
		RoleCode:  req.RoleCode,
		CreatedBy: sql.NullString{String: user_name, Valid: user_name != ""},
		Details:   req.Details,
	}

	userUUID, err := s.userRepository.SaveUser(tx, userEntity)
	if err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return uuid.Nil, errors.New(customErr.Message, customErr.StatusCode)
		}
		transactionErr = fmt.Errorf("error saving user: %w", err)
		return uuid.Nil, transactionErr
	}

	if _, ok := req.Details.(map[string]interface{}); ok {
		if err := s.saveRoleDetails(tx, userEntity, req); err != nil {
			transactionErr = fmt.Errorf("error saving role details: %w", err)
			return uuid.Nil, transactionErr
		}
	} else {
		switch userEntity.Role {
		case entity.SchoolAdmin:
			if details, ok := req.Details.(dto.SchoolAdminDetailsRequestsDTO); ok {
				req.Details = details
			} else {
				return uuid.Nil, errors.New("invalid school admin details", 400)
			}
		case entity.Driver:
			if details, ok := req.Details.(dto.DriverDetailsRequestsDTO); ok {
				req.Details = details
			} else {
				return uuid.Nil, errors.New("invalid driver details", 400)
			}
		}

		if err := s.saveRoleDetails(tx, userEntity, req); err != nil {
			transactionErr = fmt.Errorf("error saving role details: %w", err)
			return uuid.Nil, transactionErr
		}
	}

	return userUUID, nil
}

func (s *UserService) UpdateUser(id string, req dto.UserRequestsDTO, username string, detailsMap map[string]interface{}, file []byte) error {
	tx, err := s.userRepository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	exists, err := s.userRepository.CheckEmailExist(id, req.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists", 409)
	}

	exists, err = s.userRepository.CheckUsernameExist(id, req.Username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("username already exists", 409)
	}

	userData := entity.User{
		Username:  req.Username,
		Email:     req.Email,
		Role:      entity.Role(req.Role),
		RoleCode:  req.RoleCode,
		Details:   detailsMap,
		UpdatedBy: sql.NullString{String: username, Valid: username != ""},
	}

	if err := s.userRepository.UpdateUser(tx, userData, id); err != nil {
		return err
	}

	if _, ok := userData.Details.(map[string]interface{}); ok {
		if err := s.updateRoleDetails(tx, userData, req, id); err != nil {
			logger.LogError(err, "error saving role details", map[string]interface{}{})
			return fmt.Errorf("error saving role details: %w", err)
		}
	}

	return nil
}

func (service *UserService) DeleteSuperAdmin(id string, user_name string) error {
	tx, err := service.userRepository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("user not found", 404)
	}

	err = service.userRepository.DeleteSuperAdmin(tx, parsedUUID, user_name)
	if err != nil {
		return errors.New("user not found", 404)
	}

	return nil
}

func (service *UserService) DeleteSchoolAdmin(id string, user_name string) error {
	tx, err := service.userRepository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	var transactionErr error
	defer func() {
		if transactionErr != nil {
			tx.Rollback()
		} else {
			transactionErr = tx.Commit()
		}
	}()

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		transactionErr = errors.New("user not found", 404)
		return transactionErr
	}

	err = service.userRepository.DeleteSchoolAdmin(tx, parsedUUID, user_name)
	if err != nil {
		transactionErr = errors.New("user not found", 404)
		return transactionErr
	}

	return nil
}

func (service *UserService) DeleteDriver(id string, user_name string) error {
	tx, err := service.userRepository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	var transactionErr error
	defer func() {
		if transactionErr != nil {
			tx.Rollback()
		} else {
			transactionErr = tx.Commit()
		}
	}()

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		transactionErr = errors.New("user not found", 404)
		return transactionErr
	}

	err = service.userRepository.DeleteDriver(tx, parsedUUID, user_name)
	if err != nil {
		transactionErr = errors.New("user not found", 404)
		return transactionErr
	}

	return nil
}


func (service *UserService) GetSpecUserWithDetails(id string) (entity.User, error) {
	user, err := service.userRepository.FetchSpecificUser(id)
	if err != nil {
		return entity.User{}, err
	}
	switch user.RoleCode {
	case "SA":
		superAdminDetails, err := service.userRepository.FetchSuperAdminDetails(user.UUID)
		if err != nil {
			return entity.User{}, err
		}
		user.Details = superAdminDetails
		return user, nil
	case "AS":
		schoolAdminDetails, err := service.userRepository.FetchSchoolAdminDetails(user.UUID)
		if err != nil {
			return entity.User{}, err
		}
		user.Details = schoolAdminDetails
		return user, nil
	case "P":
		parentDetails, err := service.userRepository.FetchParentDetails(user.UUID)
		if err != nil {
			return entity.User{}, err
		}
		user.Details = parentDetails
		return user, nil
	case "D":
		driverDetails, err := service.userRepository.FetchDriverDetails(user.UUID)
		if err != nil {
			return entity.User{}, err
		}
		user.Details = driverDetails
		return user, nil
	default:
		return entity.User{}, errors.New("invalid role code", 0)
	}
}

func (service *UserService) GetSpecUser(id string) (entity.User, error) {
	db, err := databases.PostgresConnection()
	if err != nil {
		return entity.User{}, err
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return entity.User{}, errors.New("invalid user id", 0)
	}

	var user entity.User
	query := `
		SELECT * FROM users WHERE user_id = $1
	`

	err = db.Get(&user, query, idInt)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (service *UserService) CheckPermittedSchoolAccess(userUUID string) (string, error) {
	schoolUUID, err := service.userRepository.FetchPermittedSchoolAccess(userUUID)
	if err != nil {
		return "", err
	}

	return schoolUUID, nil
}

func (s *UserService) saveRoleDetails(tx *sqlx.Tx, userEntity entity.User, req dto.UserRequestsDTO) error {
	var params interface{}

	switch userEntity.Role {
	case entity.SuperAdmin:
		details := entity.SuperAdminDetails{
			Picture:   req.Picture,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Gender:    entity.Gender(req.Gender),
			Phone:     req.Phone,
			Address:   req.Address,
		}

		if err := s.userRepository.SaveSuperAdminDetails(tx, details, userEntity.UUID, params); err != nil {
			return err
		}
	case entity.SchoolAdmin:
		details := entity.SchoolAdminDetails{
			SchoolUUID: uuid.MustParse(req.Details.(dto.SchoolAdminDetailsRequestsDTO).SchoolUUID),
			Picture:    req.Picture,
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			Gender:     entity.Gender(req.Gender),
			Phone:      req.Phone,
			Address:    req.Address,
		}

		if err := s.userRepository.SaveSchoolAdminDetails(tx, details, userEntity.UUID, params); err != nil {
			return err
		}
	case entity.Parent:
		details := entity.ParentDetails{
			Picture:   req.Picture,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Gender:    entity.Gender(req.Gender),
			Phone:     req.Phone,
			Address:   req.Address,
		}

		if err := s.userRepository.SaveParentDetails(tx, details, userEntity.UUID, params); err != nil {
			return err
		}
	case entity.Driver:
		details := entity.DriverDetails{
			SchoolUUID:    parseSafeUUID(req.Details.(dto.DriverDetailsRequestsDTO).SchoolUUID),
			VehicleUUID:   parseSafeUUID(req.Details.(dto.DriverDetailsRequestsDTO).VehicleUUID),
			Picture:       req.Picture,
			FirstName:     req.FirstName,
			LastName:      req.LastName,
			Gender:        entity.Gender(req.Gender),
			Phone:         req.Phone,
			Address:       req.Address,
			LicenseNumber: req.Details.(dto.DriverDetailsRequestsDTO).LicenseNumber,
		}

		if err := s.userRepository.SaveDriverDetails(tx, details, userEntity.UUID, params); err != nil {
			return err
		}
	default:
		return errors.New("invalid role", 400)
	}

	return nil
}

func (s *UserService) updateRoleDetails(tx *sqlx.Tx, userEntity entity.User, req dto.UserRequestsDTO, id string) error {
	switch userEntity.Role {
	case entity.SuperAdmin:
		details := entity.SuperAdminDetails{
			Picture:   req.Picture,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Gender:    entity.Gender(req.Gender),
			Phone:     req.Phone,
			Address:   req.Address,
		}

		if err := s.userRepository.UpdateSuperAdminDetails(tx, details, id); err != nil {
			return err
		}
	case entity.SchoolAdmin:
		details := entity.SchoolAdminDetails{
			SchoolUUID: *parseSafeUUID(req.Details.(dto.SchoolAdminDetailsRequestsDTO).SchoolUUID),
			Picture:    req.Picture,
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			Gender:     entity.Gender(req.Gender),
			Phone:      req.Phone,
			Address:    req.Address,
		}

		if err := s.userRepository.UpdateSchoolAdminDetails(tx, details, id); err != nil {
			return err
		}
	case entity.Parent:
		details := entity.ParentDetails{
			Picture:   req.Picture,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Gender:    entity.Gender(req.Gender),
			Phone:     req.Phone,
			Address:   req.Address,
		}

		if err := s.userRepository.UpdateParentDetails(tx, details, id); err != nil {
			return err
		}
	case entity.Driver:
		details := entity.DriverDetails{
			SchoolUUID:    parseSafeUUID(req.Details.(dto.DriverDetailsRequestsDTO).SchoolUUID),
			VehicleUUID:   parseSafeUUID(req.Details.(dto.DriverDetailsRequestsDTO).VehicleUUID),
			Picture:       req.Picture,
			FirstName:     req.FirstName,
			LastName:      req.LastName,
			Gender:        entity.Gender(req.Gender),
			Phone:         req.Phone,
			Address:       req.Address,
			LicenseNumber: req.Details.(dto.DriverDetailsRequestsDTO).LicenseNumber,
		}

		parsedUUID, err := uuid.Parse(id)
		if err != nil {
			return fmt.Errorf("invalid UUID: %w", err)
		}

		if err := s.userRepository.UpdateDriverDetails(tx, details, parsedUUID); err != nil {
			return err
		}
	default:
		return errors.New("invalid role", 400)
	}

	return nil
}

func parseSafeUUID(id string) *uuid.UUID {
	if id == "" || id == "00000000-0000-0000-0000-000000000000" {
		return nil
	}
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil
	}
	return &parsedUUID
}