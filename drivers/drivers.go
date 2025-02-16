package drivers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	constants "github.com/siddartha999/F1_Experience_Server/config"
)

func getDriversFromRemote() {
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

func InitiateDriversInfo() {
	driversCacheFileSize, err := os.Stat(constants.DriversCacheFile)
	if err != nil || driversCacheFileSize.Size() == 0 {
		fmt.Println("Driver's cache file does not exist on the server. Retrieving from remote")
		getDriversFromRemote()
	}
}
