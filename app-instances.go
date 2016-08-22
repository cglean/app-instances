package main

import (
	"fmt"
	"encoding/json"
	"strconv"
	"strings"
	"github.com/cloudfoundry/cli/plugin"
)

// AppInstances represents Buildpack Usage CLI interface
type AppInstances struct{}

// Metadata is the data retrived from the response json
type Metadata struct {
	GUID string `json:"guid"`
}

// GetMetadata provides the Cloud Foundry CLI with metadata to provide user about how to use buildpack-usage command
func (c *AppInstances) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "app-instances",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 1,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "app-instances",
				HelpText: "Command to view the number of instances running for an app.",
				UsageDetails: plugin.Usage{
					Usage: "cf app-instances [app name]",
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(AppInstances))
}

// Run is what is executed by the Cloud Foundry CLI when the buildpack-usage command is specified
func (c AppInstances) Run(cli plugin.CliConnection, args []string) {
	if args[0] == "app-instances" {
		apps := c.GetAppData(cli)
		if len(args) == 2 {
			for _, val := range apps.Resources {
				if args[1] == val.Entity.Name {
	   			fmt.Printf("%v", val.Entity.Instances)
					return
				}
			}
		}
	}
	fmt.Printf("%s", "Applicaiton not found")
}

// AppSearchResults represents top level attributes of JSON response from Cloud Foundry API
type AppSearchResults struct {
	TotalResults int                  `json:"total_results"`
	TotalPages   int                  `json:"total_pages"`
	Resources    []AppSearchResources `json:"resources"`
}

// AppSearchResources represents resources attribute of JSON response from Cloud Foundry API
type AppSearchResources struct {
	Entity   AppSearchEntity `json:"entity"`
	Metadata Metadata        `json:"metadata"`
}

// AppSearchEntity represents entity attribute of resources attribute within JSON response from Cloud Foundry API
type AppSearchEntity struct {
	Name              string `json:"name"`
	Buildpack         string `json:"buildpack"`
	DetectedBuildpack string `json:"detected_buildpack"`
	SpaceGUID         string `json:"space_guid"`
	Instances         int    `json:"instances"`
	State             string `json:"state"`
	Memory            int    `json:"memory"`
	DiskQuota         int    `json:"disk_quota"`
}

// GetAppData requests all of the Application data from Cloud Foundry
func (c AppInstances) GetAppData(cli plugin.CliConnection) AppSearchResults {
	var res AppSearchResults
	res = c.UnmarshallAppSearchResults("/v2/apps?order-direction=asc&results-per-page=100", cli)

	if res.TotalPages > 1 {
		for i := 2; i <= res.TotalPages; i++ {
			apiUrl := fmt.Sprintf("/v2/apps?order-direction=asc&page=%v&results-per-page=100", strconv.Itoa(i))
			tRes := c.UnmarshallAppSearchResults(apiUrl, cli)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func (c AppInstances) UnmarshallAppSearchResults(apiUrl string, cli plugin.CliConnection) AppSearchResults {
	var tRes AppSearchResults
	cmd := []string{"curl", apiUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &tRes)

	return tRes
}
