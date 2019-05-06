package argparse_test

import (
	"testing"

	"github.com/skillian/argparse"
	"github.com/skillian/errors"
)

func TestArgparse(t *testing.T) {
	t.Parallel()

	p := argparse.MustNewArgumentParser(
		argparse.Description("Sample argument parser"))

	count, err := p.AddArgument(
		argparse.Action("store"),
		argparse.OptionStrings("-c", "--count"),
		argparse.Type(argparse.Int),
		argparse.Help("Set the number of items to process."))

	_ = p.MustAddArgument(
		argparse.Action("store"),
		argparse.OptionStrings("-v", "--value"),
		argparse.Type(argparse.String),
		argparse.Help("Arbitrary string value argument"))

	_ = p.MustAddArgument(
		argparse.Action("store"),
		argparse.OptionStrings("-m", "--money"),
		argparse.Type(argparse.Float64),
		argparse.Help("How much money do you want?"))

	_ = p.MustAddArgument(
		argparse.Action("store"),
		argparse.OptionStrings("--debt"),
		argparse.Type(argparse.Float64),
		argparse.Help("How much debt do you want?"))

	_ = p.MustAddArgument(
		argparse.Action("store"),
		argparse.OptionStrings("source"),
		argparse.Type(argparse.String),
		argparse.Help("Here be the source parameter that does "+
			"stuff.  In fact, it does so much stuff that I can't "+
			"even begin to tell you how amazing it is."))

	_ = p.MustAddArgument(
		argparse.Action("store"),
		argparse.OptionStrings("target"),
		argparse.Type(argparse.String),
		argparse.Help("Here be the target parameter that does "+
			"stuff.  In fact, it does so much stuff that I can't "+
			"even begin to tell you how amazing it is."))

	if err != nil {
		t.Fatal(err)
	}

	ns, err := p.ParseArgs("--count", "12345", "-h")

	if err != nil {
		t.Fatal(err)
	}

	v, ok := ns.Get(count)

	if !ok {
		t.Fatal("failed to get count argument")
	}

	if v == nil {
		t.Fatal("got nil count argument")
	}

	i, ok := v.(int)

	if !ok {
		t.Fatal(errors.NewUnexpectedType(i, v))
	}

	if i != 12345 {
		t.Fatalf("expected %d but got %d", 12345, i)
	}
}
