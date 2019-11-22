package server

import (
	"awise-socialNetwork/config"
	"awise-socialNetwork/server/middleware"
	v1 "awise-socialNetwork/server/v1"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Start for start the http server
func Start() {
	config, _ := config.GetConfig()
	r := mux.NewRouter()

	// Cors auth
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	headersOk := handlers.AllowedHeaders([]string{"authorization", "content-type"})

	// middleware
	r.Use(middleware.BasicHeader, middleware.Logger)

	// create sub router
	public := r.PathPrefix("/").Subrouter()
	private := r.PathPrefix("/").Subrouter()
	private.Use(middleware.IsGoodToken)

	// Login
	public.HandleFunc("/api/v1/login", v1.Login).Methods("POST")
	public.HandleFunc("/api/v1/login/refresh", v1.RefreshLogin).Methods("POST")

	private.HandleFunc("/api/v1/account", v1.CreateAccount).Methods("POST")
	private.HandleFunc("/api/v1/account/{id}", v1.GetAccount).Methods("GET")
	private.HandleFunc("/api/v1/account/avatar", v1.AddAccountAvatar).Methods("POST")

	private.HandleFunc("/api/v1/upload/picture", v1.UploadPicture).Methods("POST")
	private.HandleFunc("/api/v1/upload/avatar", v1.UploadAvatar).Methods("POST")

	// Ajax
	public.HandleFunc("/", nil).Methods("OPTIONS")

	log.Println("Start http server on localhost:" + strconv.Itoa(config.Port))
	srv := &http.Server{
		Handler:      handlers.CORS(originsOk, headersOk, methodsOk)(r),
		Addr:         "127.0.0.1:" + strconv.Itoa(config.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	srv.ListenAndServe()
}
