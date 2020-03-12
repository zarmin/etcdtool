package command

import (
	"os"
	"os/user"
	"time"

	"github.com/codegangsta/cli"
	"github.com/mickep76/iodatafmt"
)

// Etcdtool configuration struct.
type Etcdtool struct {
	Peers            string        `json:"peers,omitempty" yaml:"peers,omitempty" toml:"peers,omitempty"`
	Cert             string        `json:"cert,omitempty" yaml:"cert,omitempty" toml:"cert,omitempty"`
	Key              string        `json:"key,omitempty" yaml:"key,omitempty" toml:"key,omitempty"`
	CA               string        `json:"ca,omitempty" yaml:"ca,omitempty" toml:"peers,omitempty"`
	User             string        `json:"user,omitempty" yaml:"user,omitempty" toml:"user,omitempty"`
	Timeout          time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty" toml:"timeout,omitempty"`
	CommandTimeout   time.Duration `json:"commandTimeout,omitempty" yaml:"commandTimeout,omitempty" toml:"commandTimeout,omitempty"`
	Routes           []Route       `json:"routes" yaml:"routes" toml:"routes"`
	PasswordFilePath string
}

// Route configuration struct.
type Route struct {
	Regexp string `json:"regexp" yaml:"regexp" toml:"regexp"`
	Schema string `json:"schema" yaml:"schema" toml:"schema"`
}

func loadConfig(c *cli.Context) Etcdtool {
	// Enable debug
	if c.Bool("debug") {
		debug = true
	}

	// Default path for config file.
	u, _ := user.Current()
	cfgs := []string{
		"/etcd/etcdtool.json",
		"/etcd/etcdtool.yaml",
		"/etcd/etcdtool.toml",
		u.HomeDir + "/.etcdtool.json",
		u.HomeDir + "/.etcdtool.yaml",
		u.HomeDir + "/.etcdtool.toml",
	}

	// Check if we have an arg. for config file and that it exist's.
	if c.String("config") != "" {
		if _, err := os.Stat(c.String("config")); os.IsNotExist(err) {
			fatalf("Config file doesn't exist: %s", c.String("config"))
		}
		cfgs = append([]string{c.String("config")}, cfgs...)
	}

	// Check if config file exists and load it.
	e := Etcdtool{}
	for _, fn := range cfgs {
		if _, err := os.Stat(fn); os.IsNotExist(err) {
			continue
		}
		infof("Using config file: %s", fn)
		f, err := iodatafmt.FileFormat(fn)
		if err != nil {
			fatal(err.Error())
		}
		if err := iodatafmt.LoadPtr(&e, fn, f); err != nil {
			fatal(err.Error())
		}
	}

	// Override with arguments or env. variables.
	if c.String("peers") != "" {
		e.Peers = c.String("peers")
	}

	if c.String("cert") != "" {
		e.Cert = c.String("cert")
	}

	if c.String("key") != "" {
		e.Key = c.String("key")
	}

	if c.String("ca") != "" {
		e.CA = c.String("ca")
	}

	if c.String("user") != "" {
		e.User = c.String("user")
	}

	if c.Duration("timeout") != 0 {
		e.Timeout = c.Duration("timeout")
	}

	if c.Duration("command-timeout") != 0 {
		e.CommandTimeout = c.Duration("command-timeout")
	}

	// Add password file path if set
	if c.String("password-file") != "" {
		e.PasswordFilePath = c.String("password-file")
	}

	return e
}
