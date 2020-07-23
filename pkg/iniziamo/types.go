package iniziamo

import (
	"net/http"
)

type Scheme string

const (
	SchemeHTTP  Scheme = "http"
	SchemeHTTPS Scheme = "https"
)

//FIXED TYPE

type StatusValidator func(status int) error

type HeaderValidator func(value string) error

type CookieValidator func(cookie *http.Cookie) error

type BodyValidator func(contentType string, contentLength int, body []byte) error

//NON-FIXED TYPE

type HeaderParser func(string) (interface{}, error)

type CookieParser func(*http.Cookie) (interface{}, error)

type BodyParser func(contentType string, contentLength int, body []byte) (interface{}, error)
