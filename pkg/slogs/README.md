# slogs

This is a wrapper around the stdlib Golang slog package that just creates a slog logger in some standard way, with a configurable log level.

Slog was introduced in Go 1.21, so that is the minimum required Go version to use this module.

## Usage

```go
package main

import "github.com/chia-network/go-modules/pkg/slogs"

func main() {
	// Init the logger with a log-level string (debug, info, warn, error)
	// defaults to "info" if empty or unsupported string 
	// not passing any logger options
	slogs.Init("info")

	// Logs a hello world message at the info level
	slogs.Logr.Info("hello world")

	// Logs an error message at the error level 
	slogs.Logr.Error("we received an error")
}
```

In a Cobra/Viper CLI app this might look more like:

```go
package cmd

import (
	"log"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/chia-network/go-modules/pkg/slogs"
)

var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Short help message for cmd",

	Run: func(cmd *cobra.Command, args []string) {
		// Init logger
		slogs.Init(viper.GetString("log-level"))

		// Application logic below
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().String("log-level", "info", "The log-level for the application, can be one of info, warn, error, debug.")
	err := viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	if err != nil {
		log.Fatalln(err.Error()) // Have to log with standard logger until the slog logger is initialized
	}
}
```

### Logger Options

You can pass some options to the logger via the second argument to Init:

```go
slogs.Init("info", slogs.WithSourceContext(true))
```

* AddSource -- when set to `true` all logs will mention the caller (the position in the code) to the logger within the logline
