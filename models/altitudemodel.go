package models

import (
	"casaaltitudine/lang"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

// swagger:model Altitude
type Altitude struct {
	// Id of rain value
	// in: int64
	Id int64 `json:"id"`
	// Value of Altitude
	// in: int
	Value float32 `json:"valore"`
	// Timestamp of insert
	// in: time
	DateInsert time.Time `json:"data_inserimento"`
}

type Altitudes []Altitude

type ReqAddAltitude struct {
	// Value of the Altitude
	// in: int
	Value float32 `json:"valore" validate:"required"`
}

// swagger:parameters add Altitude
type ReqAltitudeBody struct {
	// - name: body
	//  in: body
	//  description: Altitude
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/ReqAddAltitude"
	//  required: true
	Body ReqAddAltitude `json:"body"`
}

// ErrHandler returns error message bassed on env debug
func ErrHandler(err error) string {
	var errmessage string
	if os.Getenv("DEBUG") == "true" {
		errmessage = err.Error()
	} else {
		errmessage = lang.Get("something_went_wrong")
	}
	return errmessage
}

func GetAltitudesSqlx(db *sql.DB) *Altitudes {
	altitudes := Altitudes{}
	rows, err := db.Query("SELECT id, valore, data_inserimento FROM altitudine order by id desc")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Altitude
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			log.Fatal(err)
		}
		altitudes = append(altitudes, p)
	}
	return &altitudes
}
func GetLastAltitudeSqlx(db *sql.DB) *Altitudes {
	altitudes := Altitudes{}
	rows, err := db.Query("SELECT id, valore, data_inserimento FROM altitudine where id = (select max(id) from altitudine)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Altitude
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			log.Fatal(err)
		}
		altitudes = append(altitudes, p)
	}
	return &altitudes
}
func GetLastHourSqlx(db *sql.DB) *Altitudes {
	altitudes := Altitudes{}

	tFine := time.Now()
	dataFine := tFine.Format("2006-01-02 15:04:05")

	tInizio := time.Now().Add(time.Duration(-1) * time.Hour)
	dataInizio := tInizio.Format("2006-01-02 15:04:05")

	sqlStatement := fmt.Sprintf("SELECT id,valore,data_inserimento FROM altitudine where data_inserimento  >= '%s' AND data_inserimento <= '%s'", dataInizio, dataFine)

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Altitude
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			log.Fatal(err)
		}
		altitudes = append(altitudes, p)
	}

	if len(altitudes) == 0 {
		elemento := GetLastAltitudeSqlx(db)
		altitudes = append(altitudes, *elemento...)
	}
	return &altitudes
}

func GetShowDataSqlx(db *sql.DB, recordNumber int) *Altitudes {
	altitudes := Altitudes{}

	sqlStatement := fmt.Sprintf("WITH t AS (SELECT id,valore,data_inserimento FROM altitudine ORDER BY data_inserimento DESC LIMIT %d) SELECT id,valore,data_inserimento FROM t ORDER BY data_inserimento ASC", recordNumber)

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Altitude
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			log.Fatal(err)
		}
		altitudes = append(altitudes, p)
	}

	return &altitudes
}

// PostAltitudeSqlx insert Altitude value
func PostAltitudeSqlx(db *sql.DB, reqAltitude *ReqAddAltitude) (*Altitude, string) {

	value := reqAltitude.Value

	var altitude Altitude

	lastInsertId := 0

	//sqlStatement := fmt.Sprintf("insert into 'pioggia' ('valore','data_inserimento') values (%d,CURRENT_TIMESTAMP) RETURNING id", value)
	sqlStatement := fmt.Sprintf("insert into altitudine (valore,data_inserimento) values (%.2f,CURRENT_TIMESTAMP) RETURNING id", value)

	err := db.QueryRow(sqlStatement).Scan(&lastInsertId)

	if err != nil {
		return &altitude, ErrHandler(err)
	}

	sqlStatement1 := fmt.Sprintf("SELECT id,valore,data_inserimento FROM altitudine where id = %d", lastInsertId)
	rows, err := db.Query(sqlStatement1)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Altitude
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		altitude = p
	}
	if err != nil {
		return &altitude, lang.Get("no_result")
	}
	return &altitude, ""
}
