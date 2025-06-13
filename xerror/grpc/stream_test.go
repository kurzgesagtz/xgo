package grpc

import (
	"context"
	"errors"
	"io"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// mockClientStream is a mock implementation of grpc.ClientStream for testing
type mockClientStream struct {
	headerErr  error
	recvMsgErr error
	sendMsgErr error
	closeSendErr error
	trailer    metadata.MD
	isServerStream bool
}

func (m *mockClientStream) Header() (metadata.MD, error) {
	return metadata.MD{}, m.headerErr
}

func (m *mockClientStream) Trailer() metadata.MD {
	return m.trailer
}

func (m *mockClientStream) CloseSend() error {
	return m.closeSendErr
}

func (m *mockClientStream) Context() context.Context {
	return context.Background()
}

func (m *mockClientStream) SendMsg(msg interface{}) error {
	return m.sendMsgErr
}

func (m *mockClientStream) RecvMsg(msg interface{}) error {
	return m.recvMsgErr
}

func TestNewStreamClientWrapper(t *testing.T) {
	mockStream := &mockClientStream{}
	mockDesc := &grpc.StreamDesc{ServerStreams: true}

	wrapper := NewStreamClientWrapper(mockStream, mockDesc)

	if wrapper == nil {
		t.Fatal("Expected non-nil wrapper")
	}

	// Check that the wrapper implements the StreamClientWrapper interface
	_, ok := wrapper.(StreamClientWrapper)
	if !ok {
		t.Error("Wrapper does not implement StreamClientWrapper interface")
	}
}

func TestSetErrorTranslator(t *testing.T) {
	mockStream := &mockClientStream{}
	mockDesc := &grpc.StreamDesc{}

	wrapper := NewStreamClientWrapper(mockStream, mockDesc)

	// Set error translator
	translator := func(err error, trailer metadata.MD) error {
		return errors.New("translated error")
	}

	wrapper.SetErrorTranslator(translator)

	// No direct way to test the translator was set, but we can test it indirectly
	// by causing an error and checking the translation in other tests
}

func TestSetOnFinished(t *testing.T) {
	mockStream := &mockClientStream{}
	mockDesc := &grpc.StreamDesc{}

	wrapper := NewStreamClientWrapper(mockStream, mockDesc)

	// Set onFinished callback
	onFinished := func(err error) {
		// In a real scenario, this would do something with the error
	}

	wrapper.SetOnFinished(onFinished)

	// No direct way to test the callback was set, but we can test it indirectly
	// by causing an error and checking the callback in other tests
}

func TestHeader(t *testing.T) {
	tests := []struct {
		name           string
		headerErr      error
		expectTranslation bool
		expectCallback bool
	}{
		{
			name:           "no error",
			headerErr:      nil,
			expectTranslation: false,
			expectCallback: false,
		},
		{
			name:           "with error",
			headerErr:      errors.New("header error"),
			expectTranslation: true,
			expectCallback: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStream := &mockClientStream{
				headerErr: tc.headerErr,
			}
			mockDesc := &grpc.StreamDesc{}

			wrapper := NewStreamClientWrapper(mockStream, mockDesc)

			var translatorCalled bool
			var callbackCalled bool

			wrapper.SetErrorTranslator(func(err error, trailer metadata.MD) error {
				translatorCalled = true
				return errors.New("translated error")
			})

			wrapper.SetOnFinished(func(err error) {
				callbackCalled = true
			})

			_, err := wrapper.Header()

			if tc.headerErr == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if tc.headerErr != nil && err == nil {
				t.Error("Expected error, got nil")
			}

			if tc.expectTranslation && !translatorCalled {
				t.Error("Expected error translator to be called")
			}

			if !tc.expectTranslation && translatorCalled {
				t.Error("Expected error translator not to be called")
			}

			if tc.expectCallback && !callbackCalled {
				t.Error("Expected onFinished callback to be called")
			}

			if !tc.expectCallback && callbackCalled {
				t.Error("Expected onFinished callback not to be called")
			}
		})
	}
}

func TestRecvMsg(t *testing.T) {
	tests := []struct {
		name           string
		recvMsgErr     error
		isServerStream bool
		expectTranslation bool
		expectCallback bool
	}{
		{
			name:           "no error, client stream",
			recvMsgErr:     nil,
			isServerStream: false,
			expectTranslation: false,
			expectCallback: true, // Callback is called for client streams when no error
		},
		{
			name:           "no error, server stream",
			recvMsgErr:     nil,
			isServerStream: true,
			expectTranslation: false,
			expectCallback: false, // Callback is not called for server streams when no error
		},
		{
			name:           "EOF error",
			recvMsgErr:     io.EOF,
			isServerStream: false,
			expectTranslation: false,
			expectCallback: true, // Callback is called with nil error for EOF
		},
		{
			name:           "other error",
			recvMsgErr:     errors.New("recv error"),
			isServerStream: false,
			expectTranslation: true,
			expectCallback: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStream := &mockClientStream{
				recvMsgErr: tc.recvMsgErr,
			}
			mockDesc := &grpc.StreamDesc{
				ServerStreams: tc.isServerStream,
			}

			wrapper := NewStreamClientWrapper(mockStream, mockDesc)

			var translatorCalled bool
			var callbackCalled bool
			var callbackErr error

			wrapper.SetErrorTranslator(func(err error, trailer metadata.MD) error {
				translatorCalled = true
				return errors.New("translated error")
			})

			wrapper.SetOnFinished(func(err error) {
				callbackCalled = true
				callbackErr = err
			})

			err := wrapper.RecvMsg(nil)

			if tc.recvMsgErr == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if tc.recvMsgErr != nil && err == nil {
				t.Error("Expected error, got nil")
			}

			if tc.expectTranslation && !translatorCalled {
				t.Error("Expected error translator to be called")
			}

			if !tc.expectTranslation && translatorCalled {
				t.Error("Expected error translator not to be called")
			}

			if tc.expectCallback && !callbackCalled {
				t.Error("Expected onFinished callback to be called")
			}

			if !tc.expectCallback && callbackCalled {
				t.Error("Expected onFinished callback not to be called")
			}

			// For EOF, callback should be called with nil error
			if tc.recvMsgErr == io.EOF && callbackErr != nil {
				t.Errorf("Expected callback with nil error for EOF, got %v", callbackErr)
			}
		})
	}
}

func TestSendMsg(t *testing.T) {
	tests := []struct {
		name           string
		sendMsgErr     error
		expectTranslation bool
		expectCallback bool
	}{
		{
			name:           "no error",
			sendMsgErr:     nil,
			expectTranslation: false,
			expectCallback: false,
		},
		{
			name:           "with error",
			sendMsgErr:     errors.New("send error"),
			expectTranslation: true,
			expectCallback: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStream := &mockClientStream{
				sendMsgErr: tc.sendMsgErr,
			}
			mockDesc := &grpc.StreamDesc{}

			wrapper := NewStreamClientWrapper(mockStream, mockDesc)

			var translatorCalled bool
			var callbackCalled bool

			wrapper.SetErrorTranslator(func(err error, trailer metadata.MD) error {
				translatorCalled = true
				return errors.New("translated error")
			})

			wrapper.SetOnFinished(func(err error) {
				callbackCalled = true
			})

			err := wrapper.SendMsg(nil)

			if tc.sendMsgErr == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if tc.sendMsgErr != nil && err == nil {
				t.Error("Expected error, got nil")
			}

			if tc.expectTranslation && !translatorCalled {
				t.Error("Expected error translator to be called")
			}

			if !tc.expectTranslation && translatorCalled {
				t.Error("Expected error translator not to be called")
			}

			if tc.expectCallback && !callbackCalled {
				t.Error("Expected onFinished callback to be called")
			}

			if !tc.expectCallback && callbackCalled {
				t.Error("Expected onFinished callback not to be called")
			}
		})
	}
}

func TestCloseSend(t *testing.T) {
	tests := []struct {
		name           string
		closeSendErr   error
		expectTranslation bool
		expectCallback bool
	}{
		{
			name:           "no error",
			closeSendErr:   nil,
			expectTranslation: false,
			expectCallback: false,
		},
		{
			name:           "with error",
			closeSendErr:   errors.New("close error"),
			expectTranslation: true,
			expectCallback: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStream := &mockClientStream{
				closeSendErr: tc.closeSendErr,
			}
			mockDesc := &grpc.StreamDesc{}

			wrapper := NewStreamClientWrapper(mockStream, mockDesc)

			var translatorCalled bool
			var callbackCalled bool

			wrapper.SetErrorTranslator(func(err error, trailer metadata.MD) error {
				translatorCalled = true
				return errors.New("translated error")
			})

			wrapper.SetOnFinished(func(err error) {
				callbackCalled = true
			})

			err := wrapper.CloseSend()

			if tc.closeSendErr == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if tc.closeSendErr != nil && err == nil {
				t.Error("Expected error, got nil")
			}

			if tc.expectTranslation && !translatorCalled {
				t.Error("Expected error translator to be called")
			}

			if !tc.expectTranslation && translatorCalled {
				t.Error("Expected error translator not to be called")
			}

			if tc.expectCallback && !callbackCalled {
				t.Error("Expected onFinished callback to be called")
			}

			if !tc.expectCallback && callbackCalled {
				t.Error("Expected onFinished callback not to be called")
			}
		})
	}
}
