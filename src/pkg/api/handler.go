package api

import (
	"fmt"
	"net/http"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("PROJECT_ID")
	if name == "" {
		name = "CBI-Project"
	}

	fmt.Fprintf(w, "CBI data collection microservices' goroutines have started for %s!\n", name)
}
