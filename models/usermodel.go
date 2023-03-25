package models

import (
	"database/sql"
	"go-auth/config"
	"go-auth/entities"
)

type UserModel struct {
	db *sql.DB
}

func NewUserModel() *UserModel {
	conn, err := config.DBConn()

	if err != nil {
		panic(err)
	}

	return &UserModel{
		db: conn,
	}
}

func (u UserModel) Where(user *entities.User, fieldName, fieldValue string) error {
	row, err := u.db.Query("select id,email,username,password from users where "+fieldName+"=? limit 1", fieldValue)

	if err != nil {
		return err
	}

	defer row.Close()

	for row.Next() {
		row.Scan(&user.Id, &user.Email, &user.Username, &user.Password)
	}

	return nil
}

func (u UserModel) Create(user entities.User) (int64, error) {
	result, err := u.db.Exec("insert into users(email,username,password) values(?,?,?)", user.Email, user.Username, user.Password)

	if err != nil {
		return 0, err
	}

	lastInsertId, _ := result.LastInsertId()

	return lastInsertId, nil
}
