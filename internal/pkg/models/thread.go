package models

import "time"

type Thread struct {
	Id int32 `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Forum string `json:"forum"`
	Message string `json:"message"`
	Votes int32 `json:"votes"`
	Slug string `json:"slug"`
	Created time.Time `json:"created"`
}

type Threads []Thread

type Vote struct {
	Nickname string `json:"nickname"`
	Voice int32 `json:"voice"`
}
