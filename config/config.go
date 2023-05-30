package config

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/onflow/flow-go/config/network"
)

const configFileName = "config.yml"

var (
	conf = viper.New()

	//go:embed config.yml
	configFile embed.FS
)

// FlowConfig Flow configuration.
type FlowConfig struct {
	NetworkConfig *network.Config `mapstructure:"network-config"`
}

// DefaultConfig initializes the flow configuration. All default values for the Flow
// configuration are stored in the config.yml file. These values can be overriden
// by node operators by setting the corresponding cli flag. DefaultConfig should be called
// before any pflags are parsed, this will allow the configuration to initialize with defaults
// from config.yml.
// Returns:
//
//	*FlowConfig: an instance of the network configuration fully initialized to the default values set in the config file
//	error: if there is any error encountered while initializing the configuration, all errors are considered irrecoverable.
func DefaultConfig() (*FlowConfig, error) {
	var flowConf FlowConfig
	err := unmarshallFlowConfig(&flowConf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall the Flow config: %w", err)
	}

	return &flowConf, nil
}

// BindPFlags binds the configuration to the cli pflag set. This should be called
// after all pflags have been parsed.
// Args:
//
//	*FlowConfig: The Flow configuration that will be used to unmarshall the configuration values into after binding pflags.
//	This needs to be done because pflags may override a configuration value.
//	error: if there is any error encountered binding pflags or unmarshalling the config struct, all errors are considered irrecoverable.
//
// Note: As configuration management is improved this func should accept the entire Flow config as the arg to unmarshall new config values into.
func BindPFlags(c *FlowConfig) error {
	if err := conf.BindPFlags(pflag.CommandLine); err != nil {
		return fmt.Errorf("failed to bind pflags: %w", err)
	}

	err := unmarshallFlowConfig(c)
	if err != nil {
		return fmt.Errorf("failed to unmarshall the Flow config: %w", err)
	}

	return nil
}

func unmarshallFlowConfig(c *FlowConfig) error {
	err := conf.Unmarshal(c)
	if err != nil {
		return fmt.Errorf("failed to unmarshal network config: %w", err)
	}

	return nil
}

func init() {
	f, err := configFile.Open(configFileName)
	if err != nil {
		panic(fmt.Errorf("failed to open config.yml: %w", err))
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
	if err != nil {
		panic(fmt.Errorf("failed to read config file into bytes buffer: %w", err))
	}

	conf.SetConfigType("yaml")

	if err = conf.ReadConfig(buf); err != nil {
		panic(fmt.Errorf("failed to initialize flow config failed to read in config file: %w", err))
	}
}
