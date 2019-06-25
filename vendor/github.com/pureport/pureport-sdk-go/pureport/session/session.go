package session

import (
	"context"
	"os"

	"github.com/op/go-logging"
	"github.com/pureport/pureport-sdk-go/pureport"
	"github.com/pureport/pureport-sdk-go/pureport/client"
	"github.com/pureport/pureport-sdk-go/pureport/credentials"
	"github.com/pureport/pureport-sdk-go/pureport/credentials/endpoint"
)

var log = logging.MustGetLogger("main_logger")

// Session contains the data for a particular request session
type Session struct {
	*credentials.Credentials
	*pureport.Configuration

	Client *client.APIClient
}

func createClient(cfg *pureport.Configuration) *client.APIClient {
	c := client.NewConfiguration()
	c.UserAgent = cfg.UserAgent
	c.BasePath = cfg.EndPoint
	//c.AddDefaultHeader()

	if hostname, err := os.Hostname(); err != nil {
		c.Host = ""
	} else {
		c.Host = hostname
	}

	return client.NewAPIClient(c)
}

func createCredentials(cfg *pureport.Configuration) *credentials.Credentials {

	cred := credentials.NewChainCredentials([]credentials.Provider{
		&credentials.StaticProvider{
			APIKey:    cfg.APIKey,
			APISecret: cfg.APISecret,
		},
		&credentials.ViperProvider{
			Profile: cfg.AuthenticationProfile,
		},
	})

	return endpoint.NewEndPointCredentials(*cfg, cfg.EndPoint, cred)
}

// NewSession creates a new request session
func NewSession(cfg *pureport.Configuration) *Session {

	log.Debug("*** Creating new Pureport Session ***")

	return &Session{
		Credentials:   createCredentials(cfg),
		Configuration: cfg,
		Client:        createClient(cfg),
	}
}

// GetSessionContext gathers the context information need to
// for communicating with the Pureport API
func (s *Session) GetSessionContext() context.Context {

	value, err := s.Credentials.Get()
	if err != nil {
		log.Criticalf("Error reading credentials: %s", err)
	}

	return context.WithValue(context.Background(), client.ContextAccessToken, value.SessionToken)
}
