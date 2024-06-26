package argparse

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/skillian/errors"
)

// Argument holds the definition of an argument.
type Argument struct {
	// parser holds a reference back to the parser that instantiated the
	// argument.
	parser *ArgumentParser

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

	// Choices holds an optional collection of allowed choices for this
	// Argument.  Choices is nil if no set of allowed values was provided.
	Choices *ArgumentChoices
}

// Bind the argument's parsed value into the given pointer.
func (a *Argument) Bind(target interface{}) error {
	return a.parser.boundArgs.bind(a, target)
}

// MustBind panics if Binding an argument fails.
func (a *Argument) MustBind(target interface{}) {
	if err := a.Bind(target); err != nil {
		panic(err)
	}
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
	var i int64
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
	key := strings.TrimSpace(strings.ToLower(v))
	act, ok := actions[key]
	if !ok {
		return func(a *Argument) error {
			return errors.Errorf(
				"unrecognized %v: %q", "Action", v,
			)
		}
	}
	return ActionFunc(act)
}

// ActionFunc lets you specify an action function value instead of just a string
// key of an action function.
func ActionFunc(f ArgumentAction) ArgumentOption {
	return func(a *Argument) error {
		a.Action = f
		switch f {
		case Store:
			if a.Nargs < 1 {
				a.Nargs = 1
			}
		case StoreTrue:
			a.Default = false
			a.Const = true
			a.Nargs = 0
		case StoreFalse:
			a.Default = true
			a.Const = false
			a.Nargs = 0
		}
		return nil
	}
}

// ArgumentOption configures an Argument.
type ArgumentOption func(a *Argument) error

// ArgumentAction is called when an argument's values are parsed from the
// command line.
type ArgumentAction interface {
	Name() string
	UpdateNamespace(a *Argument, ns Namespace, vs []interface{}) error
}

type argumentActionStruct struct {
	name            string
	updateNamespace func(a *Argument, ns Namespace, vs []interface{}) error
}

func newArgumentActionStruct(name string, f func(a *Argument, ns Namespace, vs []interface{}) error) *argumentActionStruct {
	if _, ok := actions[name]; ok {
		panic("redefinition of argument action: " + name)
	}
	s := &argumentActionStruct{name: name, updateNamespace: f}
	actions[name] = s
	return s
}

func (s argumentActionStruct) Name() string { return s.name }
func (s argumentActionStruct) UpdateNamespace(a *Argument, ns Namespace, args []interface{}) error {
	return s.updateNamespace(a, ns, args)
}

var (
	actions = make(map[string]ArgumentAction, 4)

	// Append is an ArgumentAction that appends an encountered argument to
	Append ArgumentAction = newArgumentActionStruct(
		"append",
		func(a *Argument, ns Namespace, args []interface{}) error {
			vs, err := a.defaultCreateValues(args)
			if err != nil {
				return err
			}
			ns.Append(a, getArgValueForNS(a, vs))
			return nil
		},
	)

	// Store is an ArgumentAction that sets the value associated with the
	// given argument.  If that argument already has a value in the given
	// namespace, an error is returned.
	Store ArgumentAction = newArgumentActionStruct(
		"store",
		func(a *Argument, ns Namespace, args []interface{}) error {
			if v, ok := ns.Get(a); ok {
				return errors.Errorf(
					"argument %q already defined with value %v.",
					a.Dest, v)
			}
			vs, err := a.defaultCreateValues(args)
			if err != nil {
				return err
			}
			ns.Set(a, getArgValueForNS(a, vs))
			return nil
		},
	)

	// StoreTrue is an ArgumentAction that stores the true value in the
	// given namespace for the given argument.
	StoreTrue ArgumentAction = newArgumentActionStruct(
		"store_true",
		func(a *Argument, ns Namespace, args []interface{}) error {
			if len(args) != 1 {
				return errors.Errorf(
					"one value expected for argument %q but got %d: %#v",
					a.Dest, len(args), args)
			}
			ns.Set(a, args[0])
			return nil
		},
	)

	// StoreFalse is an ArgumentAction that stores the false value in the given
	// namespace for the given argument.
	StoreFalse ArgumentAction = newArgumentActionStruct(
		"store_false",
		func(a *Argument, ns Namespace, args []interface{}) error {
			if len(args) != 1 {
				return errors.Errorf(
					"one value expected for argument %q but got %d: %#v",
					a.Dest, len(args), args)
			}
			ns.Set(a, args[0])
			return nil
		},
	)
)

func getArgValueForNS(a *Argument, vs []interface{}) interface{} {
	if a.Nargs == 1 && len(vs) == 1 {
		return vs[0]
	}
	return vs
}

// Choices sets the argument's choices.
func Choices(choices ...Choice) ArgumentOption {
	return func(a *Argument) error {
		if len(a.MetaVar) != 0 {
			return errors.Errorf("Choices take the place of a MetaVar")
		}
		a.Choices = NewChoices(choices...)
		return nil
	}
}

// ChoiceValues sets the argument's choices.
func ChoiceValues(values ...interface{}) ArgumentOption {
	return func(a *Argument) error {
		a.Choices = NewChoiceValues(values...)
		return nil
	}
}

// Const sets the Const value for the given string
func Const(v interface{}) ArgumentOption {
	return func(a *Argument) error {
		return setValue(&a.Const, "Const", v)
	}
}

// Default sets the default value of an argument.
func Default(v interface{}) ArgumentOption {
	return func(a *Argument) error {
		return setValue(&a.Default, "Default", v)
	}
}

// Dest sets the destination name in the parsed argument namespace.
func Dest(v string) ArgumentOption {
	return func(a *Argument) error {
		return setValue(&a.Dest, "Dest", v)
	}
}

// Help sets the help string of an argument.
func Help(format string, args ...interface{}) ArgumentOption {
	v := format
	if len(args) >= 0 {
		v = fmt.Sprintf(format, args...)
	}
	return func(a *Argument) error {
		return setValue(&a.Help, "Help", v)
	}
}

// MetaVar sets the help string of an argument.
func MetaVar(v ...string) ArgumentOption {
	return func(a *Argument) error {
		if a.Choices != nil {
			return errors.Errorf("Choices take the place of a MetaVar")
		}
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

func (a *Argument) defaultCreateValues(args []interface{}) (vs []interface{}, err error) {
	vs = make([]interface{}, len(args))
	if a.Choices != nil {
		for i, arg := range args {
			v, ok := a.Choices.Load(stringOf(arg))
			if !ok {
				return nil, errors.Errorf(
					"invalid choice %q for %v", v, a.Dest,
				)
			}
			vs[i] = v
		}
		return
	}
	for i, arg := range args {
		if vs[i], err = a.Type(stringOf(arg)); err != nil {
			return
		}
	}
	return
}

func stringOf(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}
