// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

type Logger interface {
	Debugf(message string, args ...interface{})
	Errorf(message string, args ...interface{})
	Noticef(message string, args ...interface{})
	Infof(message string, args ...interface{})
	Warningf(message string, args ...interface{})
}

var log Logger

func SetLogger(l Logger) {
	log = l
}

func Debug(message string, args ...interface{}) {
	if log == nil {
		return
	}

	log.Debugf(message, args...)
}

func Error(message string, args ...interface{}) {
	if log == nil {
		return
	}

	log.Errorf(message, args...)
}

func Notice(message string, args ...interface{}) {
	if log == nil {
		return
	}

	log.Noticef(message, args...)
}

func Info(message string, args ...interface{}) {
	if log == nil {
		return
	}

	log.Infof(message, args...)
}

func Warning(message string, args ...interface{}) {
	if log == nil {
		return
	}

	log.Warningf(message, args...)
}
