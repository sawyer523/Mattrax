package types

import (
	"github.com/pkg/errors"
)

// PolicyUUID is a unique identifier given to each policy
type PolicyUUID []byte

// A Policy contains instructions that are send to scoped managed device upon enrollment
type Policy struct {
	UUID        PolicyUUID
	DisplayName string
	Payload     []PolicyPayload `graphql:",optional"`
}

// PolicyPayload is a raw MDM instruction contain inside a Policy
type PolicyPayload struct {
	DisplayName string
	// Instructions map[MDMProtcol][]byte // This map is between an MDMProtocol and raw payload to be sent to the device
}

// ErrPolicyNotFound is the error returned if a user can't be found
var ErrPolicyNotFound = errors.New("Error: Policy not found")

// PolicyService contains the implemented functionality for policies
type PolicyService interface {
	GetAll() ([]Policy, error)
	Get(uuid PolicyUUID) (Policy, error)
	CreateOrEdit(uuid PolicyUUID, policy Policy) error
}
