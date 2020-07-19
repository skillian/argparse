package argparse

const (
	// Debug indicates whether or not argparse was compiled with debugging
	// enabled.  There are extremely few places where assertions are needed
	// but this is how we handle it.
	Debug = true
)

var (
	// maintainers holds a list of the maintainers of this package.
	//
	// TODO(skillian):  Is this a bad practice or just an uncommon one?
	maintainers = []string{
		"Sean Killian <skillian92@gmail.com>",
	}
)
