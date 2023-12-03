package apierror

/*

The apierror package is an error handling package that provides a few
helper features to aid with code portability, brevity and readability.

It is intented to achieve the following goals:

1. Provide a simple and consistent means of handling all (or most) errors
- Being simple to use is important to ensure that error handling
  is not a barrier in everyday development
- Being consistent is important to ensure developers can have a familiar
  way to implement error handling, and will also make it easier for
  developers who are new to a codebase to understand the error handling
- Being consistent is also important to ensure logging can have a common
  format so that logs can be ingested into 3rd party systems, and so those
  3rd party system do not require any changes over time to accomodate
  multiple logging formats
- The library currently supports a pluggable logging system, however the
  default logging is intended to not change frequently due to there being
  multiple existing 3rd party loggings system already depending on the
  current formats

2. Provide error tracing for specific API users/developers
- Tracing is very important in finding the root cause of problems
- Tracing information needs topology/code path hiding to ensure the
  security of the code/process is kept private. i.e. this will
  require some form of authentication interface

3. Lays the ground work for an out-of-band error logging/handling system
- This is intended to allow the above tracing/stack information to
  be made available to developers without having to have made the
  API call themeselves

*/

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

// APIError is the common error response struct used for all errors
type APIError struct {
	// App is the application name where this error was generated
	App string `json:"app,omitempty"`

	// ErrorCode is the equivalent HTTP error code for this error
	ErrorCode int `json:"err_code,omitempty"`
	// ErrorNumber is an application-specific unique error number (for use with i18n and error lookups)
	ErrorNumber int `json:"err_number,omitempty"`

	ErrorType ErrorType `json:"err_type,omitempty"`

	// ErrorField is the field name (if applicable of the field that is relating to the error)
	ErrorField string `json:"err_field,omitempty"`
	// ErrorMessage is the human readable error message
	ErrorMessage string `json:"err_message"`
	// OriginErrorMessage is the origin error's native error message
	OriginErrorMessage string `json:"origin_err_message,omitempty"`

	Message string `json:"message,omitempty"` // deprecated

	// SourceFile is the file name (and line number if available) of where this error was generated
	SourceFile string `json:"source_file,omitempty"`
	// SourceFunc is the package and function name where this error was generated
	SourceFunc string `json:"source_func,omitempty"`

	// Stack is the calling stack for this error
	Stack *APIError `json:"stack,omitempty"`

	// originErr contains the original error
	originErr error `json:"-"`
	// traceFunc is the function name where this error was generated
	traceFunc string `json:"-,omitempty"`
	// traceFile is the full/original file name of where this error was generated
	traceFile string `json:"-,omitempty"`
	// tracePackage is the name of the package where this error was generated
	tracePackage string       `json:"-,omitempty"`
	traceLine    int          `json:"-,omitempty"`
	tracePC      uintptr      `json:"-,omitempty"`
	traceFrames  []TraceFrame `json:"-,omitempty"`

	originPCs []uintptr `json:"-"`

	context context.Context `json:"-"`
	isNil   bool            `json:"-"`
}

type apiErrorInspect interface {
	GetJSONBytes() ([]byte, error)
	GetErrorCode() int
	GetErrorMessage() string
	GetErrorField() string
	GetErrorNumber() int
}

type apiErrorObjectInspect interface {
	GetAPIError() *APIError
}

type apiErrorProtoInspect interface {
	GetPBBytes() ([]byte, error)
}

type echoV4HTTPError interface {
	Unwrap() error
}

var (
	// NilError is a non-nil pointer used to carry a nil payload error.
	NilError = APIError{isNil: true}
	// BackProcs sets how far back the function callers go by default
	BackProcs = 1 // windows = 1
	// NumPathsToLog determines how much path information we keep/strip out
	NumPathsToLog = 1
	// DebugMode controls whether we return full stack traces or not
	DebugMode bool
	// CompactJSON controls whether or not we return compact JSON or not
	CompactJSON bool
	// ApplicationName is the application name that generated the error
	ApplicationName string
	// Logger defines a global logging interface
	Logger APILogger = &DefaultAPILoggerHandler{}
)

// NewAPIErrorCustomCallers returns a new APIError with custom caller depth
func NewAPIErrorCustomCallers(ctx context.Context, callersNum int, errType ErrorType, err error, errcode int, errfield, errmsg string, fields ...interface{}) *APIError {
	resp := APIError{
		context:            ctx,
		ErrorCode:          errcode,
		ErrorField:         errfield,
		ErrorMessage:       fmt.Sprintf(errmsg, fields...),
		ErrorType:          errType,
		App:                ApplicationName,
		OriginErrorMessage: "",
	}

	if err != nil && reflect.TypeOf(err).Kind() == reflect.Ptr && !reflect.ValueOf(err).IsNil() {
		resp.originErr = err
	} else if err != nil && reflect.TypeOf(err).Kind() == reflect.Struct {
		resp.originErr = err
	}

	// if err != nil && !reflect.ValueOf(err).IsNil() {
	// 	resp.originErr = err
	// }

	// pc := make([]uintptr, 10) // at least 1 entry needed
	// runtime.Callers(callersNum, pc)
	// runtimeFuncPtr := runtime.FuncForPC(pc[0])

	// if runtimeFuncPtr != nil {
	// 	fullfile, line := runtimeFuncPtr.FileLine(pc[0])
	// 	_, file := filepath.Split(fullfile)
	// 	resp.File = file
	// 	resp.traceLine = line
	// 	resp.Func = runtimeFuncPtr.Name()
	// }

	if pc, file, line, ok := runtime.Caller(callersNum); ok {
		runtimeFuncPtr := runtime.FuncForPC(pc)

		// Human readable versions
		// resp.Func = runtimeFuncPtr.Name()
		resp.SourceFunc = stripPath(KeepNumDirs(runtimeFuncPtr.Name(), NumPathsToLog))
		resp.SourceFile = stripGitPath(file)
		// resp.SourceFile = KeepNumDirs(file, NumPathsToLog)
		if line > 0 && resp.SourceFile != "" {
			resp.SourceFile = resp.SourceFile + ":" + strconv.Itoa(line)
		}

		// internal native versions
		resp.tracePackage, resp.traceFunc = packageAndName(runtimeFuncPtr)
		resp.traceFile = file
		resp.traceLine = line
		resp.tracePC = pc
	}

	pcs := make([]uintptr, 100)
	runtime.Callers(callersNum+1, pcs)
	resp.originPCs = pcs

	resp.traceFrames = extractFrames(pcs)
	resp.traceFrames = filterFrames(resp.traceFrames)

	// fmt.Printf("originPCs: %v\n", resp.originPCs)

	// resp.Func = KeepNumDirs(resp.Func, NumPathsToLog)
	// resp.File = KeepNumDirs(resp.File, NumPathsToLog)

	revisedErr := resp.ApplyForeignError(resp.originErr, errcode, errfield)
	if revisedErr == nil {
		return nil
	}

	resp.ApplyFieldInheritance()

	// (optional) Logging
	if Logger != nil {
		Logger.HandleLogMessageWithContext(ctx, &resp)
	}

	// if ctx != nil {
	// 	if md, ok := metadata.FromOutgoingContext(ctx); ok {
	// 		grpc.SendHeader(ctx, md)
	// 	} else {
	// 		grpc.SendHeader(ctx, nil)
	// 	}
	// }

	return &resp
}

func stripPath(in string) string {
	if lastslash := strings.LastIndex(in, "/"); lastslash >= 0 {
		return in[lastslash+1:]
	}
	return in
}

func stripGitPath(in string) string {
	stripGitUsername := false
	if beginOfBuildPath := strings.Index(in, "/build-"); beginOfBuildPath >= 0 {
		in = in[beginOfBuildPath+7:]
		if nextSlashIndex := strings.Index(in, "/"); nextSlashIndex >= 0 {
			in = in[nextSlashIndex+1:]
		}
	}
	if beginOfGitPath := strings.LastIndex(in, "atlassian/pipelines/agent/"); beginOfGitPath >= 0 {
		in = in[beginOfGitPath+26:]
		stripGitUsername = true
	}
	if beginOfGitPath := strings.LastIndex(in, "atlassian/pipelines/build/"); beginOfGitPath >= 0 {
		in = in[beginOfGitPath+26:]
		stripGitUsername = true
	}
	if beginOfGitPath := strings.LastIndex(in, "atlassian/pipelines/"); beginOfGitPath >= 0 {
		in = in[beginOfGitPath+20:]
		stripGitUsername = true
	}
	if beginOfGitPath := strings.LastIndex(in, "github.com/"); beginOfGitPath >= 0 {
		in = in[beginOfGitPath+11:]
		stripGitUsername = true
	}
	if beginOfGitPath := strings.LastIndex(in, "bitbucket.org/"); beginOfGitPath >= 0 {
		in = in[beginOfGitPath+14:]
		stripGitUsername = true
	}
	if stripGitUsername {
		if firstslash := strings.Index(in, "/"); firstslash >= 0 {
			return in[firstslash+1:]
		}
	}

	return in
}

func packageAndName(fn *runtime.Func) (string, string) {
	name := fn.Name()
	pkg := ""

	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//  runtime/debug.*T·ptrmethod
	// and want
	//  *T.ptrmethod
	// Since the package path might contains dots (e.g. code.google.com/...),
	// we first remove the path prefix if there is one.
	if lastslash := strings.LastIndex(name, "/"); lastslash >= 0 {
		pkg += name[:lastslash] + "/"
		name = name[lastslash+1:]
	}
	if period := strings.Index(name, "."); period >= 0 {
		pkg += name[:period]
		name = name[period+1:]
	}

	name = strings.Replace(name, "·", ".", -1)
	return pkg, name
}

func keepNumDirs(str string, lastn int, startat int) string {
	numFound := strings.Count(str[startat:], "/")
	if numFound > lastn {
		return keepNumDirs(str, lastn, 1+startat+strings.Index(str[startat:], "/"))
	}
	return str[startat:]
}

// KeepNumDirs returns the lastn number of directories in a path+filename combo
func KeepNumDirs(str string, lastn int) string {
	return keepNumDirs(str, lastn, 0)
}

// New is a short form of NewAPIError with an empty errfield.
func New(err error, errcode int, errmsg string, fields ...interface{}) *APIError {
	return NewAPIErrorCustomCallers(nil, BackProcs+1, ErrorType_ERROR, err, errcode, "", errmsg, fields...)
}

// NewAPIError returns a new APIError
func NewAPIError(err error, errcode int, errfield, errmsg string, fields ...interface{}) *APIError {
	return NewAPIErrorCustomCallers(nil, BackProcs+1, ErrorType_ERROR, err, errcode, errfield, errmsg, fields...)
}

// NewAPIErrorIfError returns a new APIError, if err is not nil
func NewAPIErrorIfError(err error, errcode int, errfield, errmsg string, fields ...interface{}) error {
	if err == nil {
		return nil
	}
	return NewAPIErrorCustomCallers(nil, BackProcs+1, ErrorType_ERROR, err, errcode, errfield, errmsg, fields...)
}

// NewAPIErrorWithContext returns a new APIError
func NewAPIErrorWithContext(ctx context.Context, err error, errcode int, errfield, errmsg string, fields ...interface{}) *APIError {
	return NewAPIErrorCustomCallers(ctx, BackProcs+1, ErrorType_ERROR, err, errcode, errfield, errmsg, fields...)
}

// NewAPIErrorWithContextIfError returns a new APIError, if err is not nil
func NewAPIErrorWithContextIfError(ctx context.Context, err error, errcode int, errfield, errmsg string, fields ...interface{}) error {
	if err == nil {
		return nil
	}
	return NewAPIErrorCustomCallers(ctx, BackProcs+1, ErrorType_ERROR, err, errcode, errfield, errmsg, fields...)
}

// NewAPIWarning returns a new APIError
func NewAPIWarning(err error, errcode int, errfield, errmsg string, fields ...interface{}) *APIError {
	return NewAPIErrorCustomCallers(nil, BackProcs+1, ErrorType_WARNING, err, errcode, errfield, errmsg, fields...)
}

// NewAPIWarningIfError returns a new APIError, if err is not nil
func NewAPIWarningIfError(err error, errcode int, errfield, errmsg string, fields ...interface{}) error {
	if err == nil {
		return nil
	}
	return NewAPIErrorCustomCallers(nil, BackProcs+1, ErrorType_WARNING, err, errcode, errfield, errmsg, fields...)
}

// NewAPIWarningWithContext returns a new APIError
func NewAPIWarningWithContext(ctx context.Context, err error, errcode int, errfield, errmsg string, fields ...interface{}) *APIError {
	return NewAPIErrorCustomCallers(ctx, BackProcs+1, ErrorType_WARNING, err, errcode, errfield, errmsg, fields...)
}

// NewAPIWarningWithContextIfError returns a new APIError
func NewAPIWarningWithContextIfError(ctx context.Context, err error, errcode int, errfield, errmsg string, fields ...interface{}) error {
	if err == nil {
		return nil
	}
	return NewAPIErrorCustomCallers(ctx, BackProcs+1, ErrorType_WARNING, err, errcode, errfield, errmsg, fields...)
}

// NewAPIInfo returns a new APIError
func NewAPIInfo(err error, errcode int, errfield, errmsg string, fields ...interface{}) *APIError {
	return NewAPIErrorCustomCallers(nil, BackProcs+1, ErrorType_INFO, err, errcode, errfield, errmsg, fields...)
}

// NewAPIInfoIfError returns a new APIError
func NewAPIInfoIfError(err error, errcode int, errfield, errmsg string, fields ...interface{}) error {
	if err == nil {
		return nil
	}
	return NewAPIErrorCustomCallers(nil, BackProcs+1, ErrorType_INFO, err, errcode, errfield, errmsg, fields...)
}

// NewAPIInfoWithContext returns a new APIError
func NewAPIInfoWithContext(ctx context.Context, err error, errcode int, errfield, errmsg string, fields ...interface{}) *APIError {
	return NewAPIErrorCustomCallers(ctx, BackProcs+1, ErrorType_INFO, err, errcode, errfield, errmsg, fields...)
}

// NewAPIInfoWithContextIfError returns a new APIError
func NewAPIInfoWithContextIfError(ctx context.Context, err error, errcode int, errfield, errmsg string, fields ...interface{}) error {
	if err == nil {
		return nil
	}
	return NewAPIErrorCustomCallers(ctx, BackProcs+1, ErrorType_INFO, err, errcode, errfield, errmsg, fields...)
}

// NewAPIDebug returns a new APIError
func NewAPIDebug(err error, errcode int, errfield, errmsg string, fields ...interface{}) *APIError {
	return NewAPIErrorCustomCallers(nil, BackProcs+1, ErrorType_DEBUG, err, errcode, errfield, errmsg, fields...)
}

// NewAPIDebugIfError returns a new APIError
func NewAPIDebugIfError(err error, errcode int, errfield, errmsg string, fields ...interface{}) error {
	if err == nil {
		return nil
	}
	return NewAPIErrorCustomCallers(nil, BackProcs+1, ErrorType_DEBUG, err, errcode, errfield, errmsg, fields...)
}

// NewAPIDebugWithContext returns a new APIError
func NewAPIDebugWithContext(ctx context.Context, err error, errcode int, errfield, errmsg string, fields ...interface{}) *APIError {
	return NewAPIErrorCustomCallers(ctx, BackProcs+1, ErrorType_DEBUG, err, errcode, errfield, errmsg, fields...)
}

// NewAPIDebugWithContextIfError returns a new APIError
func NewAPIDebugWithContextIfError(ctx context.Context, err error, errcode int, errfield, errmsg string, fields ...interface{}) error {
	if err == nil {
		return nil
	}
	return NewAPIErrorCustomCallers(ctx, BackProcs+1, ErrorType_DEBUG, err, errcode, errfield, errmsg, fields...)
}

// NewAPIErrorFromJSONBytes takes a JSON []byte slice and returns an APIError
func NewAPIErrorFromJSONBytes(jsonBytes []byte) (*APIError, error) {
	newAPIError := APIError{}
	if err := json.Unmarshal(jsonBytes, &newAPIError); err != nil {
		// JSON format error
		return nil, &APIError{originErr: err, ErrorMessage: "Unable to decode APIError JSON"}
	}
	if newAPIError.IsEmpty() {
		// JSON Decoded, but it didn't decode into an APIError struct
		return nil, &APIError{ErrorMessage: "Empty APIError"}
	}
	return &newAPIError, nil
}

// IsEmpty checks if an APIError is empty
func (ae *APIError) IsEmpty() bool {
	if ae == nil {
		return true
	}

	if ae.ErrorCode == 0 &&
		ae.originErr == nil &&
		ae.ErrorNumber == 0 &&
		ae.ErrorField == "" &&
		ae.ErrorMessage == "" &&
		ae.Message == "" {
		return true
	}

	return false
}

// GetContext return the context or nil
func (ae *APIError) GetContext() context.Context {
	return ae.context
}

// Unwrap gets the error at the top of the stack
func (ae *APIError) Unwrap() error {
	if ae == nil {
		return nil
	}

	if ae.originErr != nil {
		return ae.originErr
	}
	// fmt.Printf("APIError.Unwrap: %v\n", ae.GetJSONStringPretty())
	// fmt.Printf("APIError.Unwrap.Stack: %v\n", ae.Stack)
	return ae.Stack
}

// Cause gets the error at the top of the stack
func (ae *APIError) Cause() error {
	if ae == nil {
		return nil
	}

	// fmt.Printf("APIError.Cause: %v\n", ae.GetJSONStringPretty())
	// fmt.Printf("APIError.Cause.Stack: %v\n", ae.Stack)
	return ae.Stack
}

// GetTracePC
func (ae *APIError) GetTracePC() uintptr {
	if ae == nil {
		return 0
	}
	return ae.tracePC
}

// GetOriginPCs
func (ae *APIError) GetOriginPCs() []uintptr {
	if ae == nil {
		return nil
	}
	return ae.originPCs
}

type StackFrame struct {
	// The path to the file containing this ProgramCounter
	File string
	// The LineNumber in that file
	LineNumber int
	// The Name of the function that contains this ProgramCounter
	Name string
	// The Package that contains this function
	Package string
	// The underlying ProgramCounter
	ProgramCounter uintptr
}

func (ae *APIError) GetStackFrame() *StackFrame {
	if ae == nil {
		return nil
	}
	sf := &StackFrame{
		File:       ae.traceFile,
		LineNumber: ae.traceLine,
		Name:       ae.traceFunc,
		Package:    ae.tracePackage,
		// ProgramCounter: ae.originPC,
	}
	fmt.Printf("GetStackFrame originPCs: %v\n", ae.originPCs)
	if len(ae.originPCs) > 0 {
		sf.ProgramCounter = ae.originPCs[0]
	}

	return sf
}

const maxErrorDepth = 100

// StackTrace gets the stack trace
func (ae *APIError) StackTrace() []uintptr {
	var resp []uintptr
	var nextErr *APIError
	nextErr = ae
	for i := 0; i < maxErrorDepth && nextErr != nil; i++ {
		fmt.Printf("StackTrace item: %v\n", nextErr)
		fmt.Printf("StackTrace item.originPCs: %v\n", nextErr.originPCs)
		// runtimeFuncPtr := runtime.FuncForPC(nextErr.originPC)
		// pkgname, funcname := packageAndName(runtimeFuncPtr)
		// fmt.Printf("StackTrace pkgname(%s) funcname(%s)\n", pkgname, funcname)

		if len(nextErr.originPCs) > 0 {
			resp = append(resp, nextErr.originPCs...)
		}

		// resp = append(resp, nextErr.originPC)
		nextErr = nextErr.Stack
	}

	fmt.Printf("StackTrace resp: %v\n", resp)
	return resp
}

// StackTrace gets the stack trace
func (ae *APIError) StackTrace2() []StackFrame {
	if ae == nil {
		fmt.Println("APIError.StackTrace nil")
		return nil
	}
	fmt.Printf("ae.StackTrace: %v\n", ae.GetJSONStringPretty())

	var nextErr *APIError
	nextErr = ae

	var stackList []StackFrame

	for i := 0; i < maxErrorDepth && nextErr != nil; i++ {
		fmt.Printf("StackTrace item: %v\n", nextErr)
		thisStack := nextErr.GetStackFrame()
		if thisStack == nil {
			fmt.Println("APIError.StackTrace break")
			break
		}
		stackList = append(stackList, *thisStack)
		nextErr = nextErr.Stack
	}

	fmt.Printf("APIError.StackTrace %v\n", stackList)
	return stackList
}
