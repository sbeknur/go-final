package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Instructors struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Age       int32  `json:"age,omitempty"`
}

type InstructorsModel struct {
	DB *sql.DB
}

func (a InstructorsModel) Get(id int64) (*Instructors, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT * FROM instructors WHERE id = $1`

	var instructor Instructors
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.DB.QueryRowContext(ctx, query, id).Scan(
		&instructor.ID,
		&instructor.FirstName,
		&instructor.LastName,
		&instructor.Age,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &instructor, nil
}

func (a InstructorsModel) Insert(instructors *Instructors) error {
	query :=
		`INSERT INTO instructors (firstName, lastName, age)
		 VALUES ($1, $2, $3)
		 RETURNING id`
	args := []any{instructors.FirstName, instructors.LastName, instructors.Age}

	return a.DB.QueryRow(query, args...).Scan(&instructors.ID)
}

func (a InstructorsModel) GetAll(firstName string, lastName string, filters Filters) ([]*Instructors, Metadata, error) {
	query := fmt.Sprintf(` 
		SELECT COUNT(*) OVER(), id, firstName, lastName, age
		FROM instructors
		WHERE (to_tsvector('simple', firstName) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', lastName) @@ plainto_tsquery('simple', $2) OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{firstName, lastName, filters.limit(), filters.offset()}

	rows, err := a.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	instructors := []*Instructors{}
	for rows.Next() {
		var instructor Instructors
		err := rows.Scan(
			&totalRecords,
			&instructor.ID,
			&instructor.FirstName,
			&instructor.LastName,
			&instructor.Age,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		instructors = append(instructors, &instructor)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return instructors, metadata, nil
}

func (a InstructorsModel) Update(instructor *Instructors) error {
	query :=
		`UPDATE instructors
		 SET first_name = $1, last_name = $2, age = $3
		 WHERE id = $4`

	args := []any{
		&instructor.ID,
		&instructor.FirstName,
		&instructor.LastName,
		&instructor.Age,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.DB.QueryRowContext(ctx, query, args...).Scan(&instructor.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (a InstructorsModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM instructors WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := a.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}