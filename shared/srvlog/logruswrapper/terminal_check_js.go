//go:build js
// +build js

package logruswrapper

func isTerminal(fd int) bool {
	return false
}
