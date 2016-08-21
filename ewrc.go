// Copyright 2016 Alexander Zaytsev <thebestzorro@yandex.ru>
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

// Package main implements main methods of EWRC service.
package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/z0rr0/ewrc/conf"
)

const (
	defaulPort uint = 10001
	defaulDb        = "db.sqlite3"
	timeout         = 30 * time.Second

	// Name is a program name
	Name = "EWRC"
)

var (
	// Version is LUSS version
	Version = "0.0.1"
	// Revision is revision number
	Revision = "git:000000"
	// BuildDate is build date
	BuildDate = "2016-08-21_08:57:00UTC"
	// GoVersion is runtime Go language version
	GoVersion = runtime.Version()
)

func interrupt() error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	return fmt.Errorf("signal %v", <-c)
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("abnormal termination [%v]: %v\n", Version, r)
		}
	}()
	version := flag.Bool("v", false, "show version info")
	host := flag.String("h", "", "host")
	port := flag.Uint("p", defaulPort, "port")
	//database := flag.String("d", defaulDb, "database")
	//trace := flag.Bool("t", false, "turn on traces")

	flag.Parse()
	if *version {
		fmt.Printf("%v: %v\n\trevision: %v %v\n\tbuild date: %v\n", Name, Version, Revision, GoVersion, BuildDate)
		return
	}

	errChan := make(chan error)
	go func() {
		errChan <- interrupt()
	}()

	listener := net.JoinHostPort(*host, fmt.Sprint(*port))
	server := &http.Server{
		Addr:           listener,
		Handler:        http.DefaultServeMux,
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
		MaxHeaderBytes: 1 << 20,
		ErrorLog:       conf.Logger.Info,
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := "/"
		if r.URL.Path != url {
			url = strings.TrimRight(r.URL.Path, "/")
		}
		start, code := time.Now(), http.StatusOK
		defer func() {
			conf.Logger.Info.Printf("%-5v %v\t%-12v\t%v",
				r.Method,
				code,
				time.Since(start),
				r.URL.String(),
			)
		}()
		code = http.StatusNotFound
		http.NotFound(w, r)
	})
	conf.Logger.Info.Printf("\nListen %v\n", listener)
	go func() {
		errChan <- server.ListenAndServe()
	}()
	conf.Logger.Info.Printf("%v termination, reason[%v]: %v [%v]\n", Name, <-errChan, Version, Revision)

}
