package apiserver

import (
	"fmt"
	forumHandler "github.com/Arkadiyche/bd_techpark/internal/pkg/forum/delivery"
	forumRep "github.com/Arkadiyche/bd_techpark/internal/pkg/forum/repository"
	forumUC "github.com/Arkadiyche/bd_techpark/internal/pkg/forum/usecase"
	postHandler	"github.com/Arkadiyche/bd_techpark/internal/pkg/post/delivery"
	postRep	"github.com/Arkadiyche/bd_techpark/internal/pkg/post/repository"
	postUC "github.com/Arkadiyche/bd_techpark/internal/pkg/post/usecase"
	serviceHandler "github.com/Arkadiyche/bd_techpark/internal/pkg/service"
	threadHandler "github.com/Arkadiyche/bd_techpark/internal/pkg/thread/delivery"
	threadRep "github.com/Arkadiyche/bd_techpark/internal/pkg/thread/repository"
	threadUC "github.com/Arkadiyche/bd_techpark/internal/pkg/thread/usecase"
	userHandler "github.com/Arkadiyche/bd_techpark/internal/pkg/user/delivery"
	userRep "github.com/Arkadiyche/bd_techpark/internal/pkg/user/repository"
	userUC "github.com/Arkadiyche/bd_techpark/internal/pkg/user/usecase"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"log"
	"net/http"
)

type APIServer struct {
	config    *Config
	router    *mux.Router
	store     *pgx.ConnPool
}

func New(config *Config) *APIServer {
	return &APIServer{
		config:    config,
		router:    mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	config := pgx.ConnPoolConfig{
		ConnConfig:     pgx.ConnConfig{
			Host:                 "localhost",
			Port:                 5432,
			Database:             "docker",
			User:                 "docker",
			Password:             "docker",
			TLSConfig:            nil,
			UseFallbackTLS:       false,
			FallbackTLSConfig:    nil,
			Logger:               nil,
			LogLevel:             0,
			Dial:                 nil,
			RuntimeParams:        nil,
			OnNotice:             nil,
			CustomConnInfo:       nil,
			CustomCancel:         nil,
			PreferSimpleProtocol: false,
			TargetSessionAttrs:   "",
		},
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}
	connPool, err := pgx.NewConnPool(config)
	if err != nil {
		log.Fatal(err)
	}
	s.store = connPool
    defer s.store.Close()

	s.configureRouter()

	fmt.Println("starting api server")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}


func (s *APIServer) InitHandler() (user userHandler.UserHandler, forum forumHandler.ForumHandler, thread threadHandler.ThreadHandler, post postHandler.PostHandler, service serviceHandler.ServiceHandler){
	UserRep := userRep.NewUserRepository(s.store)
	UserUC := userUC.NewUserUseCase(UserRep)
	UserHandler := userHandler.UserHandler{
		UseCase: UserUC,
	}

	ForumRep := forumRep.NewForumRepository(s.store)
	ForumUC := forumUC.NewForumUseCase(ForumRep, UserRep)
	ForumHandler := forumHandler.ForumHandler{
		UseCase: ForumUC,
	}

	ThreadRep := threadRep.NewThreadRepository(s.store)
	ThreadUC := threadUC.NewThreadUseCase(ThreadRep, ForumRep)
	ThreadHandler := threadHandler.ThreadHandler{
		UseCase: ThreadUC,
	}

	PostRep := postRep.NewPostRepository(s.store)
	PostUC := postUC.NewPostUseCase(PostRep, ThreadRep, ForumRep, UserRep)
	PostHandler := postHandler.PostHandler{
		UseCase: PostUC,
	}

	ServiceHandler := serviceHandler.NewServiceHandler(s.store)

	return UserHandler, ForumHandler, ThreadHandler, PostHandler, ServiceHandler
}

func (s *APIServer) configureRouter() {
	user, forum, thread, post, service := s.InitHandler()
	//User routes ...
	s.router.HandleFunc("/api/user/{nickname}/create", user.Create)
	s.router.HandleFunc("/api/user/{nickname}/profile", user.GetProfile).Methods("GET")
	s.router.HandleFunc("/api/user/{nickname}/profile", user.UpdateProfile).Methods("POST")
	//Forum routes ...
	s.router.HandleFunc("/api/forum/create", forum.Create)
	s.router.HandleFunc("/api/forum/{slug}/details", forum.Details)
	s.router.HandleFunc("/api/forum/{slug}/users", forum.Users)
	//Thread routes ...
	s.router.HandleFunc("/api/forum/{slug}/create", thread.Create)
	s.router.HandleFunc("/api/forum/{slug}/threads", thread.ForumThreads)
	s.router.HandleFunc("/api/thread/{slug_or_id}/vote", thread.Vote)
	s.router.HandleFunc("/api/thread/{slug_or_id}/details", thread.GetThread).Methods("GET")
	s.router.HandleFunc("/api/thread/{slug_or_id}/details", thread.UpdateThread).Methods("POST")
	//Post routes ...
	s.router.HandleFunc("/api/thread/{slug_or_id}/create", post.Create)
	s.router.HandleFunc("/api/post/{id}/details", post.Get).Methods("GET")
	s.router.HandleFunc("/api/post/{id}/details", post.Update).Methods("POST")
	s.router.HandleFunc("/api/thread/{slug_or_id}/posts", post.PostThread)
	//Service routes ...
	s.router.HandleFunc("/api/service/status", service.Status)
	s.router.HandleFunc("/api/service/clear", service.Clear)
}
