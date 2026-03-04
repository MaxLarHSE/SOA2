package auth

import "database/sql"

type AuthServer struct {
	Db *sql.DB
}
