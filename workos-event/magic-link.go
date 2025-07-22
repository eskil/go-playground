package main

import (
	"net/http"
	"os"
	"fmt"
	"github.com/workos/workos-go/pkg/sso"
)

func main() {
	apiKey := os.Getenv("WORKOS_API_KEY")
	clientID := os.Getenv("WORKOS_CLIENT_ID")

	sso.Configure(apiKey, clientID)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		opts := sso.GetProfileAndTokenOptions{
			Code: r.URL.Query().Get("code"),
		}

		profileAndToken, err := sso.GetProfileAndToken(r.Context(), opts)

		if err != nil {
			fmt.Printf("error: %s", err)
			// Handle the error...
		}

		// Use the information in `profile` for further business logic.
		profile := profileAndToken.Profile

		http.Redirect(w, r, "/", http.StatusSeeOther)

		fmt.Printf("profile %s", profile)
	})
}
