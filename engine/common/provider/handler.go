package provider

import (
	"github.com/dapperlabs/flow-go/model/flow"
	"github.com/dapperlabs/flow-go/model/messages"
)

type Handler struct {
	Resource messages.Resource
	Selector flow.IdentityFilter
	Retrieve RetrieveFunc
}

type RetrieveFunc func(flow.Identifier) (flow.Entity, error)
