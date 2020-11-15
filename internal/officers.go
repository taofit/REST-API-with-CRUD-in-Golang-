package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Officer struct {
	ID   int    `json:"id"`
	NAME string `json:"name"`
}

func respondWithJSON(w http.ResponseWriter, code int, message interface{}) {
	response, _ := json.Marshal(message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func GetOfficers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := DbConn()
	selResult, err := db.Query("SELECT * FROM officers ORDER BY id DESC")
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	officers := []Officer{}
	officer := Officer{}
	for selResult.Next() {
		var id int
		var name string
		err = selResult.Scan(&id, &name)
		if err != nil {
			panic(err.Error())
		}
		officer.ID = id
		officer.NAME = name
		officers = append(officers, officer)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(&officers)
}

func GetOfficer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	db := DbConn()
	selResult := db.QueryRow("SELECT * FROM officers WHERE id=?", id)
	officer := Officer{}

	var name string
	err = selResult.Scan(&id, &name)

	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	defer db.Close()

	officer.ID = id
	officer.NAME = name

	json.NewEncoder(w).Encode(&officer)
}

func UpdateOfficer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid officer ID: " + params["id"]})
		return
	}
	var officer Officer
	_ = json.NewDecoder(r.Body).Decode(&officer)

	if officer.NAME == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request officer"})
		return
	}

	officer.ID = id

	db := DbConn()
	updateResult, err := db.Prepare("UPDATE officers SET name=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	updateResult.Exec(officer.NAME, officer.ID)
	defer db.Close()

	message := "UPDATE officer: " + officer.NAME
	log.Println(message)
	respondWithJSON(w, http.StatusAccepted, map[string]string{"success": message, "name": officer.NAME, "id": strconv.Itoa(officer.ID)})
}

func CreateOfficer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var officer Officer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&officer); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request officer"})
		return
	}

	db := DbConn()
	insert, err := db.Prepare("INSERT INTO officers(name) VALUES(?)")
	if err != nil {
		panic(err.Error())
	}
	_, err = insert.Exec(officer.NAME)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	message := "Officer '" + officer.NAME + "' is created"
	log.Println(message)
	respondWithJSON(w, http.StatusAccepted, map[string]string{"success": message})
}

func DeleteOfficer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Invalue officer ID")
		return
	}

	db := DbConn()
	delResult, err := db.Prepare("DELETE FROM officers WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delResult.Exec(id)
	defer db.Close()

	message := "officer with ID: " + params["id"] + " is deleted."
	log.Println(message)
	respondWithJSON(w, http.StatusAccepted, message)
}
