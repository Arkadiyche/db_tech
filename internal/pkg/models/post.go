package models

import "time"

type Post struct {
	Author string `json:"author"`
	Created time.Time `json:"created"`
	Forum string `json:"forum"`
	Id int64 `json:"id"`
	IsEdited bool `json:"isEdited"`
	Message string `json:"message"`
	Parent int64 `json:"parent"`
	Thread int32 `json:"thread"`
}

type PostFull struct {
	Author *User `json:"author,omitempty"`
	Forum *Forum `json:"forum,omitempty"`
	Post *Post `json:"post"`
	Thread *Thread `json:"thread,omitempty"`
}

type Posts []Post