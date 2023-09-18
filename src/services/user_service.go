package services

import (
	"clean_api/src/api/dto"
	"clean_api/src/data/db"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB *sql.DB
}

func NewUserService() *UserService{
	return &UserService{
		DB: db.GetDB(),
	}
}


func (u *UserService) GetByEmail(email string) (*dto.UserResponse, error) {
	q := `
	SELECT id, created_at, name, email, activated, version, password
	FROM users
	WHERE email = $1`
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	response := &dto.UserResponse{}

	err := u.DB.QueryRowContext(ctx, q, email).Scan(
		&response.ID,
		&response.CreatedAt,
		&response.Name,
		&response.Email,
		&response.Activated,
		&response.Version,
		&response.Password,
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}


func (u *UserService) Create(req *dto.CreateUser) (*dto.UserResponse, error) {
	q := `
	INSERT INTO users (name, email, password, activated)
	VALUES ($1, $2, $3, $4)
	RETURNING email`
	bs, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	req.Password = string(bs)

	args := []interface{}{
		req.Name,
		req.Email,
		req.Password,
		false,
	}
	tx, err := u.DB.Begin()
	if err != nil {
		return nil, err
	}
	var email string
	err = tx.QueryRow(q, args...).Scan(&email)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return u.GetByEmail(email)

}

func (u *UserService) Update(req *dto.UpdateUser, id int, version int) (*dto.UserResponse, error) {
	q := `
	UPDATE users
	SET name = $1, email = $2, activated = $3, version = version + 1
	WHERE id = $4 AND version = $5
	RETURNING id`

	args := []interface{}{
		req.Name,
		req.Email,
		req.Activated,
		id, 
		version,
	}
	
	lol := 0
	err := u.DB.QueryRow(q, args...).Scan(&lol)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil,errors.New("no rows")
		}
		return nil, err
	}
	fmt.Println(lol)
	return u.GetById(id)
}

func (u *UserService) GetById(id int) (*dto.UserResponse, error){
	q := `
	SELECT id, created_at, name, email, activated, version 
	FROM users
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	response := &dto.UserResponse{}
	err := u.DB.QueryRowContext(ctx, q, id).Scan(
		&response.ID,
		&response.CreatedAt,
		&response.Name,
		&response.Email,
		&response.Activated,
		&response.Version,
	)
	if err != nil {
		return nil, err
	}
	return response, nil

}


