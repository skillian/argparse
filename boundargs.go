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
			if err := reflectSetValue(b.Target, reflect.Zero(b.Target.Type())); err != nil {
				return err
			}
			continue
		}
		if err := reflectSetValue(b.Target, reflect.ValueOf(i)); err != nil {
			return err
		}
	}
	return nil
}

func reflectSetValue(target, value reflect.Value) error {
	logger.Verbose(
		"assigning to %v (type: %v) from %v (type: %v)",
		target, target.Type(), value, value.Type(),
	)
	tt, vt := target.Type(), value.Type()
	switch {
	case vt.ConvertibleTo(tt):
		value = value.Convert(tt)
		fallthrough
	case vt.AssignableTo(tt):
		target.Set(value)
	case vt.Kind() == reflect.Slice && tt.Kind() == reflect.Slice:
		length := value.Len()
		ts := target
		if ts.Cap() < length {
			ts = reflect.MakeSlice(tt, 0, value.Cap())
		} else {
			ts = ts.Slice(0, 0)
		}
		tz := reflect.Zero(tt.Elem())
		for i := 0; i < length; i++ {
			ts = reflect.Append(ts, tz)
			if err := reflectSetValue(
				ts.Index(i),
				value.Index(i).Elem(),
			); err != nil {
				return err
			}
		}
		target.Set(ts)
	default:
		return errors.Errorf(
			"cannot assign value %[1]v (type: %[1]T) to "+
				"target of type: %[2]v",
			value.Interface(), target,
		)
	}
	return nil
}
