package summarize

import (
	"fmt"
	"sort"
	"strings"
)

// Set is a set of integers
type Set map[int]struct{}

func (set Set) String() string {
	values := make([]string, 0, len(set))
	for elem := range set {
		values = append(values, fmt.Sprintf("%d", elem))
	}
	sort.Strings(values)
	return fmt.Sprintf("Set{%s}", strings.Join(values, ", "))
}

// Add adds an element to the set
func (set Set) Add(i int) {
	set[i] = struct{}{}
}
