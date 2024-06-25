//go:build unix

package main

import (
	"fmt"
)

var timesDir = os.Getenv("HOME") + "/.local/share/times"
