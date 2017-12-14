package routes

import (
	"net/http"

	"zenithar.org/go/common/web/utils"
)

//------------------------------------------------------------------------------
// Resources
//------------------------------------------------------------------------------

// HelloResponse contains the result of the Hello request.
type HelloResponse struct {
	Message string `json:"message"`
	Version string `json:"version"`
}

//------------------------------------------------------------------------------
// Handlers
//------------------------------------------------------------------------------

// Hello shows basic information about the API on its frontpage.
func Hello(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, 200, &HelloResponse{
		Message: "Scheduler API",
		// Version: version.Version,
	})
}
