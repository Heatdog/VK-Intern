// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	film_model "github.com/Heater_dog/Vk_Intern/internal/models/film"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// FilmsRepository is an autogenerated mock type for the FilmsRepository type
type FilmsRepository struct {
	mock.Mock
}

// DeleteActorsFormFilm provides a mock function with given fields: ctx, filmId
func (_m *FilmsRepository) DeleteActorsFormFilm(ctx context.Context, filmId uuid.UUID) error {
	ret := _m.Called(ctx, filmId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteActorsFormFilm")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, filmId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteFilm provides a mock function with given fields: ctx, filmId
func (_m *FilmsRepository) DeleteFilm(ctx context.Context, filmId uuid.UUID) error {
	ret := _m.Called(ctx, filmId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteFilm")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, filmId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetFilms provides a mock function with given fields: ctx, order, orderType
func (_m *FilmsRepository) GetFilms(ctx context.Context, order string, orderType string) ([]film_model.Film, error) {
	ret := _m.Called(ctx, order, orderType)

	if len(ret) == 0 {
		panic("no return value specified for GetFilms")
	}

	var r0 []film_model.Film
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]film_model.Film, error)); ok {
		return rf(ctx, order, orderType)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []film_model.Film); ok {
		r0 = rf(ctx, order, orderType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]film_model.Film)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, order, orderType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFilmsWithActor provides a mock function with given fields: ctx, userID
func (_m *FilmsRepository) GetFilmsWithActor(ctx context.Context, userID uuid.UUID) ([]film_model.Film, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetFilmsWithActor")
	}

	var r0 []film_model.Film
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]film_model.Film, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []film_model.Film); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]film_model.Film)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertFilm provides a mock function with given fields: ctx, film
func (_m *FilmsRepository) InsertFilm(ctx context.Context, film *film_model.FilmInsert) (string, error) {
	ret := _m.Called(ctx, film)

	if len(ret) == 0 {
		panic("no return value specified for InsertFilm")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *film_model.FilmInsert) (string, error)); ok {
		return rf(ctx, film)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *film_model.FilmInsert) string); ok {
		r0 = rf(ctx, film)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *film_model.FilmInsert) error); ok {
		r1 = rf(ctx, film)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchFilms provides a mock function with given fields: ctx, searchQuery
func (_m *FilmsRepository) SearchFilms(ctx context.Context, searchQuery string) ([]film_model.Film, error) {
	ret := _m.Called(ctx, searchQuery)

	if len(ret) == 0 {
		panic("no return value specified for SearchFilms")
	}

	var r0 []film_model.Film
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]film_model.Film, error)); ok {
		return rf(ctx, searchQuery)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []film_model.Film); ok {
		r0 = rf(ctx, searchQuery)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]film_model.Film)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, searchQuery)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateFilmActors provides a mock function with given fields: ctx, filmId, actorsID
func (_m *FilmsRepository) UpdateFilmActors(ctx context.Context, filmId uuid.UUID, actorsID []film_model.Id) error {
	ret := _m.Called(ctx, filmId, actorsID)

	if len(ret) == 0 {
		panic("no return value specified for UpdateFilmActors")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, []film_model.Id) error); ok {
		r0 = rf(ctx, filmId, actorsID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateFilmDescription provides a mock function with given fields: ctx, filmId, description
func (_m *FilmsRepository) UpdateFilmDescription(ctx context.Context, filmId uuid.UUID, description string) error {
	ret := _m.Called(ctx, filmId, description)

	if len(ret) == 0 {
		panic("no return value specified for UpdateFilmDescription")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) error); ok {
		r0 = rf(ctx, filmId, description)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateFilmRating provides a mock function with given fields: ctx, filmId, rating
func (_m *FilmsRepository) UpdateFilmRating(ctx context.Context, filmId uuid.UUID, rating string) error {
	ret := _m.Called(ctx, filmId, rating)

	if len(ret) == 0 {
		panic("no return value specified for UpdateFilmRating")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) error); ok {
		r0 = rf(ctx, filmId, rating)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateFilmReleaseDate provides a mock function with given fields: ctx, filmId, date
func (_m *FilmsRepository) UpdateFilmReleaseDate(ctx context.Context, filmId uuid.UUID, date string) error {
	ret := _m.Called(ctx, filmId, date)

	if len(ret) == 0 {
		panic("no return value specified for UpdateFilmReleaseDate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) error); ok {
		r0 = rf(ctx, filmId, date)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateFilmTitle provides a mock function with given fields: ctx, filmId, title
func (_m *FilmsRepository) UpdateFilmTitle(ctx context.Context, filmId uuid.UUID, title string) error {
	ret := _m.Called(ctx, filmId, title)

	if len(ret) == 0 {
		panic("no return value specified for UpdateFilmTitle")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) error); ok {
		r0 = rf(ctx, filmId, title)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewFilmsRepository creates a new instance of FilmsRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFilmsRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *FilmsRepository {
	mock := &FilmsRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
