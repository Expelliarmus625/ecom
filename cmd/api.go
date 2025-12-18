package main

import (
	"log"
	"net/http"
	"time"

	repo "github.com/expelliarmus625/ecom/internal/adapters/postgresql/sqlc"
	"github.com/expelliarmus625/ecom/internal/products"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

// Instance of our api
type application struct {
	config config
	// logger
	db *pgx.Conn
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	//Middleware
	r.Use(middleware.RequestID) //used for rate limiting
	r.Use(middleware.RealIP) //also important for rate-limiting and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second)) //timeout after 60 seconds. TODO: Set timeout from config
	
	r.Get("/health", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("I'm alive"))
	})

	productService := products.NewService(repo.New(app.db))
	productHandler := products.NewHandler(productService)
	r.Get("/products", productHandler.ListProducts)
	r.Get("/products/{id}", productHandler.FindProductByID)

	// http.ListenAndServe(app.config.addr, r)
	return r
}


//run
func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr: app.config.addr,
		Handler: h,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 10,
		IdleTimeout: time.Minute,
	}
	log.Printf("Server has started at addr %s", app.config.addr)
	return srv.ListenAndServe()
}

// Configuration object
type config struct {
	addr string
	db dbConfig
}


type dbConfig struct {
	dsn string
}
