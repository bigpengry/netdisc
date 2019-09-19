package handler

import (
	dblayer "github.com/bigpengry/netdisc/filestore-server/db"
	"github.com/bigpengry/netdisc/filestore-server/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pwd_salt   = "*#890"
	token_salt = "_tokensalt"
)

// SignUpHandler : 处理用户注册请求
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		if len(username) < 3 || len(password) < 5 {
			w.Write([]byte("invaild parameter"))
			return
		}
		encPwd := util.Sha1([]byte(password + pwd_salt))
		suc := dblayer.UserSignUp(username, encPwd)
		if suc {
			w.Write([]byte("Success"))
			return
		}
		w.Write([]byte("Failed"))
		return
	}
}

// SignInHandler : 处理用户登录请求
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// data, err := ioutil.ReadFile("./static/view/signin.html")
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
		// w.Write(data)
		http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPwd := util.Sha1([]byte(password + pwd_salt))
	//1.校验用户名和密码
	pwdChecked := dblayer.UserSignIn(username, encPwd)
	if !pwdChecked {
		w.Write([]byte("sign in failed"))
		return
	}
	//2.生成访问凭证（token）
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("sign in failed"))
		return
	}
	//3.登录成功后重定向到首页
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}

// QuerUserInfoHandler : 查寻用户信息接口
func QuerUserInfoHandler(w http.ResponseWriter,r *http.Request){
	//1.解析请求参数
	r.ParseForm()
	username:=r.Form.Get("username")
	//token:=r.Form.Get("username")
	// //2.验证token是否有效
	// if !IsTokenVaild(token){
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }
	//3.查询用户信息
	user,err:=dblayer.GetUserInfo(username)
	if err!=nil{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	resp:=util.NewRespMsg(0,"OK",user)
	w.Write(resp.JSONBytes())

}


// GenToken : 生成访问凭证
func GenToken(username string) string {
	//md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + token_salt))

	return tokenPrefix + ts[:8]
}

