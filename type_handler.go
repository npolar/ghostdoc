package ghostdoc

// TypeHandler interface used to check filetypes
type TypeHandler interface {
	rawInput(string) bool
	supportedFile(string) bool
}
