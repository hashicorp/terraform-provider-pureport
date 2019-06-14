package pureport

import (
	"strings"
	"time"

	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var (
	vip *viper.Viper
	log = logging.MustGetLogger("main_logger")
)

const (
	endpointEnvStr = "endpoint"
)

func init() {
	vip = viper.New()
	vip.SetConfigName("configuration")
	vip.AddConfigPath("$HOME/.pureport")
	vip.AddConfigPath(".")

	vip.SetEnvPrefix("pureport")
	vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	_ = vip.BindEnv(endpointEnvStr)

	vip.SetDefault(endpointEnvStr, "https://api.pureport.com")
}

// Configuration - Application level configuration data
type Configuration struct {

	// Timeout used for any HTTP Requests made
	Timeout time.Duration

	// The endpoint to request Authentication Tokens
	EndPoint string

	// UserAgent to use for the http client
	UserAgent string

	// Pureport API Key
	APIKey string

	// Pureport Secret Access Key
	APISecret string

	// Pureport Authentication Profile
	AuthenticationProfile string
}

// NewConfiguration creates a new configuration to application
// level configuration data.
func NewConfiguration() *Configuration {

	if err := vip.ReadInConfig(); err != nil {
		log.Warningf("Error reading configuration file: %s", err)
	}

	endpoint := vip.GetString(endpointEnvStr)

	return &Configuration{
		Timeout:  (time.Minute * 2),
		EndPoint: endpoint,
	}
}
