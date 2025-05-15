// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Main connection against ESL - Gotta add more description here
type SocketConnection struct {
	net.Conn
	err chan error
	m   chan *Message
	mtx sync.Mutex
}

// Dial - Will establish timed-out dial against specified address. In this case, it will be FreeSWITCH server
func (sc *SocketConnection) Dial(network string, addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

// Send - Will send raw message to open net connection
func (sc *SocketConnection) Send(cmd string) error {

	if strings.Contains(cmd, "\r\n") {
		return fmt.Errorf(EInvalidCommandProvided, cmd)
	}

	// lock mutex
	sc.mtx.Lock()
	defer sc.mtx.Unlock()

	_, err := io.WriteString(sc, cmd)
	if err != nil {
		return err
	}

	_, err = io.WriteString(sc, "\r\n\r\n")
	if err != nil {
		return err
	}

	return nil
}

// SendMany - Will loop against passed commands and return 1st error if error happens
func (sc *SocketConnection) SendMany(cmds []string) error {

	for _, cmd := range cmds {
		if err := sc.Send(cmd); err != nil {
			return err
		}
	}

	return nil
}

// SendEvent - Will loop against passed event headers
func (sc *SocketConnection) SendEvent(eventHeaders []string) error {
	if len(eventHeaders) <= 0 {
		return fmt.Errorf(ECouldNotSendEvent, len(eventHeaders))
	}

	// lock mutex to prevent event headers from conflicting
	sc.mtx.Lock()
	defer sc.mtx.Unlock()

	_, err := io.WriteString(sc, "sendevent ")
	if err != nil {
		return err
	}

	for _, eventHeader := range eventHeaders {
		_, err := io.WriteString(sc, eventHeader)
		if err != nil {
			return err
		}

		_, err = io.WriteString(sc, "\r\n")
		if err != nil {
			return err
		}

	}

	_, err = io.WriteString(sc, "\r\n")
	if err != nil {
		return err
	}

	return nil
}

// Execute - Helper to execute commands with its args and sync/async mode
func (sc *SocketConnection) Execute(command, args string, sync bool) (m *Message, err error) {
	return sc.SendMsg(map[string]string{
		"call-command":     "execute",
		"execute-app-name": command,
		"execute-app-arg":  args,
		"event-lock":       strconv.FormatBool(sync),
	}, "", "")
}

// ExecuteUUID - Helper to execute uuid specific commands with its args and sync/async mode
func (sc *SocketConnection) ExecuteUUID(uuid string, command string, args string, sync bool) (m *Message, err error) {
	return sc.SendMsg(map[string]string{
		"call-command":     "execute",
		"execute-app-name": command,
		"execute-app-arg":  args,
		"event-lock":       strconv.FormatBool(sync),
	}, uuid, "")
}

// SendMsg - Basically this func will send message to the opened connection
func (sc *SocketConnection) SendMsg(msg map[string]string, uuid, data string) (m *Message, err error) {

	b := bytes.NewBufferString("sendmsg")

	if uuid != "" {
		if strings.Contains(uuid, "\r\n") {
			return nil, fmt.Errorf(EInvalidCommandProvided, msg)
		}

		b.WriteString(" " + uuid)
	}

	b.WriteString("\n")

	for k, v := range msg {
		if strings.Contains(k, "\r\n") {
			return nil, fmt.Errorf(EInvalidCommandProvided, msg)
		}

		if v != "" {
			if strings.Contains(v, "\r\n") {
				return nil, fmt.Errorf(EInvalidCommandProvided, msg)
			}

			b.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
	}

	b.WriteString("\n")

	if msg["content-length"] != "" && data != "" {
		b.WriteString(data)
	}

	// lock mutex
	sc.mtx.Lock()
	_, err = b.WriteTo(sc)
	if err != nil {
		sc.mtx.Unlock()
		return nil, err
	}
	sc.mtx.Unlock()

	select {
	case err := <-sc.err:
		return nil, err
	case m := <-sc.m:
		return m, nil
	}
}

// OriginatorAdd - Will return originator address known as net.RemoteAddr()
// This will actually be a FreeSWITCH address
func (sc *SocketConnection) OriginatorAddr() net.Addr {
	return sc.RemoteAddr()
}

// ReadMessage - Will read message from channels and return them back accordingly.
//  If error is received, error will be returned. If not, message will be returned back!
func (sc *SocketConnection) ReadMessage() (*Message, error) {
	Debug("Waiting for connection message to be received ...")

	select {
	case err := <-sc.err:
		return nil, err
	case msg := <-sc.m:
		return msg, nil
	}
}

// Handle - Will handle new messages and close connection when there are no messages left to process
func (sc *SocketConnection) Handle() {

	done := make(chan bool)

	rbuf := bufio.NewReaderSize(sc, ReadBufferSize)

	go func() {
		for {
			msg, err := newMessage(rbuf, true)

			if err != nil {
				sc.err <- err
				done <- true
				break
			}

			sc.m <- msg
		}
	}()

	<-done

	// Closing the connection now as there's nothing left to do ...
	_ = sc.Close()
}

// Close - Will close down net connection and return error if error happen
func (sc *SocketConnection) Close() error {
	if err := sc.Conn.Close(); err != nil {
		return err
	}

	return nil
}
