package dispatcher

import (
	"encoding/json"
	"gateway/auth"
	"gateway/model"
	"io/ioutil"
	"net/http"
	"pkg"
)

const PREFIX = "/user"

// Dispense requests of user, whose prefix is /user/
// POST /user/login         login
// POST /user/logout        logout
func DispenseUser(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	suffix := path[len(PREFIX):]
	method := string([]byte(req.Method))
	if method == "OPTIONS" || method == "HEAD" {
		// do nothing
		return
	}
	pkg.Logd("to dispatch: path=" + path + " suffix=" + suffix)
	if suffix == "/login" && method == "POST" {
		handleLogin(w, req)
	} else if suffix == "/logout" && method == "POST" {
		handleLogout(w, req)
	} else {
		pkg.Logd("unknown request: " + method + " " + path)
		w.WriteHeader(404)
	}
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Nickname string `json:"name"`
	Account  string `json:"account"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

func handleLogin(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	var form LoginForm
	if len(body) == 0 {
		pkg.BadRequest(w, "没有请求体")
		return
	}
	if err := json.Unmarshal(body, &form); err != nil {
		pkg.BadRequest(w, "请求体json参数解析错误!")
		pkg.Logi("parse login form error:" + err.Error())
		return
	}
	email := form.Email
	password := form.Password
	if email == "" || password == "" {
		pkg.BadRequest(w, "缺少参数!")
		return
	}
	pkg.Logi("handleLogin: email=" + email)
	user, _ := model.GetUserByEmail(email)

	if user == nil || user.Password != password {
		// fail
		pkg.BadRequest(w, "账号或密码错误!")
	} else {
		token := auth.GenerateToken(user.Id, user.Email)
		response := LoginResponse{Nickname: user.Nickname, Account: user.Account, Email: user.Email, Token: token}
		pkg.WriteJsonOfResponce(w, response)
		model.UpdateLastLoginTimeOfUser(user.Id)
		pkg.Logi("Login succussfully: email=" + email)
	}
}

func handleLogout(w http.ResponseWriter, req *http.Request) {

}
