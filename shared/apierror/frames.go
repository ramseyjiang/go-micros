package apierror

import (
	"go/build"
	"path/filepath"
	"runtime"
	"strings"
)

// TraceFrame represents a function call and it's metadata. Frames are associated with a Stacktrace.
// type TraceFrame struct {
// 	Function    string                 `json:"function,omitempty"`
// 	Symbol      string                 `json:"symbol,omitempty"`
// 	Module      string                 `json:"module,omitempty"`
// 	Package     string                 `json:"package,omitempty"`
// 	Filename    string                 `json:"filename,omitempty"`
// 	AbsPath     string                 `json:"abs_path,omitempty"`
// 	Lineno      int                    `json:"lineno,omitempty"`
// 	Colno       int                    `json:"colno,omitempty"`
// 	PreContext  []string               `json:"pre_context,omitempty"`
// 	ContextLine string                 `json:"context_line,omitempty"`
// 	PostContext []string               `json:"post_context,omitempty"`
// 	InApp       bool                   `json:"in_app,omitempty"`
// 	Vars        map[string]interface{} `json:"vars,omitempty"`
// }

// NewFrame assembles a stacktrace frame out of runtime.Frame.
func NewFrame(f runtime.Frame) TraceFrame {
	var abspath, relpath string
	// NOTE: f.File paths historically use forward slash as path separator even
	// on Windows, though this is not yet documented, see
	// https://golang.org/issues/3335. In any case, filepath.IsAbs can work with
	// paths with either slash or backslash on Windows.
	switch {
	case f.File == "":
		relpath = "unknown"
		// Leave abspath as the empty string to be omitted when serializing
		// event as JSON.
		abspath = ""
	case filepath.IsAbs(f.File):
		abspath = f.File
		relpath = ""
	default:
		// f.File is a relative path. This may happen when the binary is built
		// with the -trimpath flag.
		relpath = f.File
		// Omit abspath when serializing the event as JSON.
		abspath = ""
	}

	function := f.Function
	var pkg string

	if function != "" {
		pkg, function = splitQualifiedFunctionName(function)
	}

	frame := TraceFrame{
		AbsPath:  abspath,
		Filename: relpath,
		Lineno:   int32(f.Line),
		Module:   pkg,
		Function: function,
	}

	frame.InApp = isInAppFrame(frame)

	return frame
}

// splitQualifiedFunctionName splits a package path-qualified function name into
// package name and function name. Such qualified names are found in
// runtime.Frame.Function values.
func splitQualifiedFunctionName(name string) (pkg string, fun string) {
	pkg = packageName(name)
	fun = strings.TrimPrefix(name, pkg+".")
	return
}

func extractFrames(pcs []uintptr) []TraceFrame {
	var frames []TraceFrame
	callersFrames := runtime.CallersFrames(pcs)

	for {
		callerFrame, more := callersFrames.Next()

		frames = append([]TraceFrame{
			NewFrame(callerFrame),
		}, frames...)

		if !more {
			break
		}
	}
	// spew.Dump(frames)
	return frames
}

func isInAppFrame(frame TraceFrame) bool {
	if strings.HasPrefix(frame.AbsPath, build.Default.GOROOT) ||
		strings.Contains(frame.Module, "vendor") ||
		strings.Contains(frame.Module, "third_party") {
		return false
	}

	return true
}

// packageName returns the package part of the symbol name, or the empty string
// if there is none.
// It replicates https://golang.org/pkg/debug/gosym/#Sym.PackageName, avoiding a
// dependency on debug/gosym.
func packageName(name string) string {
	// A prefix of "type." and "go." is a compiler-generated symbol that doesn't belong to any package.
	// See variable reservedimports in cmd/compile/internal/gc/subr.go
	if strings.HasPrefix(name, "go.") || strings.HasPrefix(name, "type.") {
		return ""
	}

	pathend := strings.LastIndex(name, "/")
	if pathend < 0 {
		pathend = 0
	}

	if i := strings.Index(name[pathend:], "."); i != -1 {
		return name[:pathend+i]
	}
	return ""
}

// baseName returns the symbol name without the package or receiver name.
// It replicates https://golang.org/pkg/debug/gosym/#Sym.BaseName, avoiding a
// dependency on debug/gosym.
func baseName(name string) string {
	if i := strings.LastIndex(name, "."); i != -1 {
		return name[i+1:]
	}
	return name
}

// filterFrames filters out stack frames that are not meant to be reported to
// Sentry. Those are frames internal to the SDK or Go.
func filterFrames(frames []TraceFrame) []TraceFrame {
	if len(frames) == 0 {
		return nil
	}

	filteredFrames := make([]TraceFrame, 0, len(frames))

	for _, frame := range frames {
		// Skip Go internal frames.
		if frame.Module == "runtime" || frame.Module == "testing" {
			continue
		}
		// Skip Sentry internal frames, except for frames in _test packages (for
		// testing).
		if strings.HasPrefix(frame.Module, "bitbucket.org/iqhive/apierror") &&
			!strings.HasSuffix(frame.Module, "_test") {
			continue
		}
		filteredFrames = append(filteredFrames, frame)
	}

	return filteredFrames
}
