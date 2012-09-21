package main

import (
	"github.com/davecheney/gmx"
)

var (
	active   = gmx.NewGuage("socksie.connections.active")
	accepted = gmx.NewCounter("socksie.connections.accepted")
)
