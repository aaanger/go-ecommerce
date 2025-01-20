package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aaanger/ecommerce/internal/user/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := NewUserRepository(db)

	testCases := []struct {
		name     string
		mock     func()
		user     *model.User
		expected *model.User
		err      bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("test", "test", "user").WillReturnRows(rows)
			},
			user: &model.User{
				Email:    "test",
				Password: "test",
				Role:     "user",
			},
			expected: &model.User{
				ID:       1,
				Email:    "test",
				Password: "test",
				Role:     "user",
			},
			err: false,
		},
		{
			name: "Empty fields",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("", "test", "user").WillReturnRows(rows)
			},
			user: &model.User{
				Email:    "",
				Password: "test",
				Role:     "user",
			},
			err: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			got, err := r.CreateUser(testCase.user.Email, testCase.user.Password, testCase.user.Role)
			if testCase.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}

func TestRepository_AuthUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := NewUserRepository(db)

	testCases := []struct {
		name     string
		mock     func()
		user     *model.User
		expected *model.User
		err      bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "password_hash", "role"}).AddRow(1, "test", "user")
				mock.ExpectQuery("SELECT id, password_hash, role FROM users").
					WithArgs("test").WillReturnRows(rows)
			},
			user: &model.User{
				Email: "test",
			},
			expected: &model.User{
				ID:       1,
				Email:    "test",
				Password: "test",
				Role:     "user",
			},
			err: false,
		},
		{
			name: "Empty fields",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "password_hash", "role"})
				mock.ExpectQuery("SELECT id, password_hash, role FROM users").
					WithArgs("").WillReturnRows(rows)
			},
			user: &model.User{
				Email: "",
			},
			err: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			got, err := r.AuthUser(testCase.user.Email, testCase.user.Password)
			if testCase.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}
