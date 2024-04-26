package argparse

import (
	"strings"

	"github.com/skillian/errors"
	"github.com/skillian/textwrap"
)

type helpingState struct {
	// parser holds a reference to the parser whose help output is being
	// generated
	parser *ArgumentParser

	opts []*Argument

	// columns is the number of columns wide output should be.
	columns int

	// colspcs is a precomputed slice of spaces for padding the middles of
	// strings.
	colspcs string

	// coli is the current column index in the builder.
	coli int

	// indent holds the number of columns that the help should be indented.
	indent int

	// builder builds the help string.
	builder strings.Builder
}

func (s *helpingState) init(p *ArgumentParser, columns int) {
	s.parser = p
	s.opts = p.getOptionals(true)
	s.columns = columns
	s.colspcs = strings.Repeat(" ", s.columns)
	s.indent = 16
}

func (s *helpingState) format() (v string, err error) {
	defer func() {
		if x := recover(); x != nil {
			if e, ok := x.(error); ok {
				err = errors.CreateError(e, nil, err, 0)
			} else {
				err = errors.ErrorfWithContext(err, "%v", x)
			}
		}
	}()
	s.addUsage()
	if s.parser.Description != "" {
		s.writeStrings(
			textwrap.String(
				s.parser.Description,
				s.columns,
			),
			"\n\n",
		)
	}
	s.addArguments(
		"positional arguments:",
		s.parser.Positionals,
		func(a *Argument, sb *strings.Builder) {
			sb.WriteString(a.Dest)
		})
	s.addArguments(
		"optional arguments:",
		s.opts,
		func(a *Argument, sb *strings.Builder) {
			for i, opt := range a.OptionStrings {
				if i > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(opt)
				if len(a.MetaVar) > 0 {
					sb.WriteByte(' ')
					for j, mv := range a.MetaVar {
						if j > 0 {
							sb.WriteByte(' ')
						}
						sb.WriteString(mv)
					}
				}
			}
			if a.Choices != nil {
				for j, limit := 0, a.Choices.Len(); j < limit; j++ {
					ch := a.Choices.At(j)
					if j == 0 {
						sb.WriteString(" [ ")
					} else {
						sb.WriteString(" | ")
					}
					sb.WriteString(ch.Key)
					if j == limit-1 {
						sb.WriteString(" ]")
					}
				}
			}
		})
	if len(s.parser.Epilog) > 0 {
		s.builder.WriteByte('\n')
		s.builder.WriteString(
			textwrap.String(s.parser.Epilog, s.columns),
		)
	}
	return s.builder.String(), nil
}

func (s *helpingState) addUsage() {
	s.writeStrings("usage: ", s.parser.Prog, " ")
	s.coli = s.builder.Len()
	width := s.columns - s.coli
	if width <= 0 {
		s.writeStrings("\n")
		s.coli = s.indent
		width = s.columns - s.coli
	}
	var usages []string
	for _, a := range s.opts {
		usages = append(usages, s.argUsage(a))
	}
	for _, a := range s.parser.Positionals {
		usages = append(usages, s.argUsage(a))
	}
	s.writeStrings(
		strings.Join(
			textwrap.SliceLines(usages, width, " "),
			"\n"+s.colspcs[:s.columns-width]),
		"\n\n")
}

func (s *helpingState) addArguments(prefix string, args []*Argument, sel helpHeaderSelector) {
	if len(args) == 0 {
		return
	}
	s.writeStrings(prefix, "\n")
	s.coli = 0
	for _, a := range args {
		s.writeStrings("  ")
		beforeHead := s.builder.Len()
		sel(a, &s.builder)
		s.coli = 2 + (s.builder.Len() - beforeHead)
		if s.coli <= s.indent-2 {
			s.writeStrings(s.colspcs[:s.indent-s.coli])
		} else {
			s.writeStrings("\n", s.colspcs[:s.indent])
		}
		s.coli = s.indent
		for _, v := range strings.Split(textwrap.String(a.Help, s.columns-s.indent), "\n") {
			s.writeStrings(s.colspcs[:s.indent-s.coli], v, "\n")
			s.coli = 0
		}
		if a.Choices != nil {
			s.writeSpaces(s.indent)
			s.writeString("choices:\n")
			choiceIndent := 2 * s.indent
			for i, limit := 0, a.Choices.Len(); i < limit; i++ {
				c := a.Choices.At(i)
				s.writeSpaces(s.indent)
				s.writeString(c.Key)
				s.coli = s.indent + len(c.Key)
				if s.coli < choiceIndent {
					s.writeSpaces(choiceIndent - s.coli)
				} else {
					s.writeByte('\n')
					s.writeSpaces(choiceIndent)
				}
				s.coli = choiceIndent
				for _, v := range strings.Split(textwrap.String(
					c.Help, s.columns-choiceIndent,
				), "\n") {
					s.writeSpaces(choiceIndent - s.coli)
					s.writeString(v)
					s.writeByte('\n')
					s.coli = 0
				}
			}
		}
	}
	s.writeStrings("\n")
}

type helpHeaderSelector func(a *Argument, sb *strings.Builder)

func (s *helpingState) argUsage(a *Argument) string {
	var parts []string
	if a.Optional() {
		parts = append(parts, "[", getShortestArgOptionString(a))
		parts = append(parts, a.MetaVar...)
		if a.Choices != nil {
			for i, limit := 0, a.Choices.Len(); i < limit; i++ {
				if i > 0 {
					parts = append(parts, "|")
				}
				parts = append(parts, a.Choices.At(i).Key)
			}
		}
		parts = append(parts, "]")
	} else {
		parts = a.MetaVar
	}
	return strings.Join(parts, " ")
}

// TODO: name these write* methods mustWrite* because they panic

func (s *helpingState) writeByte(b byte) {
	if err := s.builder.WriteByte(b); err != nil {
		panic(err)
	}
}

func (s *helpingState) writeSpaces(n int) {
	s.builder.Grow(n)
	for i := 0; i < n; i++ {
		if err := s.builder.WriteByte(' '); err != nil {
			panic(err)
		}
	}
}

func (s *helpingState) writeString(v string) {
	if _, err := s.builder.WriteString(v); err != nil {
		panic(err)
	}
}

func (s *helpingState) writeStrings(vs ...string) {
	{
		n := 0
		for _, v := range vs {
			n += len(v)
		}
		s.builder.Grow(n)
	}
	for _, v := range vs {
		if _, err := s.builder.WriteString(v); err != nil {
			panic(err)
		}
	}
}

func getShortestArgOptionString(a *Argument) string {
	switch len(a.OptionStrings) {
	case 0:
		return ""
	case 1:
		return a.OptionStrings[0]
	default:
		short := a.OptionStrings[0]
		for _, s := range a.OptionStrings[1:] {
			if len(s) < len(short) {
				short = s
			}
		}
		return short
	}
}
