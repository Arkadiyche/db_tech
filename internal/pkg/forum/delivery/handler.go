package delivery

import (
	"encoding/json"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/forum"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/models"
	"github.com/gorilla/mux"
	"net/http"
)

type ForumHandler struct {
	UseCase   forum.UseCase
}

func (fh *ForumHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//fmt.Println("forum Create")
	forum := models.Forum{}
	if err := json.NewDecoder(r.Body).Decode(&forum); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//fmt.Println(forum)
	f, err := fh.UseCase.Create(forum)
	//fmt.Println(err)
	if err != nil {
		switch err.Message {
		case models.NotExist.Error():
			res, err := json.Marshal(err)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(404)
			w.Write(res)
		case models.Exist.Error():
			res, err := json.Marshal(f)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(409)
			w.Write(res)
		default:
			http.Error(w, err.Message, http.StatusInternalServerError)
			return
		}
	} else {
		res, err := json.Marshal(f)
		if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
		}
		w.WriteHeader(201)
		w.Write(res)
	}
}

func (fh *ForumHandler) Details(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//fmt.Println("forum Details")
	vars := mux.Vars(r)
	slug := vars["slug"]
	f, err := fh.UseCase.Details(slug)
	if err != nil {
		res, err := json.Marshal(err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(404)
		w.Write(res)
	} else {
		res, err := json.Marshal(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(200)
		w.Write(res)
	}
}

func (fh *ForumHandler) Users(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//fmt.Println("forum Users")
	vars := mux.Vars(r)
	slug := vars["slug"]
	users, err := fh.UseCase.Users(slug, *r.URL)
	if err != nil {
		res, err := json.Marshal(err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(404)
		w.Write(res)
	} else {
		res, err1 := json.Marshal(users)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(200)
		w.Write(res)
	}
}