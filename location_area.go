package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationAreaResp struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func fetchLocationAreas(url string) (LocationAreaResp, error) {
	resp, err := http.Get(url)
	if err != nil {
		return LocationAreaResp{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaResp{}, err
	}

	var data LocationAreaResp
	err = json.Unmarshal(body, &data)
	if err != nil {
		return LocationAreaResp{}, err
	}

	for _, area := range data.Results {
		fmt.Println(area.Name)
	}

	return data, nil
}
