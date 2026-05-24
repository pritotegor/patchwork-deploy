package inventory

// FilterOptions controls how hosts are selected from an inventory.
type FilterOptions struct {
	// Tags restricts results to hosts that have ALL of the given tags.
	Tags []string
	// Group restricts results to hosts belonging to the named group.
	Group string
}

// FilterByTag returns hosts that carry every tag in the provided list.
// If tags is empty, all hosts are returned.
func FilterByTag(hosts []Host, tags []string) []Host {
	if len(tags) == 0 {
		return hosts
	}
	out := make([]Host, 0, len(hosts))
	for _, h := range hosts {
		if hasAllTags(h, tags) {
			out = append(out, h)
		}
	}
	return out
}

// FilterByGroup returns hosts whose Group field matches the given name.
// An empty group string returns all hosts.
func FilterByGroup(hosts []Host, group string) []Host {
	if group == "" {
		return hosts
	}
	out := make([]Host, 0, len(hosts))
	for _, h := range hosts {
		if h.Group == group {
			out = append(out, h)
		}
	}
	return out
}

// Apply applies FilterOptions to a slice of hosts, chaining tag and group
// filters in order.
func Filter(hosts []Host, opts FilterOptions) []Host {
	result := FilterByGroup(hosts, opts.Group)
	result = FilterByTag(result, opts.Tags)
	return result
}

// hasAllTags reports whether h contains every tag in required.
func hasAllTags(h Host, required []string) bool {
	index := make(map[string]struct{}, len(h.Tags))
	for _, t := range h.Tags {
		index[t] = struct{}{}
	}
	for _, r := range required {
		if _, ok := index[r]; !ok {
			return false
		}
	}
	return true
}
