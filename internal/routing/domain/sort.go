package domain

import "sort"

func SortTargets(targets []Target) []Target {
	out := append([]Target(nil), targets...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Priority == out[j].Priority {
			return out[i].ID < out[j].ID
		}
		return out[i].Priority < out[j].Priority
	})
	return out
}
func EnabledTargets(targets []Target) []Target {
	out := make([]Target, 0, len(targets))
	for _, t := range targets {
		if t.Enabled {
			out = append(out, t)
		}
	}
	return SortTargets(out)
}
