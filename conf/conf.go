// Copyright 2016 Alexander Zaytsev <thebestzorro@yandex.ru>
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

// Package conf contains common types and methods.
package conf

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
)

var (
	// ErrNotFound is error for not found handlers
	ErrNotFound = errors.New("not found")
	// Logger is common logger
	Logger = logger{
		Info:  log.New(os.Stdout, "INFO [ewrc]: ", log.Ldate|log.Ltime|log.Lshortfile),
		Debug: log.New(ioutil.Discard, "DEBUG [ewrc]: ", log.Ldate|log.Lmicroseconds|log.Llongfile),
	}
)

// logger is common logger structure.
type logger struct {
	Debug *log.Logger
	Info  *log.Logger
}
