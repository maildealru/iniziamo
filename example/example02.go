package example

//go:generate iniziamo

import (
	"context"

	"github.com/maildealru/iniziamo/pkg/iniziamo"
)

//https://oauth.mail.ru/docs

//iniziamo:Client
type mailruOAuth2Client interface {
	//iniziamo:Call:{"method":"POST","path":"/token"}
	ExchangeAuthCodeToToken(
		//iniziamo:Context
		context.Context,

		//iniziamo:Request
		struct {
			//iniziamo:FormParam:{"name":"code"}
			Code string

			//iniziamo:FormParam:{"name":"grant_type","const":"authorization_code"}
			GrantType string

			//iniziamo:FormParam:{"name":"redirect_uri","config":"#"}
			RedirectURL string

			//iniziamo:BasicAuth:{"config":"#"}
			ClientID, ClientSecret string
		},

		//iniziamo:Response
		struct {
			//iniziamo:ResponseBodyParser
			BodyParser func(
				//iniziamo:Context
				context.Context,

				//iniziamo:Status
				int,

				//iniziamo:ResponseBody
				[]byte,
			) (
				map[string]string, error,
			)
		},
	)

	//iniziamo:Call:{"path":"/userinfo"}
	GetUserInfo(
		//iniziamo:Context
		context.Context,

		//iniziamo:Request
		struct {
			//iniziamo:QueryParam:{"name":"access_token"}
			AccessToken string
		},

		//iniziamo:Response
		struct {
			//iniziamo:ResponseBody:{"Content-Type":"JSON"}
			Body struct {
				ExpiresIn   uint32 `json:"expires_in"`
				AccessToken string `json:"access_token"`
			}
		},
	)

	//iniziamo:Call:{"method":"POST","path":"/token"}
	RenewToken(
		//iniziamo:Context
		context.Context,

		//iniziamo:Request
		struct {
			//iniziamo:FormParam:{"name":"client_id","config":"#"}
			ClientID string

			//iniziamo:FormParam:{"name":"grant_type","const":"refresh_token"}
			GrantType string

			//iniziamo:FormParam:{"name":"refresh_token"}
			RefreshToken string
		},

		//iniziamo:Response
		struct {
			//iniziamo:StatusValidator:{"status":200}
			StatusValidator iniziamo.StatusValidator

			//iniziamo:ResponseBody:{"Content-Type":"JSON"}
			Body struct {
				ExpiresIn   uint32 `json:"expires_in"`
				AccessToken string `json:"access_token"`
			}
		},
	)
}
