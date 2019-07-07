package argparse

import "github.com/skillian/logging"

var (
	logger = logging.GetLogger("argparse")
)

func init() {
	h := new(logging.ConsoleHandler)
	h.SetFormatter(logging.DefaultFormatter{})
	h.SetLevel(logging.DebugLevel)
	logger.AddHandler(h)
	logger.SetLevel(logging.DebugLevel)
}
