package grpc

import (
	"errors"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type StreamClientErrorTranslator func(err error, trailer metadata.MD) error
type StreamClientOnFinished func(err error)

var _ grpc.ClientStream = &streamClientWrapper{}

type StreamClientWrapper interface {
	grpc.ClientStream
	SetErrorTranslator(fn StreamClientErrorTranslator)
	SetOnFinished(fn StreamClientOnFinished)
}

type streamClientWrapper struct {
	grpc.ClientStream
	desc *grpc.StreamDesc

	// xerror translator
	errorTranslator StreamClientErrorTranslator
	finished        StreamClientOnFinished
}

func (sw *streamClientWrapper) SetErrorTranslator(fn StreamClientErrorTranslator) {
	sw.errorTranslator = fn
}
func (sw *streamClientWrapper) SetOnFinished(fn StreamClientOnFinished) {
	sw.finished = fn
}

func (sw *streamClientWrapper) Header() (metadata.MD, error) {
	md, err := sw.ClientStream.Header()
	if err != nil {
		newErr := sw.errorConvertor(err)
		sw.callFinished(err)
		// xlog and metric
		return md, newErr
	}
	return md, nil
}

func (sw *streamClientWrapper) RecvMsg(m interface{}) error {
	err := sw.ClientStream.RecvMsg(m)
	if err != nil {
		if errors.Is(err, io.EOF) {
			sw.callFinished(nil)
		} else {
			err = sw.errorConvertor(err)
			sw.callFinished(err)
		}
	} else if !sw.desc.ServerStreams {
		sw.callFinished(nil)
	}
	return err
}

func (sw *streamClientWrapper) SendMsg(m interface{}) error {
	err := sw.ClientStream.SendMsg(m)
	if err != nil {
		err = sw.errorConvertor(err)
		sw.callFinished(err)
	}
	return err
}

func (sw *streamClientWrapper) CloseSend() error {
	err := sw.ClientStream.CloseSend()
	if err != nil {
		err = sw.errorConvertor(err)
		sw.callFinished(err)
	}
	return err
}

func (sw *streamClientWrapper) errorConvertor(err error) error {
	if err != nil && sw.errorTranslator != nil {
		return sw.errorTranslator(err, sw.Trailer())
	}
	return err
}

func (sw *streamClientWrapper) callFinished(err error) {
	if sw.finished != nil {
		sw.finished(err)
	}
}

func NewStreamClientWrapper(cs grpc.ClientStream, desc *grpc.StreamDesc) StreamClientWrapper {
	return &streamClientWrapper{
		ClientStream: cs,
		desc:         desc,
	}
}
