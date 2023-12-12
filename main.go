package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}


func main() {
	mux := http.NewServeMux()
	logmux := middlewareLog(mux)
	corsMux := middlewareCors(logmux)
	apiConfig := apiConfig{
		fileserverHits: 0,
	}
	

	server := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}

	mux.Handle("/app/", http.StripPrefix("/app", apiConfig.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
    //mux.Handle("/assets", http.FileServer(http.Dir(".")))

	mux.HandleFunc("/healthz", serverHealth)
	mux.HandleFunc("/metrics", apiConfig.totalHits)
	mux.HandleFunc("/reset", apiConfig.resetHits)
    log.Printf("Serving files from webserver on port: %v\n", 8080)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
	//fmt.Println("server started on port 8080")
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func serverHealth(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) totalHits(w http.ResponseWriter, r *http.Request){
	x := cfg.fileserverHits
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	responseString := fmt.Sprintf("Hits: %d", x)
	w.Write([]byte(responseString))
}

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request){
	cfg.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	//responseString := fmt.Sprint("Reset Successful")
	w.Write([]byte("Reset Successful"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits = cfg.fileserverHits + 1
	    log.Printf("Hits: %v", cfg.fileserverHits)
		next.ServeHTTP(w, r)
	})
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}


