package function

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Gist struct {
	Description string             `json:"description"`
	Public      bool               `json:"public"`
	Files       map[string]Content `json:"files"`
}

type Content struct {
	Content string `json:"content"`
}

// Handle a serverless request
func Handle(payload []byte) string {

	if os.Getenv("Http_Method") != "POST" {
		fmt.Fprintf(os.Stderr, `You must post a body to this function to be stored. Use the HTTP header "X-Filename" to override the filename`)
		os.Exit(1)
	}

	url := "https://api.github.com/gists"

	filename := "body.base64"
	if val, ok := os.LookupEnv("Http_X_Filename"); ok && len(val) > 0 {
		filename = val
	}

	if strings.Contains(filename, "..") {
		fmt.Fprintf(os.Stderr, `You must not specify a relative path.`)
		os.Exit(1)
	}

	if strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		fmt.Fprintf(os.Stderr, `You must not specify a path`)
		os.Exit(1)
	}

	gist := Gist{
		Description: fmt.Sprintf("Encoded %d bytes as base64", len(payload)),
		Public:      true,
		Files: map[string]Content{
			filename: Content{
				Content: base64.StdEncoding.EncodeToString(payload),
			},
		},
	}

	jsonStr, jerr := json.Marshal(gist)
	if jerr != nil {
		return jerr.Error()
	}

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
		res, getErr := http.Get(resp.Header.Get("Location"))
		if getErr != nil {
			fmt.Fprintf(os.Stderr, getErr.Error())
			os.Exit(1)
		}

		bytesOut, _ := ioutil.ReadAll(res.Body)
		gistResult := GistResult{}
		json.Unmarshal(bytesOut, &gistResult)
		return gistResult.HtmlURL
	}

	resBody, _ := ioutil.ReadAll(resp.Body)

	fmt.Fprintf(os.Stderr, fmt.Sprintf("Couldn't create file %d %s\n", resp.StatusCode, string(resBody)))
	os.Exit(1)

	return ""
}

type GistResult struct {
	HtmlURL string `json:"html_url"`
}

func readSecret() string {
	val, err := ioutil.ReadFile("/var/openfaas/secrets/github-token")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	return strings.TrimSpace(string(val))
}
