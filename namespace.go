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

// GetStrings is a helper function to get an argument's associated values as
// a slice of strings.
func (ns Namespace) GetStrings(a *Argument) ([]string, error) {
	v := ns.MustGet(a)
	vs, ok := v.([]interface{})
	if !ok {
		return nil, errors.Errorf(
			"%v (type: %T) is not %v (type: %T)", v, v, vs, vs)
	}
	ss := make([]string, len(vs))
	for i, v := range vs {
		ss[i], ok = v.(string)
		if !ok {
			return nil, errors.Errorf(
				"index %d of argument %v is %v (type: %T), "+
					"not type %T",
				i, a, v, v, "")
		}
	}
	return ss, nil
}

// MustGetStrings gets the arguments associated with a as a slice of strings.
// This function panics if a's values are not a slice of strings.
func (ns Namespace) MustGetStrings(a *Argument) []string {
	ss, err := ns.GetStrings(a)
	if err != nil {
		panic(err)
	}
	return ss
}

// Set a value in the namespace for the given Arg.
func (ns Namespace) Set(a *Argument, v interface{}) {
	ns[a.Dest] = v
}
