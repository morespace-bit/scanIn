package postgres

import (
	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/models"
)

// CreateActivity inserts a new activity into the database and returns its generated UUID
func (p *PostgresDB) CreateActivity(a *models.Activity) error {
	query := `INSERT INTO activities (event_id, name, type, start_time, end_time) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return p.sql.QueryRow(query, a.EventID, a.Name, a.Type, a.StartTime, a.EndTime).Scan(&a.ID)
}

// GetActivity fetches a single activity by UUID
func (p *PostgresDB) GetActivity(id uuid.UUID) (*models.Activity, error) {
	a := &models.Activity{}
	query := `SELECT id, event_id, name, type, start_time, end_time FROM activities WHERE id = $1`
	err := p.sql.QueryRow(query, id).Scan(&a.ID, &a.EventID, &a.Name, &a.Type, &a.StartTime, &a.EndTime)
	return a, err
}

// UpdateActivity updates an existing activity
func (p *PostgresDB) UpdateActivity(a *models.Activity) error {
	query := `UPDATE activities SET event_id=$1, name=$2, type=$3, start_time=$4, end_time=$5 WHERE id=$6`
	_, err := p.sql.Exec(query, a.EventID, a.Name, a.Type, a.StartTime, a.EndTime, a.ID)
	return err
}

// DeleteActivity deletes an activity by UUID
func (p *PostgresDB) DeleteActivity(id uuid.UUID) error {
	_, err := p.sql.Exec(`DELETE FROM activities WHERE id=$1`, id)
	return err
}


// GetActivities By event_id
func (p *PostgresDB) GetActivitiesByEvent(eventID uuid.UUID) ([]models.Activity, error) {
	activities := []models.Activity{}
	query := `SELECT id, event_id, name, type, start_time, end_time FROM activities WHERE event_id = $1`
	rows, err := p.sql.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Activity
		if err := rows.Scan(&a.ID, &a.EventID, &a.Name, &a.Type, &a.StartTime, &a.EndTime); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}
	return activities, nil
}
