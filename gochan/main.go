package main

import (
	"log"
	"os"
)

type CatalogResponse struct {
	PreviousPage int
	NextPage     int
	BoardLetter  string
	Page         int `json:"page"`
	Threads      []struct {
		Sub     string `json:"sub"`
		No      int    `json:"no"`
		Tim     int    `json:"tim"`
		Sticky  int    `json:"sticky"`
		Closed  int    `json:"closed"`
		Name    string `json:"name"`
		Com     string `json:"com"`
		Replies int    `json:"replies"`
		Images  int    `json:"images"`
		Time    string `json:"now"`
		// For full res images
		Ext string `json:"ext"`
	} `json:"threads"`
	LastMod int `json:"last_modified"`
}

type ThreadResponse struct {
	BoardLetter string
	ThreadNo    int
	Posts       []struct {
		Sub     string `json:"sub"`
		No      int    `json:"no"`
		Tim     int    `json:"tim"`
		Time    string `json:"now"`
		Name    string `json:"name"`
		Com     string `json:"com"`
		Replies int    `json:"replies"`
		Images  int    `json:"images"`
		// For full res images
		Ext string `json:"ext"`
	}
}

const STATIC_DIR string = "../static"
const IMAGES_DIR string = STATIC_DIR + "/images"
const API_URL = "https://a.4cdn.org"
const IMAGE_API_URL = "https://i.4cdn.org"

func main() {
	srvHost := "0.0.0.0"
	srvPort := os.Getenv("PORT")

	if srvPort == "" {
		log.Print("Server port not specified in environment")
		log.Print("Setting port to default: 4433")
		srvPort = "4433"
	}

	// Build our full server host string
	srvHost = srvHost + ":" + srvPort
	runServer(srvHost)
}
