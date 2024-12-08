package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserDataOnLoginDTO struct {
	UserID    int64  `json:"user_id"`
	UserUUID  string `json:"user_uuid"`
	RoleCode  string `json:"user_role_code"`
	Password  string `json:"user_password"`
	FirstName string `json:"user_first_name"`
	LastName  string `json:"user_last_name"`
}
