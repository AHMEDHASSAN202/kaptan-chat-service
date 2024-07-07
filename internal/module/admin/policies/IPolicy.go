package policies

import "samm/internal/module/admin/domain"

type Model struct {
	Name   string
	Object interface{}
}

type IPolicy struct {
	Models []Model
}

func NewIPolicy() *IPolicy {
	return IPolicy{
		Models: []Model{{Name: "Admin"}},
	}
}
