package mattrax

import (
	"github.com/mattrax/Mattrax/internal/types"
)

// Server holds the global server state
type Server struct {
	Config        Config
	UserService   types.UserService
	PolicyService types.PolicyService
}

// Config holds the global server config
type Config struct {
	TenantName    string
	PrimaryDomain string
}
