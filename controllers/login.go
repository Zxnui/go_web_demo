package controllers

import (
	"go-web-game/models"
	"html/template"
	"net/http"
	"strings"
	"time"
)

//登录
func Login(w http.ResponseWriter, r *http.Request) {
	tel, err := template.ParseFiles("views/login.html", "views/common/head.tpl")
	if pageNotFound(w, err) {
		return
	}
	isExit := strings.EqualFold(r.FormValue("exit"), "true")
	if isExit {
		cookie, _ := r.Cookie("Uname")
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	_, err = r.Cookie("Uname")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data := make(map[string]string)
	tel.Execute(w, data)

}

//保存登录数据
func LoginPost(w http.ResponseWriter, r *http.Request) {
	submit := r.PostFormValue("submit")
	if strings.EqualFold(submit, "1") {
		http.Redirect(w, r, "/loginInfo", http.StatusFound)
		return
	}
	data := make(map[string]string)
	login, err := template.ParseFiles("views/login.html", "views/common/head.tpl")
	if pageNotFound(w, err) {
		return
	}

	uname := r.PostFormValue("uname")
	pwd := r.PostFormValue("pwd")
	checkbox := r.PostFormValue("autoLogin") == "on"

	user, has, err := models.GetUserInfo(uname)
	_ = CheckError("查询用户失败:", err)

	if has {
		if strings.EqualFold(user.Passwd, GetMd5String(pwd)) {
			if checkbox {
				newCookie := http.Cookie{
					Name:   http.CanonicalHeaderKey("uname"),
					Value:  uname,
					MaxAge: 1<<31 - 1,
				}
				http.SetCookie(w, &newCookie)
			} else {
				newCookie := http.Cookie{
					Name:   http.CanonicalHeaderKey("uname"),
					Value:  uname,
					MaxAge: 0,
				}
				http.SetCookie(w, &newCookie)
			}
			http.Redirect(w, r, "/", http.StatusFound)
			return
		} else {
			data["error"] = "密码错误"
			data["uname"] = uname
			login.Execute(w, data)
			return
		}
	} else {
		data["error"] = "用户不存在"
		login.Execute(w, data)
		return
	}
}

//注册
func LoginInfo(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	user := new(models.User)
	t, err := template.ParseFiles("views/login_info.html", "views/common/head.tpl")
	if pageNotFound(w, err) {
		return
	}
	l, err := template.ParseFiles("views/login.html", "views/common/head.tpl")
	if pageNotFound(w, err) {
		return
	}
	if r.Method == "GET" {
		data["user"] = user
		data["error"] = ""
		t.Execute(w, data)
		return
	} else {
		user.Uname = r.PostFormValue("uname")
		user.Email = r.PostFormValue("email")
		user.Passwd = GetMd5String(r.PostFormValue("pwd"))
		user.Status = 1
		user.Subtime = time.Now()
		s := models.GetUserNameOrEmail(user.Uname, user.Email)
		if len(s) > 0 {
			data["error"] = s
			data["user"] = user
			t.Execute(w, data)
			return
		}
		err = models.InsertUser(user)
		if CheckError("数据库连接失败:", err) {
			data["user"] = user
			data["error"] = "网站连接失败，请稍后重试"
			t.Execute(w, data)
			return
		}
		data["success"] = "注册成功，请登录"
		l.Execute(w, data)
		return
	}

}

//发送邮件，找回密码
func LoginForget(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	user := new(models.User)
	t, err := template.ParseFiles("views/login_forget.html", "views/common/head.tpl")
	if pageNotFound(w, err) {
		return
	}
	l, err := template.ParseFiles("views/login.html", "views/common/head.tpl")
	if pageNotFound(w, err) {
		return
	}
	if r.Method == "GET" {
		data["user"] = user
		data["error"] = ""
		t.Execute(w, data)
		return
	} else {
		user.Uname = r.PostFormValue("uname")
		user.Email = r.PostFormValue("email")
		s := models.GetUserNameAndEmail(user.Uname, user.Email)
		if len(s) > 0 {
			data["error"] = s
			data["user"] = user
			t.Execute(w, data)
			return
		}

		//发送邮件
		err = SendToMail(user.Email, user.Uname)
		if CheckError("邮件发送失败:", err) {
			data["error"] = "邮件发送失败，稍后再试"
			data["user"] = user
			t.Execute(w, data)
			return
		}
		data["success"] = "请登录邮箱，修改密码"
		l.Execute(w, data)
		return
	}
}

//修改密码
func LoginChangePasswd(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	user := new(models.User)
	t, err := template.ParseFiles("views/login_passwd_change.html", "views/common/head.tpl")
	if pageNotFound(w, err) {
		return
	}
	l, err := template.ParseFiles("views/login.html", "views/common/head.tpl")
	if pageNotFound(w, err) {
		return
	}
	if r.Method == "GET" {
		emailCode := r.FormValue("emailCode")
		user.Email_change = emailCode
		has, _, _ := models.SelectUser(user)
		if has {
			data["emailCode"] = emailCode
			data["user"] = user
			data["error"] = ""
			t.Execute(w, data)
			return
		} else {
			data["error"] = "密码修改失败,请重试"
			l.Execute(w, data)
			return
		}

	} else {
		pwd := r.PostFormValue("pwd")
		user.Email_change = r.PostFormValue("emailCode")
		has, userInfo, _ := models.SelectUser(user)
		if has {
			id := userInfo.Id
			user = new(models.User)
			user.Passwd = GetMd5String(pwd)
			user.Email_change = GetMd5String(userInfo.Uname + userInfo.Email + userInfo.Passwd)
			err := models.UpdateUser(id, user)
			if CheckError("密码修改失败", err) {
				data["error"] = "密码修改失败,请重试"
				l.Execute(w, data)
				return
			} else {
				data["success"] = "密码修改成功，请登录"
				l.Execute(w, data)
				return
			}
		} else {
			data["error"] = "密码修改失败,请重试"
			l.Execute(w, data)
			return
		}
	}
}
