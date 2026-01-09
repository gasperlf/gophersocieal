package main

import (
	"net/http"
)

type HelthCheckResponse struct {
	Status  string `json:"status"`
	Env     string `json:"env"`
	Version string `json:"version"`
}

// healthcheck godoc
//
//	@Summary		Health check
//	@Description	Health check
//	@Tags			ops
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	HelthCheckResponse
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	response := HelthCheckResponse{
		Status:  "ok",
		Env:     app.config.env,
		Version: version,
	}

	if err := writeJSON(w, http.StatusOK, response); err != nil {
		app.internalServerError(w, r, err)
	}
}
