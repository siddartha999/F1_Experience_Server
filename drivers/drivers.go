package drivers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	constants "github.com/siddartha999/F1_Experience_Server/config"
)

type Driver struct {
	DriverNumber  uint8  `json:"driver_number"`
	BroadCastName string `json:"broadcast_name"`
	FullName      string `json:"full_name"`
	NameAcronym   string `json:"name_acronym"`
	TeamName      string `json:"team_name"`
	TeamColor     string `json:"team_colour"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	HeadShotURL   string `json:"headshot_url"`
	CountryCode   string `json:"country_code"`
}

// Fetches F1 driver's info from a remote server
func getDriversInfoFromRemote() {
	cacheFile, err := os.Create(constants.DriversCacheFile)
	if err != nil {
		fmt.Println("Unable to Create driver's cache file")
	}
	defer cacheFile.Close()

	resp, err := http.Get("https://api.openf1.org/v1/drivers")
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Unable to retrieve drivers info from remote")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unable to read the response body")
	}

	err = os.WriteFile(constants.DriversCacheFile, body, 0644)
	if err != nil {
		fmt.Println("Error caching driver's info")
	}
}

// Parses the cached driver's information
func parseDriversCache() []Driver {
	driversCacheBytes, err := os.ReadFile(constants.DriversCacheFile)
	if err != nil {
		fmt.Println("Error parsing driver's cache file")
	}

	allDrivers := []Driver{}
	err = json.Unmarshal(driversCacheBytes, &allDrivers)
	if err != nil {
		fmt.Println("Error parsing driver's cache bytes", err)
	}

	//Extract unique drivers
	uniqueDriverEntries := []Driver{}
	isDriverEntryUnique := map[string]bool{}
	for idx := 0; idx < len(allDrivers); idx++ {
		currentDriverEntry := allDrivers[idx]
		uniqueIdentifier := string(currentDriverEntry.DriverNumber) + ":" + currentDriverEntry.NameAcronym + ":" + currentDriverEntry.FullName
		exists := isDriverEntryUnique[uniqueIdentifier]
		if !exists {
			uniqueDriverEntries = append(uniqueDriverEntries, currentDriverEntry)
			isDriverEntryUnique[uniqueIdentifier] = true
		}
	}
	fmt.Printf("Found %d unique driver entries", len(uniqueDriverEntries))
	return uniqueDriverEntries
}

func InitiateDriversInfo() []Driver {
	driversCacheFileSize, err := os.Stat(constants.DriversCacheFile)
	if err != nil || driversCacheFileSize.Size() == 0 {
		fmt.Println("Driver's cache file does not exist on the server. Retrieving from remote")
		getDriversInfoFromRemote()
	}
	fmt.Println("Parsing cached driver's info")
	return parseDriversCache()
}
