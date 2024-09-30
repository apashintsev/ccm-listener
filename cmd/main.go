package main

import (
	"photon-listener/internal/app"
	"photon-listener/internal/scanner"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := app.InitApp(); err != nil {
		return err
	}

	scanner, err := scanner.NewScanner()
	if err != nil {
		return err
	}

	scanner.Listen()

	return nil
}
