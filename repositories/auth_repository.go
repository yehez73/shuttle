package repositories

import (
	"context"
	"shuttle/models/entity"
	"time"

	"github.com/jmoiron/sqlx"
)

type AuthRepositoryInterface interface {
	Login(email string) (entity.UserDataOnLogin, error)
	CheckRefreshTokenData(userID int64, token string) (entity.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, userID int64) error
	UpdateUserStatus(userUUID string, status string, lastActive time.Time) error
}

type authRepository struct {
	DB *sqlx.DB
}

func NewAuthRepository(DB *sqlx.DB) AuthRepositoryInterface {
	return &authRepository{
		DB: DB,
	}
}

func (r *authRepository) Login(email string) (entity.UserDataOnLogin, error) {
	var user entity.UserDataOnLogin
	query := `SELECT u.user_id, u.user_uuid, u.user_role_code, u.user_password,
              COALESCE(sa.user_first_name, sca.user_first_name, p.user_first_name, d.user_first_name) AS user_first_name,
              COALESCE(sa.user_last_name, sca.user_last_name, p.user_last_name, d.user_last_name) AS user_last_name
              FROM users u
              LEFT JOIN super_admin_details sa ON u.user_uuid = sa.user_uuid
              LEFT JOIN school_admin_details sca ON u.user_uuid = sca.user_uuid
              LEFT JOIN parent_details p ON u.user_uuid = p.user_uuid
              LEFT JOIN driver_details d ON u.user_uuid = d.user_uuid
              WHERE u.user_email = $1`

	row := r.DB.QueryRow(query, email)

	if err := row.Scan(&user.ID, &user.UUID, &user.RoleCode, &user.Password, &user.FirstName, &user.LastName); err != nil {
		return entity.UserDataOnLogin{}, err
	}

	return user, nil
}

func (r *authRepository) CheckRefreshTokenData(userID int64, token string) (entity.RefreshToken, error) {
	query := `
		SELECT refresh_token, expired_at, is_revoked 
		FROM refresh_tokens 
		WHERE user_id = $1 AND refresh_token = $2
	`

	var tokenData entity.RefreshToken
	err := r.DB.Get(&tokenData, query, userID, token)
	if err != nil {
		return tokenData, err
	}

	return tokenData, nil
}

func SaveRefreshToken(db sqlx.DB, refreshToken entity.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, refresh_token, expired_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id)
		DO UPDATE SET refresh_token = $3, issued_at = CURRENT_TIMESTAMP, expired_at = $4, is_revoked = false 
	`
	_, err := db.Exec(query, refreshToken.ID, refreshToken.UserID, refreshToken.RefreshToken, refreshToken.ExpiredAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *authRepository) DeleteRefreshToken(ctx context.Context, userID int64) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE user_id = $1
	`

	_, err := r.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *authRepository) UpdateUserStatus(userUUID string, status string, lastActive time.Time) error {
	query := `
		UPDATE users
		SET user_status = $1, user_last_active = $2
		WHERE user_uuid = $3
	`

	_, err := r.DB.Exec(query, status, lastActive, userUUID)
	if err != nil {
		return err
	}

	return nil
}