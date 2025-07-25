package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/database"
	"github.com/koiraladarwin/scanin/models"
)

func (p *PostgresDB) CreateUser(u *models.User) (*models.User,error) {
	var lastAutoID int
	err := p.sql.QueryRow(`SELECT COALESCE(MAX(auto_id), 0) FROM users WHERE role = $1`, u.Role).Scan(&lastAutoID)
	if err != nil {
		return nil,fmt.Errorf("failed to fetch latest auto_id: %w", err)
	}

	u.AutoId = lastAutoID + 1

	query := `
		INSERT INTO users (auto_id, full_name, image_url, position, company, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err = p.sql.QueryRow(
		query,
		u.AutoId,
		u.FullName,
		u.Image_url,
		u.Position,
		u.Company,
		u.Role,
	).Scan(&u.ID)

	if isUniqueViolationError(err) {
		return nil,db.ErrAlreadyExists
	}

	return u,err
}

func (p *PostgresDB) GetUser(id uuid.UUID) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, full_name, auto_id, image_url, position, company ,role FROM users WHERE id=$1`
	err := p.sql.QueryRow(query, id).Scan(&u.ID, &u.FullName, &u.AutoId, &u.Image_url, &u.Position, &u.Company, &u.Role)
	return u, err
}

func (p *PostgresDB) GetUserByAttendeeid(Attendeeid uuid.UUID) (*models.User, error) {
	u := &models.User{}

	a, err := p.GetAttendee(Attendeeid)
	if err != nil {
		return nil, err
	}
	query := `SELECT id, full_name, auto_id, image_url, position, company,role FROM users WHERE id=$1`
	err = p.sql.QueryRow(query, a.UserID).Scan(&u.ID, &u.FullName, &u.AutoId, &u.Image_url, &u.Position, &u.Company, &u.Role)
	return u, err
}

// func (p *PostgresDB) UpdateUser(u *models.User) error {
// 	query := `UPDATE users SET full_name=$1, email=$2, phone=$3 WHERE id=$5`
// 	_, err := p.sql.Exec(query, u.FullName, u.Email, u.Phone, u.ID)
// 	return err
// }
//
// func (p *PostgresDB) DeleteUser(id uuid.UUID) error {
// 	_, err := p.sql.Exec(`DELETE FROM users WHERE id=$1`, id)
// 	return err
// }

func (p *PostgresDB) GetUsersByEvent(eventID uuid.UUID) ([]models.UserWithRole, error) {
	var users []models.UserWithRole

	rows, err := p.sql.Query(`
		SELECT u.id, u.full_name, u.auto_id, u.position, u.company,u.image_url, u.role, a.id
		FROM attendees a
		JOIN users u ON u.id = a.user_id
		WHERE a.event_id = $1
	`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.UserWithRole
		if err := rows.Scan(&u.ID, &u.FullName, &u.AutoId, &u.Position, &u.Company, &u.Image_url, &u.Role, &u.AttendeeId); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (p *PostgresDB) GetAllUsers() ([]models.User, error) {
	rows, err := p.sql.Query(`SELECT id, full_name, auto_id, image_url, position, company,role FROM users `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.FullName, &u.AutoId, &u.Image_url, &u.Position, &u.Company, &u.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
