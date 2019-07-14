package environment

import "flag"

var Listener *bool

func init() {
	Listener = flag.Bool("l", false, "run application as listener")
	flag.Parse()
}

func IsListener() bool {
	return Listener != nil && *Listener
}