// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	environment "github.com/telia-oss/aws-env"
)

type FakeSMClient struct {
	GetSecretValueStub        func(*secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
	getSecretValueMutex       sync.RWMutex
	getSecretValueArgsForCall []struct {
		arg1 *secretsmanager.GetSecretValueInput
	}
	getSecretValueReturns struct {
		result1 *secretsmanager.GetSecretValueOutput
		result2 error
	}
	getSecretValueReturnsOnCall map[int]struct {
		result1 *secretsmanager.GetSecretValueOutput
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeSMClient) GetSecretValue(arg1 *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	fake.getSecretValueMutex.Lock()
	ret, specificReturn := fake.getSecretValueReturnsOnCall[len(fake.getSecretValueArgsForCall)]
	fake.getSecretValueArgsForCall = append(fake.getSecretValueArgsForCall, struct {
		arg1 *secretsmanager.GetSecretValueInput
	}{arg1})
	stub := fake.GetSecretValueStub
	fakeReturns := fake.getSecretValueReturns
	fake.recordInvocation("GetSecretValue", []interface{}{arg1})
	fake.getSecretValueMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeSMClient) GetSecretValueCallCount() int {
	fake.getSecretValueMutex.RLock()
	defer fake.getSecretValueMutex.RUnlock()
	return len(fake.getSecretValueArgsForCall)
}

func (fake *FakeSMClient) GetSecretValueCalls(stub func(*secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)) {
	fake.getSecretValueMutex.Lock()
	defer fake.getSecretValueMutex.Unlock()
	fake.GetSecretValueStub = stub
}

func (fake *FakeSMClient) GetSecretValueArgsForCall(i int) *secretsmanager.GetSecretValueInput {
	fake.getSecretValueMutex.RLock()
	defer fake.getSecretValueMutex.RUnlock()
	argsForCall := fake.getSecretValueArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeSMClient) GetSecretValueReturns(result1 *secretsmanager.GetSecretValueOutput, result2 error) {
	fake.getSecretValueMutex.Lock()
	defer fake.getSecretValueMutex.Unlock()
	fake.GetSecretValueStub = nil
	fake.getSecretValueReturns = struct {
		result1 *secretsmanager.GetSecretValueOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeSMClient) GetSecretValueReturnsOnCall(i int, result1 *secretsmanager.GetSecretValueOutput, result2 error) {
	fake.getSecretValueMutex.Lock()
	defer fake.getSecretValueMutex.Unlock()
	fake.GetSecretValueStub = nil
	if fake.getSecretValueReturnsOnCall == nil {
		fake.getSecretValueReturnsOnCall = make(map[int]struct {
			result1 *secretsmanager.GetSecretValueOutput
			result2 error
		})
	}
	fake.getSecretValueReturnsOnCall[i] = struct {
		result1 *secretsmanager.GetSecretValueOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeSMClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getSecretValueMutex.RLock()
	defer fake.getSecretValueMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeSMClient) recordInvocation(key string, args []interface{}) {
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

var _ environment.SMClient = new(FakeSMClient)