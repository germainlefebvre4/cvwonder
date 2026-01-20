package cmdConvert

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/germainlefebvre4/cvwonder/internal/linkedin"
	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	clientID     string
	clientSecret string
	redirectURI  string
	outputFile   string
	profileUser  string
)

func ConvertCmd() *cobra.Command {
	var cobraCmd = &cobra.Command{
		PreRun:  utils.ToggleDebug,
		Use:     "convert",
		Aliases: []string{"c", "conv"},
		Short:   "Convert LinkedIn profile to CV",
		Long:    `Convert your LinkedIn profile to CVWonder YAML format using OAuth2 authentication`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			logrus.Info("CV Wonder - LinkedIn Converter")
			logrus.Info("  Client ID: ", maskString(clientID))
			logrus.Info("  Redirect URI: ", redirectURI)
			logrus.Info("  Output file: ", outputFile)
			logrus.Info("")

			// Validate required parameters
			if clientID == "" || clientSecret == "" {
				logrus.Fatal("Client ID and Client Secret are required. Use --client-id and --client-secret flags.")
			}

			if profileUser == "" {
				logrus.Fatal("Profile username is required. Use --profile flag (e.g., --profile germainlefebvre4)")
			}

			if redirectURI == "" {
				redirectURI = "http://localhost:8080/callback"
			}

			if outputFile == "" {
				outputFile = "cv.yml"
			} // Create services
			authService, err := linkedin.NewAuthService()
			utils.CheckError(err)

			profileService, err := linkedin.NewProfileService()
			utils.CheckError(err)

			converterService, err := linkedin.NewConverterService()
			utils.CheckError(err)

			// Step 1: Get authorization URL
			state := fmt.Sprintf("state_%d", time.Now().Unix())
			authURL := authService.GetAuthorizationURL(clientID, redirectURI, state)

			logrus.Info("Step 1: Authorize the application")
			logrus.Info("Please open the following URL in your browser:")
			fmt.Println()
			fmt.Println(authURL)
			fmt.Println()
			logrus.Info("Waiting for authorization callback...")
			fmt.Println()

			// Step 2: Start local server to receive callback
			authCode, err := waitForCallback(redirectURI, state)
			utils.CheckError(err)

			logrus.Info("")
			logrus.Info("Step 2: Exchanging authorization code for access token...")

			// Step 3: Exchange code for access token
			accessToken, err := authService.ExchangeCodeForToken(ctx, clientID, clientSecret, authCode, redirectURI)
			utils.CheckError(err)

			logrus.Info("✓ Successfully obtained access token")
			logrus.Info("")

			// Step 4: Fetch LinkedIn profile
			logrus.Info("Step 3: Fetching LinkedIn profile for user: ", profileUser)
			profile, err := profileService.GetProfileByUsername(ctx, accessToken, profileUser)
			utils.CheckError(err)

			logrus.Info("✓ Successfully fetched profile for: ", profile.FirstName, " ", profile.LastName)
			logrus.Info("")

			// Step 5: Convert to YAML
			logrus.Info("Step 4: Converting profile to YAML...")
			yamlData, err := converterService.ConvertToYAML(profile)
			utils.CheckError(err)

			// Step 6: Write to file
			err = os.WriteFile(outputFile, yamlData, 0644)
			utils.CheckError(err)

			logrus.Info("✓ Successfully converted and saved to: ", outputFile)
			logrus.Info("")
			logrus.Info("You can now use 'cvwonder generate' to create your CV from this file.")
		},
	}

	cobraCmd.Flags().StringVar(&clientID, "client-id", "", "LinkedIn Application Client ID (required)")
	cobraCmd.Flags().StringVar(&clientSecret, "client-secret", "", "LinkedIn Application Client Secret (required)")
	cobraCmd.Flags().StringVar(&redirectURI, "redirect-uri", "http://localhost:8080/callback", "OAuth2 redirect URI")
	cobraCmd.Flags().StringVarP(&outputFile, "output", "o", "cv.yml", "Output YAML file")
	cobraCmd.Flags().StringVarP(&profileUser, "profile", "p", "", "LinkedIn profile username - REQUIRED (e.g., 'germainlefebvre4' for linkedin.com/in/germainlefebvre4)")

	return cobraCmd
}

// maskString masks a string for display, showing only first and last 3 characters
func maskString(s string) string {
	if len(s) <= 6 {
		return "***"
	}
	return s[:3] + "..." + s[len(s)-3:]
}

// waitForCallback starts a temporary HTTP server to receive the OAuth callback
func waitForCallback(redirectURI, expectedState string) (string, error) {
	// Parse the redirect URI to get the address to listen on
	var listenAddr string
	if strings.Contains(redirectURI, "localhost:8080") {
		listenAddr = ":8080"
	} else if strings.Contains(redirectURI, "localhost:") {
		// Extract port from redirectURI
		parts := strings.Split(redirectURI, ":")
		if len(parts) >= 3 {
			port := strings.Split(parts[2], "/")[0]
			listenAddr = ":" + port
		} else {
			listenAddr = ":8080"
		}
	} else {
		listenAddr = ":8080"
	}

	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	server := &http.Server{Addr: listenAddr}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		errorParam := r.URL.Query().Get("error")
		errorDesc := r.URL.Query().Get("error_description")

		if errorParam != "" {
			errChan <- fmt.Errorf("authorization error: %s - %s", errorParam, errorDesc)
			fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head><title>Authorization Failed</title></head>
<body>
	<h1>Authorization Failed</h1>
	<p>Error: %s</p>
	<p>%s</p>
	<p>You can close this window.</p>
</body>
</html>`, errorParam, errorDesc)
			go func() {
				time.Sleep(1 * time.Second)
				server.Shutdown(context.Background())
			}()
			return
		}

		if state != expectedState {
			errChan <- fmt.Errorf("state mismatch: expected %s, got %s", expectedState, state)
			fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head><title>Authorization Failed</title></head>
<body>
	<h1>Authorization Failed</h1>
	<p>State parameter mismatch. This may be a security issue.</p>
	<p>You can close this window.</p>
</body>
</html>`)
			go func() {
				time.Sleep(1 * time.Second)
				server.Shutdown(context.Background())
			}()
			return
		}

		if code == "" {
			errChan <- fmt.Errorf("no authorization code received")
			fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head><title>Authorization Failed</title></head>
<body>
	<h1>Authorization Failed</h1>
	<p>No authorization code received.</p>
	<p>You can close this window.</p>
</body>
</html>`)
			go func() {
				time.Sleep(1 * time.Second)
				server.Shutdown(context.Background())
			}()
			return
		}

		codeChan <- code
		fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head><title>Authorization Successful</title></head>
<body>
	<h1>✓ Authorization Successful!</h1>
	<p>You can close this window and return to the terminal.</p>
	<script>setTimeout(function() { window.close(); }, 3000);</script>
</body>
</html>`)

		go func() {
			time.Sleep(1 * time.Second)
			server.Shutdown(context.Background())
		}()
	})

	go func() {
		logrus.Debug("Starting callback server on ", listenAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("failed to start callback server: %w", err)
		}
	}()

	// Wait for either code or error, with timeout
	select {
	case code := <-codeChan:
		return code, nil
	case err := <-errChan:
		return "", err
	case <-time.After(5 * time.Minute):
		server.Shutdown(context.Background())
		return "", fmt.Errorf("authorization timeout after 5 minutes")
	}
}
