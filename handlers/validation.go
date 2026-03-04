package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"

	"SOA2/handlers/errorhandler"
	"SOA2/pkg/api"
)

func NewValidationMiddleware() func(http.Handler) http.Handler {
	spec, err := openapi3.NewLoader().LoadFromFile("api.yaml")
	if err != nil {
		log.Fatalf("failed to load openapi spec: %v", err)
	}

	if err := spec.Validate(context.Background()); err != nil {
		log.Fatalf("invalid openapi spec: %v", err)
	}

	spec.Servers = nil

	return nethttpmiddleware.OapiRequestValidatorWithOptions(spec,
		&nethttpmiddleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
			},
			ErrorHandler: validationErrorHandler,
		},
	)
}

func validationErrorHandler(w http.ResponseWriter, message string, statusCode int) {
	fieldErrors := parseValidationMessage(message)

	details := map[string]interface{}{
		"validation_errors": fieldErrors,
	}

	errorhandler.WriteError(
		w,
		http.StatusBadRequest,
		api.VALIDATIONERROR,
		"Invalid request parameters",
		&details,
	)
}

func parseValidationMessage(message string) []string {
	var errors []string
	for _, line := range strings.Split(message, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && line != "|" {
			errors = append(errors, line)
		}
	}
	if len(errors) == 0 {
		errors = append(errors, message)
	}
	return errors
}
