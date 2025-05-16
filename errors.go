// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

import "fmt"

var (
	ErrInvalidCommandProvided  = fmt.Errorf("invalid command provided")
	ErrCouldNotReadMIMEHeaders = fmt.Errorf("error while reading MIME headers")
	ErrInvalidContentLength    = fmt.Errorf("unable to get size of content-length")
	ErrUnsuccessfulReply       = fmt.Errorf("got error while reading from reply command")
	ErrCouldNotReadyBody       = fmt.Errorf("got error while reading reader body")
	ErrUnsupportedMessageType  = fmt.Errorf("unsupported message type")
	ErrCouldNotStartListener   = fmt.Errorf("got error while attempting to start listener")
	ErrInvalidServerAddr       = fmt.Errorf("invalid server address")
	ErrUnexpectedAuthHeader    = fmt.Errorf("unexpected auth/request content type")
	ErrInvalidPassword         = fmt.Errorf("invalid password")
	ErrCouldNotCreateMessage   = fmt.Errorf("error while creating new message")
	ErrCouldNotSendEvent       = fmt.Errorf("must send at least one event header")
)
