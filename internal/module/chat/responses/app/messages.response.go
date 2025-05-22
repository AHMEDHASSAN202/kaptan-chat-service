package app

import "kaptan/pkg/database/mysql"

type MessagesResponse struct {
	Docs []*MessageResponse `json:"docs"`
	Meta *mysql.Pagination  `json:"meta"`
}
