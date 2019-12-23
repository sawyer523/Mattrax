package mattrax

import (
	"github.com/mattrax/Mattrax/internal/types"
)

// Server holds the global server state
type Server struct {
	Config          Config
	UserService     types.UserService
	PolicyService   types.PolicyService
	SettingsService types.SettingsService
}

// Config holds the global server config
type Config struct {
	Port            int
	Domains         []string
	CertFile        string
	KeyFile         string
	DevelopmentMode bool
}
