package argparse

import (
	"reflect"

	"github.com/skillian/errors"
)

// boundArg binds an argument to a pointer to a value that is set after
// all arguments are parsed.
type boundArg struct {
	*Argument
	Target reflect.Value
}

// boundArgs is a collection of bound arguments.
type boundArgs []boundArg

func (bs *boundArgs) bind(a *Argument, t interface{}) error {
	if err := bs.ensureNotAlreadyBound(a); err != nil {
		return err
	}
	v := reflect.ValueOf(t)
	if v.Kind() != reflect.Ptr {
		return errors.Errorf(
			"target must be a pointer, not %v (type: %T)",
			v.Kind(), t,
		)
	}
	v = v.Elem()
	*bs = append(*bs, boundArg{a, v})
	return nil
}

func (bs *boundArgs) ensureNotAlreadyBound(a *Argument) error {
	for _, b := range *bs {
		if b.Argument == a {
			return errors.Errorf(
				"rebinding of arguments is not yet "+
					"supported.\n\nIf you want "+
					"this, please tell %v what "+
					"your use case is.",
				maintainers,
			)
		}
	}
	return nil
}

func (bs boundArgs) setValues(ns Namespace) error {
	for _, b := range bs {
		i, ok := ns[b.Dest]
		if !ok {
			// TODO: Is this an error, or should we just use a zero
			// value or something?
			return errors.Errorf(
				"unable to get value for argument %v",
				b.Argument,
			)
		}
		v := reflect.ValueOf(i)
		if !v.Type().AssignableTo(b.Target.Type()) {
			return errors.Errorf(
				"cannot assign value %v (type: %T) to "+
					"target of type: %v",
				i, i, b.Target.Type(),
			)
		}
		b.Target.Set(v)
	}
	return nil
}
