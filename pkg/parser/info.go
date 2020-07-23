package parser

type Info struct {
	Dir     string
	Name    string
	PkgName string
	Imports []Import
	Clients []Client
}

type Import struct {
	Name string
	Path string
}

type Client struct {
	GoName           string
	GoIntName        string
	GoImplName       string
	GoImplConstrName string
	Config           ClientConfig
	Calls            []Call
	ConfigVars       []ConfigVar
}

type ClientConfig struct {
	GoName string
}

type Call struct {
	GoName    string
	Method    string
	PathFmt   string
	GoCtxType string
	Request   Request
	Response  Response
}

type Request struct {
	GoName       string
	GoConstrName string
	PathParams   []Param
	QueryParams  []Param
	FormParams   []Param
	HeaderParams []Param
}

type Param struct {
	Name          string
	GoName        string
	GoType        string
	GoFieldName   string
	GoSetFlagName string
	GoSetterName  string
	Optional      bool
	DefaultValue  interface{}
	ConstValue    interface{}
	ConfigVarName string
}

type Response struct {
	GoName       string
	GoConstrName string
}

type ConfigVar struct {
	GoName string
	GoType string
}
