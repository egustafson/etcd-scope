package main

import ()

const (
	ConfigBasename  = "etcd-scope"
	ConfigExtension = "yml"

	DefaultEndpoint = "localhost:2379"
)

// Configuration is the structure of the configuration.
type Configuration struct {
	Endpoint string `yaml:"endpoint"`
}

// Config is the default configuration
var Config = Configuration{
	Endpoint: DefaultEndpoint,
}

func initConfig() error {
	Config.Endpoint = endpointFlag // just set config from the flag (or default)
	return nil
}

// ignoreInitConfig is WIP that is not necessary until more
// configuration becomes necessary
func ignoreInitConfig() error {
	cfgPath, err := locateConfigFile()
	if err != nil {
		return err
	}
	if len(cfgPath) > 0 {
		log.Debug().Str("path", cfgPath).Msg("config file selected")
		//
		// TODO: load config file
		//
	} else {
		log.Debug().Msg("no config file located, using defaults")
	}

	// TODO:  implement ENV VAR substitution

	// TODO:  implement config validation

	return nil
}

func locateConfigFile() (path string, err error) {

	// TODO - implement as soon as config other than etcd-endpoint is needed
	return "", nil
}
