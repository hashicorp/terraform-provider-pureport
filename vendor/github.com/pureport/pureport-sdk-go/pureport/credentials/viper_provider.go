package credentials

import (
	"fmt"

	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

const (
	// ViperProviderName is the name of this provider
	ViperProviderName = "ViperProvider"

	configAPIKeyStr     = "api_key"
	configAPISecretStr  = "api_secret"
	configAPIProfileStr = "profile"
)

var (
	vip *viper.Viper
	log = logging.MustGetLogger("main_logger")
)

// ViperProvider is a credentials provider using Viper as the configuration source.
type ViperProvider struct {
	retrieved bool

	// Filename - path to the configuration file (Optional)
	Filename string

	// Profile - the profile to use from the configuration file (Optional)
	Profile string
}

func init() {
	vip = viper.New()
	vip.SetConfigName("credentials")
	vip.AddConfigPath("$HOME/.pureport")
	vip.AddConfigPath(".")
	vip.SetEnvPrefix("pureport")
	vip.BindEnv(configAPIKeyStr)
	vip.BindEnv(configAPISecretStr)
	vip.BindEnv(configAPIProfileStr)
}

// NewViperCredentials creates a new credentials provider using Viper
func NewViperCredentials(profile string) *Credentials {
	return NewCredentials(&ViperProvider{
		retrieved: false,
		Profile:   profile,
	})
}

// Retrieve Provider.Retrieve()
func (p *ViperProvider) Retrieve() (Value, error) {
	p.retrieved = false

	if err := vip.ReadInConfig(); err != nil {
		log.Debugf("Configuration file not available: %s", err)
	}

	// Check environment first
	key := vip.GetString(configAPIKeyStr)
	secret := vip.GetString(configAPISecretStr)

	if profile := vip.GetString(configAPIProfileStr); profile != "" {
		p.Profile = profile
	}

	if p.Profile == "" {
		p.Profile = "default"
	}

	// Read from configuration file
	if key == "" || secret == "" {
		key = vip.GetString(fmt.Sprintf("profiles.%s.%s", p.Profile, configAPIKeyStr))
		secret = vip.GetString(fmt.Sprintf("profiles.%s.%s", p.Profile, configAPISecretStr))
	}

	if key == "" || secret == "" {
		return Value{ProviderName: ViperProviderName}, fmt.Errorf("API Key and/or Secret not found")
	}

	p.retrieved = true
	return Value{
		ProviderName: ViperProviderName,
		APIKey:       key,
		Secret:       secret,
	}, nil
}

// IsExpired Provider.IsExpired()
func (p *ViperProvider) IsExpired() bool {
	return !p.retrieved
}
