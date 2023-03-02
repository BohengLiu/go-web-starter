package model

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:user_tbl"`

	ID     int64  `bun:"id,pk,notnull"`
	Name   string `bun:"name"`
	Signup string `bun:"signup"`
}
