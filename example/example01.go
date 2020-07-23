package example

//go:generate iniziamo

import (
	"context"
)

//iniziamo:Client
type authServiceClient interface {
	//iniziamo:Call:{"method":"POST","path":"/api/{version}/login"}
	Login(
		//iniziamo:Context
		context.Context,
		//iniziamo:Request
		struct {
			//iniziamo:PathParam:{"name":"version"}
			Version string

			//iniziamo:QueryParam:{"name":"continue","required":false}
			ContinueURL string

			//iniziamo:FormParam:{"name":"login"}
			Login string

			//iniziamo:FormParam:{"name":"password"}
			Password string

			//iniziamo:Header:{"name":"X-API-Version"}
			XAPIVersion uint32
		},
		//iniziamo:Response
		struct {
			//iniziamo:Status
			Status int

			//iniziamo:Cookie:{"name":"x_session_cookie"}
			SessionCookie string

			//iniziamo:Body:{"encoding":"json"}
			UserInfo struct {
				UserID  uint64                 `json:"user_id"`
				Profile map[string]interface{} `json:"profile"`
			}
		},
	)

	//iniziamo:Call:{"method":"DELETE","path":"/api/{version}/logout"}
	Logout(
		//iniziamo:Context
		context.Context,
		//iniziamo:Request
		struct {
			//iniziamo:Cookie:{"name":"x_session_cookie"}
			SessionCookie string
		},
		//iniziamo:Response
		struct {
			//iniziamo:Status
			Status int
		},
	)
}
