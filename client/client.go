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
func (c *Client) Execute() {
	if c.opts.DeployID == "" {
		c.deploy()
		return
	}
	c.getStatus()
}

// deploy submits a deploy request to the server.
func (c *Client) deploy() {
	// Read in template file.
	tf, err := ioutil.ReadFile(c.opts.TemplateFilePath)
	if err != nil {
		PrintErr(err.Error())
		return
	}
	tmpl := string(tf[:])

	// Read in etcd2 key/values.
	keys := make(map[string]string)
	if c.opts.Etcd2FilePath != "" {
		file, err := os.Open(c.opts.Etcd2FilePath)
		if err != nil {
			PrintErr(err.Error())
			return
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
			PrintErr(err.Error())
			return
		}
		file.Close()
	}

	// Create the payload.
	payload, err := json.Marshal(server.NewServiceRequest(c.opts.Name, c.opts.Version, c.opts.NumInstances, tmpl, keys))
	if err != nil {
		PrintErr(err.Error())
		return
	}

	// Send the request.
	req, err := http.NewRequest(httpPost, fmt.Sprintf("%s%s", c.opts.Url, httpRouteV1Deploy), bytes.NewBuffer([]byte(payload)))
	if err != nil {
		PrintErr(err.Error())
		return
	}
	c.sendRequest(req)
}

// getStatus prints out the status of a previous deploy.
func (c *Client) getStatus() {
	req, err := http.NewRequest(httpGet, fmt.Sprintf("%s%s%s", c.opts.Url, httpRouteV1Status, c.opts.DeployID), nil)
	if err != nil {
		PrintErr(err.Error())
		return
	}
	c.sendRequest(req)
}

// sendRequest sends a request to the server and prints the result.
func (c *Client) sendRequest(req *http.Request) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.opts.Token))
	cl := &http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		PrintErr(err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		PrintErr(err.Error())
		return
	}
	fmt.Println(string(body))
}
