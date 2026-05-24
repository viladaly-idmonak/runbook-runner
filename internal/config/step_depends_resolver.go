package config

import "fmt"

// DependencyOrder returns step names in topological order based on their
// declared dependencies. Returns an error if a cycle is detected or a
// referenced dependency does not exist in the known step set.
func DependencyOrder(stepNames []string, deps map[string][]string) ([]string, error) {
	// Build a set of known steps for fast lookup.
	known := make(map[string]struct{}, len(stepNames))
	for _, n := range stepNames {
		known[n] = struct{}{}
	}

	// Validate all referenced deps exist.
	for step, requires := range deps {
		for _, dep := range requires {
			if _, ok := known[dep]; !ok {
				return nil, fmt.Errorf("step %q depends on unknown step %q", step, dep)
			}
		}
	}

	// Kahn's algorithm.
	inDegree := make(map[string]int, len(stepNames))
	for _, n := range stepNames {
		inDegree[n] = 0
	}
	for _, requires := range deps {
		for _, dep := range requires {
			_ = dep // edges go dep -> step; dep must finish before step
		}
	}
	// Build adjacency: dep -> list of steps that depend on dep
	adj := make(map[string][]string)
	for step, requires := range deps {
		for _, dep := range requires {
			adj[dep] = append(adj[dep], step)
			inDegree[step]++
		}
	}

	queue := []string{}
	for _, n := range stepNames {
		if inDegree[n] == 0 {
			queue = append(queue, n)
		}
	}

	ordered := make([]string, 0, len(stepNames))
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		ordered = append(ordered, cur)
		for _, next := range adj[cur] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	if len(ordered) != len(stepNames) {
		return nil, fmt.Errorf("dependency cycle detected among steps")
	}
	return ordered, nil
}
