package apiserver

import (
	"fmt"
	"github.com/gorilla/mux"
	"githun.com/Arkadiyche/bd_techpark/internal/pkg/store"
	"net/http"
)

type APIServer struct {
	config    *Config
	router    *mux.Router
	store     *store.Store
}

func New(config *Config) *APIServer {
	return &APIServer{
		config:    config,
		router:    mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {

	if err := s.configureStore(); err != nil {
		return err
	}

	//s.configureRouter()

	fmt.Println("starting api server")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st

	return nil
}