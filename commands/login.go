package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/urfave/cli"
)

const (
	AuthServiceUrl = "https://rezkit-staging-cli-auth.fly.dev"
	MaxAttempts    = 60
)

type loginResponse struct {
	AuthURL   string `json:"url"`
	SessionID string `json:"id"`
}

func authServiceUrl() string {
	if url := os.Getenv("REZKIT_AUTH_SERVICE_URL"); url != "" {
		return url
	}

	return AuthServiceUrl
}

func Login(cli *cli.Context) error {

	// Initialize a new session with the auth service
	resp, err := http.Post(authServiceUrl()+"/sessions/new", "text/plain", nil)

	if err != nil {
		return err
	}

	response := loginResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	fmt.Println("Attempting to open browser to authenticate")

	if runtime.GOOS == "linux" {
		os.StartProcess("xdg-open", []string{response.AuthURL}, nil)
	}

	fmt.Println("Please visit this URL to authenticate: ", response.AuthURL)

	ticker := time.NewTicker(2 * time.Second)

	go func() {
		for attempts := 0; attempts < MaxAttempts; attempts++ {
			// Wait for the ticker...
			<-ticker.C

			resp, err := http.Get(authServiceUrl() + "/token?id=" + response.SessionID)

			if err != nil {
				break
			}

			if resp.StatusCode == 200 {

			}
		}
	}()

	fmt.Println("Waiting for Authentication")

	return nil
}

func getTokens(ctx context.Context) {

}
