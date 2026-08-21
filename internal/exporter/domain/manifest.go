package domain

import (
	"sort"
	"strings"
)

func CanonicalPaths(paths []string) []string {
	out := paths
	sort.Strings(out)
	return out
}
func ManifestText(paths []string) string { return strings.Join(CanonicalPaths(paths), "\n") + "\n" }
