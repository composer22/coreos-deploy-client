// coreos-deploy-client is a CLI interface for submitting deploys to a CoreOS cluster running coreos-deploy.
package main

import (
	"flag"
	"strings"

	"github.com/composer22/coreos-deploy-client/client"
)

// main is the main entry point for the client.
func main() {
	opts := &client.Options{}
	var showVersion bool

	flag.StringVar(&opts.Name, "n", "", "Name of the service.")
	flag.StringVar(&opts.Name, "name", "", "Name of the service.")
	flag.StringVar(&opts.Version, "r", "", "Version of the service.")
	flag.StringVar(&opts.Version, "service_version", "", "Version of the service.")
	flag.StringVar(&opts.ImageVersion, "k", "latest", "Version of the docker image.")
	flag.StringVar(&opts.ImageVersion, "image_version", "latest", "Version of the docker image.")
	flag.IntVar(&opts.NumInstances, "i", 2, "Number of instances to instantiate.")
	flag.IntVar(&opts.NumInstances, "instances", 2, "Number of instances to instantiate.")
	flag.StringVar(&opts.TemplateFilePath, "t", "", "path/file to the unit .service file.")
	flag.StringVar(&opts.TemplateFilePath, "template_filepath", "", "path/file to the unit .service file.")
	flag.StringVar(&opts.Etcd2FilePath, "e", "", "path/file to the etcd2 key/value file.")
	flag.StringVar(&opts.Etcd2FilePath, "etcd2_filepath", "", "path/file to the etcd2 key/value file.")

	flag.StringVar(&opts.Token, "b", "", "API authorization token")
	flag.StringVar(&opts.Token, "bearer_token", "", "API authorization token.")
	flag.StringVar(&opts.Url, "u", "", "URL of the load balancer to the coreos-deploy server.")
	flag.StringVar(&opts.Url, "url", "", "URL of the load balancer to the coreos-deploy server.")

	flag.StringVar(&opts.DeployID, "p", "", "Deploy ID to lookup.")
	flag.StringVar(&opts.DeployID, "deploy_id", "", "Deploy ID to lookup.")

	flag.BoolVar(&opts.Debug, "d", false, "Enable debugging output.")
	flag.BoolVar(&opts.Debug, "debug", false, "Enable debugging output.")
	flag.BoolVar(&showVersion, "V", false, "Show version.")
	flag.BoolVar(&showVersion, "version", false, "Show version.")
	flag.Usage = client.PrintUsageAndExit
	flag.Parse()

	// Version flag request?
	if showVersion {
		client.PrintVersionAndExit()
	}

	// Check additional params beyond the flags.
	for _, arg := range flag.Args() {
		switch strings.ToLower(arg) {
		case "version":
			client.PrintVersionAndExit()
		case "help":
			client.PrintUsageAndExit()
		}
	}

	// Validate the mandatory options.
	if err := opts.Validate(); err != nil {
		client.PrintErr(err.Error())
		return
	}

	// Run the deploy.
	c := client.New(opts)
	c.Execute()
}
