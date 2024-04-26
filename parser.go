package argparse

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/skillian/errors"
)

// ArgumentParser collects allowed program arguments and parses them into a
// collection.
type ArgumentParser struct {
	// Optionals is a mapping from any of the option strings to the
	// arguments defined through AddArgument.
	Optionals map[string]*Argument

	// Positionals holds the arguments deduced by their positions in the
	// command line.  They must follow all Optional arguments.
	Positionals []*Argument

	// Prog is the name of the program
	Prog string

	// Usage describes the program's usage.  Is usually generated from the
	// arguments added to the parser.
	Usage string

	// Description is the brief description under the usage showing what
	// the command is for.
	Description string

	// Epilog is trailing text added after the argument help.
	Epilog string

	// Subparsers holds a slice of sub-parsers when your top-level parser
	// has different sub-commands.
	Subparsers []*ArgumentParser

	// Parents includes a collection of ArgumentParser objects whose
	// arguments should be included in this ArgumentParser.  We're keeping
	// it simple for now, though.
	//Parents []*ArgumentParser

	//FormatterClass reflect.Type
	//PrefixChars []rune
	//FromFilePrefixChars []rune
	//ArgumentDefault *Argument
	//ConflictHandler interface{}

	// NoHelp is false when the ArgumentParser should add the -h/--help
	// arguments to generate help output.  It is analogous to the add_help
	// attribute on the ArgumentParser class in Python.
	NoHelp bool

	// boundArgs is a collection of arguments and their bound targets
	// which are set after parsing arguments.
	boundArgs
}

// NewArgumentParser constructs a new argument parser.
func NewArgumentParser(options ...ArgumentParserOption) (*ArgumentParser, error) {
	p := new(ArgumentParser)
	p.Optionals = make(map[string]*Argument)
	for _, o := range options {
		if err := o(p); err != nil {
			return nil, errors.ErrorfWithCause(
				err,
				"error initializing %[1]v "+
					"(type: %[1]T)", p,
			)
		}
	}
	// defaults:
	if p.Prog == "" {
		p.Prog = filepath.Base(os.Args[0])
	}
	return p, nil
}

// MustNewArgumentParser creates an argument parser and panics if creation fails.
func MustNewArgumentParser(options ...ArgumentParserOption) *ArgumentParser {
	p, err := NewArgumentParser(options...)
	if err != nil {
		panic(err)
	}
	return p
}

// AddArgument adds an argument to the argument parser.
func (p *ArgumentParser) AddArgument(options ...ArgumentOption) (*Argument, error) {
	a := &Argument{parser: p}
	for _, o := range options {
		if err := o(a); err != nil {
			return nil, err
		}
	}
	// defaults:
	if a.Action == nil {
		a.Action = Store
	}
	if a.Type == nil {
		a.Type = String
	}
	if a.Dest == "" {
		var dest string
		for _, op := range a.OptionStrings {
			parts := alphaNumRegexp.FindAllString(op, -1)
			full := strings.Join(parts, "")
			if len(full) > len(dest) {
				dest = full
			}
		}
		a.Dest = dest
	}
	if len(a.MetaVar) == 0 && a.Nargs != 0 && a.Choices == nil {
		upper := strings.ToUpper(a.Dest)
		if a.Nargs < 0 || a.Nargs == 1 {
			a.MetaVar = []string{upper}
		} else {
			a.MetaVar = make([]string, a.Nargs)
			for i := range a.MetaVar {
				a.MetaVar[i] = upper
			}
		}

	}
	// add to parser:
	if a.Optional() {
		for _, op := range a.OptionStrings {
			if _, ok := p.Optionals[op]; ok {
				return nil, errors.Errorf(
					"redefinition of option: %q", op)
			}
		}
		for _, op := range a.OptionStrings {
			p.Optionals[op] = a
		}
	} else {
		p.Positionals = append(p.Positionals, a)
	}

	return a, nil
}

// MustAddArgument adds an argument or panics if argument creation fails.
func (p *ArgumentParser) MustAddArgument(options ...ArgumentOption) *Argument {
	a, err := p.AddArgument(options...)
	if err != nil {
		panic(err)
	}
	return a
}

// ParseArgs parses the given args (or os.Args[1:], if none specified) to create
// a namespace from those args.  If any arguments were bound from an Argument,
// those targets are assigned to.
func (p *ArgumentParser) ParseArgs(args ...string) (Namespace, error) {
	s := parsingState{}
	if len(args) == 0 {
		args = os.Args[1:]
	}
	p.handleHelp(args)
	s.init(p, args)
	var err error
	if err = s.parse(); err != nil {
		return nil, err
	}
	if err = p.boundArgs.setValues(s.ns); err != nil {
		return nil, err
	}
	return s.ns, nil
}

// MustParseArgs must parse its arguments or it will panic.
func (p *ArgumentParser) MustParseArgs(args ...string) Namespace {
	ns, err := p.ParseArgs(args...)
	if err != nil {
		panic(err)
	}
	return ns
}

func (p *ArgumentParser) getOptionals(sorted bool) []*Argument {
	// might as well allocate enough...
	args := make([]*Argument, 0, len(p.Optionals))
	already := make(map[*Argument]struct{})
	for _, a := range p.Optionals {
		if _, ok := already[a]; ok {
			continue
		}
		args = append(args, a)
		already[a] = struct{}{}
	}
	if sorted {
		sort.Slice(args, func(i, j int) bool {
			return strings.Compare(args[i].Dest, args[j].Dest) < 0
		})
	}
	return args
}

func (p *ArgumentParser) handleHelp(args []string) {
	if p.NoHelp {
		return
	}
	for _, arg := range args {
		// TODO: Handle checking for help within subcommands.  Make
		// this more like Python's ArgumentParser in which the help
		// argument is just another argument in the set.
		if arg != "-h" && arg != "--help" {
			continue
		}
		v, err := p.FormatHelp()
		if err != nil {
			v = err.Error()
		}
		fmt.Fprintln(os.Stderr, v)
		os.Exit(1)
	}
}

// FormatHelp builds the help output into a string and returns it.
func (p *ArgumentParser) FormatHelp() (string, error) {
	s := helpingState{}
	s.init(p, 80)
	return s.format()
}

// ArgumentParserOption is a function that applies changes to the
// ArgumentParser during construction.
type ArgumentParserOption func(p *ArgumentParser) error

// Prog sets the Program name of the ArgumentParser during its construction.
func Prog(v string) ArgumentParserOption {
	return func(p *ArgumentParser) error {
		return setValue(&p.Prog, "Prog", v)
	}
}

// Usage sets the argument parser's usage string.
func Usage(v string) ArgumentParserOption {
	return func(p *ArgumentParser) error {
		return setValue(&p.Usage, "Usage", v)
	}
}

// Description sets the argument parser's description string.
func Description(v string) ArgumentParserOption {
	return func(p *ArgumentParser) error {
		return setValue(&p.Description, "Description", v)
	}
}

// Epilog sets the argument parser's description string.
func Epilog(v string) ArgumentParserOption {
	return func(p *ArgumentParser) error {
		return setValue(&p.Epilog, "Epilog", v)
	}
}

func setValue(p interface{}, name string, i interface{}) error {
	pv := reflect.ValueOf(p)
	if pv.Kind() != reflect.Ptr {
		return errors.Errorf(
			"unexpected kind: %s", pv.Kind())
	}
	t := pv.Elem()
	s := reflect.ValueOf(i)
	if !s.Type().AssignableTo(t.Type()) {
		return errors.Errorf(
			"mismatched types: %v vs. %v",
			t.Kind(), s.Kind())
	}
	t.Set(s)
	return nil
}
