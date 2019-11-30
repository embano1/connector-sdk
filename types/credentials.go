package types

import (
	"log"
	"os"

	"github.com/openfaas/faas-provider/auth"
)

// GetCredentials returns a pointer to basic auth credentials using environment
// variable lookup for runtime secrets retrieved via "secret_mount_path"
// environment variable.
//
// It panics if credentials cannot be retrieved.
func GetCredentials() *auth.BasicAuthCredentials {
	var credentials *auth.BasicAuthCredentials

	if val, ok := os.LookupEnv("basic_auth"); ok && len(val) > 0 {
		if val == "true" || val == "1" {

			reader := auth.ReadBasicAuthFromDisk{}

			if val, ok := os.LookupEnv("secret_mount_path"); ok && len(val) > 0 {
				reader.SecretMountPath = os.Getenv("secret_mount_path")
			}

			res, err := reader.Read()
			if err != nil {
				log.Fatalln(err)
			}
			credentials = res
		}
	}
	return credentials
}
