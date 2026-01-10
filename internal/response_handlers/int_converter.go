package reponsehandlers

import (
	"cmd/internal/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ParseIdResponse(r *http.Request, w http.ResponseWriter, fieldName string) (int64, bool) {
	strId := chi.URLParam(r, fieldName)
	num, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		log.Printf("failed to parse %s: Error: %v\n", fieldName, err)
		json.WriteError(w, http.StatusBadRequest, "Invalid Request")
		return 0, false
	}
	return int64(num), true
}

func ParseId32Response(r *http.Request, w http.ResponseWriter, fieldName string) (int32, bool) {
	strId := chi.URLParam(r, fieldName)
	num, err := strconv.ParseInt(strId, 10, 32)
	if err != nil {
		log.Printf("failed to parse %s: Error: %v\n", fieldName, err)
		json.WriteError(w, http.StatusBadRequest, "Invalid Request")
		return 0, false
	}
	return int32(num), true

}

func ParseSymbolName(r *http.Request, w http.ResponseWriter, fieldName string) string {
	strId := chi.URLParam(r, fieldName)
	return strId
}
