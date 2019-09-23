package handler

import(
	"net/http"
	"netdisc/filestore-server/util"
)

// HTTPInterceptor : http请求拦截器
func HTTPInterceptor(hf http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter,r *http.Request){
			r.ParseForm()
			username := r.Form.Get("username")
			token := r.Form.Get("token")
			if len(username) < 3 || !IsTokenValid(token) {
				// token校验失败则跳转到直接返回失败提示
				resp := util.NewRespMsg(-1,"token无效",nil)
				w.Write(resp.JSONBytes())
				return
			}
			hf(w,r)
	})
}

// IsTokenValid : 检验token是否有效
func IsTokenValid(token string)bool{
	if len(token)!=40{
		return false
	}
	//1.检验token是否超时
	//2.查询token是否存在
	//3.比较token是否一致
	return true
}