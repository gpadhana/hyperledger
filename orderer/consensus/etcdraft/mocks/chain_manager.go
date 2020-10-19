// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/hyperledger/fabric/orderer/common/types"
	"github.com/hyperledger/fabric/orderer/consensus"
	"github.com/hyperledger/fabric/orderer/consensus/etcdraft"
)

type ChainManager struct {
	CreateChainStub        func(string)
	createChainMutex       sync.RWMutex
	createChainArgsForCall []struct {
		arg1 string
	}
	GetConsensusChainStub        func(string) consensus.Chain
	getConsensusChainMutex       sync.RWMutex
	getConsensusChainArgsForCall []struct {
		arg1 string
	}
	getConsensusChainReturns struct {
		result1 consensus.Chain
	}
	getConsensusChainReturnsOnCall map[int]struct {
		result1 consensus.Chain
	}
	ReportRelationAndStatusMetricsStub        func(string, types.ClusterRelation, types.Status)
	reportRelationAndStatusMetricsMutex       sync.RWMutex
	reportRelationAndStatusMetricsArgsForCall []struct {
		arg1 string
		arg2 types.ClusterRelation
		arg3 types.Status
	}
	SwitchChainToFollowerStub        func(string)
	switchChainToFollowerMutex       sync.RWMutex
	switchChainToFollowerArgsForCall []struct {
		arg1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChainManager) CreateChain(arg1 string) {
	fake.createChainMutex.Lock()
	fake.createChainArgsForCall = append(fake.createChainArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("CreateChain", []interface{}{arg1})
	fake.createChainMutex.Unlock()
	if fake.CreateChainStub != nil {
		fake.CreateChainStub(arg1)
	}
}

func (fake *ChainManager) CreateChainCallCount() int {
	fake.createChainMutex.RLock()
	defer fake.createChainMutex.RUnlock()
	return len(fake.createChainArgsForCall)
}

func (fake *ChainManager) CreateChainCalls(stub func(string)) {
	fake.createChainMutex.Lock()
	defer fake.createChainMutex.Unlock()
	fake.CreateChainStub = stub
}

func (fake *ChainManager) CreateChainArgsForCall(i int) string {
	fake.createChainMutex.RLock()
	defer fake.createChainMutex.RUnlock()
	argsForCall := fake.createChainArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ChainManager) GetConsensusChain(arg1 string) consensus.Chain {
	fake.getConsensusChainMutex.Lock()
	ret, specificReturn := fake.getConsensusChainReturnsOnCall[len(fake.getConsensusChainArgsForCall)]
	fake.getConsensusChainArgsForCall = append(fake.getConsensusChainArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetConsensusChain", []interface{}{arg1})
	fake.getConsensusChainMutex.Unlock()
	if fake.GetConsensusChainStub != nil {
		return fake.GetConsensusChainStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getConsensusChainReturns
	return fakeReturns.result1
}

func (fake *ChainManager) GetConsensusChainCallCount() int {
	fake.getConsensusChainMutex.RLock()
	defer fake.getConsensusChainMutex.RUnlock()
	return len(fake.getConsensusChainArgsForCall)
}

func (fake *ChainManager) GetConsensusChainCalls(stub func(string) consensus.Chain) {
	fake.getConsensusChainMutex.Lock()
	defer fake.getConsensusChainMutex.Unlock()
	fake.GetConsensusChainStub = stub
}

func (fake *ChainManager) GetConsensusChainArgsForCall(i int) string {
	fake.getConsensusChainMutex.RLock()
	defer fake.getConsensusChainMutex.RUnlock()
	argsForCall := fake.getConsensusChainArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ChainManager) GetConsensusChainReturns(result1 consensus.Chain) {
	fake.getConsensusChainMutex.Lock()
	defer fake.getConsensusChainMutex.Unlock()
	fake.GetConsensusChainStub = nil
	fake.getConsensusChainReturns = struct {
		result1 consensus.Chain
	}{result1}
}

func (fake *ChainManager) GetConsensusChainReturnsOnCall(i int, result1 consensus.Chain) {
	fake.getConsensusChainMutex.Lock()
	defer fake.getConsensusChainMutex.Unlock()
	fake.GetConsensusChainStub = nil
	if fake.getConsensusChainReturnsOnCall == nil {
		fake.getConsensusChainReturnsOnCall = make(map[int]struct {
			result1 consensus.Chain
		})
	}
	fake.getConsensusChainReturnsOnCall[i] = struct {
		result1 consensus.Chain
	}{result1}
}

func (fake *ChainManager) ReportRelationAndStatusMetrics(arg1 string, arg2 types.ClusterRelation, arg3 types.Status) {
	fake.reportRelationAndStatusMetricsMutex.Lock()
	fake.reportRelationAndStatusMetricsArgsForCall = append(fake.reportRelationAndStatusMetricsArgsForCall, struct {
		arg1 string
		arg2 types.ClusterRelation
		arg3 types.Status
	}{arg1, arg2, arg3})
	fake.recordInvocation("ReportRelationAndStatusMetrics", []interface{}{arg1, arg2, arg3})
	fake.reportRelationAndStatusMetricsMutex.Unlock()
	if fake.ReportRelationAndStatusMetricsStub != nil {
		fake.ReportRelationAndStatusMetricsStub(arg1, arg2, arg3)
	}
}

func (fake *ChainManager) ReportRelationAndStatusMetricsCallCount() int {
	fake.reportRelationAndStatusMetricsMutex.RLock()
	defer fake.reportRelationAndStatusMetricsMutex.RUnlock()
	return len(fake.reportRelationAndStatusMetricsArgsForCall)
}

func (fake *ChainManager) ReportRelationAndStatusMetricsCalls(stub func(string, types.ClusterRelation, types.Status)) {
	fake.reportRelationAndStatusMetricsMutex.Lock()
	defer fake.reportRelationAndStatusMetricsMutex.Unlock()
	fake.ReportRelationAndStatusMetricsStub = stub
}

func (fake *ChainManager) ReportRelationAndStatusMetricsArgsForCall(i int) (string, types.ClusterRelation, types.Status) {
	fake.reportRelationAndStatusMetricsMutex.RLock()
	defer fake.reportRelationAndStatusMetricsMutex.RUnlock()
	argsForCall := fake.reportRelationAndStatusMetricsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *ChainManager) SwitchChainToFollower(arg1 string) {
	fake.switchChainToFollowerMutex.Lock()
	fake.switchChainToFollowerArgsForCall = append(fake.switchChainToFollowerArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("SwitchChainToFollower", []interface{}{arg1})
	fake.switchChainToFollowerMutex.Unlock()
	if fake.SwitchChainToFollowerStub != nil {
		fake.SwitchChainToFollowerStub(arg1)
	}
}

func (fake *ChainManager) SwitchChainToFollowerCallCount() int {
	fake.switchChainToFollowerMutex.RLock()
	defer fake.switchChainToFollowerMutex.RUnlock()
	return len(fake.switchChainToFollowerArgsForCall)
}

func (fake *ChainManager) SwitchChainToFollowerCalls(stub func(string)) {
	fake.switchChainToFollowerMutex.Lock()
	defer fake.switchChainToFollowerMutex.Unlock()
	fake.SwitchChainToFollowerStub = stub
}

func (fake *ChainManager) SwitchChainToFollowerArgsForCall(i int) string {
	fake.switchChainToFollowerMutex.RLock()
	defer fake.switchChainToFollowerMutex.RUnlock()
	argsForCall := fake.switchChainToFollowerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ChainManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createChainMutex.RLock()
	defer fake.createChainMutex.RUnlock()
	fake.getConsensusChainMutex.RLock()
	defer fake.getConsensusChainMutex.RUnlock()
	fake.reportRelationAndStatusMetricsMutex.RLock()
	defer fake.reportRelationAndStatusMetricsMutex.RUnlock()
	fake.switchChainToFollowerMutex.RLock()
	defer fake.switchChainToFollowerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChainManager) recordInvocation(key string, args []interface{}) {
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

var _ etcdraft.ChainManager = new(ChainManager)
