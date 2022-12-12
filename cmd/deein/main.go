package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	deein "github.com/deeincom/deeincom"
	"github.com/deeincom/deeincom/config"
	"github.com/deeincom/deeincom/internal/cli/root"

	// commands
	_ "github.com/deeincom/deeincom/internal/cli/api"
	_ "github.com/deeincom/deeincom/internal/cli/version"
	_ "github.com/deeincom/deeincom/internal/cli/web"
)

func main() {
	fmt.Println("starting cmd: ", os.Args[1])
	for _, cmd := range root.CmdSet {
		cfg := cmd.String("cfg", "config.json", "path to config file")

		if os.Args[1] == cmd.Name() {
			cmd.Parse(os.Args[2:])
			cmd.App = newApp(*cfg)

			err := cmd.Exec()
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

	fmt.Println("Your cmd was not supported. See `as -h`")
	for _, cmd := range root.CmdSet {
		fmt.Println("- cmd: " + cmd.Name())
		cmd.Usage()
	}
}

func newApp(cfg string) *deein.App {
	retryNum := 0
retry:
	c, err := config.Read(cfg)
	if isMissingConfig(err) {
		retryNum = retryNum + 1
		log.Println("CONFIG: Failed", err, cfg)
		err := config.Create()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		if retryNum < 3 {
			goto retry
		}
		log.Println("CONFIG: MAX RETRY")
		os.Exit(1)

	}

	if err != nil {
		panic(err)
	}

	log.Println("CONFIG: OK")

	app, err := deein.New(c)
	if err != nil {
		panic(err)
	}
	return app
}

func isMissingConfig(err error) bool {
	var pathError *os.PathError
	return errors.As(err, &pathError)
}
