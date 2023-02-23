package main

import (
	"errors"
	"fmt"
	"net/http"

	data "github.com/sbeknur/go-final/internal/data"
	"github.com/sbeknur/go-final/internal/validator"
)

func (app *application) createInstructorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Age       int32  `json:"age"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	instructors := &data.Instructors{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Age:       input.Age,
	}

	err = app.models.Instructors.Insert(instructors)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/instructors/%d", instructors.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"instructors": instructors}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showInstructorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	instructors, err := app.models.Instructors.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"instructor": instructors}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listInstructorsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName string
		LastName  string
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.FirstName = app.readString(qs, "firstname", "")
	input.LastName = app.readString(qs, "lastname", "")

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.SortSafeList = []string{"id", "firstname", "lastname", "age", "-id", "-firstname", "-lastname", "-age"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	instructors, metadata, err := app.models.Instructors.GetAll(input.FirstName, input.LastName, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"instructors": instructors, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
