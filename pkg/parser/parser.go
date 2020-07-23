package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"strings"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseFile(path string) *Info {
	fs := token.NewFileSet()

	f, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	info := parseFileInfo(f, fs)
	info.Dir = filepath.Dir(path)
	info.Name = filepath.Base(path)

	return info
}

const (
	docPrefix = "iniziamo:"

	DocClient     = docPrefix + "Client"
	DocCall       = docPrefix + "Call"
	DocContext    = docPrefix + "Context"
	DocRequest    = docPrefix + "Request"
	DocResponse   = docPrefix + "Response"
	DocPathParam  = docPrefix + "PathParam"
	DocQueryParam = docPrefix + "QueryParam"
	DocFormParam  = docPrefix + "FormParam"
	DocHeader     = docPrefix + "Header"
	DocCookie     = docPrefix + "Cookie"
)

func parseFileInfo(f *ast.File, fs *token.FileSet) *Info {
	info := &Info{
		PkgName: f.Name.Name,
		Imports: make([]Import, 0),
		Clients: make([]Client, 0),
	}

	cmts := parseSpecialComments(f, fs)
	for _, decl := range f.Decls {
		decl, ok := decl.(*ast.GenDecl)
		if !ok || len(decl.Specs) == 0 {
			log.Printf("not a decl (skip): %+v", decl)
			continue
		}

		isImp := false
		for _, spec := range decl.Specs {
			imp, ok := spec.(*ast.ImportSpec)
			if !ok {
				log.Printf("spec is not an import (expect type): %+v", spec)
				break
			}

			impInfo := Import{
				Path: imp.Path.Value,
			}
			if imp.Name != nil {
				impInfo.Name = imp.Name.Name
			}

			isImp = true
			info.Imports = append(info.Imports, impInfo)
		}
		if isImp {
			log.Println("imports are parsed (continue)")
			continue
		}
		if decl.Doc == nil {
			log.Printf("no doc for decl (skip): %+v", decl)
			continue
		}

		cDoc := decl.Doc.Text()
		if !strings.HasPrefix(cDoc, DocClient) {
			log.Printf("decl is not a client (skip): %+v", decl)
			continue
		}
		if len(decl.Specs) != 1 {
			log.Fatal("invalid number of specs")
		}

		spec := decl.Specs[0]
		ts, ok := spec.(*ast.TypeSpec)
		if !ok {
			log.Fatalf("not a type spec: %+v", spec)
		}
		itf, ok := ts.Type.(*ast.InterfaceType)
		if !ok {
			log.Fatalf("not an interface type: %+v", ts)
		}

		name := ts.Name.Name
		if len(name) == 0 || name[0] < 'a' || name[1] > 'z' {
			log.Fatalf("invalid client name: %s", name)
		}
		if itf.Methods == nil || len(itf.Methods.List) == 0 {
			log.Fatalf("client has no methods: %+v", itf)
		}

		c := Client{
			GoName: name,
			Calls:  make([]Call, 0),
		}

		for _, mt := range itf.Methods.List {
			if mt.Doc == nil {
				log.Printf("method has no doc (skip): %+v", mt)
				continue
			}

			callDoc := mt.Doc.Text()
			if !strings.HasPrefix(callDoc, DocCall+":") {
				continue
			}
			if len(mt.Names) != 1 {
				log.Fatal("too many call names")
			}

			call := Call{
				GoName: mt.Names[0].Name,
			}

			callInfo := callInfo{}
			if err := callInfo.Parse([]byte(callDoc[len(DocCall+":"):])); err != nil {
				log.Fatalf("invalid call info: %s", err.Error())
			}

			//TODO: make fmt
			call.Method = callInfo.Method
			call.PathFmt = callInfo.PathFmt

			call.Request.PathParams = make([]Param, 0)
			call.Request.QueryParams = make([]Param, 0)
			call.Request.FormParams = make([]Param, 0)
			call.Request.HeaderParams = make([]Param, 0)

			mtType := mt.Type.(*ast.FuncType)
			if mtType.Params == nil {
				//NOTE: no params is allowed
				continue
			}
			if len(mtType.Params.List) > 3 {
				log.Fatalf("invalid number of call func params: %+v", mtType.Params.List)
			}

			for _, mtParam := range mtType.Params.List {
				cmt := cmts[fs.Position(mtParam.Pos()).Line]
				if cmt == "" {
					log.Fatalf("no comment for call func param: %+v", mtParam)
				}

				if cmt == DocContext {
					ctxType, ok := mtParam.Type.(*ast.SelectorExpr)
					if !ok {
						log.Fatalf("invalid context type: %+v", mtParam.Type)
					}
					call.GoCtxType = fmt.Sprintf(
						"%s.%s", ctxType.X.(*ast.Ident).Name, ctxType.Sel.Name,
					)
				} else {
					structType, ok := mtParam.Type.(*ast.StructType)
					if !ok {
						log.Fatalf("struct type is expected: %+v", mtParam.Type)
					}

					//NOTE: empty struct is allowed
					if structType.Fields == nil {
						continue
					}

					for _, field := range structType.Fields.List {
						if field.Doc == nil {
							log.Fatalf("no doc comment for struct field: %+v", field)
						}

						//NOTE: basic auth requires two names
						if len(field.Names) < 1 || len(field.Names) > 2 {
							log.Fatalf("invalid field name: %+v", field.Names)
						}

						goName := field.Names[0].Name
						fieldDoc := field.Doc.Text()

						docType := ""
						if strings.HasPrefix(fieldDoc, DocPathParam) {
							docType = DocPathParam
						} else if strings.HasPrefix(fieldDoc, DocQueryParam) {
							docType = DocQueryParam
						} else if strings.HasPrefix(fieldDoc, DocFormParam) {
							docType = DocFormParam
						} else if cmt == DocRequest && strings.HasPrefix(fieldDoc, DocHeader) {
							docType = DocHeader
						}
						if docType != "" {
							paramInfo := paramInfo{}
							if err := paramInfo.UnmarshalJSON([]byte(fieldDoc[len(docType+":"):])); err != nil {
								log.Fatalf("failed to parse param info: %q: %s", fieldDoc, err)
							}

							pr := Param{
								Name:         paramInfo.Name,
								GoName:       goName,
								GoType:       field.Type.(*ast.Ident).Name,
								Optional:     paramInfo.Optional,
								DefaultValue: paramInfo.DefaultValue,
							}

							switch docType {
							case DocPathParam:
								call.Request.PathParams = append(call.Request.PathParams, pr)
							case DocQueryParam:
								call.Request.QueryParams = append(call.Request.QueryParams, pr)
							case DocFormParam:
								call.Request.FormParams = append(call.Request.FormParams, pr)
							case DocHeader:
								call.Request.HeaderParams = append(call.Request.HeaderParams, pr)
							}
						}
					}
				}
			}

			c.Calls = append(c.Calls, call)
		}

		info.Clients = append(info.Clients, c)
	}

	return info
}

func parseSpecialComments(f *ast.File, fs *token.FileSet) map[int]string {
	m := make(map[int]string)
	for _, comment := range f.Comments {
		txt := strings.TrimSpace(comment.Text())
		if !(txt == DocContext || txt == DocRequest || txt == DocResponse) {
			continue
		}

		p := fs.Position(comment.Pos())
		m[p.Line+1] = txt
	}
	return m
}
