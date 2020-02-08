package main

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	ExtensionName = "terraless-extension-aws-organisation-roles"
	VERSION       = "0.1.0"
)

type ExtensionAwsOrganisationRoles struct {
	logger hclog.Logger
}

func (extension *ExtensionAwsOrganisationRoles) Info() schema.PluginInfo {
	return schema.PluginInfo{
		Name:    ExtensionName,
		Version: VERSION,
	}
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "extension-plugin",
	MagicCookieValue: "terraless",
}

var logger hclog.Logger

func main() {
	logger = hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	logrus.Info("Running Plugin: Terraless Module AWS Organisation Roles")
	extension := &ExtensionAwsOrganisationRoles{
		logger: logger,
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"Extension": &schema.ExtensionPlugin{Impl: extension},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
