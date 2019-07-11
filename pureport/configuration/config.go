package configuration

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/hashicorp/terraform/httpclient"
	"github.com/pureport/pureport-sdk-go/pureport"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	ppLog "github.com/pureport/pureport-sdk-go/pureport/logging"
	"github.com/pureport/pureport-sdk-go/pureport/session"
	"github.com/pureport/terraform-provider-pureport/version"
)

var (
	logMutex sync.Mutex
)

type Config struct {
	Session *session.Session

	APIKey                string
	APISecret             string
	AuthenticationProfile string
	EndPoint              string
}

func (c *Config) LoadAndValidate() error {

	// Lock the configuration while we update
	logMutex.Lock()
	defer logMutex.Unlock()

	// Validate that if the API Key was specified that a secret was specified as well.
	if (c.APIKey == "") != (c.APISecret == "") {
		return fmt.Errorf("API Key and Secret both need to be specified for successful authentication.")
	}

	cfg := pureport.NewConfiguration()
	cfg.APIKey = c.APIKey
	cfg.APISecret = c.APISecret
	cfg.AuthenticationProfile = c.AuthenticationProfile
	cfg.EndPoint = c.EndPoint

	logCfg := ppLog.NewLogConfig()

	// Map Terrform Log Levels to our SDK Levels
	switch os.Getenv("TF_LOG") {
	case "TRACE":
		logCfg.Level = "DEBUG"
	case "DEBUG":
		logCfg.Level = "DEBUG"
	case "INFO":
		logCfg.Level = "INFO"
	case "WARN":
		logCfg.Level = "WARNING"
	case "ERROR":
		logCfg.Level = "ERROR"
	default:
		logCfg.Level = "WARNING"
	}

	ppLog.SetupLogger(logCfg)

	terraformVersion := httpclient.UserAgentString()
	providerVersion := fmt.Sprintf("terraform-provider-pureport/%s", version.ProviderVersion)
	terraformWebsite := "(+https://www.pureport.com)"

	cfg.UserAgent = fmt.Sprintf("%s %s %s", terraformVersion, terraformWebsite, providerVersion)
	c.Session = session.NewSession(cfg)

	return nil
}

func (c *Config) getAccounts() ([]client.Account, error) {

	ctx := c.Session.GetSessionContext()

	accounts, resp, err := c.Session.Client.AccountsApi.FindAllAccounts(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("Error when Reading Pureport Account data: %v", err)
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Error Response while Reading Pureport Account data")
	}

	// Filter the results
	var filteredAccounts []client.Account
	for _, account := range accounts {
		if _, ok := account.Tags["test-acc"]; ok {
			filteredAccounts = append(filteredAccounts, account)
		}
	}

	// Sort the list
	sort.Slice(filteredAccounts, func(i int, j int) bool {
		return filteredAccounts[i].Name < filteredAccounts[j].Name
	})

	return filteredAccounts, nil
}

func (c *Config) GetAccNetworks() ([]client.Network, error) {

	accounts, err := c.getAccounts()
	if err != nil {
		return nil, fmt.Errorf("Error reading account information: %v", err)
	}

	// Filter the results
	var filteredNetworks []client.Network
	ctx := c.Session.GetSessionContext()

	for _, account := range accounts {

		networks, resp, err := c.Session.Client.NetworksApi.FindNetworks(ctx, account.Id)
		if err != nil {
			return nil, fmt.Errorf("Error when Reading Network data: %v", err)
		}

		if resp.StatusCode >= 300 {
			return nil, fmt.Errorf("Error Response while Reading Network data")
		}

		for _, network := range networks {
			if _, ok := network.Tags["test-acc"]; ok {
				filteredNetworks = append(filteredNetworks, network)
			}
		}
	}

	// Sort the list
	sort.Slice(filteredNetworks, func(i int, j int) bool {
		return filteredNetworks[i].Name < filteredNetworks[j].Name
	})

	return filteredNetworks, nil
}

func (c *Config) GetAccConnections() ([]client.Connection, error) {

	networks, err := c.GetAccNetworks()
	if err != nil {
		return nil, fmt.Errorf("Error reading Networks: %v", err)
	}

	// Filter the results
	var filteredConnections []client.Connection
	ctx := c.Session.GetSessionContext()

	for _, network := range networks {

		connections, resp, err := c.Session.Client.ConnectionsApi.GetConnections(ctx, network.Id)
		if err != nil {
			return nil, fmt.Errorf("Error when Reading Connections data: %v", err)
		}

		if resp.StatusCode >= 300 {
			return nil, fmt.Errorf("Error Response while Reading Connections data")
		}

		for _, connection := range connections {
			if _, ok := connection.Tags["test-acc"]; ok {
				filteredConnections = append(filteredConnections, connection)
			}
		}
	}

	// Sort the list
	sort.Slice(filteredConnections, func(i int, j int) bool {
		return filteredConnections[i].Name < filteredConnections[j].Name
	})

	return filteredConnections, nil
}

func (c *Config) SweepNetworks(networks []client.Network) error {

	ctx := c.Session.GetSessionContext()

	for _, network := range networks {
		if _, ok := network.Tags["sweep"]; ok {
			resp, err := c.Session.Client.NetworksApi.DeleteNetwork(ctx, network.Id)
			if err != nil {
				return fmt.Errorf("Error when Deleting Network: %v", err)
			}

			if resp.StatusCode >= 300 {
				return fmt.Errorf("Error Response while Deleting Network : id=%s", network.Id)
			}
		}
	}

	return nil
}

func (c *Config) SweepConnections(connections []client.Connection) error {

	ctx := c.Session.GetSessionContext()

	for _, connection := range connections {
		if _, ok := connection.Tags["sweep"]; ok {
			_, resp, err := c.Session.Client.ConnectionsApi.DeleteConnection(ctx, connection.Id)
			if err != nil {
				return fmt.Errorf("Error when Deleting Connection: %v", err)
			}

			if resp.StatusCode >= 300 {
				return fmt.Errorf("Error Response while Deleting Connection: id=%s", connection.Id)
			}
		}
	}

	return nil
}
