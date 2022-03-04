package controllers

import (
	"casaaltitudine/lang"
	"casaaltitudine/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// BaseHandler will hold everything that controller needs
type BaseHandlerSqlx struct {
	db *sqlx.DB
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandlerSqlx(db *sqlx.DB) *BaseHandlerSqlx {
	return &BaseHandlerSqlx{
		db: db,
	}
}

// swagger:model CommonError
type CommonError struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
}

// swagger:model CommonSuccess
type CommonSuccess struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
}

// swagger:model GetAltitudes
type GetAltitudes struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the response
	// in: string
	Message string            `json:"message"`
	Data    *models.Altitudes `json:"data"`
}

// swagger:model GetAltitude
type GetAltitude struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the response
	// in: string
	Message string `json:"message"`
	// Umidity value
	Data *models.Altitude `json:"data"`
}

// ErrHandler returns error message response
func ErrHandler(errmessage string) *CommonError {
	errresponse := CommonError{}
	errresponse.Status = 0
	errresponse.Message = errmessage
	return &errresponse
}

// swagger:route GET /altitudes listAltitude
// Get Altitude list
//
// security:
// - apiKey: []
// responses:
//  401: CommonError
//  200: GetAltitudes
func (h *BaseHandlerSqlx) GetAltitudesSqlx(w http.ResponseWriter, r *http.Request) {
	response := GetAltitudes{}

	altitudes := models.GetAltitudesSqlx(h.db.DB)

	response.Status = 1
	response.Message = lang.Get("success")
	response.Data = altitudes

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// swagger:route GET /lasthour lastHour
// Get list of last hour of altitude values .... or the last value inserted
//
// security:
// - apiKey: []
// responses:
//  401: CommonError
//  200: GetAltitudes
func (h *BaseHandlerSqlx) GetLastHourSqlx(w http.ResponseWriter, r *http.Request) {
	response := GetAltitudes{}

	altitudes := models.GetLastHourSqlx(h.db.DB)

	response.Status = 1
	response.Message = lang.Get("success")
	response.Data = altitudes

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// swagger:route POST /altitude addAltitude
// Create a new altitude value
//
// security:
// - apiKey: []
// responses:
//  401: CommonError
//  200: GetAltitude
func (h *BaseHandlerSqlx) PostAltitudeSqlx(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	response := GetAltitude{}

	decoder := json.NewDecoder(r.Body)
	var reqAltitude *models.ReqAddAltitude
	err := decoder.Decode(&reqAltitude)
	fmt.Println(err)

	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler(lang.Get("invalid_request")))
		return
	}

	rain, errmessage := models.PostAltitudeSqlx(h.db.DB, reqAltitude)
	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	response.Status = 1
	response.Message = lang.Get("insert_success")
	response.Data = rain
	json.NewEncoder(w).Encode(response)
}
