package services

import (
	"clean_api/src/data/db"
	"clean_api/src/data/models"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"errors"
	"time"
)

const (
	ScopeActivation = "activation"
	ScopeAuthentication = "authentication"

)

type Token struct {
	PlainText string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserId    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func GenerateToken(userId int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserId: int64(userId),
		Expiry: time.Now().Add(ttl),
		Scope: scope,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]	
	return token, nil
}

type TokenService struct {
	DB *sql.DB
}

func NewTokenService() *TokenService{
	return &TokenService{
		DB: db.GetDB(),
	}
}

func (t *TokenService) New(userId int, ttl time.Duration, scope string) (*Token, error) {
	q := `
	select expiry  from tokens where user_id = $1 and scope = 'authentication'
	`
	var expiry time.Time
	err := t.DB.QueryRow(q, userId).Scan(&expiry)
	if err == nil && !time.Now().After(expiry) {
		return nil, errors.New("already a token exists for user, if you lost token request a delete")
	} 
	
	token, err := GenerateToken(int(userId), ttl, scope)
	if err != nil {
		return nil, err
	}
	err = t.Insert(token)
	return token, err
}

func (t *TokenService) Insert(token *Token) (error) {
	q := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)
	`
	args := []interface{}{
		token.Hash,
		token.UserId,
		token.Expiry,
		token.Scope,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second *3)
	defer cancel()

	_, err := t.DB.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return nil
}

func (t *TokenService) DeleteAllForUser(scope string, userId int) error {
	q := `DELETE FROM tokens WHERE scope = $1 AND user_id = $2`
	_, err := t.DB.Exec(q, scope, userId)
	if err != nil {
		return err
	}
	return nil
}

func (t *TokenService) GetForToken(tokenScope, plainText string) (*models.Users, error) {
	q := `
	SELECT users.id, users.created_at, users.name, users.email, users.activated, users.version FROM users
	INNER JOIN tokens
	ON users.id = tokens.user_id
	WHERE tokens.hash = $1
	AND tokens.scope = $2
	AND tokens.expiry > $3
	`
	tokensHash := sha256.Sum256([]byte(plainText))
	args := []interface{}{
		tokensHash[:],
		tokenScope, 
		time.Now(),
	}
	var usr models.Users
	err := t.DB.QueryRow(q, args...).Scan(
		&usr.ID,
		&usr.CreatedAt,
		&usr.Name,
		&usr.Email,
		&usr.Activated,
		&usr.Version,
	)
	if err != nil {
		if err.Error() == "sql: no rows in result set"{
			return nil, errors.New("no valid token for user found (maybe user is already activated)")
		}
		return nil, err

	}
	return &usr, nil
}