package gate

import (
	"reflect"
	"samm/pkg/utils"
)

type Gate struct {
	Models map[string]interface{}
}

func NewGate() *Gate {
	return &Gate{
		Models: make(map[string]interface{}),
	}
}

func (gate *Gate) Register(name interface{}, policy interface{}) {
	gate.Models[utils.GetStructName(name)] = policy
}

func (gate *Gate) Authorize(obj interface{}, method string, args ...interface{}) bool {
	results := utils.CallMethod(gate.Models[utils.GetStructName(obj)], method, append([]interface{}{obj}, args...)...)
	if len(results) > 0 && results[0].Kind() == reflect.Bool {
		return results[0].Bool()
	}
	return false
}
