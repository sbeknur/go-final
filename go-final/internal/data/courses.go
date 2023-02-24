package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/sbeknur/go-final/internal/validator"
)

type Course struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"-"`
	Title          string    `json:"title"`
	Published_date string    `json:"published_date,omitempty"`
	Runtime        Runtime   `json:"runtime,omitempty"`
	Lectures       []string  `json:"lectures,omitempty"`
	Version        int32     `json:"version"`
}

type CourseModel struct {
	DB *sql.DB
}

func (m CourseModel) Insert(course *Course) error {
	query :=
		`INSERT INTO courses (title, published_date, runtime, lectures)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at, version`
	args := []any{course.Title, course.Published_date, course.Runtime, pq.Array(course.Lectures)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&course.ID, &course.CreatedAt, &course.Version)
}

func (m CourseModel) Get(id int64) (*Course, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT * FROM courses WHERE id = $1`

	var course Course
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&course.ID,
		&course.CreatedAt,
		&course.Title,
		&course.Published_date,
		&course.Runtime,
		pq.Array(&course.Lectures),
		&course.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &course, nil
}

func (m CourseModel) GetAll(title string, lectures []string, filters Filters) ([]*Course, Metadata, error) {
	query := fmt.Sprintf(` 
		SELECT COUNT(*) OVER(), id, created_at, title, published_date, runtime, lectures, version
		FROM courses 
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (lectures @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{title, pq.Array(lectures), filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	courses := []*Course{}
	for rows.Next() {
		var course Course
		err := rows.Scan(
			&totalRecords,
			&course.ID,
			&course.CreatedAt,
			&course.Title,
			&course.Published_date,
			&course.Runtime,
			pq.Array(&course.Lectures),
			&course.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		courses = append(courses, &course)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return courses, metadata, nil
}

func (m CourseModel) Update(course *Course) error {
	query :=
		`UPDATE courses
		 SET title = $1, published_date = $2, runtime = $3, lectures = $4, version = version + 1
		 WHERE id = $5 AND version = $6
		 RETURNING version`

	args := []any{
		course.Title,
		course.Published_date,
		course.Runtime,
		pq.Array(course.Lectures),
		course.ID,
		course.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&course.Version)
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

func (m CourseModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM courses WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
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

func ValidateCourse(v *validator.Validator, course *Course) {
	v.Check(course.Title != "", "title", "must be provided")
	v.Check(len(course.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(course.Published_date != "", "published_date", "must be provided")
	v.Check(course.Runtime != 0, "runtime", "must be provided")
	v.Check(course.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(course.Lectures != nil, "lectures", "must be provided")
	v.Check(len(course.Lectures) >= 1, "lectures", "must contain at least 1 lecture")
	v.Check(len(course.Lectures) <= 5, "lectures", "must not contain more than 5 lectures")
	v.Check(validator.Unique(course.Lectures), "lectures", "must not contain duplicate values")
}
