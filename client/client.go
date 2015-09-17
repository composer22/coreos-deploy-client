package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/composer22/coreos-deploy/server"
)

// Client represents an instance of a connection to the server.
type Client struct {
	opts *Options
}

// New is a factory function that returns a new client instance.
func New(o *Options) *Client {
	return &Client{
		opts: o,
	}
}

// PrintVersionAndExit prints the version of the server then exits.
func PrintVersionAndExit() {
	fmt.Printf("coreos-deploy-client version %s\n", version)
	os.Exit(0)
}

// Execute communicates to the server with an action.
func (c *Client) Execute() (string, error) {
	if c.opts.DeployID == "" {
		result, err := c.deploy()
		return result, err
	}
	result, err := c.getStatus()
	return result, err
}

// Text template variables to substitute in the .service template.
type ServiceTemplateVars struct {
	Name         string `json:"name"`         // The name of the service to deploy.
	Version      string `json:"version"`      // The version of the service.
	ImageVersion string `json:"imageVersion"` // The version of the docker image.
	NumInstances int    `json:"numInstances"` // The number of instances to deploy.
}

// deploy submits a deploy request to the server.
func (c *Client) deploy() (string, error) {
	// Read in template file.
	tf, err := ioutil.ReadFile(c.opts.TemplateFilePath)
	if err != nil {
		return "", err
	}

	// Fill any template variables.
	stv := &ServiceTemplateVars{
		Name:         c.opts.Name,
		Version:      c.opts.Version,
		ImageVersion: c.opts.ImageVersion,
		NumInstances: c.opts.NumInstances,
	}
	t, err := template.New("service template").Parse(string(tf[:]))
	if err != nil {
		return "", err
	}

	var tb bytes.Buffer
	err = t.Execute(&tb, stv)
	if err != nil {
		return "", err
	}
	tmpl := tb.String()

	// Read in etcd2 key/values.
	keys := make(map[string]string)
	if c.opts.Etcd2FilePath != "" {
		file, err := os.Open(c.opts.Etcd2FilePath)
		if err != nil {
			return "", err
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			words := strings.Fields(line)
			if len(words) >= 2 {
				value := strings.TrimSpace(strings.Replace(line, words[0], "", 1))
				keys[words[0]] = value
			}
		}

		if err := scanner.Err(); err != nil {
			file.Close()
			return "", err
		}
		file.Close()
	}

	// Create the payload.
	payload, err := json.Marshal(server.NewServiceRequest(c.opts.Name, c.opts.Version, c.opts.NumInstances, tmpl, keys))
	if err != nil {
		return "", err
	}

	// Send the request.
	req, err := http.NewRequest(httpPost, fmt.Sprintf("%s%s", c.opts.Url, httpRouteV1Deploy), bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", err
	}
	return c.sendRequest(req)
}

// getStatus prints out the status of a previous deploy.
func (c *Client) getStatus() (string, error) {
	req, err := http.NewRequest(httpGet, fmt.Sprintf("%s%s%s", c.opts.Url, httpRouteV1Status, c.opts.DeployID), nil)
	if err != nil {
		return "", err
	}
	return c.sendRequest(req)
}

// sendRequest sends a request to the server and prints the result.
func (c *Client) sendRequest(req *http.Request) (string, error) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.opts.Token))
	cl := &http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
