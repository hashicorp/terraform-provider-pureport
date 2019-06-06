package endpoint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/op/go-logging"
	"github.com/pureport/pureport-sdk-go/pureport"
	"github.com/pureport/pureport-sdk-go/pureport/credentials"
)

var log = logging.MustGetLogger("main_logger")

const providerName = "EndpointCredentialsProvider"

// Provider for endpoint based credential requests
type Provider struct {
	credentials.Expiry
	*credentials.Credentials

	// https://golang.org/pkg/net/http/
	Client *http.Client

	// HTTP Endpoint to query credentials from
	EndPoint string

	// ExpiryWindow allow the refresh of the authorization token
	// to be refreshed prior to expiration to ensure that we have
	// a valid token.
	ExpiryWindow time.Duration
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Status      int    `json:"status"`
	Code        string `json:"code"`
	Message     string `json:"message"`
}

func newEndpointProvider(cfg pureport.Configuration, endpoint string, cred *credentials.Credentials) credentials.Provider {

	return &Provider{
		Client: &http.Client{
			Timeout: cfg.Timeout,
		},
		// Trim any trailing slashes in the endpoint if they exist
		EndPoint:    strings.TrimRight(endpoint, "/"),
		Credentials: cred,
	}
}

// NewEndPointCredentials creates a new credentials for requesting and updating
// credentials from a remote endpoint.
func NewEndPointCredentials(cfg pureport.Configuration, endpoint string, cred *credentials.Credentials) *credentials.Credentials {
	return credentials.NewCredentials(newEndpointProvider(cfg, endpoint, cred))
}

// IsExpired - see Provider.IsExpired()
func (p *Provider) IsExpired() bool {
	return p.Expiry.IsExpired()
}

// Retrieve - see Provider.Retrieve()
func (p *Provider) Retrieve() (credentials.Value, error) {

	local, err := p.Credentials.Get()
	if err != nil {
		return credentials.Value{ProviderName: providerName}, err
	}

	// Create the body of the request
	values := map[string]string{"key": local.APIKey, "secret": local.Secret}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		return credentials.Value{ProviderName: providerName}, fmt.Errorf("Error creating credential body")
	}

	buf := bytes.NewBuffer(jsonValue)

	log.Debugf("Logging in to EndPoint: %s/login", p.EndPoint)

	// Create the HTTP Request
	resp, err := p.Client.Post(fmt.Sprintf("%s/login", p.EndPoint), "application/json", buf)
	if err != nil {
		log.Errorf("HTTP Response: %s Error: %s", resp, err)
		return credentials.Value{ProviderName: providerName}, fmt.Errorf("Error creating credentials login request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("HTTP Body: %s, Error: %s", string(body), err)
		return credentials.Value{ProviderName: providerName}, fmt.Errorf("Error reading credential request body")
	}

	var data loginResponse
	if err = json.Unmarshal(body, &data); err != nil {
		return credentials.Value{ProviderName: providerName}, fmt.Errorf("Error reading credential response")
	}

	// Check to make sure an error wasn't returned
	if data.Status != 0 {
		return credentials.Value{ProviderName: providerName}, fmt.Errorf(
			"(%v) %v, %v: %v",
			data.Status,
			data.Code,
			data.Message,
			p.EndPoint,
		)
	}

	// Initialize the Expiry
	expiresIn, err := time.ParseDuration(fmt.Sprintf("%vs", data.ExpiresIn))
	if err != nil {
		return credentials.Value{ProviderName: providerName}, fmt.Errorf("Error converting expiry time")
	}

	p.Expiry.SetExpiration(time.Now().Add(expiresIn), p.ExpiryWindow)

	return credentials.Value{
		APIKey:       local.APIKey,
		Secret:       local.Secret,
		SessionToken: data.AccessToken,
	}, nil
}
