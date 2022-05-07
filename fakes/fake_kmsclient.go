// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/aws/aws-sdk-go/service/kms"
	environment "github.com/telia-oss/aws-env"
)

type FakeKMSClient struct {
	DecryptStub        func(*kms.DecryptInput) (*kms.DecryptOutput, error)
	decryptMutex       sync.RWMutex
	decryptArgsForCall []struct {
		arg1 *kms.DecryptInput
	}
	decryptReturns struct {
		result1 *kms.DecryptOutput
		result2 error
	}
	decryptReturnsOnCall map[int]struct {
		result1 *kms.DecryptOutput
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeKMSClient) Decrypt(arg1 *kms.DecryptInput) (*kms.DecryptOutput, error) {
	fake.decryptMutex.Lock()
	ret, specificReturn := fake.decryptReturnsOnCall[len(fake.decryptArgsForCall)]
	fake.decryptArgsForCall = append(fake.decryptArgsForCall, struct {
		arg1 *kms.DecryptInput
	}{arg1})
	stub := fake.DecryptStub
	fakeReturns := fake.decryptReturns
	fake.recordInvocation("Decrypt", []interface{}{arg1})
	fake.decryptMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeKMSClient) DecryptCallCount() int {
	fake.decryptMutex.RLock()
	defer fake.decryptMutex.RUnlock()
	return len(fake.decryptArgsForCall)
}

func (fake *FakeKMSClient) DecryptCalls(stub func(*kms.DecryptInput) (*kms.DecryptOutput, error)) {
	fake.decryptMutex.Lock()
	defer fake.decryptMutex.Unlock()
	fake.DecryptStub = stub
}

func (fake *FakeKMSClient) DecryptArgsForCall(i int) *kms.DecryptInput {
	fake.decryptMutex.RLock()
	defer fake.decryptMutex.RUnlock()
	argsForCall := fake.decryptArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeKMSClient) DecryptReturns(result1 *kms.DecryptOutput, result2 error) {
	fake.decryptMutex.Lock()
	defer fake.decryptMutex.Unlock()
	fake.DecryptStub = nil
	fake.decryptReturns = struct {
		result1 *kms.DecryptOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeKMSClient) DecryptReturnsOnCall(i int, result1 *kms.DecryptOutput, result2 error) {
	fake.decryptMutex.Lock()
	defer fake.decryptMutex.Unlock()
	fake.DecryptStub = nil
	if fake.decryptReturnsOnCall == nil {
		fake.decryptReturnsOnCall = make(map[int]struct {
			result1 *kms.DecryptOutput
			result2 error
		})
	}
	fake.decryptReturnsOnCall[i] = struct {
		result1 *kms.DecryptOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeKMSClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.decryptMutex.RLock()
	defer fake.decryptMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeKMSClient) recordInvocation(key string, args []interface{}) {
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

var _ environment.KMSClient = new(FakeKMSClient)