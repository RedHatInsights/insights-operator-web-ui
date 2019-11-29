/*
Copyright © 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tisnik/insights-operator-web-ui/types"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// APIPrefix represents part of URL that is appended before the actual endpoint address
const APIPrefix = "/api/v1/"

var controllerURL = ""

func performReadRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Expected HTTP status 200 OK, got %d", response.StatusCode)
	}
	body, readErr := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if readErr != nil {
		return nil, fmt.Errorf("Unable to read response body")
	}

	return body, nil
}

func performWriteRequest(url string, method string, payload io.Reader) error {
	var client http.Client

	request, err := http.NewRequest(method, url, payload)
	if err != nil {
		return fmt.Errorf("Error creating request %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Expected HTTP status 200 OK, 201 Created or 202 Accepted, got %d", response.StatusCode)
	}
	return nil
}

func readListOfClusters(controllerURL string, apiPrefix string) ([]types.Cluster, error) {
	clusters := []types.Cluster{}

	url := controllerURL + apiPrefix + "client/cluster"
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &clusters)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}

func readListOfConfigurationProfiles(controllerURL string, apiPrefix string) ([]types.ConfigurationProfile, error) {
	profiles := []types.ConfigurationProfile{}

	url := controllerURL + apiPrefix + "client/profile"
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &profiles)
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

func readListOfConfigurations(controllerURL string, apiPrefix string) ([]types.ClusterConfiguration, error) {
	configurations := []types.ClusterConfiguration{}

	url := controllerURL + apiPrefix + "client/configuration"
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &configurations)
	if err != nil {
		return nil, err
	}
	return configurations, nil
}

func readListOfTriggers(controllerURL string, apiPrefix string, clusterName string) ([]types.Trigger, error) {
	var triggers []types.Trigger
	url := controllerURL + apiPrefix + "client/cluster/" + clusterName + "/trigger"
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &triggers)
	if err != nil {
		return nil, err
	}
	return triggers, nil
}

func readListOfAllTriggers(controllerURL string, apiPrefix string) ([]types.Trigger, error) {
	var triggers []types.Trigger
	url := controllerURL + apiPrefix + "client/trigger"
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &triggers)
	if err != nil {
		return nil, err
	}
	return triggers, nil
}

func readConfigurationProfile(controllerURL string, apiPrefix string, profileID string) (*types.ConfigurationProfile, error) {
	var profile types.ConfigurationProfile
	url := controllerURL + apiPrefix + "client/profile/" + profileID
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func readClusterConfigurationByID(controllerURL string, apiPrefix string, configurationID string) (*string, error) {
	url := controllerURL + apiPrefix + "client/configuration/" + configurationID
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	str := string(body)
	return &str, nil
}

func getContentType(filename string) string {
	// TODO: to map
	if strings.HasSuffix(filename, ".html") {
		return "text/html"
	} else if strings.HasSuffix(filename, ".js") {
		return "application/javascript"
	} else if strings.HasSuffix(filename, ".css") {
		return "text/css"
	}
	return "text/html"
}

func sendStaticPage(writer http.ResponseWriter, filename string) {
	body, err := ioutil.ReadFile(filename)
	if err == nil {
		writer.Header().Set("Server", "A Go Web Server")
		writer.Header().Set("Content-Type", getContentType(filename))
		fmt.Fprint(writer, string(body))
	} else {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
	}
}

func staticPage(filename string) func(writer http.ResponseWriter, request *http.Request) {
	log.Println("Serving static file", filename)
	return func(writer http.ResponseWriter, request *http.Request) {
		sendStaticPage(writer, filename)
	}
}

// ListClustersDynContent represents dynamic part of HTML page with list of clusters
type ListClustersDynContent struct {
	Items []types.Cluster
}

func listClusters(writer http.ResponseWriter, request *http.Request) {
	clusters, err := readListOfClusters(controllerURL, APIPrefix)
	if err != nil {
		log.Println("Error reading list of clusters", err)
		return
	}

	t, err := template.ParseFiles("html/list_clusters.html")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}

	dynData := ListClustersDynContent{Items: clusters}
	err = t.Execute(writer, dynData)
	if err != nil {
		println("Error executing template")
	}
}

// ListProfilesDynContent represents dynamic part of HTML page with list of configuration profiles
type ListProfilesDynContent struct {
	Items []types.ConfigurationProfile
}

func listProfiles(writer http.ResponseWriter, request *http.Request) {
	profiles, err := readListOfConfigurationProfiles(controllerURL, APIPrefix)
	if err != nil {
		log.Println("Error reading list of configuration profiles", err)
		return
	}

	t, err := template.ParseFiles("html/list_profiles.html")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}

	dynData := ListProfilesDynContent{Items: profiles}
	err = t.Execute(writer, dynData)
	if err != nil {
		println("Error executing template")
	}
}

// ListConfigurationsDynContent represents dynamic part of HTML page with list of configurations
type ListConfigurationsDynContent struct {
	Items []types.ClusterConfiguration
}

// ListTriggersDynContent represents dynamic part of HTML page with list of triggers
type ListTriggersDynContent struct {
	Items []types.Trigger
}

var epoch = time.Unix(0, 0).Format(time.RFC1123)

var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

func listConfigurations(writer http.ResponseWriter, request *http.Request) {
	configurations, err := readListOfConfigurations(controllerURL, APIPrefix)
	// NoCache headers
	for k, v := range noCacheHeaders {
		writer.Header().Set(k, v)
	}

	if err != nil {
		log.Println("Error reading list of cluster configurations", err)
		return
	}

	t, err := template.ParseFiles("html/list_configurations.html")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}

	dynData := ListConfigurationsDynContent{Items: configurations}
	err = t.Execute(writer, dynData)
	if err != nil {
		println("Error executing template")
	}
}

func listTriggers(writer http.ResponseWriter, request *http.Request) {
	clusterName, ok := request.URL.Query()["clusterName"]
	var triggers []types.Trigger
	var err error

	if !ok {
		triggers, err = readListOfAllTriggers(controllerURL, APIPrefix)
	} else {
		triggers, err = readListOfTriggers(controllerURL, APIPrefix, clusterName[0])
	}

	// NoCache headers
	for k, v := range noCacheHeaders {
		writer.Header().Set(k, v)
	}

	if err != nil {
		log.Println("Error reading list of triggers", err)
		return
	}

	log.Println(triggers)
	t, err := template.ParseFiles("html/list_triggers.html")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}

	dynData := ListTriggersDynContent{Items: triggers}
	err = t.Execute(writer, dynData)
	if err != nil {
		println("Error executing template")
	}
}

// DescribeConfigurationDynContent represents dynamic part of HTML page with configuration description
type DescribeConfigurationDynContent struct {
	Configuration types.ConfigurationProfile
}

func describeConfiguration(writer http.ResponseWriter, request *http.Request) {
	configID, ok := request.URL.Query()["configuration"]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}

	configuration, err := readConfigurationProfile(controllerURL, APIPrefix, configID[0])
	fmt.Println(configuration)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}

	t, err := template.ParseFiles("html/describe_configuration.html")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Error parsing template")
		return
	}

	dynData := DescribeConfigurationDynContent{Configuration: *configuration}
	err = t.Execute(writer, dynData)
	if err != nil {
		println("Error executing template")
	}
}

func storeProfile(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	form := request.Form

	username := form.Get("username")
	description := form.Get("description")
	configuration := form.Get("configuration")

	log.Println("username", username)
	log.Println("description", description)
	log.Println("configuration", configuration)

	query := "username=" + url.QueryEscape(username) + "&description=" + url.QueryEscape(description)
	url := controllerURL + APIPrefix + "client/profile?" + query

	err := performWriteRequest(url, "POST", strings.NewReader(configuration))
	if err != nil {
		log.Println("Error communicating with the service", err)
		http.Redirect(writer, request, "/profile-not-created", 301)
	} else {
		log.Println("Configuration profile has been created")
		http.Redirect(writer, request, "/profile-created", 301)
	}
}

func storeConfiguration(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	form := request.Form

	username := form.Get("username")
	cluster := form.Get("cluster")
	reason := form.Get("reason")
	description := form.Get("description")
	configuration := form.Get("configuration")

	log.Println("username", username)
	log.Println("cluster", cluster)
	log.Println("reason", reason)
	log.Println("description", description)
	log.Println("configuration", configuration)

	query := "username=" + url.QueryEscape(username) + "&reason=" + url.QueryEscape(reason) + "&description=" + url.QueryEscape(description)
	url := controllerURL + APIPrefix + "client/cluster/" + url.PathEscape(cluster) + "/configuration?" + query

	err := performWriteRequest(url, "POST", strings.NewReader(configuration))
	if err != nil {
		log.Println("Error communicating with the service", err)
		http.Redirect(writer, request, "/configuration-not-created", 301)
	} else {
		log.Println("Configuration has been created")
		http.Redirect(writer, request, "/configuration-created", 301)
	}
}

func enableConfiguration(writer http.ResponseWriter, request *http.Request) {
	configurationID, ok := request.URL.Query()["id"]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}
	url := controllerURL + APIPrefix + "client/configuration/" + configurationID[0] + "/enable"
	err := performWriteRequest(url, "PUT", nil)
	if err != nil {
		fmt.Println("Error communicating with the service", err)
		return
	}

	// everything is ok, configuration has been enabled
	fmt.Println("Configuration " + configurationID[0] + " has been enabled")
	http.Redirect(writer, request, "/list-configurations", 307)
}

func disableConfiguration(writer http.ResponseWriter, request *http.Request) {
	configurationID, ok := request.URL.Query()["id"]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}
	url := controllerURL + APIPrefix + "client/configuration/" + configurationID[0] + "/disable"
	err := performWriteRequest(url, "PUT", nil)
	if err != nil {
		fmt.Println("Error communicating with the service", err)
		return
	}

	// everything is ok, configuration has been disabled
	fmt.Println("Configuration " + configurationID[0] + " has been disabled")
	http.Redirect(writer, request, "/list-configurations", 307)
}

func activateTrigger(writer http.ResponseWriter, request *http.Request) {
	triggerID, ok := request.URL.Query()["id"]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}
	url := controllerURL + APIPrefix + "client/trigger/" + triggerID[0] + "/activate"

	err := performWriteRequest(url, "PUT", nil)
	if err != nil {
		fmt.Println("Error communicating with the service", err)
		return
	}

	// everything is ok, trigger has been activated
	fmt.Println("Trigger " + triggerID[0] + " has been activated")
	http.Redirect(writer, request, "/list-triggers", 307)
}

func deactivateTrigger(writer http.ResponseWriter, request *http.Request) {
	triggerID, ok := request.URL.Query()["id"]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}
	url := controllerURL + APIPrefix + "client/trigger/" + triggerID[0] + "/deactivate"

	err := performWriteRequest(url, "PUT", nil)
	if err != nil {
		fmt.Println("Error communicating with the service", err)
		return
	}

	// everything is ok, trigger has been deactivated
	fmt.Println("Trigger " + triggerID[0] + " has been deactivated")
	http.Redirect(writer, request, "/list-triggers", 307)
}

func triggerMustGatherConfiguration(writer http.ResponseWriter, request *http.Request) {
	clusterID, ok := request.URL.Query()["clusterID"]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}
	id, err := strconv.Atoi(clusterID[0])
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}

	clusterName, ok := request.URL.Query()["clusterName"]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Not found!")
		return
	}

	t, err := template.ParseFiles("html/trigger_must_gather.html")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "Error parsing template")
		return
	}
	dynData := types.Cluster{ID: id, Name: clusterName[0]}
	err = t.Execute(writer, dynData)
	if err != nil {
		println("Error executing template")
	}
}

// POST must-gather to REST API
func triggerMustGather(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	form := request.Form

	clusterID := form.Get("clusterid")
	clusterName := form.Get("clustername")
	username := form.Get("username")
	reason := form.Get("reason")
	link := form.Get("link")

	log.Println("clusterID", clusterID)
	log.Println("clusterName", clusterName)
	log.Println("username", username)
	log.Println("reason", reason)
	log.Println("link", link)

	query := "username=" + url.QueryEscape(username) + "&reason=" + url.QueryEscape(reason) + "&link=" + url.QueryEscape(link)
	log.Println(query)
	url := controllerURL + APIPrefix + "client/cluster/" + url.PathEscape(clusterName) + "/trigger/must-gather?" + query
	log.Println(url)

	err := performWriteRequest(url, "POST", nil)
	if err != nil {
		log.Println("Error communicating with the service", err)
		http.Redirect(writer, request, "/trigger-not-created", 301)
	} else {
		log.Println("Trigger has been created")
		http.Redirect(writer, request, "/trigger-created", 301)
	}
}

func startHTTPServer(address string) {
	http.HandleFunc("/", staticPage("html/index.html"))
	http.HandleFunc("/bootstrap.min.css", staticPage("html/bootstrap.min.css"))
	http.HandleFunc("/bootstrap.min.js", staticPage("html/bootstrap.min.js"))
	http.HandleFunc("/ccx.css", staticPage("html/ccx.css"))
	http.HandleFunc("/configuration-created", staticPage("html/configuration_created.html"))
	http.HandleFunc("/configuration-not-created", staticPage("html/configuration_not_created.html"))
	http.HandleFunc("/profile-created", staticPage("html/profile_created.html"))
	http.HandleFunc("/profile-not-created", staticPage("html/profile_not_created.html"))
	http.HandleFunc("/list-clusters", listClusters)
	http.HandleFunc("/list-profiles", listProfiles)
	http.HandleFunc("/list-configurations", listConfigurations)
	http.HandleFunc("/list-all-triggers", listTriggers)
	http.HandleFunc("/list-triggers", listTriggers)
	http.HandleFunc("/describe-configuration", describeConfiguration)
	http.HandleFunc("/new-profile", staticPage("html/new_profile.html"))
	http.HandleFunc("/new-configuration", staticPage("html/new_configuration.html"))
	http.HandleFunc("/store-profile", storeProfile)
	http.HandleFunc("/store-configuration", storeConfiguration)
	http.HandleFunc("/enable-configuration", enableConfiguration)
	http.HandleFunc("/disable-configuration", disableConfiguration)
	http.HandleFunc("/activate-trigger", activateTrigger)
	http.HandleFunc("/deactivate-trigger", deactivateTrigger)
	http.HandleFunc("/trigger-must-gather-configuration", triggerMustGatherConfiguration)
	http.HandleFunc("/trigger-must-gather", triggerMustGather)
	http.HandleFunc("/trigger-created", staticPage("html/trigger_created.html"))
	http.HandleFunc("/trigger-not-created", staticPage("html/trigger_not_created.html"))
	http.ListenAndServe(address, nil)
}

func main() {
	log.Println("Reading configuration")
	configFile, specified := os.LookupEnv("INSIGHTS_WEB_UI_CONFIG_FILE")
	if specified {
		// we need to separate the directory name and filename without extension
		directory, basename := filepath.Split(configFile)
		file := strings.TrimSuffix(basename, filepath.Ext(basename))
		// parse the configuration
		viper.SetConfigName(file)
		viper.AddConfigPath(directory)
	} else {
		// parse the configuration
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	controllerURL = viper.GetString("controller_url")
	address := viper.GetString("address")

	log.Println("Starting the service at address: " + address)
	startHTTPServer(address)
}
