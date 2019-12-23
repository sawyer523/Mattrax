package types

// PolicyUUID is a unique identifier given to each policy
type PolicyUUID []byte

// A Policy contains instructions that are send to scoped managed device upon enrollment
type Policy struct {
	UUID        PolicyUUID `sqlgen:",primary"`
	DisplayName string
	Payload     []PolicyPayload `graphql:",optional"`
}

// PolicyPayload is a raw MDM instruction contain inside a Policy
type PolicyPayload struct {
	DisplayName  string
	Instructions map[MDMProtcol][]byte // This map is between an MDMProtocol and raw payload to be sent to the device
}

// PolicyService contains the implemented functionality for policies
type PolicyService interface {
	GetAll() ([]Policy, error)
	Get(uuid PolicyUUID) (Policy, error)
	CreateOrEdit(uuid PolicyUUID, policy Policy) error
	GenerateUUID() (PolicyUUID, error)
}
