package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type Dataset struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string            `json:"type"`
	Geometry   Geometry          `json:"geometry"`
	Properties datasetProperties `json:"properties"`
}

type datasetProperties struct {
	ID          string
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Address     string `json:"eventAddress"`
}

func MapboxDataset(index int, artist string) {
	indexInt := strconv.Itoa(index)
	datasetName := strings.Replace(artist, " ", "-", -1) + "-tourdates_std"
	// Step 1: Create a new dataset
	datasetID, err := createDataset(datasetName, "Tourdata for "+artist)
	if err != nil {
		fmt.Println("Error creating dataset:", err)
		return
	}
	fmt.Println("Created dataset with ID:", datasetID)

	// Step 2: Add features to the dataset from the JSON file
	features, err := readFeaturesFromFile("db/mapbox_std/" + indexInt + ".geojson")
	if err != nil {
		fmt.Println("Error reading features:", err)
		return
	}

	for _, feature := range features.Features {
		featureID := feature.Properties.ID
		if featureID == "" {
			featureID = uuid.New().String()
		}
		err = addFeature(datasetID, featureID, feature)
		if err != nil {
			fmt.Println("Error adding feature:", err)
		} else {
			fmt.Println("Feature added successfully:", featureID)
		}
	}
}

func createDataset(name, description string) (string, error) {
	url := fmt.Sprintf("https://api.mapbox.com/datasets/v1/%s?access_token=%s", os.Getenv("MAPBOX_USERNAME"), os.Getenv("MAPBOX_ACCESS_TOKEN"))

	dataset := Dataset{Name: name, Description: description}
	data, err := json.Marshal(dataset)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	//if resp.StatusCode != http.StatusCreated {
	//	body, _ := io.ReadAll(resp.Body)
	//	return "", fmt.Errorf("failed to create dataset: %s", body)
	//}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result["id"].(string), nil
}

func readFeaturesFromFile(filePath string) (*FeatureCollection, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var features FeatureCollection
	err = json.NewDecoder(file).Decode(&features)
	if err != nil {
		return nil, err
	}

	return &features, nil
}

func addFeature(datasetID, featureID string, feature Feature) error {
	url := fmt.Sprintf("https://api.mapbox.com/datasets/v1/%s/%s/features/%s?access_token=%s", os.Getenv("MAPBOX_USERNAME"), datasetID, featureID, os.Getenv("MAPBOX_ACCESS_TOKEN"))

	data, err := json.Marshal(feature)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to add feature: %s", body)
	}

	return nil
}
