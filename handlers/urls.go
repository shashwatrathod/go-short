package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shashwatrathod/url-shortner/core"
	"github.com/shashwatrathod/url-shortner/middleware"
)

// CreateUrlAliasRequest defines the request body for creating a URL alias.
//
// @Description Request body for creating a URL alias.
type CreateUrlAliasRequest struct {
	OriginalUrl string `json:"originalUrl" validate:"required,url" example:"https://example.com/very/long/url/to/shorten"`
}

// CreateUrlAliasResponse defines the response body for a created URL alias.
//
// @Description Response body for a created URL alias.
type CreateUrlAliasResponse struct {
	UrlAlias string `json:"urlAlias" example:"aBcDeFg1"`
}

// CreateUrlAliasHandler handles HTTP requests for creating a new URL alias
// or retrieving an existing one for a given original URL.
// It expects a CreateUrlRequest in the request body.
//
// On success, it responds with a JSON object containing the aliased URL.
// If an error occurs during processing (e.g., issues with application environment,
// database operations, or URL generation), it logs the error and responds with
// an HTTP 500 Internal Server Error.
//
// @Summary Create or get a URL alias
// @Description Creates a new URL alias for a given original URL or returns an existing one.
// @Tags urls
// @Accept json
// @Produce json
// @Param request body CreateUrlAliasRequest true "Request body to create a URL alias"
// @Success 200 {object} CreateUrlAliasResponse "Successfully created or retrieved alias"
// @Failure 400 {object} middleware.ValidationError "Invalid request payload (validation error)"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /create [post]
func CreateUrlAliasHandler(w http.ResponseWriter, r *http.Request, req CreateUrlAliasRequest) {
	appEnv, ok := r.Context().Value(middleware.ContextAppEnvKey).(*middleware.AppEnv)

	if !ok || appEnv == nil {
		log.Printf("CreateUrlAliasHandler: Error accessing AppEnv.")

		SendInternalServerError(w, "CreateUrlAliasHandler: Error accessing AppEnv.")
		return
	}

	existingAlias, err := appEnv.UrlAliasDao.FindByOriginalUrl(r.Context(), req.OriginalUrl)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		SendInternalServerError(w, "CreateUrlAliasHandler: Unexpected error while processing request.")
		return
	}

	if existingAlias != nil {
		log.Printf("Found an existing alias for %s : %s", req.OriginalUrl, existingAlias.Alias)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&CreateUrlAliasResponse{
			UrlAlias: existingAlias.Alias,
		})
		return
	}

	// no existing urls in db - create new short url.
	shortUrl, err := core.GenerateAlias(req.OriginalUrl, appEnv.AliasingStrategy)

	if err != nil {
		log.Printf("CreateShortUrlHandler: Unexpected error while generating alias : %s.", err)
		SendInternalServerError(w, "CreateShortUrlHandler: Unexpected error while generating alias.")
		return
	}

	urlAlias, err := appEnv.UrlAliasDao.CreateUrlAlias(r.Context(), shortUrl, req.OriginalUrl)

	if err != nil {
		log.Printf("CreateShortUrlHandler: Unexpected error while saving alias : %s.", err)
		SendInternalServerError(w, "CreateShortUrlHandler: Unexpected error while saving alias.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&CreateUrlAliasResponse{
		UrlAlias: urlAlias.Alias,
	})
}

// GetUrlAliasHandler handles HTTP requests to retrieve and redirect to an original URL
// based on a given alias.
//
// If the alias is found, it redirects the client to the original URL (HTTP 302).
// If the alias is not found, it responds with an HTTP 404 Not Found.
// If an internal error occurs, it responds with an HTTP 500 Internal Server Error.
// @Summary Redirect to original URL
// @Description Retrieves the original URL for a given alias and redirects to it.
// @Tags urls
// @Produce html
// @Param alias path string true "URL Alias" example:"aBcDeFg1"
// @Success 302 "Redirects to the original URL (Location header will be set)"
// @Failure 404 {object} ErrorResponse "Alias not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /{alias} [get]
func GetUrlAliasHandler(w http.ResponseWriter, r *http.Request) {
	appEnv, ok := r.Context().Value(middleware.ContextAppEnvKey).(*middleware.AppEnv)

	if !ok || appEnv == nil {
		log.Printf("GetUrlAliasHandler: Error accessing AppEnv.")
		SendInternalServerError(w, "GetUrlAliasHandler: Error accessing AppEnv.")
		return
	}

	vars := mux.Vars(r)
	alias := vars["alias"]

	existingAlias, err := appEnv.UrlAliasDao.FindByAlias(r.Context(), alias)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		SendInternalServerError(w, "GetUrlAliasHandler: Unexpected error while processing request.")
		return
	}

	if existingAlias != nil {
		http.Redirect(w, r, existingAlias.OriginalURL, http.StatusFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(ErrorResponse{Error: "Not Found", Message: "The requested alias was not found."})
}
