package pureport

import (
	"fmt"

	"github.com/hashicorp/terraform/httpclient"
	"github.com/pureport/pureport-sdk-go/pureport"
	ppLog "github.com/pureport/pureport-sdk-go/pureport/logging"
	"github.com/pureport/pureport-sdk-go/pureport/session"
	"github.com/pureport/terraform-provider-pureport/version"
)

type Config struct {
	Session *session.Session

	userAgent string
}

func (c *Config) LoadAndValidate() error {

	cfg := pureport.NewConfiguration()

	logCfg := ppLog.NewLogConfig()

	ppLog.SetupLogger(logCfg)

	terraformVersion := httpclient.UserAgentString()
	providerVersion := fmt.Sprintf("terraform-provider-pureport/%s", version.ProviderVersion)
	terraformWebsite := "(+https://www.pureport.com)"
	userAgent := fmt.Sprintf("%s %s %s", terraformVersion, terraformWebsite, providerVersion)

	cfg.UserAgent = userAgent
	c.Session = session.NewSession(cfg)

	return nil
}
