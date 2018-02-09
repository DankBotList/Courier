package courier

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"errors"

	"path/filepath"

	"github.com/satori/go.uuid"
)

// Data abstract containing methods for saving and loading.

// Data an abstract struct used for it's functions to save and load config files.
type Data struct{}

func (d *Data) save(saveLoc string, inter interface{}) error {
	// Make all the directories
	if err := os.MkdirAll(filepath.Dir(saveLoc), os.ModeDir|0775); err != nil {
		return err
	}

	data, err := json.MarshalIndent(inter, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(saveLoc, data, 0660)
}

func (d *Data) load(saveLoc string, inter interface{}) error {

	if _, err := os.Stat(saveLoc); os.IsNotExist(err) {
		return DefaultConfigSavedError
	} else if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(saveLoc)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, inter); err != nil {
		return err
	}

	return nil

}

// ========== Main configuration.

// ConfigSaveLocation the location to save the config to.
var ConfigSaveLocation = "courier/config.json"

// DefaultConfigSavedError an error returned if the default config is saved.
var DefaultConfigSavedError = errors.New("the default config has been saved, please edit it")

// DefaultConfig the default configuration to save.
var DefaultConfig = Config{
	Data:              Data{},
	Repo:              "gogs@git.cory.red:DankBotList/Site.git",
	Ref:               "refs/heads/master",
	PollTimeSeconds:   20,
	WebSocketPath:     "/courier/ws",
	AuthenticationKey: uuid.Must(uuid.NewV4()).String() + "-" + uuid.Must(uuid.NewV4()).String(),
	DataFolder:        "courier/data/",
}

// Config the main configuration.
type Config struct {
	Data              `json:"-"`
	Repo              string `json:"repo"`
	Ref               string `json:"ref"`
	PollTimeSeconds   int    `json:"poll_time_seconds"`
	WebSocketPath     string `json:"web_socket_path"`
	AuthenticationKey string `json:"authentication_key"`
	DataFolder        string `json:"data_folder"`
}

// Save saves the config.
func (c *Config) Save() error {
	saveLoc, envThere := os.LookupEnv("COURIER_CONFIG")
	if !envThere {
		saveLoc = ConfigSaveLocation
	}

	return c.save(saveLoc, c)
}

// Load loads the config.
func (c *Config) Load() error {

	saveLoc, envThere := os.LookupEnv("COURIER_CONFIG")
	if !envThere {
		saveLoc = ConfigSaveLocation
	}

	if err := c.load(saveLoc, c); err == DefaultConfigSavedError {
		if err := DefaultConfig.Save(); err != nil {
			return err
		}
		return DefaultConfigSavedError
	} else if err != nil {
		return err
	}

	return nil

}

// ========== Hosts.
type Host struct {
}

// HostsData the data for hosts, their locations and ID's.
type HostsData struct {
	Data     `json:"-"`
	DataFile string `json:"-"`
	Hosts    []Host `json:"hosts"`
}

func (c *HostsData) SetDataFile(dataFile string) {
	c.DataFile = dataFile
}

// Save saves the config.
func (c *HostsData) Save() error {
	if c.DataFile == "" {
		return c.save("courier/data/hosts.json", c)
	} else {
		return c.save(c.DataFile, c)
	}
}

// Load loads the config.
func (c *HostsData) Load() error {
	saveLoc := "courier/data/hosts.json"
	if c.DataFile != "" {
		saveLoc = c.DataFile
	}

	if err := c.load(saveLoc, c); err == DefaultConfigSavedError {
		if err := DefaultConfig.Save(); err != nil {
			return err
		}
		return DefaultConfigSavedError
	} else if err != nil {
		return err
	}

	return nil
}
