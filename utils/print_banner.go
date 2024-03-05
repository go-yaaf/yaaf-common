package utils

import (
	"fmt"
	. "github.com/go-yaaf/yaaf-common/config"
)

// PrintBanner displays the application's banner and lists all configuration variables along with their values in alphabetical order.
// This function is designed to provide a clear overview of the application's current configuration state at startup or when needed,
// facilitating easier debugging and configuration verification. It ensures that the configuration variables are presented in a sorted manner
// for quick reference and readability.
func PrintBanner(banner string) {
	fmt.Print(banner)
	allConfigVars := Get().GetAllVars()
	for _, key := range Get().GetAllKeysSorted() {
		fmt.Printf("%s: %s\n", key, allConfigVars[key])
	}
}
