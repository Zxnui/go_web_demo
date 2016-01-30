package controllers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"go-web-game/conf"
	"go-web-game/models"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
)

var FuncMap = template.FuncMap{
	"AddOne": AddOne,
	"SubStr": SubStr,
}

//生成32位Md5码
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func loginCookie(w http.ResponseWriter, r *http.Request) bool {
	_, err := r.Cookie("Uname")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return true
	}
	return false
}

func pageNotFound(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Println("找不到页面,错误代码:" + err.Error())
		t, _ := template.ParseFiles("views/error/404.html")
		t.Execute(w, nil)
		return true
	}
	return false
}

func CheckError(s string, err error) bool {
	if err != nil {
		log.Println(s + err.Error())
		return true
	}
	return false
}

//string转化为int64
func StringToInt64(s string) (i int64) {
	if s == "" {
		i = 0
	} else {
		st, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			i = 0
		} else {
			i = st
		}
	}
	return
}

//string转化为float64
func StringToFloat64(s string) (i float64) {
	if s == "" {
		i = 0
	} else {
		st, err := strconv.ParseFloat(s, 64)
		if err != nil {
			i = 0
		} else {
			i = st
		}
	}
	fmt.Println(s, i)
	return
}

//template自定义模板方法
func AddOne(a int) int {
	return a + 1
}

//为模板增加方法
func TemParseFiles(t *template.Template, filenames ...string) (*template.Template, error) {
	if len(filenames) == 0 {
		return nil, fmt.Errorf("html/template: no files named in call to ParseFiles")
	}
	for _, filename := range filenames {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		s := string(b)
		_, err = t.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

//截取string
func SubStr(s string, a int) string {
	return string([]byte(s)[0:a]) + "..."
}

//发送邮件
func SendToMail(toEmail, name string) error {
	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

	email, _ := conf.Cfg.GetValue("email", "email")
	password, _ := conf.Cfg.GetValue("email", "passwd")
	host, _ := conf.Cfg.GetValue("email", "host")
	subject, _ := conf.Cfg.GetValue("email", "subject")
	sendName, _ := conf.Cfg.GetValue("email", "sendName")

	from := mail.Address{sendName, email}
	to := mail.Address{name, toEmail}

	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", email, password, hp[0])

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", b64.EncodeToString([]byte(subject)))
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"
	header["Content-Transfer-Encoding"] = "base64"

	//获取用户信息
	userInfo, _, _ := models.GetUserInfo(name)
	EmailMd5 := GetMd5String(toEmail + name + userInfo.Passwd)
	newUser := new(models.User)
	newUser.Email_change = EmailMd5
	err := models.UpdateUser(userInfo.Id, newUser)
	if CheckError("数据库读取失败", err) {
		return err
	}
	body := `<strong>修改密码</strong>
	<p>访问以下网址：<a href="http://` + conf.HttpWeb + `/login/changePasswd?emailCode=` + EmailMd5 + `">http://` + conf.HttpWeb + `/login/changePasswd?emailCode=` + EmailMd5 + `<a><p>
	<p>如果以上链接无法访问，请将该网址复制并粘贴至新的浏览器窗口中。</p>`

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + b64.EncodeToString([]byte(body))

	send_to := strings.Split(toEmail, ";")
	err = smtp.SendMail(host, auth, email, send_to, []byte(message))
	return err
}
