package argparse

import "github.com/skillian/errors"

// Namespace maps argument destination names with their values.  Values
// are of the type the Argument's Type function converts them to (string, by
// default).  If an argument's Nargs are >1, then the value is a slice of
// interface{} with the elements being the type set by the argument's Type
// function.
type Namespace map[string]interface{}

// Append a set of values to the namespace.
func (ns Namespace) Append(a *Argument, vs ...interface{}) {
	var values []interface{}
	existing, ok := ns[a.Dest]
	if ok {
		values, ok = existing.([]interface{})
		if !ok {
			values = make([]interface{}, 1, len(vs)+1)
			values[0] = existing
		}
	}
	values = append(values, vs...)
	ns[a.Dest] = values
}

// Get the value from the Namespace associated with the given argument's Dest.
func (ns Namespace) Get(a *Argument) (v interface{}, ok bool) {
	v, ok = ns[a.Dest]
	return
}

// MustGet retrieves an argument from the given namespace.  It panics if the
// argument wasn't found in the namespace.
func (ns Namespace) MustGet(a *Argument) interface{} {
	v, ok := ns.Get(a)
	if !ok {
		panic(errors.Errorf("failed to get argument %q", a.Dest))
	}
	return v
}

// Set a value in the namespace for the given Arg.
func (ns Namespace) Set(a *Argument, v interface{}) {
	ns[a.Dest] = v
}
