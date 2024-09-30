package main

import (
	"log"

	"gototp/internal/gototp"
	"gototp/internal/view"
)

func main() {
	_gototp, err := gototp.New()
	if err != nil {
		log.Fatalln(err)
	}
	_view := view.New(
		_gototp,
	)
	_view.Run()
}
