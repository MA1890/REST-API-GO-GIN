package database

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	Id      int `json:"id"`
	UserId  int `json:"user_id"`
	EventId int `json:"event-id"`
}

func (m *AttendeeModel) Insert(attendee *Attendee) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "INSERT INTO attendees(user_id, event_id) VALUES (?, ?) RETURNING id"
	err := m.DB.QueryRowContext(ctx, query, &attendee.UserId, &attendee.EventId).Scan(&attendee.Id)
	if err != nil {
		return nil, err
	}
	return attendee, nil
}
func (m *AttendeeModel) GetByEventAndAttendee(evenId, userId int) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var attendee Attendee
	query := "SELECT id, user_id, event_id FROM attendees WHERE user_id = ? AND event_id = ?"
	err := m.DB.QueryRowContext(ctx, query, userId, evenId).Scan(&attendee.Id, &attendee.UserId, &attendee.EventId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &attendee, nil
}
func (m *AttendeeModel) GetAttendeesByEvent(eventId int) ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var users []User
	query := "SELECT u.id, u.name, u.email FROM users u JOIN attendees a ON u.id = a.user_id WHERE a.event_id = ?"
	rows, err := m.DB.QueryContext(ctx, query, eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
func (m *AttendeeModel) Delete(userId, eventId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := "DELETE FROM attendees WHERE user_id = ? AND event_id = ?"
	_, err := m.DB.ExecContext(ctx, query, userId, eventId)
	if err != nil {
		return err
	}
	return nil
}
