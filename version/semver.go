package version

import (
	"errors"
	"strconv"
	"strings"
)

var ErrMalformed = errors.New("malformed version format, must be <major>.<minor>.<patch>")

type Version struct {
	Major int
	Minor int
	Patch int
	Input string
}

func (v Version) String() string {
	return v.Input
}

func Parse(s string) (Version, error) {
	p := strings.Split(strings.TrimPrefix(s, "v"), ".")

	if len(p) < 3 {
		return Version{}, ErrMalformed
	}

	major, err := strconv.ParseUint(p[0], 10, 64)
	if err != nil {
		return Version{}, ErrMalformed
	}

	minor, err := strconv.ParseUint(p[1], 10, 64)
	if err != nil {
		return Version{}, ErrMalformed
	}

	patch, err := strconv.ParseUint(p[2], 10, 64)
	if err != nil {
		return Version{}, ErrMalformed
	}

	return Version{
		Major: int(major),
		Minor: int(minor),
		Patch: int(patch),
		Input: s,
	}, nil
}

//  -1 —> v is less than other
//  0  —> v is equal to other
//  1  —> v is greather than other
func (v Version) Compare(other Version) int {
	if v.Major != other.Major {
		if v.Major > other.Major {
			return 1
		}
		return -1
	}
	if v.Minor != other.Minor {
		if v.Minor > other.Minor {
			return 1
		}
		return -1
	}
	if v.Patch != other.Patch {
		if v.Patch > other.Patch {
			return 1
		}
		return -1
	}
	return 0
}
