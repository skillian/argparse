package argparse

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/skillian/errors"
)

// Argument holds the definition of an argument.
type Argument struct {
	// Action holds the action to perform after successful parsing of
	// values associated with the given argument.
	Action ArgumentAction

	// Const holds the value associated with this argument when the
	// argument is present.
	Const interface{}

	// Default is the value associated with the argument when a specific
	// value is not otherwise provided.
	Default interface{}

	// Dest is the string key that the argument can be retrieved by.
	Dest string

	// Help is the help text associated with the argument.
	Help string

	// MetaVar is the variable that the argument is represented with when
	// displaying its usage.  It is a slice in case Nargs is non-zero.
	MetaVar []string

	// Nargs is the number of values that this argument can accept.  It
	// should be a positive int unless it is one of the sentinel values:
	// ZeroOrOne, ZeroOrMore, or OneOrMore.
	Nargs int

	// OptionStrings are the possible string values that the argument can
	// be matched against.
	OptionStrings []string

	// Required determines if the argument is required or not.
	Required bool

	// Type holds a function that can be used to parse a string value into
	// the type desired by this argument.
	Type ValueParser
}

// Optional returns whether or not this is an optional (flag) argument.  If
// it is not, then it is a positional argument.
func (a *Argument) Optional() bool {
	for _, s := range a.OptionStrings {
		if strings.HasPrefix(s, "-") {
			return true
		}
	}
	return false
}

const (
	// OneOrMore means that one or more argument values are accepted by
	// the argument.
	OneOrMore int = -1 - iota

	// ZeroOrMore indicates that zero or more arguments are accepted.
	ZeroOrMore

	// ZeroOrOne indicates that zero or one argument is allowed
	ZeroOrOne
)

// isValidNarg is a helper function that can tell if a Nargs value is either a
// valid number of arguments or valid sentinel value.
func isValidNarg(v int) bool {
	return v >= ZeroOrOne
}

// ValueParser can parse a string value into a Go value.
type ValueParser func(v string) (interface{}, error)

// Bool converts the given string into a boolean value.
// It implements the ValueParser interface.
func Bool(v string) (interface{}, error) {
	if strings.EqualFold(v, "true") {
		return true, nil
	}
	if strings.EqualFold(v, "false") {
		return false, nil
	}
	return nil, errors.NewUnexpectedType(false, v)
}

// Float32 converts the given string into a float32 value.
// It implements the ValueParser interface.
func Float32(v string) (interface{}, error) {
	var f float32
	err := sscanf(v, "%f", &f)
	return f, err
}

// Float64 converts the given string into a float64 value.
// It implements the ValueParser interface.
func Float64(v string) (interface{}, error) {
	var f float64
	err := sscanf(v, "%f", &f)
	return f, err
}

// Int converts the given string into a int value.
// It implements the ValueParser interface.
func Int(v string) (interface{}, error) {
	var i int
	err := sscanf(v, "%d", &i)
	return i, err
}

// Int8 converts the given string into a int8 value.
// It implements the ValueParser interface.
func Int8(v string) (interface{}, error) {
	var i int8
	err := sscanf(v, "%d", &i)
	return i, err
}

// Int16 converts the given string into a int16 value.
// It implements the ValueParser interface.
func Int16(v string) (interface{}, error) {
	var i int16
	err := sscanf(v, "%d", &i)
	return i, err
}

// Int32 converts the given string into a int32 value.
// It implements the ValueParser interface.
func Int32(v string) (interface{}, error) {
	var i int32
	err := sscanf(v, "%d", &i)
	return i, err
}

// Int64 converts the given string into a int value.
// It implements the ValueParser interface.
func Int64(v string) (interface{}, error) {
	var i int
	err := sscanf(v, "%d", &i)
	return i, err
}

// Uint converts the given string into a uint value.
// It implements the ValueParser interface.
func Uint(v string) (interface{}, error) {
	var i uint
	err := sscanf(v, "%u", &i)
	return i, err
}

// Uint8 converts the given string into a uint8 value.
// It implements the ValueParser interface.
func Uint8(v string) (interface{}, error) {
	var i uint8
	err := sscanf(v, "%u", &i)
	return i, err
}

// Uint16 converts the given string into a uint16 value.
// It implements the ValueParser interface.
func Uint16(v string) (interface{}, error) {
	var i uint16
	err := sscanf(v, "%u", &i)
	return i, err
}

// Uint32 converts the given string into a uint32 value.
// It implements the ValueParser interface.
func Uint32(v string) (interface{}, error) {
	var i uint32
	err := sscanf(v, "%u", &i)
	return i, err
}

// Uint64 converts the given string into a uint64 value.
// It implements the ValueParser interface.
func Uint64(v string) (interface{}, error) {
	var i uint64
	err := sscanf(v, "%u", &i)
	return i, err
}

// String is a "dummy" ValueParser filled in automatically by AddArgument if
// no other ValueParser is used.
func String(v string) (interface{}, error) {
	return v, nil
}

func sscanf(s, f string, p interface{}) error {
	n, err := fmt.Sscanf(s, f, p)
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.Errorf("%d != 1", n)
	}
	return nil
}

// Action takes the name of an action instead of the action function.
// it works similarly to Python's argparse.ArgumentParser.add_argument's
// action parameter when set to a string value.
func Action(v string) ArgumentOption {
	return func(a *Argument) error {
		switch strings.ToLower(v) {
		case "append":
			a.Nargs = OneOrMore
			a.Action = Append
		case "store":
			a.Nargs = 1
			a.Action = Store
		case "store_true":
			a.Default = false
			a.Const = true
			a.Nargs = 0
			a.Action = StoreTrue
		case "store_false":
			a.Default = true
			a.Const = false
			a.Nargs = 0
			a.Action = StoreFalse
		default:
			return errors.Errorf(
				"unrecognized value")
		}
		return nil
	}
}

// ArgumentOption configures an Argument.
type ArgumentOption func(a *Argument) error

// ArgumentAction is called when an argument's values are parsed from the
// command line.
type ArgumentAction func(a *Argument, ns Namespace, vs []interface{}) error

// Append is an ArgumentAction that appends an encountered argument to
func Append(a *Argument, ns Namespace, vs []interface{}) error {
	ns.Append(a, vs...)
	return nil
}

// Store is an ArgumentAction that sets the value associated with the given
// argument.  If that argument already has a value in the given namespace,
// an error is returned.
func Store(a *Argument, ns Namespace, vs []interface{}) error {
	if v, ok := ns.Get(a); ok {
		return errors.Errorf(
			"argument %q already defined with value %v.",
			a.Dest, v)
	}
	var v interface{}
	if a.Nargs == 1 && len(vs) == 1 {
		v = vs[0]
	} else {
		v = vs
	}
	ns.Set(a, v)
	return nil
}

// StoreTrue is an ArgumentAction that stores the true value in the given
// namespace for the given argument.
func StoreTrue(a *Argument, ns Namespace, vs []interface{}) error {
	if len(vs) > 0 {
		return errors.Errorf(
			"no values expected for argument %q but got %d",
			a.Dest, len(vs))
	}
	ns.Set(a, true)
	return nil
}

// StoreFalse is an ArgumentAction that stores the false value in the given
// namespace for the given argument.
func StoreFalse(a *Argument, ns Namespace, vs []interface{}) error {
	if len(vs) > 0 {
		return errors.Errorf(
			"no values expected for argument %q but got %d",
			a.Dest, len(vs))
	}
	ns.Set(a, false)
	return nil
}

// Const sets the Const value for the given string
func Const(v interface{}) ArgumentOption {
	return func(a *Argument) error {
		return setValue(&a.Const, "Const", v)
	}
}

// Default sets the default value of an argument.
func Default(v string) ArgumentOption {
	return func(a *Argument) error {
		return setValue(&a.Default, "Default", v)
	}
}

// Help sets the help string of an argument.
func Help(v string) ArgumentOption {
	return func(a *Argument) error {
		return setValue(&a.Help, "Help", v)
	}
}

// MetaVar sets the help string of an argument.
func MetaVar(v ...string) ArgumentOption {
	return func(a *Argument) error {
		return setValue(&a.MetaVar, "MetaVar", v)
	}
}

// Nargs sets the number of values the argument can accept.
func Nargs(v int) ArgumentOption {
	return func(a *Argument) error {
		if !isValidNarg(v) {
			return errors.Errorf(
				"%d is not a valid number of arguments", v)
		}
		a.Nargs = v
		return nil
	}
}

var (
	alphaNumRegexp = regexp.MustCompile("[0-9A-Za-z]+")
)

// OptionStrings sets the arg strings.
func OptionStrings(ops ...string) ArgumentOption {
	return func(a *Argument) error {
		if len(ops) == 0 {
			return errors.Errorf("no option strings specified")
		}
		var positional, optional bool
		for _, op := range ops {
			if len(op) > 0 && op[0] == '-' {
				optional = true
			} else {
				positional = true
			}
		}
		if optional == positional {
			return errors.Errorf(
				"cannot determine if argument %s is "+
					"optional or positional",
				ops[0])
		}
		err := setValue(&a.OptionStrings, "OptionStrings", ops)
		if err != nil {
			return err
		}
		return nil
	}
}

// Required flags the Argument as required.
func Required(a *Argument) error {
	a.Required = true
	return nil
}

// Type sets the Type (actually a ValueParser function)
// of the argument.
func Type(t ValueParser) ArgumentOption {
	return func(a *Argument) error {
		if a.Type != nil {
			return errors.Errorf(
				"type already set!")
		}
		a.Type = t
		return nil
	}
}
