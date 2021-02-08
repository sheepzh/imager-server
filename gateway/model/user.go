package model

import (
	"database/sql"
	"fmt"
	"pkg"
	"time"
)

type User struct {
	Id               int64        `json:"id"`
	Nickname         string       `json:"nickname"`
	Email            string       `json:"email"`
	Account          string       `json:"account"`
	Password         string       `json:"-"`
	LastLoginTimeSql sql.NullTime `json:"-"`
	LastLoginTime    int64        `json:"lastLoginTime"`
	CreatedTimeSql   sql.NullTime `json:"-"`
	CreatedTime      int64        `json:"createdTime"`
}

var DB *sql.DB

func InitUserDb() error {
	// init the db connection
	// username:password@tcp(ip:port)/database
	connSettings := pkg.GetArgRequired("-ud", "--user-databse")
	if connSettings == "" {
		pkg.Loge("db connection format: username:password@tcp(ip:port)/database")
	}

	db, err := sql.Open("mysql", connSettings+"?parseTime=true")
	DB = db
	if err != nil {
		pkg.Loge("connected to mysql failed: " + err.Error())
		return err
	}
	return nil
}

// query by email
// @param email      the email for the user
// @return nil if not found, else the user info
// @throws while the db error happens
func GetUserByEmail(email string) (*User, error) {
	user := new(User)
	row := DB.QueryRow(`select id, nickname, email, account, password, last_login_time, created_time from user where email = ? and deleted_flag is false`, email)
	if err := row.Scan(&user.Id, &user.Nickname, &user.Email, &user.Account, &user.Password, &user.LastLoginTimeSql, &user.CreatedTimeSql); err != nil {
		pkg.Logi("db error while query user by email:" + err.Error())
		return nil, nil
	}

	if user.LastLoginTimeSql.Valid {
		user.LastLoginTime = user.LastLoginTimeSql.Time.Unix()
	}

	if user.CreatedTimeSql.Valid {
		user.CreatedTime = user.CreatedTimeSql.Time.Unix()
	}

	return user, nil
}

// update the last login time of user
// @param userId     id of the user
func UpdateLastLoginTimeOfUser(userId int64) {
	_, err := DB.Exec(fmt.Sprintf("update user set last_login_time = '%s' where id = %d and deleted_flag is false", time.Now().Local().Format("2006-01-02 15:04:05"), userId))
	if err != nil {
		pkg.Logw(err.Error())
	}
}
