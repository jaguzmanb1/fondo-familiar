package main

import (
	"context"
	"database/sql"
	"fmt"
	"fondo-mod/auth"
	"fondo-mod/data"
	"fondo-mod/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gohandlers "github.com/gorilla/handlers"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
)

var _ = godotenv.Load(".env")
var (
	//ConnectionString cadena de conexión a la base de datos
	ConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("user"),
		os.Getenv("pass"),
		os.Getenv("host"),
		os.Getenv("port"),
		os.Getenv("db_name"))
)

func main() {
	l := hclog.New(&hclog.LoggerOptions{
		Name:  "fondo-app",
		Level: hclog.LevelFromString("DEBUG"),
	})
	handlerLogger := l.Named("Handler")
	serviceLogger := l.Named("Service")
	authLogger := l.Named("Auth")

	db, err := sql.Open("mysql", ConnectionString)
	if err != nil {
		l.Error("Can't connect to mysql database", "error", err)
	}

	// JSON validator
	v := data.NewValidation()

	// Token validator handler
	auth := auth.New(authLogger)

	// New user service
	us := data.NewUserService(db, serviceLogger)

	// New user handler
	uha := handlers.New(us, handlerLogger, v)

	// Router creationg
	sm := mux.NewRouter()

	sm.Use(auth.MiddlewareTokenValidationRol3)

	postAportesR1 := sm.Methods(http.MethodPost).Subrouter()
	postAportesR1.Use(uha.MiddlewareValidateAporte)
	postAportesR1.Use(auth.MiddlewareTokenValidationRol1)
	postAportesR1.HandleFunc("/aportes", uha.CreateAporte)

	getAllR3ID := sm.Methods(http.MethodGet).Subrouter()
	getAllR3ID.Use(uha.MiddlewareCheckUserIDCall)
	getAllR3ID.HandleFunc("/usuarios/{id:[0-9]+}/aportes", uha.GetAllAportesByID)
	getAllR3ID.HandleFunc("/usuarios/{id:[0-9]+}/aportes/sum", uha.GetSumAportesByID)
	getAllR3ID.HandleFunc("/usuarios/{id:[0-9]+}/creditos", uha.GetAllCreditosByUserID)

	getAllR1 := sm.Methods(http.MethodGet).Subrouter()
	getAllR1.Use(auth.MiddlewareTokenValidationRol1)
	getAllR1.HandleFunc("/aportes", uha.GetAllAportes)
	getAllR1.HandleFunc("/creditos", uha.GetAllCreditos)

	postCreditosR1 := sm.Methods(http.MethodPost).Subrouter()
	postCreditosR1.Use(uha.MiddlewareValidateCredito)
	postCreditosR1.Use(auth.MiddlewareTokenValidationRol1)
	postCreditosR1.HandleFunc("/creditos", uha.CreateCredito)

	postCreditosPagosR1 := sm.Methods(http.MethodPost).Subrouter()
	postCreditosPagosR1.Use(uha.MiddlewareValidatePago)
	postCreditosPagosR1.HandleFunc("/pago", uha.CreatePago)

	getCreditosR3 := sm.Methods(http.MethodGet).Subrouter()
	getCreditosR3.Use(uha.MiddlewareValidateCredito)
	getCreditosR3.HandleFunc("/creditos/proyeccion", uha.GetProyeccionCredito)

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	s := http.Server{
		Addr:         os.Getenv("bindAddress"),                         // configure the bind address
		Handler:      ch(sm),                                           // set the default handler
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}), // set the logger for the server
		ReadTimeout:  5 * time.Second,                                  // max time to read request from the client
		WriteTimeout: 10 * time.Second,                                 // max time to write response to the client
		IdleTimeout:  120 * time.Second,                                // max time for connections using TCP Keep-Alive
	}

	go func() {
		l.Debug("[main] Starting server on", "port", os.Getenv("bindAddress"))

		err := s.ListenAndServeTLS("cert/server.crt", "cert/server.key")
		if err != nil {
			l.Error("[main] Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
