package client

import (
	"fmt"
	"os"
)

const usageText = `
Description: coreos-deploy-client is a CLI for deploying services to the coreos-deploy API.

Usage: coreos-deploy-client [options...]

Server options:
    -n, --name NAME                  NAME of the service (mandatory).
    -r, --service_version VERSION    VERSION of the service (mandatory).
    -n, --instances INSTANCES        Number of INSTANCES to deploy. (default: 2).
    -t, --template_filepath TEMPLATE Path and filename to the unit .service TEMPLATE (mandatory).
    -e, --etcd2_filepath ETCD2FILE   Path and filename to the etcd2 key/value ETCD2FILE.
    -b, --bearer_token TOKEN         The API authorization TOKEN for the server.
    -u, --url URL                    URL of the coreos-deploy server.

    -p, --deploy_id ID               Lookup the status of a previous deployment.

    -d, --debug                      Enable debugging output (default: false)

Common options:
    -h, --help                       Show this message
    -V, --version                    Show version

Examples:

   # Deploy a service and return a deploy ID...
    coreos-deploy-client -n my-application -r 1.0.1 -n 2 -t /path/to/my-application@.service \
	  -e /path/to/my-application.etcd2 -b AP1T0K3N -u http://coreos-dev.example.com:80

	# Check the status of a recent deploy...
	coreos-deploy-client -b AP1T0K3N -u http://coreos-dev.example.com:80 -p DC8D9C2E-8161-4FC0-937F-4CA7037970D5
`

// PrintUsageAndExit is used to print out command line options.
func PrintUsageAndExit() {
	fmt.Printf("%s\n", usageText)
	os.Exit(0)
}
