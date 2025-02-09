package repository

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	email := "test@example.com"
	password := "password123"
	role := "user"

	
	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(strings.ToLower(email), sqlmock.AnyArg(), role).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	
	user, err := repo.CreateUser(email, password, role)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, strings.ToLower(email), user.Email)
	assert.Equal(t, role, user.Role)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_HashError(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	_, err = repo.CreateUser("test@example.com", "", "user") 
	assert.Error(t, err)
}

func TestCreateUser_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	email := "test@example.com"
	password := "password123"
	role := "user"

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(strings.ToLower(email), sqlmock.AnyArg(), role). 
		WillReturnError(errors.New("db error"))

	_, err = repo.CreateUser(email, password, role)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet()) 
}

func TestAuthUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	email := "test@example.com"
	password := "password123"
	role := "user"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mock.ExpectQuery(`SELECT id, password_hash, role FROM users WHERE email=`).
		WithArgs(strings.ToLower(email)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "password_hash", "role"}).
			AddRow(1, string(hashedPassword), role))

	user, err := repo.AuthUser(email, password)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, role, user.Role)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthUser_InvalidEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectQuery(`SELECT id, password_hash, role FROM users WHERE email=`).
		WithArgs("invalid@example.com").
		WillReturnError(errors.New("sql: no rows in result set"))

	_, err = repo.AuthUser("invalid@example.com", "password")
	assert.Error(t, err)
	assert.Equal(t, "invalid email", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthUser_WrongPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	email := "test@example.com"
	wrongPassword := "wrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	mock.ExpectQuery(`SELECT id, password_hash, role FROM users WHERE email=`).
		WithArgs(strings.ToLower(email)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "password_hash", "role"}).
			AddRow(1, string(hashedPassword), "user"))

	_, err = repo.AuthUser(email, wrongPassword)
	assert.Error(t, err)
	assert.Equal(t, "wrong password", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}
