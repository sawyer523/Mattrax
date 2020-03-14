package mattrax

import (
	"github.com/alexflint/go-arg"
	"github.com/mattrax/Mattrax/pkg/types"
)

// Config holds the read only server config
// These values *should* not need to change after deployment
// These values are set via command line flags
type Config struct {
	Port            int    `arg:"-p" help:"the port for the HTTPS webserver to listen on" placeholder:"443" default:"443"`
	Domain          string `arg:"-d" help:"the domain name the server is accessible on" placeholder:"mdm.example.com"`
	Database        string `help:"the database provider followed by ":" and then the connection address" placeholder:"boltdb:/var/mattrax.db" default:"boltdb:/var/mattrax.db"`
	CertFile        string `arg:"--cert" help:"the path to the https certificate for the HTTPS webserver" placeholder:"/your-cert.pem"`
	KeyFile         string `arg:"--key" help:"the path to the https private key for the HTTPS webserver" placeholder:"/your-key.pem"`
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

	if config.Database == "" {
		if config.DevelopmentMode {
			config.Database = "boltdb:./mattrax.db"
		} else {
			p.Fail("you must provide a database path.")
		}
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
	return "Mattrax MDM Server. Created By Oscar Beaumont. Version: " + Version
}
