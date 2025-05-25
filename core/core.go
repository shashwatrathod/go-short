package core

// desired length of the short url. can support 62^8 unique strings
const SHORT_URL_LEN = 8;

// generates the short url from the original url using the provided strategy
func GenerateShortUrl(original_url string, shorteningStrategy ShorteningStrategy) (string, error) {
	return shorteningStrategy.Shorten(original_url, SHORT_URL_LEN), nil
}