package entity

import "time"

type UserDataOnLogin struct {
	ID        int64  `db:"user_id"`
	UUID      string `db:"user_uuid"`
	RoleCode  string `db:"user_role_code"`
	Password  string `db:"user_password"`
	FirstName string `db:"user_first_name"`
	LastName  string `db:"user_last_name"`
}

type RefreshToken struct {
	ID           int64     `db:"id"`
	UserID       int64     `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	IssuedAt     time.Time `db:"issued_at"`
	ExpiredAt    time.Time `db:"expired_at"`
	Revoked      bool      `db:"is_revoked"`
	LastUsedAt   time.Time `db:"last_used_at"`
}
