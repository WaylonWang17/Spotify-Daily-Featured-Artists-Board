package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type AuthData struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func requestAccessToken() AuthData {
	/*
		flow for api requests
		1. http request
		2. add header
		3. do request and get response
		4. check status code
		5. read response body
	*/
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	endpoint := "https://accounts.spotify.com/api/token"

	//encode clientid:clientsecret in base64
	//used for authentication
	authStr := clientID + ":" + clientSecret
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authStr)) //encoding from bytes to string

	//form - url.Values{} is a built in go type
	form := url.Values{}
	form.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		panic(err) //hard brake stops the entire program
	}

	//url, headers, form, json
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+encodedAuth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		print(resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var auth AuthData

	err = json.Unmarshal(body, &auth)
	if err != nil {
		panic(err)
	}

	return auth
}

func getAccess() AuthData {
	file, err := os.Open("authentication.json")
	if err != nil {
		fmt.Println("Error opening the file", err)
		panic(err)
	}

	defer file.Close()

	//decode json into struct
	var auth AuthData

	err = json.NewDecoder(file).Decode(&auth) //newdecoder is just reading the file we passed in and decode is filling the authdata struct with matching variables from the file
	if err != nil {
		panic(err)
	}

	return auth

}

func getArtist(auth AuthData) {
	artistID := "1ok4DP80jKsX7GZZ6yr2xR?si=DljzMFSiQ-6D4BFNIzJRGQ"
	url := fmt.Sprintf("https://api.spotify.com/v1/artists/%s", artistID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	//set header
	req.Header.Set("Authorization", auth.TokenType+" "+auth.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}

func getTrack(auth AuthData) {
	trackID := "2plbrEY59IikOBgBGLjaoe?si=8fbdac58be9c4b91"
	url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", trackID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", auth.TokenType+" "+auth.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}

func main() {
	auth := requestAccessToken()
	// auth := getAccess()
	getArtist(auth)
	// getTrack(auth)

}
