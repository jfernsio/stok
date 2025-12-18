package marketux

import (
	"encoding/json"
	"fmt"
	"io"

	// "io/ioutil"
	"net/http"

	"github.com/fatih/color"
	// "net/url"
)

type Response struct {
	Meta struct {
		Found    int `json:"found"`
		Returned int `json:"returned"`
		Limit    int `json:"limit"`
		Page     int `json:"page"`
	} `json:"meta"`
	Data []struct {
		Title       string `json:"title"`
		Snippet     string `json:"snippet"`
		Url         string `json:"url"`
		ImageUrl    string `json:"image_url"`
		PublishedAt string `json:"published_at"`
		Source      string `json:"source"`
	} `json:"data"`
}

func Run() {
	baseURL := "https://api.marketaux.com/v1/news/all?symbols=TSLA&api_token=mujnLdZZuIS6KimqhvLMBUJSmhKCAi6Iw5jHyH7B"
	// api_token := "mujnLdZZuIS6KimqhvLMBUJSmhKCAi6Iw5jHyH7B"
	// endpoint := "symbols=TSLA"

	// url := baseURL + endpoint + "&api_token=" + api_token

	resp, err := http.Get(baseURL)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		fmt.Printf("API error %d: %s\n", resp.StatusCode, string(body))
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	var res Response
	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Println("JSON parse error:", err)
		return
	}
	color.Cyan("Found: %d, Returned: %d\n", res.Meta.Found, res.Meta.Returned)
	for _, item := range res.Data {
		color.Yellow("Title: %s\n", item.Title)
		color.White("Snippet: %s\n", item.Snippet)
		color.Blue("URL: %s\n", item.Url)
		color.Green("Published At: %s\n", item.PublishedAt)
		color.Red("Source: %s\n", item.Source)
		fmt.Println("---")
	}
}


