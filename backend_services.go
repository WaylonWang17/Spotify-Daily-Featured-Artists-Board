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
	"log"
	"math/rand"
	"time"

	"github.com/joho/godotenv"
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
	//checks to see if we can load the env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

func getArtist(auth AuthData, artistID string) map[string]any{
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
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("spotify_data.json", body, 0644)
    if err != nil {
        panic(err)
    }

	var artist map[string]any //variable artist where each key is a string and value can be "any"
	json.Unmarshal(body, &artist)
	return artist
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

func getRandomArtist(auth AuthData, genre string) map[string]any {
	url := fmt.Sprintf("https://api.spotify.com/v1/search?q=genre:%s&type=artist&limit=50", url.QueryEscape(genre))

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+auth.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]any
	json.Unmarshal(body, &result)

	artists, ok := result["artists"].(map[string]any)["items"].([]any)
	if !ok || len(artists) == 0 {
		fmt.Println("⚠️ No artists found for genre:", genre)
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	randomArtist := artists[rand.Intn(len(artists))].(map[string]any)

	// ✅ Extract the first available image (usually largest)
	var imageURL string
	if imgs, ok := randomArtist["images"].([]any); ok && len(imgs) > 0 {
		image := imgs[0].(map[string]any)
		imageURL = image["url"].(string)
	}

	topTrack := getTopTrack(auth, randomArtist["id"].(string))


	// ✅ Return simplified artist object including the image
	return map[string]any{
		"name":   randomArtist["name"],
		"id":     randomArtist["id"],
		"genres": randomArtist["genres"],
		"href":   randomArtist["external_urls"].(map[string]any)["spotify"],
		"image":  imageURL,
		"genre":  genre,
		"followers": randomArtist["followers"].(map[string]any)["total"],
		"top_track": topTrack,
	}
}

func getTopTrack(auth AuthData, artistID string) map[string]any {
	url := fmt.Sprintf("https://api.spotify.com/v1/artists/%s/top-tracks?market=US", artistID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil{
		panic(err)
	}

	req.Header.Add("Authorization", "Bearer "+auth.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error fetching top track:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil{
		panic(err)
	}

	var result map[string]any
	json.Unmarshal(body, &result)

	tracks, ok := result["tracks"].([]any)
	if !ok || len(tracks) == 0 {
		return nil
	}

	topTrack := tracks[0].(map[string]any) // take the first (most popular)
	return map[string]any{
		"name":   topTrack["name"],
		"href":   topTrack["external_urls"].(map[string]any)["spotify"],
		"preview": topTrack["preview_url"], // might be nil for some songs
	}
}

func main() {
	auth := requestAccessToken()

	genres := []string{
		"pop",
		"hip hop",
		"rap",
		"r&b",
		"rock",
		"alternative",
		"indie",
		"country",
		"jazz",
		"classical",
		"metal",
		"edm",
		"dance",
		"house",
		"techno",
		"folk",
		"soul",
		"punk",
		"blues",
		"reggae",
		"latin",
		"k-pop",
		"j-pop",
	}

	//deprecated consider swapping out later
	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len(genres), func(i, j int) { genres[i], genres[j] = genres[j], genres[i] })
	selectedGenres := genres[:3]
	fmt.Println("🎧 Selected genres:", selectedGenres)

	

	var all []map[string]any //declares dynamic array called all

	//loop through artistids and append artist into all 
	for _, genre := range selectedGenres {
		artist := getRandomArtist(auth, genre)
		if artist != nil{
			all = append(all, artist)
		}
	}

	data, _ := json.MarshalIndent(all, "", "  ") //converts all slice into pretty printed json where "" => no prefix for each line and " " => indent each nested level with 2 spaces
	os.WriteFile("spotify_data.json", data, 0644)
	fmt.Println("✅ Wrote 3 artists to spotify_data.json")
}
