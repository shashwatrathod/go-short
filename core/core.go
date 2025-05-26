package core

// desired length of the url alias. can support 62^8 unique strings
const ALIAS_LEN = 8;

// generates the url alias from the original url using the provided strategy
func GenerateAlias(original_url string, aliasingStrategy AliasingStrategy) (string, error) {
	return aliasingStrategy.Alias(original_url, ALIAS_LEN), nil
}