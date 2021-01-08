package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sachinsu/chaosengg/internal/app"
	"github.com/sachinsu/chaosengg/internal/setup"
)

// Server holding config details
type Server struct {
	DbConnectionString string
	Port               string
	BaseURL            string
	Logger             *log.Logger
	Router             *httprouter.Router
}

// NewServer creates new instance
func NewServer(dbConn string, port string, baseurl string) (*Server, error) {

	s := Server{Logger: log.New(os.Stdout, "Logger: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile),
		DbConnectionString: dbConn,
		Port:               port,
		BaseURL:            baseurl,
		Router:             httprouter.New(),
	}

	s.Router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", r.Header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
		}
		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})

	return &s, nil

}

func (s *Server) SetupRoutes() {
	s.Router.GET("/", Index)

	s.Router.GET(s.BaseURL+"/cities", s.ListCities)
	s.Router.GET(s.BaseURL+"/cities/latency", s.ListCitiesWithLatency)
}

// Run the http server
func (s *Server) Run() error {
	//ref: https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/

	s.SetupRoutes()

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         s.Port,
		Handler:      s.Router,
	}

	s.Logger.Printf("Starting HTTP Server, %s with base %s\n", s.Port, s.BaseURL)
	return srv.ListenAndServe()
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// ListCities route
func (s *Server) ListCities(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//ref: https://github.com/julienschmidt/httprouter/issues/9s
	queryValues := r.URL.Query()
	start := 0
	limit := 5000

	if len(queryValues.Get("startid")) > 0 {
		start, _ = strconv.Atoi(queryValues.Get("startid"))
	}

	if len(queryValues.Get("limit")) > 0 {
		limit, _ = strconv.Atoi(queryValues.Get("limit"))
	}

	citilist, err := app.GetCities(r.Context(), s.DbConnectionString, start, limit)

	if err != nil {
		s.Logger.Fatalf("Error while getting cities %v", err)
		RespondErr(w, r, http.StatusInternalServerError)
	} else {
		RespondJSON(w, r, http.StatusOK, citilist)
	}

}

// ListCitiesWithLatency route
func (s *Server) ListCitiesWithLatency(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	queryValues := r.URL.Query()
	delay := 100

	if len(queryValues.Get("delay")) > 0 {
		delay, _ = strconv.Atoi(queryValues.Get("delay"))
	}

	proxy := setup.InjectLatency("pgsql", delay)
	defer proxy.RemoveToxic("latency_downstream")
	s.ListCities(w, r, p)
}

// RespondJson refer Use https://github.com/unrolled/render for rendering JSON/HTML/Text Content or from what Ryer recommends
func RespondJSON(w http.ResponseWriter, r *http.Request,
	status int, data interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if data != nil {
		EncodeBody(w, r, data)
	}
}

// EncodeBody for JSON generation
func EncodeBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

// RespondErr sends error
func RespondErr(w http.ResponseWriter, r *http.Request,
	status int, args ...interface{}) {
	RespondJSON(w, r, status, map[string]interface{}{
		"error": map[string]interface{}{"message": fmt.Sprint(args...)},
	})
}
