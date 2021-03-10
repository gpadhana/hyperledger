/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mgmt

import (
	"reflect"
	"sync"

	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/msp/cache"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// FIXME: AS SOON AS THE CHAIN MANAGEMENT CODE IS COMPLETE,
// THESE MAPS AND HELPER FUNCTIONS SHOULD DISAPPEAR BECAUSE
// OWNERSHIP OF PER-CHAIN MSP MANAGERS WILL BE HANDLED BY IT;
// HOWEVER IN THE INTERIM, THESE HELPER FUNCTIONS ARE REQUIRED

var (
	m         sync.Mutex
	localMsp  msp.MSP
	mspMap    = make(map[string]msp.MSPManager)
	mspLogger = flogging.MustGetLogger("msp")
)

// TODO - this is a temporary solution to allow the peer to track whether the
// MSPManager has been setup for a channel, which indicates whether the channel
// exists or not
type mspMgmtMgr struct {
	msp.MSPManager
	// track whether this MSPManager has been setup successfully
	up bool
}

func (mgr *mspMgmtMgr) DeserializeIdentity(serializedIdentity []byte) (msp.Identity, error) {
	if !mgr.up {
		return nil, errors.New("channel doesn't exist")
	}
	return mgr.MSPManager.DeserializeIdentity(serializedIdentity)
}

func (mgr *mspMgmtMgr) Setup(msps []msp.MSP) error {
	err := mgr.MSPManager.Setup(msps)
	if err == nil {
		mgr.up = true
	}
	return err
}

// GetManagerForChain returns the msp manager for the supplied
// chain; if no such manager exists, one is created
func GetManagerForChain(chainID string) msp.MSPManager {
	m.Lock()
	defer m.Unlock()

	mspMgr, ok := mspMap[chainID]
	if !ok {
		mspLogger.Debugf("Created new msp manager for channel `%s`", chainID)
		mspMgmtMgr := &mspMgmtMgr{msp.NewMSPManager(), false}
		mspMap[chainID] = mspMgmtMgr
		mspMgr = mspMgmtMgr
	} else {
		// check for internal mspManagerImpl and mspMgmtMgr types. if a different
		// type is found, it's because a developer has added a new type that
		// implements the MSPManager interface and should add a case to the logic
		// above to handle it.
		if !(reflect.TypeOf(mspMgr).Elem().Name() == "mspManagerImpl" || reflect.TypeOf(mspMgr).Elem().Name() == "mspMgmtMgr") {
			panic("Found unexpected MSPManager type.")
		}
		mspLogger.Debugf("Returning existing manager for channel '%s'", chainID)
	}
	return mspMgr
}

// GetManagers returns all the managers registered
func GetDeserializers() map[string]msp.IdentityDeserializer {
	m.Lock()
	defer m.Unlock()

	clone := make(map[string]msp.IdentityDeserializer)

	for key, mspManager := range mspMap {
		clone[key] = mspManager
	}

	return clone
}

// XXXSetMSPManager is a stopgap solution to transition from the custom MSP config block
// parsing to the channelconfig.Resources interface, while preserving the problematic singleton
// nature of the MSP manager
func XXXSetMSPManager(chainID string, manager msp.MSPManager) {
	m.Lock()
	defer m.Unlock()

	mspMap[chainID] = &mspMgmtMgr{manager, true}
}

// GetLocalMSP returns the local msp (and creates it if it doesn't exist)
func GetLocalMSP(cryptoProvider bccsp.BCCSP) msp.MSP {
	m.Lock()
	defer m.Unlock()

	if localMsp != nil {
		return localMsp
	}

	localMsp = loadLocalMSP(cryptoProvider)

	return localMsp
}

func loadLocalMSP(bccsp bccsp.BCCSP) msp.MSP {
	// determine the type of MSP (by default, we'll use bccspMSP)
	mspType := viper.GetString("peer.localMspType")
	if mspType == "" {
		mspType = msp.ProviderTypeToString(msp.FABRIC)
	}

	newOpts, found := msp.Options[mspType]
	if !found {
		mspLogger.Panicf("msp type " + mspType + " unknown")
	}

	mspInst, err := msp.New(newOpts, bccsp)
	if err != nil {
		mspLogger.Fatalf("Failed to initialize local MSP, received err %+v", err)
	}
	switch mspType {
	case msp.ProviderTypeToString(msp.FABRIC):
		mspInst, err = cache.New(mspInst)
		if err != nil {
			mspLogger.Fatalf("Failed to initialize local MSP, received err %+v", err)
		}
	case msp.ProviderTypeToString(msp.IDEMIX):
		// Do nothing
	default:
		panic("msp type " + mspType + " unknown")
	}

	mspLogger.Debugf("Created new local MSP")

	return mspInst
}

// GetIdentityDeserializer returns the IdentityDeserializer for the given chain
func GetIdentityDeserializer(chainID string, cryptoProvider bccsp.BCCSP) msp.IdentityDeserializer {
	if chainID == "" {
		return GetLocalMSP(cryptoProvider)
	}

	return GetManagerForChain(chainID)
}
