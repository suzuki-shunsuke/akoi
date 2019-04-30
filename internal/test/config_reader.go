package test

// Don't edit this file.
// This file is generated by gomic 0.5.2.
// https://github.com/suzuki-shunsuke/gomic

import (
	testing "testing"

	test "github.com/suzuki-shunsuke/akoi/internal/domain"
	gomic "github.com/suzuki-shunsuke/gomic/gomic"
)

type (
	// ConfigReader is a mock.
	ConfigReader struct {
		t                      *testing.T
		name                   string
		callbackNotImplemented gomic.CallbackNotImplemented
		impl                   struct {
			Read func(p0 string) (test.Config, error)
		}
	}
)

// NewConfigReader returns ConfigReader .
func NewConfigReader(t *testing.T, cb gomic.CallbackNotImplemented) *ConfigReader {
	return &ConfigReader{
		t: t, name: "ConfigReader", callbackNotImplemented: cb}
}

// Read is a mock method.
func (mock ConfigReader) Read(p0 string) (test.Config, error) {
	methodName := "Read" // nolint: goconst
	if mock.impl.Read != nil {
		return mock.impl.Read(p0)
	}
	if mock.callbackNotImplemented != nil {
		mock.callbackNotImplemented(mock.t, mock.name, methodName)
	} else {
		gomic.DefaultCallbackNotImplemented(mock.t, mock.name, methodName)
	}
	return mock.fakeZeroRead(p0)
}

// SetFuncRead sets a method and returns the mock.
func (mock *ConfigReader) SetFuncRead(impl func(p0 string) (test.Config, error)) *ConfigReader {
	mock.impl.Read = impl
	return mock
}

// SetReturnRead sets a fake method.
func (mock *ConfigReader) SetReturnRead(r0 test.Config, r1 error) *ConfigReader {
	mock.impl.Read = func(string) (test.Config, error) {
		return r0, r1
	}
	return mock
}

// fakeZeroRead is a fake method which returns zero values.
func (mock ConfigReader) fakeZeroRead(p0 string) (test.Config, error) {
	var (
		r0 test.Config
		r1 error
	)
	return r0, r1
}
