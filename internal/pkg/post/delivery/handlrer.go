package delivery

import (
	"encoding/json"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/models"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/post"
	"github.com/gorilla/mux"
	"net/http"
)

type PostHandler struct {
	UseCase   post.UseCase
}

func (ph *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	posts := models.Posts{}
	vars := mux.Vars(r)
	if err := json.NewDecoder(r.Body).Decode(&posts); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	SlugOrId := vars["slug_or_id"]
	err := ph.UseCase.Create(SlugOrId, &posts)
	if err != nil {
		res, err1 := json.Marshal(err)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		switch err.Message {
		case models.NotExist.Error():
			w.WriteHeader(404)
			w.Write(res)
		case models.Exist.Error():
			w.WriteHeader(409)
			w.Write(res)
		default:
			w.WriteHeader(409)
			w.Write(res)
		}
	} else {
		res, err1 := json.Marshal(posts)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(201)
		w.Write(res)
	}
}

func (ph *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	post, err := ph.UseCase.Get(id, *r.URL)
	if err != nil {
		res, err1 := json.Marshal(err)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(409)
		w.Write(res)
	} else {
		res, err1 := json.Marshal(post)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(201)
		w.Write(res)
	}
}

func (ph *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	post := models.Post{}
	vars := mux.Vars(r)
	id := vars["id"]
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := ph.UseCase.Update(id, &post)
	if err != nil {
		res, err1 := json.Marshal(err)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(404)
		w.Write(res)
	} else {
		res, err1 := json.Marshal(&post)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(200)
		w.Write(res)
	}
}

func (ph *PostHandler) PostThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["slug_or_id"]
	post, err := ph.UseCase.GetThreadPosts(id, *r.URL)
	if err != nil {
		res, err1 := json.Marshal(err)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(404)
		w.Write(res)
	} else {
		res, err1 := json.Marshal(post)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(200)
		w.Write(res)
	}
}