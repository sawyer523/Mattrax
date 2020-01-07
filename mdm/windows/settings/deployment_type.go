package wsettings

// TODO: Probally global setting for both MDM

type DeploymentType int

func (dt DeploymentType) String() string {
	if dt == DeploymentStandalone {
		return "Standalone"
	}
	if dt == DeploymentAzureAD {
		return "AzureAD"
	}
	return "" // TODO: This should probs be an error
}

const (
	DeploymentStandalone DeploymentType = iota
	DeploymentAzureAD
)
