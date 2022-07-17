// Code generated by counterfeiter. DO NOT EDIT.
package mock

import (
	"sync"

	"github.com/hyperledger/fabric/v2/core/ledger"
	"github.com/hyperledger/fabric/v2/core/ledger/cceventmgmt"
)

type ChaincodeInfoProvider struct {
	GetDeployedChaincodeInfoStub        func(string, *cceventmgmt.ChaincodeDefinition) (*ledger.DeployedChaincodeInfo, error)
	getDeployedChaincodeInfoMutex       sync.RWMutex
	getDeployedChaincodeInfoArgsForCall []struct {
		arg1 string
		arg2 *cceventmgmt.ChaincodeDefinition
	}
	getDeployedChaincodeInfoReturns struct {
		result1 *ledger.DeployedChaincodeInfo
		result2 error
	}
	getDeployedChaincodeInfoReturnsOnCall map[int]struct {
		result1 *ledger.DeployedChaincodeInfo
		result2 error
	}
	RetrieveChaincodeArtifactsStub        func(*cceventmgmt.ChaincodeDefinition) (bool, []byte, error)
	retrieveChaincodeArtifactsMutex       sync.RWMutex
	retrieveChaincodeArtifactsArgsForCall []struct {
		arg1 *cceventmgmt.ChaincodeDefinition
	}
	retrieveChaincodeArtifactsReturns struct {
		result1 bool
		result2 []byte
		result3 error
	}
	retrieveChaincodeArtifactsReturnsOnCall map[int]struct {
		result1 bool
		result2 []byte
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChaincodeInfoProvider) GetDeployedChaincodeInfo(arg1 string, arg2 *cceventmgmt.ChaincodeDefinition) (*ledger.DeployedChaincodeInfo, error) {
	fake.getDeployedChaincodeInfoMutex.Lock()
	ret, specificReturn := fake.getDeployedChaincodeInfoReturnsOnCall[len(fake.getDeployedChaincodeInfoArgsForCall)]
	fake.getDeployedChaincodeInfoArgsForCall = append(fake.getDeployedChaincodeInfoArgsForCall, struct {
		arg1 string
		arg2 *cceventmgmt.ChaincodeDefinition
	}{arg1, arg2})
	fake.recordInvocation("GetDeployedChaincodeInfo", []interface{}{arg1, arg2})
	fake.getDeployedChaincodeInfoMutex.Unlock()
	if fake.GetDeployedChaincodeInfoStub != nil {
		return fake.GetDeployedChaincodeInfoStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getDeployedChaincodeInfoReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ChaincodeInfoProvider) GetDeployedChaincodeInfoCallCount() int {
	fake.getDeployedChaincodeInfoMutex.RLock()
	defer fake.getDeployedChaincodeInfoMutex.RUnlock()
	return len(fake.getDeployedChaincodeInfoArgsForCall)
}

func (fake *ChaincodeInfoProvider) GetDeployedChaincodeInfoCalls(stub func(string, *cceventmgmt.ChaincodeDefinition) (*ledger.DeployedChaincodeInfo, error)) {
	fake.getDeployedChaincodeInfoMutex.Lock()
	defer fake.getDeployedChaincodeInfoMutex.Unlock()
	fake.GetDeployedChaincodeInfoStub = stub
}

func (fake *ChaincodeInfoProvider) GetDeployedChaincodeInfoArgsForCall(i int) (string, *cceventmgmt.ChaincodeDefinition) {
	fake.getDeployedChaincodeInfoMutex.RLock()
	defer fake.getDeployedChaincodeInfoMutex.RUnlock()
	argsForCall := fake.getDeployedChaincodeInfoArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *ChaincodeInfoProvider) GetDeployedChaincodeInfoReturns(result1 *ledger.DeployedChaincodeInfo, result2 error) {
	fake.getDeployedChaincodeInfoMutex.Lock()
	defer fake.getDeployedChaincodeInfoMutex.Unlock()
	fake.GetDeployedChaincodeInfoStub = nil
	fake.getDeployedChaincodeInfoReturns = struct {
		result1 *ledger.DeployedChaincodeInfo
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeInfoProvider) GetDeployedChaincodeInfoReturnsOnCall(i int, result1 *ledger.DeployedChaincodeInfo, result2 error) {
	fake.getDeployedChaincodeInfoMutex.Lock()
	defer fake.getDeployedChaincodeInfoMutex.Unlock()
	fake.GetDeployedChaincodeInfoStub = nil
	if fake.getDeployedChaincodeInfoReturnsOnCall == nil {
		fake.getDeployedChaincodeInfoReturnsOnCall = make(map[int]struct {
			result1 *ledger.DeployedChaincodeInfo
			result2 error
		})
	}
	fake.getDeployedChaincodeInfoReturnsOnCall[i] = struct {
		result1 *ledger.DeployedChaincodeInfo
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeInfoProvider) RetrieveChaincodeArtifacts(arg1 *cceventmgmt.ChaincodeDefinition) (bool, []byte, error) {
	fake.retrieveChaincodeArtifactsMutex.Lock()
	ret, specificReturn := fake.retrieveChaincodeArtifactsReturnsOnCall[len(fake.retrieveChaincodeArtifactsArgsForCall)]
	fake.retrieveChaincodeArtifactsArgsForCall = append(fake.retrieveChaincodeArtifactsArgsForCall, struct {
		arg1 *cceventmgmt.ChaincodeDefinition
	}{arg1})
	fake.recordInvocation("RetrieveChaincodeArtifacts", []interface{}{arg1})
	fake.retrieveChaincodeArtifactsMutex.Unlock()
	if fake.RetrieveChaincodeArtifactsStub != nil {
		return fake.RetrieveChaincodeArtifactsStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.retrieveChaincodeArtifactsReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *ChaincodeInfoProvider) RetrieveChaincodeArtifactsCallCount() int {
	fake.retrieveChaincodeArtifactsMutex.RLock()
	defer fake.retrieveChaincodeArtifactsMutex.RUnlock()
	return len(fake.retrieveChaincodeArtifactsArgsForCall)
}

func (fake *ChaincodeInfoProvider) RetrieveChaincodeArtifactsCalls(stub func(*cceventmgmt.ChaincodeDefinition) (bool, []byte, error)) {
	fake.retrieveChaincodeArtifactsMutex.Lock()
	defer fake.retrieveChaincodeArtifactsMutex.Unlock()
	fake.RetrieveChaincodeArtifactsStub = stub
}

func (fake *ChaincodeInfoProvider) RetrieveChaincodeArtifactsArgsForCall(i int) *cceventmgmt.ChaincodeDefinition {
	fake.retrieveChaincodeArtifactsMutex.RLock()
	defer fake.retrieveChaincodeArtifactsMutex.RUnlock()
	argsForCall := fake.retrieveChaincodeArtifactsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ChaincodeInfoProvider) RetrieveChaincodeArtifactsReturns(result1 bool, result2 []byte, result3 error) {
	fake.retrieveChaincodeArtifactsMutex.Lock()
	defer fake.retrieveChaincodeArtifactsMutex.Unlock()
	fake.RetrieveChaincodeArtifactsStub = nil
	fake.retrieveChaincodeArtifactsReturns = struct {
		result1 bool
		result2 []byte
		result3 error
	}{result1, result2, result3}
}

func (fake *ChaincodeInfoProvider) RetrieveChaincodeArtifactsReturnsOnCall(i int, result1 bool, result2 []byte, result3 error) {
	fake.retrieveChaincodeArtifactsMutex.Lock()
	defer fake.retrieveChaincodeArtifactsMutex.Unlock()
	fake.RetrieveChaincodeArtifactsStub = nil
	if fake.retrieveChaincodeArtifactsReturnsOnCall == nil {
		fake.retrieveChaincodeArtifactsReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 []byte
			result3 error
		})
	}
	fake.retrieveChaincodeArtifactsReturnsOnCall[i] = struct {
		result1 bool
		result2 []byte
		result3 error
	}{result1, result2, result3}
}

func (fake *ChaincodeInfoProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getDeployedChaincodeInfoMutex.RLock()
	defer fake.getDeployedChaincodeInfoMutex.RUnlock()
	fake.retrieveChaincodeArtifactsMutex.RLock()
	defer fake.retrieveChaincodeArtifactsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChaincodeInfoProvider) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
