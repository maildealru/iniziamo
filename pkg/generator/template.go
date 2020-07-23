package generator

const (
	IniziamoTemplate = `package {{.PkgName}}
// Code generated by iniziamo. DO NOT EDIT.

import (
	"fmt"
	"strconv"

	{{range $import := .Imports}}
		{{$import.Name}} {{$import.Path}}
	{{end}}

	"github.com/go-resty/resty/v2"
)

{{range $client := .Clients}}
	type {{$client.Config.GoName}} struct {
		Scheme iniziamo.Scheme
		Host   string
		Port   int
	}

	type {{$client.GoIntName}} interface {
		{{range $call := $client.Calls}}
			{{$call.GoName}}(
				{{$call.GoCtxType}}, *{{$call.Request.GoName}},
			) (
				*{{$call.Response.GoName}}, error,
			)
		{{end}}
	}

	type {{$client.GoImplName}} struct {
		conf {{$client.Config.GoName}}

		baseURL string
	}

	func {{$client.GoImplConstrName}}(conf {{$client.Config.GoName}}) {{$client.GoIntName}} {
		return &{{$client.GoImplName}}{
			conf: conf,
			baseURL: fmt.Sprintf(
				"%s://%s%s%s",
				conf.Scheme, conf.Host,
				func() string {
					if conf.Port == 0 {
						return ""
					}
					return ":"
				}(),
				func() string {
					if conf.Port == 0 {
						return ""
					}
					return strconv.Itoa(conf.Port)
				}(),
			),
		}
	}


	{{range $call := $client.Calls}}
		type {{$call.Request.GoName}} struct {
			{{range $param := $call.Request.QueryParams}}
				{{$param.GoFieldName}} {{$param.GoType}}
			{{- end}}

			{{range $param := $call.Request.QueryParams}}
				{{if $param.Optional -}}
					{{$param.GoSetFlagName}} bool
				{{- end}}
			{{- end}}
		}

		func {{$call.Request.GoConstrName}}() *{{$call.Request.GoName}} {
			return &{{$call.Request.GoName}}{
				//TODO
			}
		}

		{{range $param := $call.Request.QueryParams}}
			func (r *{{$call.Request.GoName}}) {{$param.GoSetterName}}(value {{$param.GoType}}) *{{$call.Request.GoName}} {
				r.{{$param.GoFieldName}} = value
				{{if $param.Optional -}}
					r.{{$param.GoSetFlagName}} = true
				{{- end}}
				return r
			}
		{{- end}}

		type {{$call.Response.GoName}} struct {
			//TODO
		}

		func {{$call.Response.GoConstrName}}() *{{$call.Response.GoName}} {
			return &{{$call.Response.GoName}}{
				//TODO
			}
		}

		func (c *{{$client.GoImplName}}) {{$call.GoName}}(
			ctx {{$call.GoCtxType}}, request *{{$call.Request.GoName}},
		) (
			response *{{$call.Response.GoName}}, err error,
		) {
			_, err = resty.
				New().R().
				SetContext(ctx).
				Execute("{{$call.Method}}", c.baseURL + "{{$call.PathFmt}}")
			if err != nil {
				return nil, err
			}

			//TODO
			return nil, nil
		}
	{{end}}
{{end}}
`
)
