package main

import (
	//"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
    //mux.Handle("/assets", http.FileServer(http.Dir(".")))

	mux.HandleFunc("/healthz", serverHealth)
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
