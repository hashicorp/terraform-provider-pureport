package credentials

import (
	"errors"
	"fmt"
	"os"
	"os/user"

	"gopkg.in/ini.v1"
)

const (
	// FileProviderName the name of this provider
	FileProviderName = "FileProvider"

	// CredConfigFileEnvStr - environment variable for the path to the file
	// containing credentials information.
	CredConfigFileEnvStr = "PUREPORT_CREDENTIALS_FILE"

	// CredConfigProfileEnvStr - environment variable for the credentials
	// profile to use.
	CredConfigProfileEnvStr = "PUREPORT_PROFILE"

	// CredConfigAPIKeyStr - INI key for the API Key Value
	CredConfigAPIKeyStr = "api_key"

	// CredConfigAPISecretStr - INI key for the API Secret Value
	CredConfigAPISecretStr = "api_secret"
)

var (

	// ErrorFileAPIKeyNotFound is returned when the API Key can not be found in the configuration file
	ErrorFileAPIKeyNotFound = errors.New("api_key not found in configuration file")

	// ErrorFileAPISecretNotFound is returned when the API Secret can not be found in the configuration file
	ErrorFileAPISecretNotFound = errors.New("api_secret not found in configuration file")

	// ErrorFileConfigurationFileNotFound is return when the specified configuration file
	// can't be found.
	ErrorFileConfigurationFileNotFound = errors.New("No available credentials file found")

	// ErrorFileHomeDirectoryNotFound is returned when the users home directory can't be found.
	ErrorFileHomeDirectoryNotFound = errors.New("No available home directory found for current user")
)

// FileProvider holds information for fetch credentials information from
// configurations files.
type FileProvider struct {

	// The location of the configuration file
	// Linux/OSX "$HOME/.pureport/credentials"
	Filename string

	// The Profile for the credentials that we want to use
	Profile string

	// Have the credentials already been retrieved
	retrieved bool
}

// NewFileCredentials constructs a Credentials object for the specified file and profile
func NewFileCredentials(filename string, profile string) *Credentials {
	return NewCredentials(&FileProvider{
		Filename: filename,
		Profile:  profile,
	})
}

// Retrieve Provider.Retrieve()
func (p *FileProvider) Retrieve() (Value, error) {
	p.retrieved = false

	filename, err := p.filename()
	if err != nil {
		return Value{ProviderName: FileProviderName}, err
	}

	credentials, err := loadProfile(filename, p.profile())
	if err != nil {
		return Value{ProviderName: FileProviderName}, err
	}

	p.retrieved = true
	return credentials, nil
}

// IsExpired Provider.IsExpired()
func (p *FileProvider) IsExpired() bool {
	return !p.retrieved
}

// Helper method to get the profile name if it hasn't
// already been specified either from the environment
// or use the default.
func (p *FileProvider) profile() string {
	if p.Profile == "" {
		p.Profile = os.Getenv(CredConfigProfileEnvStr)
	}

	if p.Profile == "" {
		p.Profile = "default"
	}
	return p.Profile
}

func (p *FileProvider) filename() (string, error) {
	if len(p.Filename) != 0 {
		return p.Filename, nil
	}

	if p.Filename = os.Getenv(CredConfigFileEnvStr); len(p.Filename) != 0 {
		return p.Filename, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", ErrorFileHomeDirectoryNotFound
	}

	defaultPath := usr.HomeDir + "/.pureport/credentials"
	_, err = os.Stat(defaultPath)
	if err != nil {
		return "", ErrorFileConfigurationFileNotFound
	}

	return defaultPath, nil
}

func loadProfile(filename string, profile string) (Value, error) {
	config, err := ini.Load(filename)
	if err != nil {
		return Value{ProviderName: FileProviderName}, fmt.Errorf("Failed to load configuration file: %s", filename)
	}

	sec, err := config.GetSection(profile)
	if err != nil {
		return Value{ProviderName: FileProviderName}, fmt.Errorf("Failed to find specified profile: %s", profile)
	}

	key, err := sec.GetKey(CredConfigAPIKeyStr)
	if err != nil {
		return Value{ProviderName: FileProviderName}, ErrorFileAPIKeyNotFound
	}

	secret, err := sec.GetKey(CredConfigAPISecretStr)
	if err != nil {
		return Value{ProviderName: FileProviderName}, ErrorFileAPISecretNotFound
	}

	return Value{
		ProviderName: FileProviderName,
		APIKey:       key.String(),
		Secret:       secret.String(),
	}, nil
}
