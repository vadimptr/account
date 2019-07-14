package environment

import "flag"

var Listener *bool
var Method *string
const DefaultMethod = "serializable"

func init() {
	Listener = flag.Bool("l", false, "run application as listener")
	Method = flag.String("m", DefaultMethod, "[serializable or optimistic")
	flag.Parse()
}

func IsListener() bool {
	return Listener != nil && *Listener
}

func GetMethod() string {
	if Method != nil {
		return *Method
	}
	return DefaultMethod
}
