package apiserver

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"strings"
)

type IBuild interface {
	Build() *Config
}

type Config struct {
	host, port string
}

func (c *Config) GetAddr() string {
	return c.host + ":" + c.port
}

type ConfigBuilder struct {
	config *Config
}

func (c *ConfigBuilder) setHostPort(host, port string) error {
	if host == "" || port == "" {
		return errors.New("Error encode file")
	}

	c.config.host, c.config.port = host, port
	return nil
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: &Config{},
	}
}

func (c *ConfigBuilder) Parse(configPath string) (IBuild, error) {
	var f func([]byte) (IBuild, error)

	switch true {
	case strings.HasSuffix(configPath, ".xml"):
		f = c.fromXML
	case strings.HasSuffix(configPath, ".yml"):
		f = c.fromYML
	case strings.HasSuffix(configPath, ".json"):
		f = c.fromJSON
	default:
		return nil, errors.New("Invalid type file: " + configPath)
	}

	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	return f(bytes)
}

//----------------------------------------------------------------------------//

// ConfigJSON JSON configure
type ConfigJSON struct {
	ConfigBuilder
}

func (c *ConfigBuilder) fromJSON(b []byte) (IBuild, error) {
	temp := struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}{}

	err := json.Unmarshal(b, &temp)
	if err != nil {
		return nil, err
	}
	err = c.setHostPort(temp.Host, temp.Port)
	return &ConfigJSON{*c}, err
}

func (c *ConfigJSON) Build() *Config {
	return c.config
}

//----------------------------------------------------------------------------//

// ConfigXML XML Configure
type ConfigXML struct {
	ConfigBuilder
}

func (c *ConfigBuilder) fromXML(b []byte) (IBuild, error) {
	temp := struct {
		XMLName xml.Name `xml:"config"`
		Host    string   `xml:"host"`
		Port    string   `xml:"port"`
	}{}

	err := xml.Unmarshal(b, &temp)
	if err != nil {
		return nil, err
	}

	err = c.setHostPort(temp.Host, temp.Port)
	return &ConfigXML{*c}, err
}

func (c *ConfigXML) Build() *Config {
	return c.config
}

//----------------------------------------------------------------------------//

// ConfigYML YML configure
type ConfigYML struct {
	ConfigBuilder
}

func (c *ConfigBuilder) fromYML(b []byte) (IBuild, error) {
	return nil, nil
}
