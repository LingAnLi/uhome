package handler

import (
	"net/http"
	Error "uhome/UHomeWeb/error"
)

func ToObtainSessionId(r *http.Request)(sessionId string,err error){
	cookie,err:=r.Cookie("userlogin")
	 if err!=nil|| cookie.Value=="" {
		err=Error.NewError(404,"缓存为空")

	}
	sessionId=cookie.Value
return
}
