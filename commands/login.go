package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/rezkit/cli/internal/config"
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

type oAuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
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
		os.StartProcess("xdg-open", []string{response.AuthURL}, &os.ProcAttr{})
	}

	fmt.Println("Please visit this URL to authenticate: ", response.AuthURL)

	ticker := time.NewTicker(2 * time.Second)

	fmt.Println("Waiting for Authentication")

	tokenData := oAuthResponse{}
	attempts := 0

	for attempts = 0; attempts < MaxAttempts; attempts++ {
		// Wait for the ticker...
		<-ticker.C

		fmt.Println("Checking auth state")

		resp, e := http.Get(authServiceUrl() + "/token?id=" + response.SessionID)

		if e != nil {
			err = e
			break
		}

		if resp.StatusCode == 200 {

			if e := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
				err = e
				break
			}

			if tokenData.AccessToken != "" {
				break
			}
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to authenticate. Please try again")
		return err
	}

	if attempts >= MaxAttempts {
		fmt.Fprintln(os.Stderr, "Authentication timed out after", MaxAttempts, "attempts.")
		return nil
	}

	expires := time.Now().Add(time.Duration(tokenData.ExpiresIn) * time.Second)

	// Store the auth data in the config
	config.GetConfig().Set("authentication.access_token", tokenData.AccessToken)
	config.GetConfig().Set("authentication.refresh_token", tokenData.RefreshToken)
	config.GetConfig().Set("authentication.expires", expires)

	return nil
}
