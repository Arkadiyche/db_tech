package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/models"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/user"
	"github.com/gorilla/mux"
	"net/http"
)

type UserHandler struct {
	UseCase   user.UseCase
}

func (uh *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("user Create")
	user := models.User{}
	vars := mux.Vars(r)
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.Nickname = vars["nickname"]
	users, err := uh.UseCase.Create(user)
	fmt.Println(err)
	switch err {
	case nil:
		fmt.Println(users)
		res, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(201)
		w.Write(res)
	case models.Exist:
		res, err := json.Marshal(users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(409)
		w.Write(res)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (uh *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("user GetProfile")
	vars := mux.Vars(r)
	nickname := vars["nickname"]
	user, err := uh.UseCase.GetProfile(nickname)
	fmt.Println(err)
	if err == nil {
		res, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(200)
		w.Write(res)
	} else {
		res, err := json.Marshal(err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(404)
		w.Write(res)
	}
}


func (uh *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("user UpdateProfile")
	user := models.User{}
	vars := mux.Vars(r)
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(user)
	user.Nickname = vars["nickname"]
	err := uh.UseCase.UpdateProfile(&user)
	fmt.Println(err)
	if err != nil {
		res, err1 := json.Marshal(err)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			return
		}
		if err.Message == models.NotExist.Error() {
			w.WriteHeader(404)
			w.Write(res)
		} else {
			w.WriteHeader(409)
			w.Write(res)
		}
	} else {
		res, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(200)
		w.Write(res)
	}
}