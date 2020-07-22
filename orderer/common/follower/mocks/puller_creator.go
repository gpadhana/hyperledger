// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/orderer/common/follower"
)

type BlockPullerFactory struct {
	BlockPullerStub        func(*common.Block) (follower.ChannelPuller, error)
	blockPullerMutex       sync.RWMutex
	blockPullerArgsForCall []struct {
		arg1 *common.Block
	}
	blockPullerReturns struct {
		result1 follower.ChannelPuller
		result2 error
	}
	blockPullerReturnsOnCall map[int]struct {
		result1 follower.ChannelPuller
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *BlockPullerFactory) BlockPuller(arg1 *common.Block) (follower.ChannelPuller, error) {
	fake.blockPullerMutex.Lock()
	ret, specificReturn := fake.blockPullerReturnsOnCall[len(fake.blockPullerArgsForCall)]
	fake.blockPullerArgsForCall = append(fake.blockPullerArgsForCall, struct {
		arg1 *common.Block
	}{arg1})
	fake.recordInvocation("BlockPuller", []interface{}{arg1})
	fake.blockPullerMutex.Unlock()
	if fake.BlockPullerStub != nil {
		return fake.BlockPullerStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.blockPullerReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *BlockPullerFactory) BlockPullerCallCount() int {
	fake.blockPullerMutex.RLock()
	defer fake.blockPullerMutex.RUnlock()
	return len(fake.blockPullerArgsForCall)
}

func (fake *BlockPullerFactory) BlockPullerCalls(stub func(*common.Block) (follower.ChannelPuller, error)) {
	fake.blockPullerMutex.Lock()
	defer fake.blockPullerMutex.Unlock()
	fake.BlockPullerStub = stub
}

func (fake *BlockPullerFactory) BlockPullerArgsForCall(i int) *common.Block {
	fake.blockPullerMutex.RLock()
	defer fake.blockPullerMutex.RUnlock()
	argsForCall := fake.blockPullerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *BlockPullerFactory) BlockPullerReturns(result1 follower.ChannelPuller, result2 error) {
	fake.blockPullerMutex.Lock()
	defer fake.blockPullerMutex.Unlock()
	fake.BlockPullerStub = nil
	fake.blockPullerReturns = struct {
		result1 follower.ChannelPuller
		result2 error
	}{result1, result2}
}

func (fake *BlockPullerFactory) BlockPullerReturnsOnCall(i int, result1 follower.ChannelPuller, result2 error) {
	fake.blockPullerMutex.Lock()
	defer fake.blockPullerMutex.Unlock()
	fake.BlockPullerStub = nil
	if fake.blockPullerReturnsOnCall == nil {
		fake.blockPullerReturnsOnCall = make(map[int]struct {
			result1 follower.ChannelPuller
			result2 error
		})
	}
	fake.blockPullerReturnsOnCall[i] = struct {
		result1 follower.ChannelPuller
		result2 error
	}{result1, result2}
}

func (fake *BlockPullerFactory) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.blockPullerMutex.RLock()
	defer fake.blockPullerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *BlockPullerFactory) recordInvocation(key string, args []interface{}) {
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
