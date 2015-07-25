package configuration

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Bowery/prompt"
)

var (
	homeDir               = os.Getenv("HOME")
	configurationFilePath = homeDir + "/.gost"
)

// ErrNoConfigFound is thrown when no configuration file is found
var ErrNoConfigFound = errors.New("no config file found")

const (
	authURL  = "https://api.github.com/authorizations"
	authJSON = `{"scopes": "gist", "note": "gost cli"}`
)

// Configuration holds Github the config
type Configuration struct {
	Username string
	Token    string
	Private  bool
}

type authResponse struct {
	Token string `json:"token"`
}

// NewConfiguration prompts the user for the gitub credentials and creates a
// configuration with an access token
func NewConfiguration() (*Configuration, error) {
	u, err := prompt.BasicDefault("Username: ", os.Getenv("USER"))
	if err != nil {
		return nil, err
	}
	pw, err := prompt.Password("Password: ")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", authURL, bytes.NewBufferString(authJSON))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(u, pw)
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	if res.StatusCode > 400 {
		return nil, fmt.Errorf("could not get token, code: %d, body: %s", res.StatusCode, body)
	}

	r := new(authResponse)
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	c := &Configuration{
		Username: u,
		Private:  true,
		Token:    r.Token,
	}
	return c, nil
}

//LoadConfiguration load the configuration from a file
func LoadConfiguration() (*Configuration, error) {
	c := new(Configuration)
	if _, err := os.Stat(configurationFilePath); os.IsNotExist(err) {
		return c, ErrNoConfigFound
	}
	file, err := ioutil.ReadFile(configurationFilePath)
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(file, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}

//Save saves the configuration to a file
func (c *Configuration) Save() error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configurationFilePath, b, 0600)
	if err != nil {
		return err
	}
	return nil
}
