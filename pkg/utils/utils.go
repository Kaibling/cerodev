package utils

import (
	cryptoRand "crypto/rand"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

// GenerateULID generates and returns a ULID string.
func GenerateULID() string {
	t := time.Now().UTC()
	entropy := ulid.Monotonic(cryptoRand.Reader, 0)

	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

func ContainerName(userID, gitRepo string) string {
	return "cd-" + userID + "-" + strings.Split(gitRepo, "/")[len(strings.Split(gitRepo, "/"))-1]
}

func PrintPrettyJSON(a any) {
	json, _ := json.MarshalIndent(a, "", "  ") //nolint: errchkjson
	fmt.Println(string(json))                  //nolint: forbidigo
}

func RemoveEmptyStrings(s []string) []string {
	var result []string

	for _, str := range s {
		if str != "" {
			result = append(result, str)
		}
	}

	return result
}
