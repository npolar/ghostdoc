package context

import (
	"time"

	"github.com/codegangsta/cli"
)

// GhostContext is a interface to codegangsta cli
type GhostContext interface {
	Int(name string) int
	Duration(name string) time.Duration
	Float64(name string) float64
	Bool(name string) bool
	BoolT(name string) bool
	String(name string) string
	StringSlice(name string) []string
	IntSlice(name string) []int
	Generic(name string) interface{}
	GlobalInt(name string) int
	GlobalDuration(name string) time.Duration
	GlobalBool(name string) bool
	GlobalString(name string) string
	GlobalStringSlice(name string) []string
	GlobalIntSlice(name string) []int
	GlobalGeneric(name string) interface{}
	NumFlags() int
	IsSet(name string) bool
	GlobalIsSet(name string) bool
	FlagNames() (names []string)
	GlobalFlagNames() (names []string)
	Args() cli.Args
	Cli() *cli.Context
}
