// Copyright (c) 2021, Geert JM Vanderkelen

/*
Package xlog implements a sofisticated logging package.

### Quick Start

When importing xlog, logging is can be used right away:

	package main

	import "github.com/geertjanvdk/xkit/xlog"

	func main() {
		xlog.Info("I am an informational.")
	}

It is possible to add fields as well as error or a scope:

	package main

	import "github.com/geertjanvdk/xkit/xlog"

	func main() {
		if err := createUser(user); err != nil {
			xlog.WithError(err).Errorf("creating user for %s", appName)
		} else {
			xlog.WithScope("users").WithField("user", user.ID).Infof("user created")
		}
	}


### Setting Logging Levels

xlog does not support the classic way of setting the level of logging. Normally,
when you set log level WARN, nothing informational is written, and only errors
are reported.

In xlog levels are activated, so that the application can choose which are
logged in a dynamic manner.

For example: an application is running and has the defaults set. This means that
informational, warning, and error messages are logged, as well as fatal or
terminating events. When, however, the developers need for a shor time to have
debugging enabled (for example in staging environments), they can do so by
activating the debug level.

	xlog.ActivateLevel(xlog.DebugLevel)

How this is dynamically done is up to the application. It could be possible to
have a process signal, an API call, or as simple as creating a file in a
specific folder.

Deactivating is simply done using using the Deactivate method:

	xlog.DeactivateLevel(xlog.DebugLevel)


### Custom Loggers

xlog comes with a logger. It is however more practical for applications to
instantiate their own, and customize it.

	package main

	import "github.com/geertjanvdk/xkit/xlog"

	func main() {
		logger := xlog.New()
		logger.DeactivateLevels(xlog.WarnLevel, xlog.InfoLevel)
		logger.Scope = "my-app"
		logger.Errorf("only errors are logged (and fatal events)")
	}

*/
package xlog
