package main

import (
	"casaaltitudine/config"
	"casaaltitudine/controllers"
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
	})
	dbsqlx := config.ConnectDBSqlx()
	hsqlx := controllers.NewBaseHandlerSqlx(dbsqlx)

	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// documentation for developers
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	r.Handle("/docs", sh)

	// documentation for share
	opts1 := middleware.RedocOpts{SpecURL: "/swagger.yaml", Path: "docs"}
	sh1 := middleware.Redoc(opts1, nil)
	r.Handle("/docs", sh1)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"hello\": \"world\"}"))
	})

	altitudine := r.PathPrefix("/").Subrouter()
	altitudine.HandleFunc("/altitude", hsqlx.PostAltitudeSqlx).Methods("POST")
	altitudine.HandleFunc("/altitudes", hsqlx.GetAltitudesSqlx).Methods("GET")
	altitudine.HandleFunc("/lasthour", hsqlx.GetLastHourSqlx).Methods("GET")
	altitudine.HandleFunc("/showdata/{recordNumber}", hsqlx.GetShowDataSqlx).Methods("GET")

	http.Handle("/", r)
	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", "", "5558"),
		Handler: cors.Default().Handler(r),
	}
	s.ListenAndServe()
}
