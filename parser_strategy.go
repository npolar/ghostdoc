package ghostdoc

import "github.com/npolar/ghostdoc/context"

//ParserStrategy interface for supported file types
type ParserStrategy interface {
	getContext() context.GhostContext
	isRawInput(string) bool
	isSupportedFile(string) bool
	parse(*rawFile, chan *dataFile)
}
