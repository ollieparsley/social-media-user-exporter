package platform

import "os"

//Interface ...
type Interface interface {
	Fetch() error
}

// getEnv Get an env vairable and set a default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv("SMUE_" + key); ok {
		return value
	}
	return fallback
}
