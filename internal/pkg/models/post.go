package models

type Post struct {
	Author string `json:"author"`
	Created string `json:"created"`
	Forum string `json:"forum"`
	Id string `json:"id"`
	IsEdited bool `json:"isEdited"`
	Message string `json:"message"`
	Parent int64 `json:"parent"`
	Thread int32 `json:"thread"`
}

type PostFull struct {
	Author User `json:"author"`
	Forum Forum `json:"forum"`
	Post Post `json:"post"`
	Thread Thread `json:"thread"`
}