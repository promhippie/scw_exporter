package config

import (
	"time"
)

// Server defines the general server configuration.
type Server struct {
	Addr    string
	Path    string
	Timeout time.Duration
}

// Logs defines the level and color for log configuration.
type Logs struct {
	Level  string
	Pretty bool
}

// Target defines the target specific configuration.
type Target struct {
	AccessKey string
	SecretKey string
	Org       string
	Project   string
	Region    string
	Zone      string
	Timeout   time.Duration
}

// Collector defines the collector specific configuration.
type Collector struct {
	Dashboard      bool
	SecurityGroups bool
	Servers        bool
	Baremetal      bool
	Snapshots      bool
	Volumes        bool
}

// Config is a combination of all available configurations.
type Config struct {
	Server    Server
	Logs      Logs
	Target    Target
	Collector Collector
}

// Load initializes a default configuration struct.
func Load() *Config {
	return &Config{}
}
