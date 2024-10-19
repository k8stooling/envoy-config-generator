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

	err = tmpl.Execute(os.Stdout, map[string]interface{}{
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
