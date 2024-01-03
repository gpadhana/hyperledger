/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package orderers

import (
	"github.com/hyperledger/fabric/common/flogging"
)

type ConnectionSourcer interface {
	RandomEndpoint() (*Endpoint, error)
	ShuffledEndpoints() []*Endpoint
	Update(globalAddrs []string, orgs map[string]OrdererOrg)
}

type ConnectionSourceCreator interface {
	CreateConnectionSource(logger *flogging.FabricLogger) ConnectionSourcer
}

type ConnectionSourceFactory struct {
	Overrides map[string]*Endpoint
}

func (f *ConnectionSourceFactory) CreateConnectionSource(logger *flogging.FabricLogger) ConnectionSourcer {
	return NewConnectionSource(logger, f.Overrides)
}
