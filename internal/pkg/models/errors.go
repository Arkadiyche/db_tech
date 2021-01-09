package models

import "errors"

type Error struct {
	Message string `json:"message"`
}

var Exist error = errors.New("exist")
var NotExist error = errors.New("doesn't exist")
