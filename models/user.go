package models

import (
	"time"
)

type User struct {
	Id           int64
	Uname        string
	Passwd       string
	Status       int
	Email        string
	Subtime      time.Time
	Email_change string
}

func SelectUser(user *User) (bool, *User, error) {
	has, err := Xorm.Get(user)
	return has, user, err
}

func UpdateUser(id int64, user *User) error {
	_, err := Xorm.Id(id).Update(user)
	return err
}

func GetUserInfo(uname string) (*User, bool, error) {
	user := new(User)
	has, err := Xorm.Where("uname=?", uname).Get(user)
	return user, has, err
}

func GetUserNameOrEmail(uname, email string) string {
	s := ""
	user := new(User)
	hasname, _ := Xorm.Where("uname=?", uname).Get(user)
	user = new(User)
	hasemail, _ := Xorm.Where("email=?", email).Get(user)
	if hasname {
		s = s + "用户名已存在"
	}
	if hasemail {
		if len(s) == 0 {
			s = "邮箱已存在"
		} else {
			s = s + ",邮箱已存在"
		}

	}

	return s
}

func GetUserNameAndEmail(uname, email string) string {
	s := ""
	user := new(User)
	hasname, _ := Xorm.Where("uname=?", uname).Get(user)
	user = new(User)
	hasemail, _ := Xorm.Where("email=?", email).Get(user)
	if !hasname {
		s = s + "用户名不存在"
	}
	if !hasemail {
		if len(s) == 0 {
			s = "邮箱不存在"
		} else {
			s = s + ",邮箱不存在"
		}

	}

	return s
}

func InsertUser(user *User) error {
	_, err := Xorm.Insert(user)
	return err
}
