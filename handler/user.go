package handler

import (
	"fmt"
	"go-cloud/db"
	"go-cloud/util"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pwd_salt = "#890"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("invalid parameter"))
		return
	}

	encPwd := util.Sha1([]byte(passwd + pwd_salt))
	suc := db.UserSignup(username, encPwd)
	if suc {
		w.Write([]byte("success"))
	} else {
		w.Write([]byte("failed"))
	}

}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPwd := util.Sha1([]byte(password + pwd_salt))
	pwdChecked := db.UserSignin(username, encPwd)
	if !pwdChecked {
		w.Write([]byte("failed"))
		return
	}
	token := GenToken(username)
	upRes := db.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("failed"))
		return
	}
	w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
}

func GenToken(username string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}
