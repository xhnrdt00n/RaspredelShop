package db

import (
	"context"
	"log"
)

type User struct {
	Id       string `json:"-"`
	Login    string `json:"Login"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

func (db *PgxCon) AddUser(u User) (string, error) {
	connCtx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	var id string

	_ = db.pgConn.QueryRow(connCtx, "SELECT id from users WHERE id=$1", u.Id).Scan(&id)
	if id != "" {
		return id, nil
	}

	tx, _ := db.pgConn.Begin(connCtx)
	err := tx.QueryRow(connCtx,
		"INSERT INTO ShopUser (LoginName,Passhash,Email) VALUES ($1,$2,$3) returning id",
		u.Login, u.Password, u.Email).Scan(&id)

	if err != nil {
		tx.Rollback(connCtx)
		return "", err
	}

	tx.Commit(connCtx)
	return id, nil
}

func (db *PgxCon) GetUser(login, email string) (*User, error) {
	var user User
	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()
	err := db.pgConn.QueryRow(connCtx, "SELECT id,LoginName,Passhash,Email FROM ShopUser WHERE LoginName=$1 AND Email=$2", login, email).
		Scan(&user.Id, &user.Login, &user.Password, &user.Email)
	log.Println(user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *PgxCon) SetRoleForUser(uid string, roleid string) error {
	connCtx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	var roid string

	_ = db.pgConn.QueryRow(connCtx, "SELECT role_id FROM RoleUsers WHERE user_id=$1", uid).Scan(&roid)
	if roid != "" {
		return nil
	}

	tx, _ := db.pgConn.Begin(connCtx)
	err := tx.QueryRow(connCtx,
		"INSERT INTO RoleUsers (user_id,role_id) VALUES ($1,$2) returning id",
		uid, roleid).Scan(&roid)

	if err != nil {
		tx.Rollback(connCtx)
		return err
	}

	tx.Commit(connCtx)
	return nil
}

func (db *PgxCon) GetRoleByUserID(id string) (int, error) {
	var roleID int
	connCtx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()

	err := db.pgConn.QueryRow(connCtx, "SELECT role_id FROM RoleUsers WHERE user_id=$1", id).
		Scan(&roleID)
	log.Println(roleID)
	if err != nil {
		return 0, err
	}
	return roleID, nil
}
