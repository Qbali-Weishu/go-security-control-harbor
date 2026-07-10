package version

import (
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Parts [3]int
}

func Parse(raw string) (Version, error) {
	fields := strings.Split(raw, ".")
	if len(fields) != 3 {
		return Version{}, fmt.Errorf("expected semantic version with three numeric parts: %q", raw)
	}
	var v Version
	for i, part := range fields {
		n, err := strconv.Atoi(part)
		if err != nil {
			return Version{}, fmt.Errorf("invalid numeric segment %q in %q: %w", part, raw, err)
		}
		v.Parts[i] = n
	}
	return v, nil
}

func Compare(left, right string) (int, error) {
	l, err := Parse(left)
	if err != nil {
		return 0, err
	}
	r, err := Parse(right)
	if err != nil {
		return 0, err
	}
	for i := 0; i < len(l.Parts); i++ {
		if l.Parts[i] < r.Parts[i] {
			return -1, nil
		}
		if l.Parts[i] > r.Parts[i] {
			return 1, nil
		}
	}
	return 0, nil
}

// InRange reports whether candidate falls within [min, maxExclusive).
func InRange(candidate, min, maxExclusive string) (bool, error) {
	if min != "" {
		cmp, err := Compare(candidate, min)
		if err != nil {
			return false, err
		}
		if cmp < 0 {
			return false, nil
		}
	}
	if maxExclusive != "" {
		cmp, err := Compare(candidate, maxExclusive)
		if err != nil {
			return false, err
		}
		if cmp > 0 {
			return false, nil
		}
	}
	return true, nil
}
