package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
)

func run() {
    baseURL := "https://marketaux.com"

    baseURL.Path += "v1/news/all"

    params := url.Values{}
    params.Add("api_token", "YOUR_API_TOKEN")
    params.Add("symbols", "aapl,tsla")
    params.Add("search", "ipo")
    params.Add("limit", "50")

    baseURL.RawQuery = params.Encode()

    req, _ := http.NewRequest("GET", baseURL.String(), nil)

    res, _ := http.DefaultClient.Do(req)

    defer res.Body.Close()

    body, _ := ioutil.ReadAll(res.Body)

    fmt.Println(string(body))
}