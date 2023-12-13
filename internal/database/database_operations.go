package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error){
	db := DB{
		path: path,
	}
	return &db , nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error){

	err := db.ensureDB()
	if err != nil {
		return Chirp{}, err
	}
	chirpData, err := db.loadDB()
	if err != nil {
		return Chirp{} , err
	}
	currentLength := len(chirpData.Chirps)
	newChirp := Chirp{
		Body: body,
		ID: currentLength + 1,
	}
	chirpData.Chirps[currentLength + 1] = newChirp
	err = db.writeDB(chirpData)
	if err != nil {
		return Chirp{} , err
	}
	return newChirp , nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error){

	chirpData, err := db.loadDB()
	if err != nil {
		return []Chirp{} , err
	}
	
	chirps := make([]Chirp , 0)
	for _ , val := range chirpData.Chirps {
		chirps = append(chirps, val)
	}
	return chirps , nil

}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error{
	db.mux.RLock()
	defer db.mux.RUnlock()
	file1, err := os.Open(db.path)
	if err != nil {
		if os.IsNotExist(err) {
			file , err1 := os.Create(db.path)
			if err1 != nil {
				return err1
			}
			defer file.Close()
			return nil
		} else {
			return err
		}
	}
	defer file1.Close()
    return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error){
    db.mux.RLock()
	defer db.mux.RUnlock()
	file, err := os.Open(db.path)

	if err != nil {
		return DBStructure{}, err
	}
	defer file.Close()
	var chirpData DBStructure
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&chirpData)
	if err != nil {
		return DBStructure{}, fmt.Errorf("error decoding JSON: %v", err)
	}

	return chirpData , nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.RLock()
	defer db.mux.RUnlock()
	dat , err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, dat, 0644)
	if err != nil {
		log.Printf("Error writing to file: %v\n", err)
		return err
	}
	return nil
}





