package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cobra"
)

type DoorDashAccessKey struct {
	DeveloperId   string `json:"developer_id"`
	KeyId         string `json:"key_id"`
	SigningSecret string `json:"signing_secret"`
}

var DoorDashAccessKeyInput string
var DoorDashAccessKeyFilePathInput string
var DurationInput int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "make-doordash-jwt",
	Short: "Make a DoorDash Developer JWT",
	Long:  `make-doordash-jwt is a CLI library for generating a DoorDash Developer JWT.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate the key inputs
		if DoorDashAccessKeyInput == "" && DoorDashAccessKeyFilePathInput == "" {
			color.Red("Error: no DoorDash Access Key JSON provided. Provide an Access Key with the -o or -f flags\n\n")
			cmd.Help()
			os.Exit(1)
		}

		if DoorDashAccessKeyInput != "" && DoorDashAccessKeyFilePathInput != "" {
			color.Red("Error: provide only the -o or -f flag, not both\n\n")
			cmd.Help()
			os.Exit(1)
		}

		// Unmarshal the Access Key JSON
		accessKey := DoorDashAccessKey{}
		if DoorDashAccessKeyInput != "" {
			if err := json.Unmarshal([]byte(DoorDashAccessKeyInput), &accessKey); err != nil {
				color.Red("Couldn't parse the Access Key JSON you provided.\n    %v\n", err)
				os.Exit(1)
			}
		} else if DoorDashAccessKeyFilePathInput != "" {
			file, err := os.Open(DoorDashAccessKeyFilePathInput)
			if err != nil {
				color.Red("Couldn't open the file you provided\n    %v\n", err)
				os.Exit(1)
			}
			defer file.Close()

			if err := json.NewDecoder(file).Decode(&accessKey); err != nil {
				color.Red("Couldn't parse the Access Key JSON you provided\n    %v\n", err)
				os.Exit(1)
			}
		}

		// Validate the duration
		if DurationInput < 1 {
			color.Red("Durations less than 1 minute aren't supported. Defaulting to 30 minutes.")
			DurationInput = 30
		}
		if DurationInput > 30 {
			color.Red("Durations greater than 30 minutes aren't supported. Defaulting to 30 minutes.")
			DurationInput = 30
		}

		// Get the JWT
		jwt, err := GetJWT(accessKey)
		if err != nil {
			fmt.Printf("Couldn't generate a JWT\n    %v\n", err)
			os.Exit(1)
		}
		fmt.Println(jwt)
	},
}

func GetJWT(accessKey DoorDashAccessKey) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": accessKey.DeveloperId,
		"aud": "doordash",
		"exp": time.Now().Add(time.Minute * time.Duration(DurationInput)).Unix(),
		"iat": time.Now().Unix(),
		"kid": accessKey.KeyId,
	})

	t.Header["dd-ver"] = "DD-JWT-V1"

	decodedSigningSecret, err := base64.RawURLEncoding.DecodeString(accessKey.SigningSecret)
	if err != nil {
		return "", err
	}

	jwt, err := t.SignedString(decodedSigningSecret)
	if err != nil {
		return "", err
	}
	return jwt, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&DoorDashAccessKeyInput, "object", "o", "", "JSON object with developer_id, key_id, and signing_secret copied from the DoorDash Developer Portal.")
	rootCmd.Flags().StringVarP(&DoorDashAccessKeyFilePathInput, "file", "f", "", "path to a file containing JSON object with developer_id, key_id, and signing_secret copied from the DoorDash Developer Portal.")
	rootCmd.Flags().IntVarP(&DurationInput, "duration", "d", 30, "duration in minutes for which the JWT should be valid.")
}
