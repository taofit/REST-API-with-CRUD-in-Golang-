package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type CaseOfficer struct {
	CASE    int `json:case`
	OFFICER int `json:officer`
}

func AssignCases() {
	var availableOfficerIds = getIdsOfAvailableOfficer()
	lenOfAvilableOfficer := len(availableOfficerIds)
	lenOfAvailableCase := getNumOfAvailableCase()

	updateNumOfCase := lenOfAvailableCase
	if updateNumOfCase > lenOfAvilableOfficer {
		updateNumOfCase = lenOfAvilableOfficer
	}

	for _, id := range availableOfficerIds[:updateNumOfCase] {
		updateBikeTheft(id)
	}
}

func getIdsOfAvailableOfficer() []int {
	db := dbConn()
	rows, err := db.Query(`SELECT o.id FROM officers o
		LEFT JOIN bike_thefts bt
		ON o.id = bt.officer AND bt.solved = 0
		WHERE bt.id IS NULL`)

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var officerIds []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			panic(err.Error())
		}
		officerIds = append(officerIds, id)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return officerIds
}

func getNumOfAvailableCase() int {
	db := dbConn()
	availableCaseNum := 0

	err := db.QueryRow(
		`SELECT COUNT(id) FROM bike_thefts
		WHERE solved = 0 AND officer IS NULL`).Scan(&availableCaseNum)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	return availableCaseNum
}

func updateBikeTheft(officerId int) {
	db := dbConn()
	updateResult, err := db.Prepare("UPDATE bike_thefts SET officer=? WHERE solved = 0 AND officer IS NULL LIMIT 1")
	if err != nil {
		panic(err.Error())
	}
	updateResult.Exec(officerId)
	log.Println("UPDATE: bike_theft table with offficer id")
	defer db.Close()
}

func AssignCaseToEnOfficer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var caseOfficer CaseOfficer

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&caseOfficer); err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	if err, message := checkOfficer(w, caseOfficer.OFFICER); err {
		respondWithJSON(w, http.StatusBadRequest, message)
		return
	}

	db := dbConn()
	updateResult, err := db.Prepare(`UPDATE bike_thefts 
										SET officer = ?
										WHERE 
										NOT EXISTS(
										SELECT 1 FROM (SELECT solved,officer FROM bike_thefts) AS temp
										WHERE solved = 0 AND officer = ?
										) AND solved = 0 AND officer IS NULL AND id = ?
										LIMIT 1`)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	sqlResult, err := updateResult.Exec(caseOfficer.OFFICER, caseOfficer.OFFICER, caseOfficer.CASE)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	updatedRow, _ := sqlResult.RowsAffected()
	defer db.Close()

	var message string
	if updatedRow == 1 {
		message = fmt.Sprintf("Successfully assign case id: %d to officer id: %d", caseOfficer.CASE, caseOfficer.OFFICER)
	} else {
		message = fmt.Sprintf("Case id: %d is not assigned to officer id: %d", caseOfficer.CASE, caseOfficer.OFFICER)

	}
	respondWithJSON(w, http.StatusAccepted, message)
}

func checkOfficer(w http.ResponseWriter, id int) (bool, string) {
	db := dbConn()
	officerId := 0
	err := db.QueryRow(
		`SELECT id FROM officers
		WHERE id = ?`, id).Scan(&officerId)
	if err != nil {
		return true, err.Error()
	}
	defer db.Close()

	return false, ""
}
