package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func removeAll(path string) error {
	contents, err := filepath.Glob(path)
	if err != nil {
		return err
	}
	for _, item := range contents {
		err = os.Remove(item)
		if err != nil {
			return err
		}
	}
	return nil
}

func downloadImage(tnDownloadUrl string, imagePath string) error {
	// Download image
	resp, err := http.Get(tnDownloadUrl)
	if err != nil {
		log.Print("Failed to download thread thumbnail image from: ", tnDownloadUrl)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Response gave status code: ", resp.StatusCode)
	}
	defer resp.Body.Close()

	// Create image file
	file, err := os.Create(imagePath)
	if err != nil {
		log.Print("Failed to create thread thumbnail image file: ", imagePath)
		return err
	}
	defer file.Close()

	// Copy image content to file
	file.ReadFrom(resp.Body)
	if err != nil {
		log.Print("Failed to save thread thumbnail image: ", imagePath)
		return err
	}
	return nil
}

// PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func isValidBoard(board string) bool {
	validBoards := []string{
		"a", "c", "w", "m", "cgl",
		"cm", "cm", "n", "jp", "vt", "v",
		"vg", "vmg", "vm", "vp", "vr", "vrpg",
		"vst", "co", "g", "tv", "k", "o",
		"an", "tg", "sp", "xs", "pw", "sci",
		"his", "int", "out", "toy", "i", "po",
		"p", "ck", "ic", "wg", "lit", "mu",
		"fa", "3", "gd", "diy", "wsg", "qst",
		"biz", "trv", "fit", "x", "adv", "lgbt",
		"mlp", "news", "wsr", "vip", "b", "r9k",
		"pol", "bant", "soc", "s4s", "s", "hc",
		"hm", "h", "e", "u", "d", "y",
		"t", "hr", "gif", "aco", "r",
	}

	for _, curBoard := range validBoards {
		if curBoard == board {
			return true
		}
	}
	return false
}
