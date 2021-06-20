package main

import (
  cmdline  "github.com/galdor/go-cmdline"
  context  "context"
  internal "github.com/gpenaud/needys-api-strategy/internal"
  log      "github.com/sirupsen/logrus"
  os       "os"
  signal   "os/signal"
  syscall  "syscall"
)

func registerCliConfiguration(a *internal.Application) {
  cmdline := cmdline.New()

  a.Config = &internal.Configuration{}

  // application configuration flags
  cmdline.AddOption("e", "environment", "ENVIRONMENT", "the current environment (development, integration, production)")
  cmdline.SetOptionDefault("environment", "production")

  cmdline.AddOption("v", "verbosity", "LEVEL", "verbosity for log-level (error, warning, info, debug)")
  cmdline.SetOptionDefault("verbosity", "info")

  cmdline.AddOption("l", "log-format", "FORMAT", "log format (text, json)")
  cmdline.SetOptionDefault("log-format", "unset")

  cmdline.AddFlag("", "log-healthcheck", "log healthcheck queries")

  // application server configuration flags
  cmdline.AddOption("", "server.host", "HOST", "host of application")
  cmdline.SetOptionDefault("server.host", "localhost")

  cmdline.AddOption("", "server.port", "PORT", "port of application")
  cmdline.SetOptionDefault("server.port", "8011")

  cmdline.Parse(os.Args)

  // application general configuration
  a.Config.Environment    = cmdline.OptionValue("environment")
  a.Config.Verbosity      = cmdline.OptionValue("verbosity")
  a.Config.LogFormat      = cmdline.OptionValue("log-format")
  a.Config.LogHealthcheck = cmdline.IsOptionSet("log-healthcheck")

  // a server configuration values
  a.Config.Server.Host = cmdline.OptionValue("server.host")
  a.Config.Server.Port = cmdline.OptionValue("server.port")
}

var BuildTime = "unset"
var Commit 		= "unset"
var Release 	= "unset"

func registerVersion(a *internal.Application) {
  a.Version = &internal.Version{BuildTime, Commit, Release}
}

var mainLog *log.Entry
var a        internal.Application

func init() {
  mainLog = log.WithFields(log.Fields{
    "_file": "cmd/needys-api-strategy-server/main.go",
    "_type": "system",
  })

  registerCliConfiguration(&a)
  registerVersion(&a)

  a.Initialize()
}

func main() {
  c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

  go func() {
		oscall := <-c

    mainLog.WithFields(log.Fields{
      "signal": oscall,
    }).Warn("received a system call")

		cancel()
	}()

  a.Run(ctx)
}
