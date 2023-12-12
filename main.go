package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

type apiConfig struct {
	fileserverHits int
}


func main() {
	
	apiConfig := apiConfig{
		fileserverHits: 0,
	}
	
	r := chi.NewRouter()
	r.Use(middlewareLog)
	r.Use(middlewareCors)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	apiRouter := chi.NewRouter()
	r.Mount("/api", apiRouter)

	// appRouter := chi.NewRouter()
	// r.Mount("/app" , appRouter)

	fsHandler := apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	r.Handle("/app/*", fsHandler)
	r.Handle("/app" , fsHandler)
	apiRouter.Get("/healthz", serverHealth)
	apiRouter.Get("/metrics", apiConfig.totalHits)
	apiRouter.Get("/reset", apiConfig.resetHits)
    apiRouter.Post("/validate_chirp", validateChirp)

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

func validateChirp(w http.ResponseWriter, r *http.Request){
	type parameters struct {
        // these tags indicate how the keys in the JSON should be mapped to the struct fields
        // the struct fields must be exported (start with a capital letter) if you want them parsed
        Body string `json:"body"`
    }

	type errMessage struct {
		Error string `json:"error"`
	}

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        // an error will be thrown if the JSON is invalid or has the wrong types
        // any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		errMessage := errMessage{
			Error: "Something went wrong",
		}
		dat, err := json.Marshal(errMessage)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			w.Write([]byte("Something went wrong"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(dat)
		
    }
	cleanedResponse := getCleanedResponse(params.Body)

	var isValid bool
    log.Printf("Body length : %v", len(params.Body))
	if len(params.Body) <= 140 {
		isValid = true
	} else {
		isValid = false
	}

	if isValid {
		type returnVals struct {
			// the key will be the name of struct field unless you give it an explicit JSON tag
			CreatedAt time.Time `json:"created_at"`
			CleanedBody string `json:"cleaned_body"`
		}
		respBody := returnVals{
			CreatedAt: time.Now(),
			CleanedBody: cleanedResponse,
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			w.Write([]byte("Something went wrong"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(dat)
	
	} else {
		errMessage := errMessage{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(errMessage)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			w.Write([]byte("Something went wrong"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
	}
}

func getCleanedResponse(msg string) string {
	//lowerMsg := strings.ToLower(msg)
     words := strings.Split(msg, " ")
	 for i, val := range words {
		lowerVal := strings.ToLower(val)
		if lowerVal == "profane" || lowerVal == "sharbert" || lowerVal == "fornax" || lowerVal == "kerfuffle" {
            words[i] = "****"
		}
	 }
	 cleanedMsg := strings.Join(words, " ")
	 return cleanedMsg
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


