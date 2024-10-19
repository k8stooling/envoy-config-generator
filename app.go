package main

import (
	_ "embed"
	"log"
	"os"
	"strings"
	"text/template"
)

//go:embed envoy_template.yaml
var envoyTemplate string

type HostInfo struct {
	Host string
	Port string
}

func main() {
	hostsStr := os.Getenv("ENVOY_HOSTS") // comma-separated "host1:port1,host2,host3:port3"
	if hostsStr == "" {
		log.Fatal("ENVOY_HOSTS environment variable is required")
	}

	hosts := parseHosts(hostsStr)

	tmpl, err := template.New("envoy").Parse(envoyTemplate)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	// Check if ENVOY_CONFIG_PATH is set
	configPath := os.Getenv("ENVOY_CONFIG_PATH")
	var output *os.File
	if configPath != "" {
		// Write to file if ENVOY_CONFIG_PATH is defined
		output, err = os.Create(configPath)
		if err != nil {
			log.Fatalf("Error creating config file: %v", err)
		}
		defer output.Close()
		log.Printf("Writing config to file: %s", configPath)
	} else {
		// Default to stdout
		output = os.Stdout
		log.Println("Writing config to stdout")
	}

	// Execute the template and write to the chosen output (file or stdout)
	err = tmpl.Execute(output, map[string]interface{}{
		"Hosts": hosts,
	})
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}

func parseHosts(hostsStr string) []HostInfo {
	var hosts []HostInfo
	for _, host := range strings.Split(hostsStr, ",") {
		parts := strings.Split(host, ":")
		hostName := parts[0]
		port := "443" // default port
		if len(parts) > 1 {
			port = parts[1]
		}
		hosts = append(hosts, HostInfo{Host: hostName, Port: port})
	}
	return hosts
}
