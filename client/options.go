package client

import (
	"encoding/json"
	"errors"
)

// Options represents parameters that are passed to the application.
type Options struct {
	Name             string `json:"name"`             // The name of the service to deploy.
	Version          string `json:"version"`          // The version of the service.
	NumInstances     int    `json:"numInstances"`     // The number of instances to deploy.
	TemplateFilePath string `json:"templateFilePath"` // The path to the unit .service file source file.
	Etcd2FilePath    string `json:"etcd2FilePath"`    // The path to the etcd2 key/value source file.
	Token            string `json:"bearerToken"`      // The API authorization token to the server.
	Url              string `json:"URL"`              // The URL to the coreos-deploy endpoint.
	DeployID         string `json:"deployID"`         // The Deploy ID to query.
	Debug            bool   `json:"debugEnabled"`     // Is debugging enabled in the application?
}

// Validate options
func (o *Options) Validate() error {
	if o.DeployID == "" {
		if o.Name == "" {
			return errors.New("Service Name is mandatory.")
		}
		if o.Version == "" {
			return errors.New("Service Version is mandatory.")
		}
		if o.TemplateFilePath == "" {
			return errors.New("Unit template file is mandatory.")
		}
	}

	if o.Token == "" {
		return errors.New("Authorization token is mandatory.")
	}
	if o.Url == "" {
		return errors.New("URL is mandatory.")
	}

	return nil
}

// String is an implentation of the Stringer interface so the structure is returned as a string
// to fmt.Print() etc.
func (o *Options) String() string {
	b, _ := json.Marshal(o)
	return string(b)
}
