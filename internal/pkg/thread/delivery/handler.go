package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/models"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/thread"
	"github.com/gorilla/mux"
	"net/http"
)

type ThreadHandler struct {
	UseCase   thread.UseCase
}

func (th *ThreadHandler) Create(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	fmt.Println("thread Create")
	thread := models.Thread{}
	vars := mux.Vars(r)
	if err := json.NewDecoder(r.Body).Decode(&thread); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	thread.Forum = vars["slug"]
	t, err := th.UseCase.Create(&thread)
	fmt.Println(err)
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
			res, err := json.Marshal(t)
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
		res, err := json.Marshal(thread)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(201)
		w.Write(res)
	}
}

func (th *ThreadHandler) ForumThreads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	slug := vars["slug"]
	ts, err := th.UseCase.ForumThreads(slug, *r.URL)
	if err != nil {
		res, err := json.Marshal(err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(404)
		w.Write(res)
	} else {
		res, err1 := json.Marshal(ts)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(200)
		w.Write(res)
	}
}

func (th *ThreadHandler) GetThread(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	SlugOrId := vars["slug_or_id"]
	t, err := th.UseCase.GetThread(SlugOrId)
	if err != nil {
		res, err := json.Marshal(err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(404)
		w.Write(res)
	} else {
		res, err1 := json.Marshal(t)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(200)
		w.Write(res)
	}
}

func (th *ThreadHandler) UpdateThread(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	thread := models.Thread{}
	vars := mux.Vars(r)
	SlugOrId := vars["slug_or_id"]
	if err := json.NewDecoder(r.Body).Decode(&thread); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t, err := th.UseCase.UpdateThread(SlugOrId, thread)
	if err != nil {
		res, err := json.Marshal(err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(404)
		w.Write(res)
	} else {
		res, err1 := json.Marshal(t)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(200)
		w.Write(res)
	}
}

func (th *ThreadHandler) Vote(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	vote := models.Vote{}
	vars := mux.Vars(r)
	if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	SlugOrId := vars["slug_or_id"]
	t, err := th.UseCase.Vote(SlugOrId, vote)
	if err != nil {
		res, err := json.Marshal(err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(404)
		w.Write(res)
	} else {
		res, err1 := json.Marshal(t)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(200)
		w.Write(res)
	}
}