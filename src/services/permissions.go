package services

import (
	"clean_api/src/data/db"
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type PermissionService struct {
	DB *sql.DB
}

func NewPermissionService() *PermissionService{
	return &PermissionService{
		DB: db.GetDB(),
	}
}


type Permissions []string


func (p *PermissionService) GetAllForUser(userId int) (Permissions, error) {
	q := `
	SELECT permissions.code
	FROM permissions
	INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
	INNER JOIN users ON users_permissions.user_id = users.id
	WHERE users.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := p.DB.QueryContext(ctx, q, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ps Permissions
	for rows.Next() {
		var permission string
		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}
		ps = append(ps, permission)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ps, nil
}

func (p *PermissionService) AddForUser(userId int, roles ...string) error {
	q := `
	INSERT INTO users_permissions
	SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)
	`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, q, userId, pq.Array(roles))
	return err
}