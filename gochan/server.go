package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"text/template"
)

func runServer(srvHost string) {
	fileServer := http.FileServer(http.Dir(STATIC_DIR))
	http.Handle("/", fileServer)
	http.HandleFunc("/get_boards", getBoards)
	http.HandleFunc("/get_thread", getThread)

	// Start server
	log.Print("Starting server on: ", srvHost)
	if err := http.ListenAndServe(srvHost, nil); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

func getBoards(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != "GET" {
		log.Print("Invalid request method received: ", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get required board
	requiredBoard := r.URL.Query()["board"][0]
	numBoards := len(r.URL.Query()["board"])

	/*
		Here we have two variables:
			requiuredPage
			pageIndex
		This is because:
		If the user requests page 1, for example. We need to access the
		first index of our JSON array (as this will be the first page), however
		arrays are indexed from 0 so page 1 is actually index 0.

		Problem is, if we use requiredPage - 1 (0) to access the array, and then
		work out the next page's number (so that the links in the webpage take us
		to the next page), we end up incrementing by 1. 0 + 1 = 1 - meaning we end
		up requesting page 1 as the next page, getting us stuck in a first-page-loop
		forever :')
	*/
	requestedPage := r.URL.Query()["page"][0]
	numPages := len(r.URL.Query()["page"])
	requiredPage, err := strconv.Atoi(requestedPage)
	pageIndex := requiredPage - 1

	if err != nil {
		log.Fatalf("Cannot convert value %s to integer", requestedPage)
	}

	// Ensure user hasn't specified more than 1 board
	if numBoards > 1 {
		log.Print("Number of boards requested: ", numBoards)
		http.Error(w, "Cannot specify more than 1 board", http.StatusBadRequest)
		return
	}

	if numPages > 1 {
		log.Print("Number of pages requested: ", numPages)
		http.Error(w, "Cannot specify more than 1 page", http.StatusBadRequest)
		return
	}

	// Ensure requested board is valid
	if !isValidBoard(requiredBoard) {
		log.Print("Board is not valid: ", requiredBoard)
		http.Error(w, "Non valid board specified: "+requiredBoard, http.StatusBadRequest)
		return
	}
	log.Print("User requested:")
	log.Print("Board: ", requiredBoard)
	log.Print("Page: ", requiredPage)
	catalogAPIEndpoint := API_URL + "/" + requiredBoard + "/catalog.json"

	// Send request
	log.Print("Retrieving data from endpoint: ", catalogAPIEndpoint)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get(catalogAPIEndpoint)
	if err != nil {
		log.Fatalf("Error getting board /%s/ catalog: %s", requiredBoard, err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Endpoint returned status code: ", resp.StatusCode)
	}
	log.Printf("Endpoint %s responsed with status code: %d", catalogAPIEndpoint, resp.StatusCode)

	// Read response
	defer resp.Body.Close()
	log.Print("Reading response from endpoint into memory")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response from board [/%s] catalog: %s",
			requiredBoard, err)
	}

	// Ensure JSON is valid
	log.Print("Ensuring JSON is valid")
	if !json.Valid(body) {
		log.Fatal("JSON is invalid")
	} else {
		log.Print("JSON is valid")
	}

	var result []CatalogResponse
	log.Print("Unmarshalling response into struct")
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal("Could not unmarshal JSON: ", err)
	}

	result[pageIndex].BoardLetter = requiredBoard

	// Remove current thumbnails in IMAGES_DIR
	err = removeAll(IMAGES_DIR + "/*")
	if err != nil {
		log.Fatalf("Failed to clear %s directory: %s", IMAGES_DIR, err)
	}

	var wg sync.WaitGroup
	// Now download all full images
	for _, thread := range result[pageIndex].Threads {
		if thread.Tim == 0 {
			// This post doesn't have a thumbnail - skip it
			continue
		}
		tnDownloadUrl := fmt.Sprintf("%s/%s/%d%s", IMAGE_API_URL, requiredBoard, thread.Tim, thread.Ext)
		imagePath := fmt.Sprintf("%s/%d%s", IMAGES_DIR, thread.Tim, thread.Ext)

		log.Print("Downloading thread full-res image: ", tnDownloadUrl)
		wg.Add(1)
		go func(url, path string) {
			defer wg.Done()
			err := downloadImage(url, path)
			if err != nil {
				log.Fatal("Failed to retrieve thumbnail image: ", err)
			}
		}(tnDownloadUrl, imagePath)
	}

	// Wait for download to complete before rendering page
	wg.Wait()

	/*
		Define our next and previous page numbers
		Since we never have a page 0 (as far as the HTML is concerned)
		We should set the page as 0 if it can't go that way
		i.e if we're on page 1 (first page), we should set previous page
		as 0. We can check for this in the template and use it to conditionally
		render the previous or next page button/link
	*/
	var nextPage, previousPage int
	if requiredPage == len(result) {
		// We're on the last page
		nextPage = 0
		previousPage = requiredPage - 1
	} else if requiredPage == 0 {
		// We're on the first page
		nextPage = requiredPage + 1
		previousPage = 0
	} else {
		// We're somewhere in between
		nextPage = requiredPage + 1
		previousPage = requiredPage - 1
	}

	tmpl := template.Must(template.ParseFiles(STATIC_DIR + "/showboard.html"))
	data := CatalogResponse{
		PreviousPage: previousPage,
		NextPage:     nextPage,
		BoardLetter:  result[pageIndex].BoardLetter,
		Page:         result[pageIndex].Page,
		Threads:      result[pageIndex].Threads,
	}
	log.Print("Rendering catalog page")
	if err := tmpl.Execute(w, data); err != nil {
		log.Fatal(err)
	}
}

func getThread(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != "GET" {
		log.Print("Invalid request method received: ", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get required thread info
	requiredBoard := r.URL.Query()["board"][0]
	requiredThread := r.URL.Query()["thread"][0]
	numThreads := len(r.URL.Query()["thread"])

	// Ensure user hasn't specified more than 1 thread
	if numThreads > 1 {
		log.Print("Number of threads requested: ", numThreads)
		http.Error(w, "Cannot specify more than 1 thread", http.StatusBadRequest)
		return
	}

	log.Print("User requested thread: ", requiredThread)
	threadAPIEndpoint := fmt.Sprintf("%s/%s/thread/%s.json",
		API_URL, requiredBoard, requiredThread)

	// Send request
	log.Print("Retrieving data from endpoint: ", threadAPIEndpoint)
	resp, err := http.Get(threadAPIEndpoint)
	if err != nil {
		log.Fatalf("Error getting thread %s: %s", requiredThread, err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Endpoint returned status code: ", resp.StatusCode)
	}
	log.Printf("Endpoint %s responsed with status code: %d",
		threadAPIEndpoint, resp.StatusCode)

	// Read response
	defer resp.Body.Close()
	log.Print("Reading response from endpoint into memory")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response from thread %s catalog: %s",
			requiredBoard, err)
	}

	// Ensure JSON is valid
	log.Print("Ensuring JSON is valid")
	if !json.Valid(body) {
		log.Fatal("JSON is invalid")
	} else {
		log.Print("JSON is valid")
	}

	var result ThreadResponse
	log.Print("Unmarshalling response into struct")
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal("Could not unmarshal JSON: ", err)
	}

	result.BoardLetter = requiredBoard
	result.ThreadNo, err = strconv.Atoi(requiredThread)
	if err != nil {
		log.Fatalf("Cannot convert value %s to integer", requiredThread)
	}

	// Remove current thumbnails in IMAGES_DIR
	err = removeAll(IMAGES_DIR + "/*")
	if err != nil {
		log.Fatal("Failed to clear "+IMAGES_DIR+" directory: ", err)
	}

	var wg sync.WaitGroup
	// Now download all full images
	for _, thread := range result.Posts {
		if thread.Tim == 0 {
			// This post doesn't have a thumbnail - skip it
			continue
		}
		tnDownloadUrl := fmt.Sprintf("%s/%s/%d%s", IMAGE_API_URL, requiredBoard, thread.Tim, thread.Ext)
		imagePath := fmt.Sprintf("%s/%d%s", IMAGES_DIR, thread.Tim, thread.Ext)

		log.Print("Downloading thread full-res image: ", tnDownloadUrl)
		wg.Add(1)
		go func(url, path string) {
			defer wg.Done()
			err := downloadImage(url, path)
			if err != nil {
				log.Fatal("Failed to retrieve thumbnail image: ", err)
			}
		}(tnDownloadUrl, imagePath)
	}

	// Wait for download to complete before rendering page
	wg.Wait()

	tmpl := template.Must(template.ParseFiles(STATIC_DIR + "/showthread.html"))
	data := ThreadResponse{
		BoardLetter: result.BoardLetter,
		ThreadNo:    result.ThreadNo,
		Posts:       result.Posts,
	}
	log.Print("Rendering thread page")
	if err := tmpl.Execute(w, data); err != nil {
		log.Fatal(err)
	}
}
