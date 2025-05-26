package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shashwatrathod/url-shortner/core"
	"github.com/shashwatrathod/url-shortner/middleware"
)


type CreateUrlRequest struct {
	OriginalUrl string `json:"originalUrl" validate:"required,url"`
}

type CreateUrlResponse struct {
	ShortUrl string `json:"shortUrl"`
}


// CreateShortUrlHandler handles HTTP requests for creating a new short URL
// or retrieving an existing one for a given original URL.
// It expects a CreateUrlRequest in the request body.
//
// On success, it responds with a JSON object containing the short URL.
// If an error occurs during processing (e.g., issues with application environment,
// database operations, or URL generation), it logs the error and responds with
// an HTTP 500 Internal Server Error.
func CreateShortUrlHandler(w http.ResponseWriter, r *http.Request, req CreateUrlRequest) {
	appEnv, ok := r.Context().Value(middleware.ContextAppEnvKey).(*middleware.AppEnv)

	if !ok || appEnv == nil {
		log.Printf("CreateShortUrlHandler: Error accessing AppEnv.")
		http.Error(w, "CreateShortUrlHandler: Error accessing AppEnv.", http.StatusInternalServerError);
		return
	}

	// 1. Check if the URL already exists in any of the shards.
	existingShortUrl, err := appEnv.ShortURLDAO.FindByOriginalUrl(r.Context(), req.OriginalUrl);

	if err != nil {
		log.Printf("Error: %s", err.Error())
		http.Error(w, "CreateShortUrlHandler: Unexpected error while processing request.", http.StatusInternalServerError);
		return
	}

	if existingShortUrl != nil {
		log.Printf("Found an existing short Url for %s : %s", req.OriginalUrl, existingShortUrl.ShortURL);
		json.NewEncoder(w).Encode(&CreateUrlResponse{
			ShortUrl: existingShortUrl.ShortURL,
		})
		return
	}

	// no existing urls in db - create new short url.
	shortUrl, err := core.GenerateShortUrl(req.OriginalUrl, appEnv.ShorteningStrategy)

	if err != nil {
		log.Printf("CreateShortUrlHandler: Unexpected error while generating ShortUrl : %s.", err)
		http.Error(w, "CreateShortUrlHandler: Unexpected error while generating ShortUrl.", http.StatusInternalServerError);
	}

	dbShortUrl, err := appEnv.ShortURLDAO.CreateShortURL(r.Context(), shortUrl, req.OriginalUrl)

	if err != nil {
		log.Printf("CreateShortUrlHandler: Unexpected error while saving ShortUrl : %s.", err)
		http.Error(w, "CreateShortUrlHandler: Unexpected error while saving ShortUrl.", http.StatusInternalServerError);
	}

	json.NewEncoder(w).Encode(&CreateUrlResponse{
		ShortUrl: dbShortUrl.ShortURL,
	})
}

func GetShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	appEnv, ok := r.Context().Value(middleware.ContextAppEnvKey).(*middleware.AppEnv)

	if !ok || appEnv == nil {
		log.Printf("GetShortUrlHandler: Error accessing AppEnv.")
		http.Error(w, "GetShortUrlHandler: Error accessing AppEnv.", http.StatusInternalServerError);
		return
	}

	vars := mux.Vars(r)
	shortUrl := vars["shortUrl"]

	existingShortUrl, err := appEnv.ShortURLDAO.FindByShortUrl(r.Context(), shortUrl);

	if err != nil {
		log.Printf("Error: %s", err.Error())
		http.Error(w, "GetShortUrlHandler: Unexpected error while processing request.", http.StatusInternalServerError);
		return
	}

	if existingShortUrl != nil {
		http.Redirect(w, r, existingShortUrl.OriginalURL, http.StatusFound)
		return
	}

	http.NotFound(w, r)
}