package service

import (
	"encoding/json"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/models"
	"github.com/jackc/pgx"
	"net/http"
)

type ServiceHandler struct {
	db *pgx.ConnPool
}

func NewServiceHandler(db *pgx.ConnPool) ServiceHandler {
	return ServiceHandler{
		db: db,
	}
}

func (sh *ServiceHandler) Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := models.Status{}
	err := sh.db.QueryRow("SELECT " +
		"(SELECT COUNT(*) FROM users), (SELECT COUNT(*) FROM forums), (SELECT COUNT(*) FROM threads), (SELECT COUNT(*) FROM posts)").Scan(
		&status.User,
		&status.Forum,
		&status.Thread,
		&status.Post)
	if err != nil {
		res, err := json.Marshal(status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(200)
		w.Write(res)
	}
	res, err := json.Marshal(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(200)
	w.Write(res)
}

func (sh *ServiceHandler) Clear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sh.db.Exec("TRUNCATE  votes, posts, forum_users, threads, forums, users CASCADE")
	w.WriteHeader(200)
}
