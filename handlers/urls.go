package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shashwatrathod/url-shortner/core"
	"github.com/shashwatrathod/url-shortner/middleware"
)


type CreateUrlAliasRequest struct {
	OriginalUrl string `json:"originalUrl" validate:"required,url"`
}

type CreateUrlAliasResponse struct {
	UrlAlias string `json:"urlAlias"`
}


// CreateUrlAliasHandler handles HTTP requests for creating a new URL alias
// or retrieving an existing one for a given original URL.
// It expects a CreateUrlRequest in the request body.
//
// On success, it responds with a JSON object containing the aliased URL.
// If an error occurs during processing (e.g., issues with application environment,
// database operations, or URL generation), it logs the error and responds with
// an HTTP 500 Internal Server Error.
func CreateUrlAliasHandler(w http.ResponseWriter, r *http.Request, req CreateUrlAliasRequest) {
	appEnv, ok := r.Context().Value(middleware.ContextAppEnvKey).(*middleware.AppEnv)

	if !ok || appEnv == nil {
		log.Printf("CreateUrlAliasHandler: Error accessing AppEnv.")
		http.Error(w, "CreateUrlAliasHandler: Error accessing AppEnv.", http.StatusInternalServerError);
		return
	}

	existingAlias, err := appEnv.UrlAliasDao.FindByOriginalUrl(r.Context(), req.OriginalUrl);

	if err != nil {
		log.Printf("Error: %s", err.Error())
		http.Error(w, "CreateUrlAliasHandler: Unexpected error while processing request.", http.StatusInternalServerError);
		return
	}

	if existingAlias != nil {
		log.Printf("Found an existing alias for %s : %s", req.OriginalUrl, existingAlias.Alias);
		json.NewEncoder(w).Encode(&CreateUrlAliasResponse{
			UrlAlias: existingAlias.Alias,
		})
		return
	}

	// no existing urls in db - create new short url.
	shortUrl, err := core.GenerateAlias(req.OriginalUrl, appEnv.AliasingStrategy)

	if err != nil {
		log.Printf("CreateShortUrlHandler: Unexpected error while generating alias : %s.", err)
		http.Error(w, "CreateShortUrlHandler: Unexpected error while generating alias.", http.StatusInternalServerError);
	}

	urlAlias, err := appEnv.UrlAliasDao.CreateUrlAlias(r.Context(), shortUrl, req.OriginalUrl)

	if err != nil {
		log.Printf("CreateShortUrlHandler: Unexpected error while saving alias : %s.", err)
		http.Error(w, "CreateShortUrlHandler: Unexpected error while saving alias.", http.StatusInternalServerError);
	}

	json.NewEncoder(w).Encode(&CreateUrlAliasResponse{
		UrlAlias: urlAlias.Alias,
	})
}

func GetUrlAliasHandler(w http.ResponseWriter, r *http.Request) {
	appEnv, ok := r.Context().Value(middleware.ContextAppEnvKey).(*middleware.AppEnv)

	if !ok || appEnv == nil {
		log.Printf("GetUrlAliasHandler: Error accessing AppEnv.")
		http.Error(w, "GetUrlAliasHandler: Error accessing AppEnv.", http.StatusInternalServerError);
		return
	}

	vars := mux.Vars(r)
	alias := vars["alias"]

	existingAlias, err := appEnv.UrlAliasDao.FindByAlias(r.Context(), alias);

	if err != nil {
		log.Printf("Error: %s", err.Error())
		http.Error(w, "GetUrlAliasHandler: Unexpected error while processing request.", http.StatusInternalServerError);
		return
	}

	if existingAlias != nil {
		http.Redirect(w, r, existingAlias.OriginalURL, http.StatusFound)
		return
	}

	http.NotFound(w, r)
}