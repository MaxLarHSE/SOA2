package errorhandler

import (
	"SOA2/pkg/api"
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, httpStatus int, errorCode api.ErrorResponseErrorCode, message string, details *map[string]interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	resp := api.ErrorResponse{
		ErrorCode: errorCode,
		Message:   message,
		Details:   details,
	}

	json.NewEncoder(w).Encode(resp)
}
func CustomErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	details := map[string]interface{}{
		"parse_error": err.Error(),
	}
	WriteError(w, http.StatusBadRequest, api.VALIDATIONERROR, "Invalid request parameters", &details)
}
