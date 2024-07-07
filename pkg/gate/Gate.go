package gate

import (
	"reflect"
	"samm/pkg/utils"
)

type Model struct {
	Name   string
	Policy interface{}
}

type Gate struct {
	Models []Model
}

func NewGate() *Gate {
	return &Gate{
		Models: make([]Model, 0),
	}
}

func (gate *Gate) Register(name interface{}, policy interface{}) {
	gate.Models = append(gate.Models, Model{Name: utils.GetStructName(name), Policy: policy})
}

func (gate *Gate) Authorize(name interface{}, method string, args ...interface{}) bool {
	results := utils.CallMethod(gate.Models[name], method, args)
	if len(results) > 0 && results[0].Kind() == reflect.Bool {
		return results[0].Bool()
	}
	return false
}
