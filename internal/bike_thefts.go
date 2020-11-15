package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type TheftCase struct {
	ID          int       `json:"id"`
	TITLE       string    `json:"title"`
	BRAND       string    `json:"brand"`
	CITY        string    `json:"city"`
	DESCRIPTION string    `json:"description"`
	REPORTED    time.Time `json:"reported"`
	UPDATED     time.Time `json:"updated"`
	SOLVED      bool      `json:"solved"`
	OFFICER     Officer
	IMAGE       []byte
}

func CreateCaseNoImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var theftCase TheftCase

	_ = json.NewDecoder(r.Body).Decode(&theftCase)

	if theftCase.TITLE == "" || theftCase.BRAND == "" || theftCase.CITY == "" || theftCase.DESCRIPTION == "" {
		respondWithJSON(w, http.StatusBadRequest, "Some fields are missing please enter them again")
		return
	}

	db := DbConn()
	insert, err := db.Prepare("INSERT INTO bike_thefts(title, brand, city, description) VALUES(?,?,?,?)")
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = insert.Exec(theftCase.TITLE, theftCase.BRAND, theftCase.CITY, theftCase.DESCRIPTION)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	defer db.Close()
	// uploadFile(w, r)

	message := "Bike theft: '" + theftCase.TITLE + "' is created"
	log.Println(message)
	respondWithJSON(w, http.StatusAccepted, message)
}

func GetCases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := DbConn()
	selResult, err := db.Query(`SELECT bt.id, bt.title, bt.brand, bt.city, bt.description, bt.reported, bt.updated, bt.solved, IFNULL(o.id, 0), IFNULL(o.name, '')
								FROM bike_thefts bt
								LEFT JOIN officers o
								ON o.id = bt.officer
								ORDER BY bt.id DESC`)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	theftCases := []TheftCase{}
	theftCase := TheftCase{}
	for selResult.Next() {
		var officerId int
		var officerName string
		err = selResult.Scan(&theftCase.ID, &theftCase.TITLE, &theftCase.BRAND,
			&theftCase.CITY, &theftCase.DESCRIPTION, &theftCase.REPORTED, &theftCase.UPDATED, &theftCase.SOLVED, &officerId, &officerName)
		if err != nil {
			panic(err.Error())
		}
		theftCase.OFFICER.ID = officerId
		theftCase.OFFICER.NAME = officerName
		theftCases = append(theftCases, theftCase)
	}

	defer db.Close()
	json.NewEncoder(w).Encode(&theftCases)
}

func GetCase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	db := DbConn()
	var officerId int
	var officerName string
	selResult := db.QueryRow(`SELECT bt.id, bt.title, bt.brand, bt.city, bt.description, bt.reported, bt.updated, bt.solved, IFNULL(bt.image, ''), IFNULL(o.id, 0), IFNULL(o.name, '')
								FROM bike_thefts bt
								LEFT JOIN officers o
								ON o.id = bt.officer
								WHERE bt.id=?`, id)

	theftCase := TheftCase{}

	err = selResult.Scan(&theftCase.ID, &theftCase.TITLE, &theftCase.BRAND,
		&theftCase.CITY, &theftCase.DESCRIPTION, &theftCase.REPORTED, &theftCase.UPDATED, &theftCase.SOLVED, &theftCase.IMAGE, &officerId, &officerName)

	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	defer db.Close()
	theftCase.OFFICER.ID = officerId
	theftCase.OFFICER.NAME = officerName
	json.NewEncoder(w).Encode(&theftCase)
}

func UpdateCase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	var theftCase TheftCase
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&theftCase); err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	db := DbConn()
	updateResult, err := db.Prepare("UPDATE bike_thefts SET solved=? WHERE id=?")
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	_, err = updateResult.Exec(theftCase.SOLVED, id)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	defer db.Close()

	resolved := "unresolved"
	if theftCase.SOLVED {
		resolved = "resolved"
	}
	message := "UPDATE bike theft ID:" + strconv.Itoa(id) + " to " + resolved
	log.Println(message)
	respondWithJSON(w, http.StatusBadRequest, message)
}

func CreateCase(w http.ResponseWriter, r *http.Request) {
	data := r.FormValue("data")
	var theftCase TheftCase
	err := json.Unmarshal([]byte(data), &theftCase)
	if err != nil {
		panic(err.Error())
	}

	if theftCase.TITLE == "" || theftCase.BRAND == "" || theftCase.CITY == "" || theftCase.DESCRIPTION == "" {
		respondWithJSON(w, http.StatusBadRequest, "Some fields are missing please enter them again")
		return
	}
	uploadFile(w, r, &theftCase)

	db := DbConn()
	insert, err := db.Prepare("INSERT INTO bike_thefts(title, brand, city, description, image) VALUES(?,?,?,?,?)")
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = insert.Exec(theftCase.TITLE, theftCase.BRAND, theftCase.CITY, theftCase.DESCRIPTION, theftCase.IMAGE)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	defer db.Close()

	message := "Bike theft: '" + theftCase.TITLE + "' is created"
	log.Println(message)
	respondWithJSON(w, http.StatusAccepted, message)
}

func uploadFile(w http.ResponseWriter, r *http.Request, theftCase *TheftCase) {
	log.Println("File Upload Endpoint Hit")
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("image")
	if err != nil {
		log.Println("error retrieving the file")
		log.Println(err)
		return
	}
	defer file.Close()
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}

	theftCase.IMAGE = fileBytes
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func GetImage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	db := DbConn()
	var img []byte
	selResult := db.QueryRow(`SELECT IFNULL(image, '') FROM bike_thefts WHERE id=?`, id)
	err = selResult.Scan(&img)

	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	defer db.Close()

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(img)))
	if _, err := w.Write(img); err != nil {
		log.Println("unable to write image.")
	}
}
