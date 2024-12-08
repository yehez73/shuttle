package repositories

import (
	"fmt"
	"shuttle/models/entity"

	"github.com/jmoiron/sqlx"
)

type SchoolRepositoryInterface interface {
	FetchAllSchools(offset, limit int, sortField, sortDirection string) ([]entity.School, map[string]entity.SchoolAdminDetails, error)
	FetchSpecSchool(uuid string) (entity.School, entity.SchoolAdminDetails, error)
	SaveSchool(entity.School) error
	UpdateSchool(entity.School) error
	DeleteSchool(entity.School) error
	CountSchools() (int, error)
}

type schoolRepository struct {
	DB *sqlx.DB
}

func NewSchoolRepository(DB *sqlx.DB) SchoolRepositoryInterface {
	return &schoolRepository{
		DB: DB,
	}
}

func (repositories *schoolRepository) FetchAllSchools(offset, limit int, sortField, sortDirection string) ([]entity.School, map[string]entity.SchoolAdminDetails, error) {
    var schools []entity.School
	var adminMap = make(map[string]entity.SchoolAdminDetails)

    query := fmt.Sprintf(`
        SELECT s.school_uuid, s.school_name, s.school_address, s.school_contact, s.school_email, s.created_at,
		u.user_uuid, sad.school_uuid, sad.user_first_name, sad.user_last_name
		FROM schools s
		LEFT JOIN school_admin_details sad ON s.school_uuid = sad.school_uuid
		LEFT JOIN users u ON sad.user_uuid = u.user_uuid
		WHERE s.deleted_at IS NULL
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, sortField, sortDirection)

	rows, err := repositories.DB.Queryx(query, limit, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var school entity.School
		var admin entity.SchoolAdminDetails

		if err := rows.Scan(&school.UUID, &school.Name, &school.Address, &school.Contact, &school.Email, &school.CreatedAt,
			&admin.UserUUID, &admin.SchoolUUID, &admin.FirstName, &admin.LastName); err != nil {
			return nil, nil, err
		}

		schools = append(schools, school)
		if admin.UserUUID.String() != "" {
			adminMap[school.UUID.String()] = admin
		}
	}

	return schools, adminMap, nil
}

func (repositories *schoolRepository) FetchSpecSchool(uuid string) (entity.School, entity.SchoolAdminDetails, error) {
	var school entity.School
	var admin entity.SchoolAdminDetails

	query := `
		SELECT s.school_uuid, s.school_name, s.school_address, s.school_contact, s.school_email, s.school_description, s.created_at,
			s.created_by, s.updated_at, s.updated_by, 
			sad.user_uuid, sad.school_uuid AS admin_school_uuid, sad.user_first_name, sad.user_last_name
		FROM schools s
		LEFT JOIN school_admin_details sad ON s.school_uuid = sad.school_uuid
		WHERE s.deleted_at IS NULL AND s.school_uuid = $1
	`

	err := repositories.DB.QueryRowx(query, uuid).Scan(
		&school.UUID, &school.Name, &school.Address, &school.Contact, &school.Email, &school.Description, &school.CreatedAt,
		&school.CreatedBy, &school.UpdatedAt, &school.UpdatedBy, &admin.UserUUID, &admin.SchoolUUID, &admin.FirstName, &admin.LastName,
	)
	if err != nil {
		return entity.School{}, entity.SchoolAdminDetails{}, err
	}

	return school, admin, nil
}

func (r *schoolRepository) SaveSchool(school entity.School) error {
	query := `INSERT INTO schools (school_id, school_uuid, school_name, school_address, school_contact, school_email, school_description, created_by)
			  VALUES (:school_id, :school_uuid, :school_name, :school_address, :school_contact, :school_email, :school_description, :created_by)`
	_, err := r.DB.NamedExec(query, school)
	if err != nil {
		return err
	}

	return nil
}

func (r *schoolRepository) UpdateSchool(school entity.School) error {
	query := `
		UPDATE schools SET school_name = :school_name, school_address = :school_address, school_contact = :school_contact, school_email = :school_email, school_description = :school_description, updated_at = :updated_at, updated_by = :updated_by 
		WHERE school_uuid = :school_uuid`
	_, err := r.DB.NamedExec(query, school)
	if err != nil {
		return err
	}

	return nil
}

func (r *schoolRepository) DeleteSchool(school entity.School) error {
	query := `UPDATE schools SET deleted_at = :deleted_at, deleted_by = :deleted_by WHERE school_uuid = :school_uuid`
	_, err := r.DB.NamedExec(query, school)
	if err != nil {
		return err
	}

	return nil
}

func (repositories *schoolRepository) CountSchools() (int, error) {
    var total int

    query := `
        SELECT COUNT(*) 
        FROM schools 
		WHERE deleted_at IS NULL
    `

    if err := repositories.DB.Get(&total, query); err != nil {
        return 0, err
    }

    return total, nil
}