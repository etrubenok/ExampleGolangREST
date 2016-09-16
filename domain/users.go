package domain

import (
  "github.com/jackc/pgx"
)

type UserEntity struct {
  Id string `json:"id"`
  Name string `json:"name"`
  Phone string `json:"phone"`
}

func GetUserById(conn DbConnection, id string) (*UserEntity, error) {
  userEntity := new(UserEntity)
  err := conn.QueryRow("select id, name, phone from users where id=$1", id).Scan(
    &userEntity.Id, 
    &userEntity.Name,
    &userEntity.Phone)

  return userEntity, err
}

type DbRow interface {
  Scan(dest ...interface{}) error
}

type DbConnection interface {
  QueryRow(query string, args ...interface{}) DbRow
}

type NewDbConnection struct {
  Db *pgx.ConnPool
}

func (dbConnection *NewDbConnection) QueryRow(query string, args ...interface{}) DbRow {
  return dbConnection.Db.QueryRow(query, args...)
}