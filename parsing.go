package argparse

import "github.com/skillian/errors"

type parsingState struct {
	// parser is the parser whose arguments are being parsed.
	parser *ArgumentParser

	// args is the slice of all arguments
	args []string

	// argi is the index of the current argument
	argi int

	// Namespace is the currently built up argument namespace.
	ns Namespace

	// posi is the index of the currently expected positional argument.
	posi int
}

func (s *parsingState) init(p *ArgumentParser, args []string) {
	s.parser = p
	s.args = args
	s.argi = 0
	s.ns = make(Namespace)
}

func (s *parsingState) parse() error {
	for s.argi < len(s.args) {
		arg := s.args[s.argi]
		a, ok := s.parser.Optionals[arg]
		if ok {
			s.argi++
		} else {
			if s.posi >= len(s.parser.Positionals) {
				return errors.Errorf(
					"unexpected argument: %q", arg)
			}
			a = s.parser.Positionals[s.posi]
			s.posi++

		}
		if err := s.handle(a); err != nil {
			return err
		}
	}
	allArgs := append(s.parser.getOptionals(false), s.parser.Positionals...)
	for _, a := range allArgs {
		if _, ok := s.ns.Get(a); !ok {
			if a.Required {
				return errors.Errorf(
					"missing required argument %q", a.Dest)
			}
			if a.Default != nil {
				if err := a.Action(a, s.ns, []interface{}{a.Default}); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *parsingState) handle(a *Argument) error {
	args, err := s.getArgs(a)
	if err != nil {
		return err
	}
	switch a.Nargs {
	case 0:
		if len(args) != 0 {
			return errors.Errorf(
				"argument %q expected 0 values, not %d",
				a.Dest, len(args))
		}
		return a.Action(a, s.ns, []interface{}{a.Const})
	case ZeroOrOne:
		if len(args) == 0 {
			return a.Action(a, s.ns, []interface{}{a.Const})
		}
		v, err := a.createValue(args[0])
		if err != nil {
			return errors.ErrorfWithCause(err, "%v failed", a.Type)
		}
		return a.Action(a, s.ns, []interface{}{v})
	case ZeroOrMore:
		if len(args) == 0 {
			return a.Action(a, s.ns, []interface{}{a.Const})
		}
		fallthrough
	case OneOrMore:
		switch len(args) {
		case 0:
			return errors.Errorf(
				"expected one or more arguments but got zero.")
		case 1:
			v, err := a.createValue(args[0])
			if err != nil {
				return errors.ErrorfWithCause(
					err, "%v failed", a.Type)
			}
			return a.Action(a, s.ns, []interface{}{v})
		}
		fallthrough
	default:
		vs := make([]interface{}, len(args))
		for i, arg := range args {
			v, err := a.createValue(arg)
			if err != nil {
				return errors.ErrorfWithCause(
					err, "%v failed", a.Type)
			}
			vs[i] = v
		}
		return a.Action(a, s.ns, vs)
	}
}

func (s *parsingState) getArgs(a *Argument) ([]string, error) {
	r := s.remainder()
	if a.Nargs > len(r) {
		return nil, errors.Errorf(
			"not enough values for argument %q", a.Dest)
	}
	switch a.Nargs {
	case 0:
		return nil, nil
	case ZeroOrOne:
		if len(r) > 0 {
			if _, ok := s.parser.Optionals[r[0]]; ok {
				return nil, nil
			}
			s.argi++
			return r[:1], nil
		}
		return nil, nil
	case ZeroOrMore:
		if len(r) == 0 {
			return nil, nil
		}
		fallthrough
	case OneOrMore:
		if len(r) == 0 {
			return nil, errors.Errorf(
				"expected at least one value for argument %q",
				a.Dest)
		}
		i := 0
		for ; i < len(r); i++ {
			if _, ok := s.parser.Optionals[r[i]]; ok {
				break
			}
		}
		s.argi += i
		return r[:i], nil
	default:
		s.argi += a.Nargs
		return r[:a.Nargs], nil
	}
}

// remainder gets the remaining args or nil if there are no remaining args.
func (s *parsingState) remainder() []string {
	if s.argi >= len(s.args) {
		return nil
	}
	return s.args[s.argi:]
}
