package function

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Handle a serverless request
func Handle(payload []byte) string {
	url := "https://api.github.com/gists"

	var jsonStr = []byte(`{
                "description": ` + fmt.Sprintf("A gist of %d bytes", len(payload)) + `,
                "public": true,
                "files": {
                        "file1.txt": {
                            "content": "` + string(payload) + `"
                        }
                    }
                }`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "token "+readSecret()) // The token
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.Status == "201 Created" {
		return resp.Header.Get("Location")
	}

	resBody, _ := ioutil.ReadAll(resp.Body)

	fmt.Fprintf(os.Stderr, fmt.Sprintf("Couldn't create file %d %s\n", resp.StatusCode, string(resBody)))
	os.Exit(1)

	return ""
}

func readSecret() string {
	val, err := ioutil.ReadFile("/var/openfaas/secrets/github-token")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	return strings.TrimSpace(string(val))
}
