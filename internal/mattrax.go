package mattrax

import (
	"github.com/alexflint/go-arg"
	"github.com/mattrax/Mattrax/internal/certificates"
	"github.com/mattrax/Mattrax/internal/devices"
	"github.com/mattrax/Mattrax/internal/settings"
	"github.com/mattrax/Mattrax/internal/types"
)

// Version contains the Mattrax server version
// This varible's correct value is injected at build time
var Version string = "0.0.0-development"

// Server holds the global server state
type Server struct {
	Version      string // The Mattrax server version
	Config       Config // Config holds static varibles defined through Command Line flags. It is read only!
	Settings     *settings.Service
	Certificates *certificates.Service
	Devices      devices.Service

	// TODO Cleanup below
	UserService   types.UserService
	PolicyService types.PolicyService
}

// Config holds the static server config
// These values are set by command line flags.
type Config struct {
	Port            int    `arg:"-p" help:"the port for the HTTPS webserver to listen on" placeholder:"443" default:"443"`
	Domain          string `arg:"-d" help:"the domain name the server is accessible on" placeholder:"mdm.example.com"`
	DBPath          string `help:"the path where the file database is stored" placeholder:"/var/mattrax.db" default:"/var/mattrax.db" graphql:"DBPath"`
	CertFile        string `arg:"--cert" help:"the path to the https certificate for the HTTPS webserver" placeholder:"/dont-put-your-cert-file-here.pem"`
	KeyFile         string `arg:"--key" help:"the path to the https certificate private key for the HTTPS webserver" placeholder:"/dont-put-your-key-file-here.pem"`
	DevelopmentMode bool   `arg:"--dev" help:"enables verbose output and loosens security measures to aid developers" default:"false"`
}

// Verify checks that the recieved values are valid and are what the server is expecting
func (config *Config) Verify(p *arg.Parser) {
	if config.Port <= 0 || config.Port > 49151 {
		p.Fail("invalid port. The port must be between 0 and 49151.")
	}

	if config.Domain == "" {
		p.Fail("you must provide a domain")
	} else if !types.IsDNSNameRegex.MatchString(config.Domain) {
		p.Fail("invalid domain name. Please ensure it doesn't start with a schema or end with a path.")
	}

	if config.DBPath == "" {
		p.Fail("you must provide a database path")
	} else if config.DBPath == "." {
		config.DBPath = "./mattrax.db"
	}

	if config.CertFile == "" {
		if config.DevelopmentMode == true {
			config.CertFile = "./certs/certificate.pem"
		} else {
			p.Fail("you must provide a certificate file")
		}
	}

	if config.KeyFile == "" {
		if config.DevelopmentMode == true {
			config.KeyFile = "./certs/privatekey.pem"
		} else {
			p.Fail("you must provide a certificate key file")
		}
	}
}

// Version is used by go-args to show a custom version string.
func (config Config) Version() string {
	return "Mattrax MDM Server. Created By Oscar Beaumont. Version " + Version
}
