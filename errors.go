// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

var (
	EInvalidCommandProvided  = "invalid command provided. Command cannot contain \\r and/or \\n. Provided command is: %s"
	ECouldNotReadMIMEHeaders = "error while reading MIME headers: %s"
	EInvalidContentLength    = "unable to get size of content-length: %s"
	EUnsuccessfulReply       = "got error while reading from reply command: %s"
	ECouldNotReadyBody       = "got error while reading reader body: %s"
	EUnsupportedMessageType  = "unsupported message type! We got '%s'. Supported types are: %v "
	ECouldNotDecode          = "could not decode/unescape message: %s"
	ECouldNotStartListener   = "got error while attempting to start listener: %s"
	EListenerConnection      = "listener connection error: %s"
	EInvalidServerAddr       = "please make sure to pass along valid address. You've passed: \"%s\""
	EUnexpectedAuthHeader    = "expected auth/request content type. Got %s"
	EInvalidPassword         = "could not authenticate against freeSWITCH with provided password: %s"
	ECouldNotCreateMessage   = "error while creating new message: %s"
	ECouldNotSendEvent       = "must send at least one event header, detected `%d` header"
)
