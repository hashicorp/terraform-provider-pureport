package pureport

import (
	"fmt"
	"os"
	"sync"

	"github.com/hashicorp/terraform/httpclient"
	"github.com/pureport/pureport-sdk-go/pureport"
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
