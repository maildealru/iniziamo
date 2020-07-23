package parser

//go:generate easyjson

import (
	"net/http"

	"github.com/pkg/errors"
)

//easyjson:json
type callInfo struct {
	Method  string `json:"method"`
	PathFmt string `json:"path"`
}

func (v *callInfo) Parse(data []byte) error {
	if err := (v).UnmarshalJSON(data); err != nil {
		return errors.Wrapf(err, "json call: unmarshalling error [%q]", data)
	}

	//NOTE: maybe default path = "/" is not a bad idea...
	if v.PathFmt == "" {
		return errors.Errorf("json call: path is empty [%q]", data)
	}

	//TODO: validate path
	if v.Method == "" {
		v.Method = http.MethodGet
	}

	knownMethods := map[string]struct{}{
		http.MethodHead: {}, http.MethodGet: {}, http.MethodPost: {},
		http.MethodPut: {}, http.MethodPatch: {}, http.MethodDelete: {},
	}
	if _, ok := knownMethods[v.Method]; !ok {
		return errors.Errorf("json call: method is unknown [%q]", data)
	}

	return nil
}

//easyjson:json
type paramInfo struct {
	Name          string      `json:"name"`
	Optional      bool        `json:"optional"`
	DefaultValue  interface{} `json:"default"`
	ConstValue    interface{} `json:"const"`
	ConfigVarName string      `json:"config"`
}
