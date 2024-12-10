package dto

import ()

type Role string
type Gender string

// type Validatable interface {
// 	Validate() error
// }

const (
	SuperAdmin  Role = "superadmin"
	SchoolAdmin Role = "schooladmin"
	Parent      Role = "parent"
	Driver      Role = "driver"

	Female Gender = "female"
	Male   Gender = "male"
)

type UserRequestsDTO struct {
	Username  string      `json:"username" validate:"required,username,min=5"`
	Email     string      `json:"email" validate:"required,email"`
	Password  string      `json:"password" validate:"required,min=8"`
	Role      Role        `json:"role" validate:"required,role"`
	RoleCode  string      `json:"role_code"`
	Picture   string      `json:"picture"`
	FirstName string      `json:"first_name" validate:"required,max=255"`
	LastName  string      `json:"last_name" validate:"required,max=255"`
	Gender    Gender      `json:"gender" validate:"required,gender"`
	Phone     string      `json:"phone" validate:"required,phone"`
	Address   string      `json:"address" validate:"required,max=255"`
	Details   interface{} `json:"details"`
}

type SchoolAdminDetailsRequestsDTO struct {
	SchoolUUID string `json:"school_uuid" validate:"required"`
}

type DriverDetailsRequestsDTO struct {
	SchoolUUID    string `json:"school_uuid"`
	VehicleUUID   string `json:"vehicle_uuid"`
	LicenseNumber string `json:"license_number" validate:"required"`
}

type UserResponseDTO struct {
	UUID       string      `json:"user_uuid"`
	Username   string      `json:"user_username"`
	Email      string      `json:"user_email"`
	Role       Role        `json:"user_role,omitempty"`
	RoleCode   string      `json:"user_role_code,omitempty"`
	Status     string      `json:"user_status"`
	LastActive string      `json:"user_last_active"`
	Details    interface{} `json:"user_details"`
	CreatedAt  string      `json:"created_at,omitempty"`
	CreatedBy  string      `json:"created_by,omitempty"`
	UpdatedAt  string      `json:"updated_at,omitempty"`
	UpdatedBy  string      `json:"updated_by,omitempty"`
}

type SuperAdminDetailsResponseDTO struct {
	Picture   string `json:"user_picture,omitempty"`
	FirstName string `json:"user_first_name"`
	LastName  string `json:"user_last_name"`
	Gender    Gender `json:"user_gender"`
	Phone     string `json:"user_phone"`
	Address   string `json:"user_address,omitempty"`
}

type SchoolAdminDetailsResponseDTO struct {
	SchoolUUID string `json:"school_uuid,omitempty"`
	SchoolName string `json:"school_name"`
	Picture    string `json:"user_picture,omitempty"`
	FirstName  string `json:"user_first_name"`
	LastName   string `json:"user_last_name"`
	Gender     Gender `json:"user_gender"`
	Phone      string `json:"user_phone"`
	Address    string `json:"user_address,omitempty"`
}

type ParentDetailsResponseDTO struct {
	Picture   string `json:"user_picture,omitempty"`
	FirstName string `json:"user_first_name"`
	LastName  string `json:"user_last_name"`
	Gender    Gender `json:"user_gender"`
	Phone     string `json:"user_phone"`
	Address   string `json:"user_address,omitempty"`
}

type DriverDetailsResponseDTO struct {
	SchoolUUID    string `json:"school_uuid,omitempty"`
	SchoolName    string `json:"school_name"`
	VehicleUUID   string `json:"vehicle_uuid,omitempty"`
	VehicleNumber string `json:"vehicle_number"`
	Picture       string `json:"user_picture,omitempty"`
	FirstName     string `json:"user_first_name"`
	LastName      string `json:"user_last_name"`
	Gender        Gender `json:"user_gender"`
	Phone         string `json:"user_phone"`
	Address       string `json:"user_address,omitempty"`
	LicenseNumber string `json:"license_number"`
}

// func (d ParentDetailsRequestsDTO) Validate() error {
// 	if err := validateGender(string(d.Gender)); err != nil {
// 		return err
// 	}
// 	if err := validatePhone(d.Phone); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (d SchoolAdminDetailsRequestsDTO) Validate() error {
// 	if err := validateGender(string(d.Gender)); err != nil {
// 		return err
// 	}
// 	if err := validatePhone(d.Phone); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (d SuperAdminDetailsRequestsDTO) Validate() error {
// 	if err := validateGender(string(d.Gender)); err != nil {
// 		return err
// 	}
// 	if err := validatePhone(d.Phone); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (d DriverDetailsRequestsDTO) Validate() error {
// 	if err := validateGender(string(d.Gender)); err != nil {
// 		return err
// 	}
// 	if err := validatePhone(d.Phone); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func validateGender(gender string) error {
// 	validGenders := map[string]bool{
// 		string(Male):   true,
// 		string(Female): true,
// 	}
// 	if !validGenders[strings.ToLower(gender)] {
// 		return errors.New("invalid gender")
// 	}
// 	return nil
// }

// func ValidatePhone(phone string) error {
// 	phoneRegex := regexp.MustCompile(`^\+?[0-9]{12,15}$`)
// 	if !phoneRegex.MatchString(phone) {
// 		return errors.New("phone number must be between 12 and 15 digits")
// 	}
// 	return nil
// }
