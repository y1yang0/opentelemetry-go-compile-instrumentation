// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package net

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
)

type testRequest struct{}

type testResponse struct{}

type netAttrsGetter struct{}

func (netAttrsGetter) GetURLScheme(_ testRequest) string {
	return "test"
}

func (netAttrsGetter) GetURLPath(_ testRequest) string {
	return "test"
}

func (netAttrsGetter) GetURLQuery(_ testRequest) string {
	return "test"
}

func (netAttrsGetter) GetNetworkType(_ testRequest, _ testResponse) string {
	return "test"
}

func (netAttrsGetter) GetNetworkTransport(_ testRequest, _ testResponse) string {
	return "test"
}

func (netAttrsGetter) GetNetworkProtocolName(_ testRequest, _ testResponse) string {
	return "test"
}

func (netAttrsGetter) GetNetworkProtocolVersion(_ testRequest, _ testResponse) string {
	return "test"
}

func (netAttrsGetter) GetNetworkLocalInetAddress(_ testRequest, _ testResponse) string {
	return "test"
}

func (netAttrsGetter) GetNetworkLocalPort(_ testRequest, _ testResponse) int {
	return 8080
}

func (netAttrsGetter) GetNetworkPeerInetAddress(_ testRequest, _ testResponse) string {
	return "test"
}

func (netAttrsGetter) GetNetworkPeerPort(_ testRequest, _ testResponse) int {
	return 8080
}

type mockClientAttributesGetter struct {
	address string
	port    int
}

func (m *mockClientAttributesGetter) GetClientAddress(_ string) string {
	return m.address
}

func (m *mockClientAttributesGetter) GetClientPort(_ string) int {
	return m.port
}

type mockAddressAndPortExtractor struct {
	address string
	port    int
}

func (m *mockAddressAndPortExtractor) Extract(_ string) AddressAndPort {
	return AddressAndPort{
		Address: m.address,
		Port:    m.port,
	}
}

func TestNoopAddressAndPortExtractorExtractShouldReturnConstantValue(t *testing.T) {
	extractor := &NoopAddressAndPortExtractor[string]{}
	actual := extractor.Extract("any request")
	assert.Equal(t, AddressAndPort{}, actual)
}

type MockGetter struct {
	Address string
	Port    int
}

func (m *MockGetter) GetServerAddress(_ any) string {
	return m.Address
}

func (m *MockGetter) GetServerPort(_ any) int {
	return m.Port
}

type MockFallbackExtractor struct {
	Extracted AddressAndPort
}

func (m *MockFallbackExtractor) Extract(_ any) AddressAndPort {
	return m.Extracted
}

type MockAddressAndPortExtractor[REQUEST AddressAndPort] struct{}

func (*MockAddressAndPortExtractor[REQUEST]) Extract(request AddressAndPort) AddressAndPort {
	return request
}

type MockClientAttributesGetter struct{}

func (MockClientAttributesGetter) GetClientAddress(request AddressAndPort) string {
	return request.Address
}

func (MockClientAttributesGetter) GetClientPort(request AddressAndPort) int {
	return request.Port
}

type MockServerAttributesGetter struct{}

func (MockServerAttributesGetter) GetServerAddress(request AddressAndPort) string {
	return request.Address
}

func (MockServerAttributesGetter) GetServerPort(request AddressAndPort) int {
	return request.Port
}

func TestClientAddressAndPortExtractorExtract(t *testing.T) {
	mockGetter := &mockClientAttributesGetter{
		address: "192.168.1.1",
		port:    8080,
	}

	mockFallback := &mockAddressAndPortExtractor{
		address: "127.0.0.1",
		port:    9090,
	}

	extractor := &ClientAddressAndPortExtractor[string]{
		getter:            mockGetter,
		fallbackExtractor: mockFallback,
	}

	result := extractor.Extract("testRequest")
	if result.Address != "192.168.1.1" || result.Port != 8080 {
		t.Errorf("Expected address and port to be '192.168.1.1:8080', got '%s:%d'", result.Address, result.Port)
	}

	mockGetter.address = ""
	mockGetter.port = 0
	result = extractor.Extract("testRequest")
	if result.Address != "127.0.0.1" || result.Port != 9090 {
		t.Errorf("Expected fallback address and port to be '127.0.0.1:9090', got '%s:%d'", result.Address, result.Port)
	}
}

func TestExtractWhenGetterReturnsDefaultsShouldUseFallbackExtractor(t *testing.T) {
	mockGetter := &MockGetter{Address: "", Port: 0}
	mockFallbackExtractor := &MockFallbackExtractor{Extracted: AddressAndPort{Address: "fallbackAddress", Port: 8080}}
	extractor := ServerAddressAndPortExtractor[any]{
		getter:            mockGetter,
		fallbackExtractor: mockFallbackExtractor,
	}
	result := extractor.Extract("testRequest")
	assert.Equal(t, "fallbackAddress", result.Address)
	assert.Equal(t, 8080, result.Port)
}

func TestExtractWhenGetterReturnsNonDefaultsShouldReturnDirectly(t *testing.T) {
	mockGetter := &MockGetter{Address: "directAddress", Port: 9090}
	mockFallbackExtractor := &MockFallbackExtractor{Extracted: AddressAndPort{Address: "fallbackAddress", Port: 8080}}
	extractor := ServerAddressAndPortExtractor[any]{
		getter:            mockGetter,
		fallbackExtractor: mockFallbackExtractor,
	}
	result := extractor.Extract("testRequest")
	assert.Equal(t, "directAddress", result.Address)
	assert.Equal(t, 9090, result.Port)
}

func TestInternalClientAttributesExtractorOnStart(t *testing.T) {
	tests := []struct {
		name           string
		address        string
		port           int
		capturePort    bool
		expectedResult []attribute.KeyValue
	}{
		{
			name:           "AddressEmpty_NoAttributesAdded",
			address:        "",
			port:           0,
			capturePort:    false,
			expectedResult: make([]attribute.KeyValue, 0),
		},
		{
			name:        "AddressNotEmpty_AttributeAdded",
			address:     "192.0.2.1",
			port:        0,
			capturePort: false,
			expectedResult: []attribute.KeyValue{
				{
					Key:   semconv.ClientAddressKey,
					Value: attribute.StringValue("192.0.2.1"),
				},
			},
		},
		{
			name:        "CapturePortTrue_PortAdded",
			address:     "192.0.2.1",
			port:        8080,
			capturePort: true,
			expectedResult: []attribute.KeyValue{
				{
					Key:   semconv.ClientAddressKey,
					Value: attribute.StringValue("192.0.2.1"),
				},
				{
					Key:   semconv.ClientPortKey,
					Value: attribute.IntValue(8080),
				},
			},
		},
		{
			name:        "CapturePortTrue_PortZero_NoPortAttribute",
			address:     "192.0.2.1",
			port:        0,
			capturePort: true,
			expectedResult: []attribute.KeyValue{
				{
					Key:   semconv.ClientAddressKey,
					Value: attribute.StringValue("192.0.2.1"),
				},
			},
		},
		{
			name:        "CapturePortFalse_NoPortAttribute",
			address:     "192.0.2.1",
			port:        8080,
			capturePort: false,
			expectedResult: []attribute.KeyValue{
				{
					Key:   semconv.ClientAddressKey,
					Value: attribute.StringValue("192.0.2.1"),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.capturePort {
				ie := CreateClientAttributesExtractor[AddressAndPort, AddressAndPort](MockClientAttributesGetter{})
				attributes, _ := ie.OnStart(context.TODO(), make([]attribute.KeyValue, 0), AddressAndPort{
					Address: test.address,
					Port:    test.port,
				})
				assert.Equal(t, test.expectedResult, attributes)
				attributes, _ = ie.OnEnd(
					context.TODO(),
					make([]attribute.KeyValue, 0),
					AddressAndPort{},
					AddressAndPort{},
					nil,
				)
				assert.Equal(t, make([]attribute.KeyValue, 0), attributes)
			} else {
				ie := &InternalClientAttributesExtractor[AddressAndPort]{
					addressAndPortExtractor: &MockAddressAndPortExtractor[AddressAndPort]{},
					capturePort:             test.capturePort,
				}
				attributes, _ := ie.OnStart(context.TODO(), make([]attribute.KeyValue, 0), AddressAndPort{
					Address: test.address,
					Port:    test.port,
				})
				assert.Equal(t, test.expectedResult, attributes)
			}
		})
	}
}

func TestInternalServerAttributesExtractorOnStart(t *testing.T) {
	tests := []struct {
		name           string
		address        string
		port           int
		capturePort    bool
		expectedResult []attribute.KeyValue
	}{
		{
			name:           "AddressEmpty_NoAttributesAdded",
			address:        "",
			port:           0,
			capturePort:    false,
			expectedResult: make([]attribute.KeyValue, 0),
		},
		{
			name:        "AddressNotEmpty_AttributeAdded",
			address:     "192.0.2.1",
			port:        0,
			capturePort: false,
			expectedResult: []attribute.KeyValue{
				{
					Key:   semconv.ServerAddressKey,
					Value: attribute.StringValue("192.0.2.1"),
				},
			},
		},
		{
			name:        "CapturePortTrue_PortAdded",
			address:     "192.0.2.1",
			port:        8080,
			capturePort: true,
			expectedResult: []attribute.KeyValue{
				{
					Key:   semconv.ServerAddressKey,
					Value: attribute.StringValue("192.0.2.1"),
				},
				{
					Key:   semconv.ServerPortKey,
					Value: attribute.IntValue(8080),
				},
			},
		},
		{
			name:        "CapturePortTrue_PortZero_NoPortAttribute",
			address:     "192.0.2.1",
			port:        0,
			capturePort: true,
			expectedResult: []attribute.KeyValue{
				{
					Key:   semconv.ServerAddressKey,
					Value: attribute.StringValue("192.0.2.1"),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ie := CreateServerAttributesExtractor[AddressAndPort, AddressAndPort](MockServerAttributesGetter{})
			attributes, _ := ie.OnStart(context.TODO(), make([]attribute.KeyValue, 0), AddressAndPort{
				Address: test.address,
				Port:    test.port,
			})
			assert.Equal(t, test.expectedResult, attributes)
			attributes, _ = ie.OnEnd(
				context.TODO(),
				make([]attribute.KeyValue, 0),
				AddressAndPort{},
				AddressAndPort{},
				nil,
			)
			assert.Equal(t, make([]attribute.KeyValue, 0), attributes)
		})
	}
}

func TestInternalNetworkAttributesExtractorOnStart(t *testing.T) {
	tests := []struct {
		name                         string
		captureProtocolAttributes    bool
		captureLocalSocketAttributes bool
		expectedResult               []attribute.KeyValue
	}{
		{
			name: "test_default",
			expectedResult: []attribute.KeyValue{
				{
					Key:   semconv.NetworkPeerAddressKey,
					Value: attribute.StringValue("test"),
				},
				{
					Key:   semconv.NetworkPeerPortKey,
					Value: attribute.IntValue(8080),
				},
			},
		},
		{
			name:                      "test_captureProtocolAttributes",
			captureProtocolAttributes: true,
			expectedResult: []attribute.KeyValue{
				{
					Key:   semconv.NetworkTransportKey,
					Value: attribute.StringValue("test"),
				},
				{
					Key:   semconv.NetworkTypeKey,
					Value: attribute.StringValue("test"),
				},
				{
					Key:   semconv.NetworkProtocolNameKey,
					Value: attribute.StringValue("test"),
				},
				{
					Key:   semconv.NetworkProtocolVersionKey,
					Value: attribute.StringValue("test"),
				},
				{
					Key:   semconv.NetworkPeerAddressKey,
					Value: attribute.StringValue("test"),
				},
				{
					Key:   semconv.NetworkPeerPortKey,
					Value: attribute.IntValue(8080),
				},
			},
		},
		{
			name:                         "AddressEmpty_NoAttributesAdded",
			captureProtocolAttributes:    true,
			captureLocalSocketAttributes: true,
			expectedResult: []attribute.KeyValue{
				{
					Key:   semconv.NetworkTransportKey,
					Value: attribute.StringValue("test"),
				},
				{
					Key:   semconv.NetworkTypeKey,
					Value: attribute.StringValue("test"),
				},
				{
					Key:   semconv.NetworkProtocolNameKey,
					Value: attribute.StringValue("test"),
				},
				{
					Key:   semconv.NetworkProtocolVersionKey,
					Value: attribute.StringValue("test"),
				},
				{
					Key:   semconv.NetworkLocalAddressKey,
					Value: attribute.StringValue("test"),
				},
				{
					Key:   semconv.NetworkLocalPortKey,
					Value: attribute.IntValue(8080),
				},
				{
					Key:   semconv.NetworkPeerAddressKey,
					Value: attribute.StringValue("test"),
				},
				{
					Key:   semconv.NetworkPeerPortKey,
					Value: attribute.IntValue(8080),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.captureProtocolAttributes && test.captureLocalSocketAttributes {
				ie := CreateNetworkAttributesExtractor[testRequest, testResponse](netAttrsGetter{})
				attributes, _ := ie.OnEnd(
					context.TODO(),
					make([]attribute.KeyValue, 0),
					testRequest{},
					testResponse{},
					nil,
				)
				assert.Equal(t, test.expectedResult, attributes)
				attributes, _ = ie.OnStart(context.TODO(), make([]attribute.KeyValue, 0), testRequest{})
				assert.Equal(t, make([]attribute.KeyValue, 0), attributes)
			} else {
				ie := &InternalNetworkAttributesExtractor[testRequest, testResponse]{
					getter:                       netAttrsGetter{},
					captureProtocolAttributes:    test.captureProtocolAttributes,
					captureLocalSocketAttributes: test.captureLocalSocketAttributes,
				}
				attributes, _ := ie.OnEnd(context.TODO(), make([]attribute.KeyValue, 0), testRequest{}, testResponse{})
				assert.Equal(t, test.expectedResult, attributes)
			}
		})
	}
}
