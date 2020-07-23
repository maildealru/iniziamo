package preprocessor

import (
	"strings"

	"github.com/maildealru/iniziamo/pkg/parser"
)

type Preprocessor struct{}

func NewPreprocessor() *Preprocessor {
	return &Preprocessor{}
}

func (p *Preprocessor) Preprocess(info *parser.Info) {
	preprocessImports(info.Imports)
	preprocessClients(info.Clients)
}

func preprocessImports(imports []parser.Import) {
	const tImpPath = `"github.com/maildealru/iniziamo/pkg/iniziamo"`
	for _, imp := range imports {
		if imp.Path == tImpPath {
			return
		}
	}
	imports = append(
		imports,
		parser.Import{
			Path: tImpPath,
		},
	)
}

func preprocessClients(clients []parser.Client) {
	for i := 0; i < len(clients); i++ {
		client := &clients[i]

		client.GoIntName = strings.ToUpper(client.GoName[:1]) + client.GoName[1:]
		client.GoImplName = client.GoName + "Impl"
		client.GoImplConstrName = "New" + client.GoIntName
		client.Config.GoName = client.GoIntName + "Config"

		for j := 0; j < len(client.Calls); j++ {
			call := &client.Calls[j]

			call.Request.GoName = client.GoName + call.GoName + "Request"
			call.Request.GoConstrName = client.GoIntName + call.GoName + "Request"

			for k := 0; k < len(call.Request.QueryParams); k++ {
				p := &call.Request.QueryParams[k]

				p.GoFieldName = strings.ToLower(p.GoName[:1]) + p.GoName[1:]
				p.GoSetFlagName = p.GoFieldName + "Set"
				p.GoSetterName = p.GoName
			}

			call.Response.GoName = client.GoName + call.GoName + "Response"
			call.Response.GoConstrName = client.GoIntName + call.GoName + "Response"
		}
	}
}
