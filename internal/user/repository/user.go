package repository

import (
	"database/sql"
	"errors"
	"github.com/aaanger/ecommerce/internal/user/model"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

//go:generate mockery --name=IUserRepository

type IUserRepository interface {
	CreateUser(email, password, role string) (*model.User, error)
	AuthUser(email, password string) (*model.User, error)
	GetEmail(userID int) string
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(email, password, role string) (*model.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Email:    strings.ToLower(email),
		Password: string(passwordHash),
		Role:     role,
	}
	row := r.db.QueryRow(`INSERT INTO users (email, password_hash, role) VALUES($1, $2, $3) RETURNING id;`, user.Email, user.Password, user.Role)
	err = row.Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) AuthUser(email, password string) (*model.User, error) {
	user := model.User{
		Email: strings.ToLower(email),
	}
	row := r.db.QueryRow(`SELECT id, password_hash, role FROM users WHERE email=$1;`, email)
	err := row.Scan(&user.ID, &user.Password, &user.Role)
	if err != nil {
		return nil, errors.New("invalid email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("wrong password")
	}

	return &user, nil
}

func (r *UserRepository) GetEmail(userID int) string {
	var email string
	row := r.db.QueryRow(`SELECT id FROM users WHERE email = $1;`, userID)
	err := row.Scan(&email)
	if err != nil {
		return ""
	}

	return email
}
