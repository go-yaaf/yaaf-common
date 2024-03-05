package utils

import (
	"fmt"
	. "github.com/go-yaaf/yaaf-common/config"
	"net/url"
	"regexp"
)

// PrintBanner displays the application's banner and lists all configuration variables along with their values in alphabetical order.
// This function is designed to provide a clear overview of the application's current configuration state at startup or when needed,
// facilitating easier debugging and configuration verification. It ensures that the configuration variables are presented in a sorted manner
// for quick reference and readability.
func PrintBanner(banner string) {
	fmt.Print(banner)
	allConfigVars := Get().GetAllVars()
	for _, key := range Get().GetAllKeysSorted() {
		// test if valid URI
		val, err := obfuscateCredentials(allConfigVars[key])
		if err == nil {
			fmt.Printf("%s: %s\n", key, val)
		} else {
			fmt.Printf("%s: %s\n", key, allConfigVars[key])
		}
	}
}

// obfuscateCredentials takes a string as input and, based on its format, either:
// - obfuscates user credentials in a URI,
// - obfuscates the first 10 characters of an API key,
// - or returns the original string if it doesn't match the above formats.
func obfuscateCredentials(input string) (string, error) {
	// Attempt to parse the input as a URI.
	parsedURI, err := url.Parse(input)
	if err == nil && parsedURI.Scheme != "" && parsedURI.Host != "" {
		// Input is a valid URI.
		if userInfo := parsedURI.User; userInfo != nil {
			username := userInfo.Username()
			obfuscatedUserInfo := username + ":*****"

			obfuscatedURI := ""
			if parsedURI.Scheme != "" {
				obfuscatedURI += parsedURI.Scheme + "://"
			}
			obfuscatedURI += obfuscatedUserInfo + "@"
			obfuscatedURI += parsedURI.Host
			if parsedURI.Path != "" {
				obfuscatedURI += parsedURI.Path
			}
			if parsedURI.RawQuery != "" {
				obfuscatedURI += "?" + parsedURI.RawQuery
			}
			if parsedURI.Fragment != "" {
				obfuscatedURI += "#" + parsedURI.Fragment
			}

			return obfuscatedURI, nil
		}
		// URI has no credentials to obfuscate.
		return input, nil
	} else {
		// Check if the input is in the form of an API key (assumed to be a hexadecimal string).
		matched, _ := regexp.MatchString(`^[a-fA-F0-9]{11,}$`, input)
		if matched && len(input) > 10 {
			// Obfuscate the first 10 characters of the API key.
			return "**********" + input[9:], nil
		}
	}

	// Input does not match any special handling criteria; return it unchanged.
	return input, nil
}
