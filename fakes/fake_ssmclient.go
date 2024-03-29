// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/aws/aws-sdk-go/service/ssm"
	environment "github.com/telia-oss/aws-env"
)

type FakeSSMClient struct {
	GetParameterStub        func(*ssm.GetParameterInput) (*ssm.GetParameterOutput, error)
	getParameterMutex       sync.RWMutex
	getParameterArgsForCall []struct {
		arg1 *ssm.GetParameterInput
	}
	getParameterReturns struct {
		result1 *ssm.GetParameterOutput
		result2 error
	}
	getParameterReturnsOnCall map[int]struct {
		result1 *ssm.GetParameterOutput
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeSSMClient) GetParameter(arg1 *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	fake.getParameterMutex.Lock()
	ret, specificReturn := fake.getParameterReturnsOnCall[len(fake.getParameterArgsForCall)]
	fake.getParameterArgsForCall = append(fake.getParameterArgsForCall, struct {
		arg1 *ssm.GetParameterInput
	}{arg1})
	stub := fake.GetParameterStub
	fakeReturns := fake.getParameterReturns
	fake.recordInvocation("GetParameter", []interface{}{arg1})
	fake.getParameterMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeSSMClient) GetParameterCallCount() int {
	fake.getParameterMutex.RLock()
	defer fake.getParameterMutex.RUnlock()
	return len(fake.getParameterArgsForCall)
}

func (fake *FakeSSMClient) GetParameterCalls(stub func(*ssm.GetParameterInput) (*ssm.GetParameterOutput, error)) {
	fake.getParameterMutex.Lock()
	defer fake.getParameterMutex.Unlock()
	fake.GetParameterStub = stub
}

func (fake *FakeSSMClient) GetParameterArgsForCall(i int) *ssm.GetParameterInput {
	fake.getParameterMutex.RLock()
	defer fake.getParameterMutex.RUnlock()
	argsForCall := fake.getParameterArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeSSMClient) GetParameterReturns(result1 *ssm.GetParameterOutput, result2 error) {
	fake.getParameterMutex.Lock()
	defer fake.getParameterMutex.Unlock()
	fake.GetParameterStub = nil
	fake.getParameterReturns = struct {
		result1 *ssm.GetParameterOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeSSMClient) GetParameterReturnsOnCall(i int, result1 *ssm.GetParameterOutput, result2 error) {
	fake.getParameterMutex.Lock()
	defer fake.getParameterMutex.Unlock()
	fake.GetParameterStub = nil
	if fake.getParameterReturnsOnCall == nil {
		fake.getParameterReturnsOnCall = make(map[int]struct {
			result1 *ssm.GetParameterOutput
			result2 error
		})
	}
	fake.getParameterReturnsOnCall[i] = struct {
		result1 *ssm.GetParameterOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeSSMClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getParameterMutex.RLock()
	defer fake.getParameterMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeSSMClient) recordInvocation(key string, args []interface{}) {
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

var _ environment.SSMClient = new(FakeSSMClient)
