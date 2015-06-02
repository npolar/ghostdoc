package context

import (
	"time"

	"github.com/codegangsta/cli"
)

// CliContext typedef
type CliContext struct {
	cli *cli.Context
}

// NewCliContext instantiates a new CliContext
func NewCliContext(c *cli.Context) *CliContext {
	return &CliContext{cli: c}
}

// Cli gets the underlying cli.Context, avoid using this directly
func (c *CliContext) Cli() *cli.Context {
	return c.cli
}

// Int @see cli.Context.Int
func (c *CliContext) Int(name string) int {
	return c.cli.Int(name)
}

// Duration @see cli.Context.Duration
func (c *CliContext) Duration(name string) time.Duration {
	return c.cli.Duration(name)
}

// Float64 @see cli.Context.Float64
func (c *CliContext) Float64(name string) float64 {
	return c.cli.Float64(name)
}

// Bool @see cli.Context.Bool
func (c *CliContext) Bool(name string) bool {
	return c.cli.Bool(name)
}

// BoolT @see cli.Context.BoolT
func (c *CliContext) BoolT(name string) bool {
	return c.cli.BoolT(name)
}

// String @see cli.Context.String
func (c *CliContext) String(name string) string {
	return c.cli.String(name)
}

// StringSlice @see cli.Context.StringSlice
func (c *CliContext) StringSlice(name string) []string {
	return c.cli.StringSlice(name)
}

// IntSlice @see cli.Context.IntSlice
func (c *CliContext) IntSlice(name string) []int {
	return c.cli.IntSlice(name)
}

// Generic @see cli.Context.Generic
func (c *CliContext) Generic(name string) interface{} {
	return c.cli.Generic(name)
}

// GlobalInt @see cli.Context.GlobalInt
func (c *CliContext) GlobalInt(name string) int {
	return c.cli.GlobalInt(name)
}

// GlobalDuration @see cli.Context.GlobalDuration
func (c *CliContext) GlobalDuration(name string) time.Duration {
	return c.cli.GlobalDuration(name)
}

// GlobalBool @see cli.Context.GlobalBool
func (c *CliContext) GlobalBool(name string) bool {
	return c.cli.GlobalBool(name)
}

// GlobalString @see cli.Context.GlobalString
func (c *CliContext) GlobalString(name string) string {
	return c.cli.GlobalString(name)
}

// GlobalStringSlice @see cli.Context.GlobalStringSlice
func (c *CliContext) GlobalStringSlice(name string) []string {
	return c.cli.GlobalStringSlice(name)
}

// GlobalIntSlice @see cli.Context.GlobalIntSlice
func (c *CliContext) GlobalIntSlice(name string) []int {
	return c.cli.GlobalIntSlice(name)
}

// GlobalGeneric @see cli.Context.GlobalGeneric
func (c *CliContext) GlobalGeneric(name string) interface{} {
	return c.cli.GlobalGeneric(name)
}

// NumFlags @see cli.Context.NumFlags
func (c *CliContext) NumFlags() int {
	return c.cli.NumFlags()
}

// IsSet @see cli.Context.IsSet
func (c *CliContext) IsSet(name string) bool {
	return c.cli.IsSet(name)
}

// GlobalIsSet @see cli.Context.GlobalIsSet
func (c *CliContext) GlobalIsSet(name string) bool {
	return c.cli.GlobalIsSet(name)
}

// FlagNames @see cli.Context.FlagNames
func (c *CliContext) FlagNames() (names []string) {
	return c.cli.FlagNames()
}

// GlobalFlagNames @see cli.Context.GlobalFlagNames
func (c *CliContext) GlobalFlagNames() (names []string) {
	return c.cli.GlobalFlagNames()
}

// Args @see cli.Context.Args
func (c *CliContext) Args() cli.Args {
	return c.cli.Args()
}
