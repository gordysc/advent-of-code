package grid

import "fmt"

// sprint is fmt.Sprint behind a tiny wrapper so grid.go stays free of fmt.
func sprint(v any) string {
	return fmt.Sprint(v)
}
