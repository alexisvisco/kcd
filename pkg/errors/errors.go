package errors

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"strings"
	"text/tabwriter"

	"github.com/sirupsen/logrus"
)

// DeactivateLocation allow to not execute the code for retrieving location of errors since this is not a
// free operation.
var DeactivateLocation = false

// Error is the represented error and underlying errors too
type Error struct {
	// location is filled by default
	location *location

	// Kind represent the type of error
	Kind Kind

	// Message is the information about the error
	Message string

	// fields is the way to wrap context to the error
	fields map[string]interface{}

	// Err is the wrapped error
	Err error
}

// New create a new error
// By default the kind is KindInternal
func New(format string, i ...interface{}) *Error {
	return &Error{
		location: retrieveLocation(),
		Kind:     KindInternal,
		Message:  fmt.Sprintf(format, i...),
		fields:   map[string]interface{}{},
		Err:      nil,
	}
}

// NewWithKind create a new error with a kind
func NewWithKind(kind Kind, format string, i ...interface{}) *Error {
	return &Error{
		location: retrieveLocation(),
		Kind:     kind,
		Message:  fmt.Sprintf(format, i...),
		fields:   map[string]interface{}{},
		Err:      nil,
	}
}

// Wrap the raw error with a descriptive message
// By default the kind is KindInternal
func Wrap(err error, format string, i ...interface{}) *Error {
	e, ok := err.(*Error)

	newErr := &Error{
		location: retrieveLocation(),
		Message:  fmt.Sprintf(format, i...),
		fields:   map[string]interface{}{},
		Err:      err,
	}
	if ok {
		newErr.Kind = e.Kind
		newErr.fields = e.fields
	} else {
		newErr.Kind = KindInternal
	}

	return newErr
}

// WrapWithKind the raw error with a descriptive message with custom kind
func WrapWithKind(kind Kind, err error, format string, i ...interface{}) *Error {
	e, ok := err.(*Error)

	newErr := &Error{
		location: retrieveLocation(),
		Message:  fmt.Sprintf(format, i...),
		fields:   map[string]interface{}{},
		Err:      err,
	}
	if ok {
		newErr.fields = e.fields
	}

	newErr.Kind = e.Kind

	return newErr
}

// Error return the error as string with the underlying errors
func (e *Error) Error() string {
	var errors []string

	if e.Message != "" {
		errors = append(errors, e.Message)
	}

	tmp := e.Err

	for tmp != nil {
		x, ok := tmp.(*Error)

		if ok {
			if x.Message != "" {
				errors = append(errors, x.Message)
			} else {
				errors = append(errors, x.Error())
			}
			tmp = x.Err
		} else {
			errors = append(errors, tmp.Error())
			break
		}
	}

	return strings.Join(errors, ": ")
}

// Stacktrace format a pretty printed stacktrace of the errors
func (e *Error) Stacktrace() string {
	st := colorRed + "┌─────────────────────────────────────────────\n│ DEBUG ERROR: "

	const space = "  "

	var tmp error = e
	isFirst := true
	for tmp != nil {
		x, ok := tmp.(*Error)

		if ok {
			if isFirst {
				st += fmt.Sprintf("%s%q%s", colorTeal, x.Message, colorRed) + "\n│" + space
				if x.location != nil {
					st += fmt.Sprintf("at %s:%d %s.%s %s%s%s",
						x.location.file,
						x.location.line,
						x.location.pkg,
						x.location.function,
						colorPurple, "<- error happened here", colorRed)
				}
				isFirst = false
			} else {
				if x.location != nil {
					st += fmt.Sprintf("at %s:%d %s.%s => %s%q%s",
						x.location.file,
						x.location.line,
						x.location.pkg,
						x.location.function,
						colorTeal, x.Message, colorRed)
				} else {
					st += fmt.Sprintf("%s%q%s", colorTeal, x.Message, colorRed)
				}
			}

			if x.Err != nil {
				st += "\n│" + space
			}

			tmp = x.Err
		} else {
			st += tmp.Error()
			break
		}
	}

	if len(e.fields) > 0 {
		st += "\n│\n│ FIELDS ATTACHED:\n"
		buffer := bytes.NewBuffer([]byte{})
		w := tabwriter.NewWriter(buffer, 0, 0, 3, ' ', tabwriter.TabIndent)
		for k, v := range e.fields {
			fmt.Fprintf(w, "│ %s%q\t%s%v%s\n", colorTeal, k, colorMagenta, v, colorRed)
		}
		w.Flush()
		st += buffer.String()
	}

	st += "\n" + colorReset

	return st
}

// WithMessage set the message
func (e *Error) WithMessage(format string, i ...interface{}) *Error {
	e.Message = fmt.Sprintf(format, i...)
	return e
}

// WithField add a field
func (e *Error) WithField(key string, value interface{}) *Error {
	e.fields[key] = value
	return e
}

// WithFields add fields
func (e *Error) WithFields(m map[string]interface{}) *Error {
	for k, v := range m {
		e.fields[k] = v
	}
	return e
}

// WithKind set the kind
func (e *Error) WithKind(kind Kind) *Error {
	e.Kind = kind
	return e
}

// WithOpHere relocalise the error where this function is called
func (e *Error) WithOpHere() *Error {
	e.location = retrieveLocation()
	return e
}

// GetField get a field by its key
func (e *Error) GetField(key string) (value interface{}, ok bool) {
	value, ok = e.fields[key]
	return value, ok
}

// locations return a list of the location of each errors
func (e *Error) locations() []string {
	if e.location == nil {
		return []string{}
	}

	var st []string

	var tmp error = e
	for tmp != nil {
		x, ok := tmp.(*Error)

		if ok {
			st = append(st,
				fmt.Sprintf(
					"%s/%s.%s:%d",
					x.location.pkg,
					x.location.file,
					x.location.function,
					x.location.line,
				))
			tmp = x.Err
		} else {
			break
		}
	}

	return st
}

// Log is a shortcut to log an error
func (e *Error) Log() *logrus.Entry {
	entry := logrus.
		WithFields(e.fields)

	if e.location != nil {
		entry = entry.WithField("stacktrace", strings.Join(e.locations(), ", "))
	}

	return entry.WithError(e)
}

type location struct {
	pkg      string
	file     string
	function string
	line     int
}

// retrieveLocation store the location of the error
// if DeactivateLocation is false : return nil
func retrieveLocation() *location {
	if DeactivateLocation {
		return nil
	}

	pc, file, line, _ := runtime.Caller(2)
	_, fileName := path.Split(file)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""
	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	return &location{
		pkg:      packageName,
		file:     fileName,
		function: funcName,
		line:     line,
	}
}
