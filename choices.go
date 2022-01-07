package argparse

import "fmt"

// Choice keeps track of choices by tracking the string representation of the
// choice and the actual value.
type Choice struct {
	Key   string
	Value interface{}
}

// ArgumentChoices keeps track of a collection of argument choices.
type ArgumentChoices struct {
	items []Choice
	index map[string]int
}

// NewChoices creates a Choices collection from the given slice.
func NewChoices(choices ...Choice) *ArgumentChoices {
	dup := make([]Choice, len(choices))
	copy(dup, choices)
	return newChoices(dup)
}

// newChoices creates the actual *Choices object by taking ownership of the
// passed-in choices slice.
func newChoices(choices []Choice) *ArgumentChoices {
	cs := &ArgumentChoices{
		items: choices,
		index: make(map[string]int, len(choices)),
	}
	for i, c := range choices {
		cs.index[c.Key] = i
	}
	return cs
}

// NewChoiceValues creates a Choices collection from the given values.  The
// string representation of each value becomes that value's key in the
// collection.
func NewChoiceValues(values ...interface{}) *ArgumentChoices {
	choices := make([]Choice, len(values))
	for i, v := range values {
		key, ok := v.(string) // micro-optimization
		if !ok {
			key = fmt.Sprint(v)
		}
		choices[i] = Choice{
			Key:   key,
			Value: v,
		}
	}
	return newChoices(choices)
}

// At returns a pointer to the Choice at the given index.  Do not mutate this
// Choice's key.
func (cs *ArgumentChoices) At(index int) *Choice {
	if index < 0 || index >= len(cs.items) {
		return nil
	}
	return &cs.items[index]
}

// Len gets the number of choices within the collection.
func (cs *ArgumentChoices) Len() int { return len(cs.items) }

// Load a value from the collection by its key.
func (cs *ArgumentChoices) Load(key string) (value interface{}, ok bool) {
	var index int
	index, ok = cs.index[key]
	if !ok {
		return
	}
	value = cs.items[index].Value
	return
}
