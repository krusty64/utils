package strtools

import (
	"unicode"
	"unicode/utf8"
)

func min3int(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		} else {
			return c
		}
	} else if b < c {
		return b
	} else {
		return c
	}
}

func min3float32(a, b, c float32) float32 {
	if a < b {
		if a < c {
			return a
		} else {
			return c
		}
	} else if b < c {
		return b
	} else {
		return c
	}
}

// This is the naive implementation based on dynamic programming. Distances
// for pairs of prefixes are computed and stored in a matrix to prevent
// recomputing them again. Loosely based on the wikipedia article:
// http://en.wikipedia.org/wiki/Levenshtein_distance#Iterative_with_full_matrix
// WARNING: not utf8 compatible
func LevenshteinDistance(s, t string) int {
	// s and t are also used as subscripts for variable names
	// the variable 'is' means i_s, the current index into string s,
	// likewise 'ls' or l_s is the length of string s
	ls := len(s) + 1
	lt := len(t) + 1
	dist := make([]int, ls*lt)
	//y * w + h for position in array
	for is := 1; is < ls; is++ {
		dist[is*lt+0] = is
	}

	for it := 1; it < lt; it++ {
		dist[0*lt+it] = it
	}

	for it := 1; it < lt; it++ {
		for is := 1; is < ls; is++ {
			if s[is-1] == t[it-1] {
				dist[is*lt+it] = dist[(is-1)*lt+(it-1)]
			} else {
				dist[is*lt+it] = min3int(
					dist[(is-1)*lt+it]+1,     //deletion
					dist[is*lt+(it-1)]+1,     //insertion
					dist[(is-1)*lt+(it-1)]+1) //substitution
			}
		}
	}

	return dist[ls*lt-1]
}

// Optimization of the algorithm, taking advantage of the fact that only one
// row (and one additional element) of the matrix needs to be in stored at
// any given time. Reduces memory consumption from O(n*m) to O(min(n,m)) and
// improves performance by a factor of 4-5 for word sized inputs.
// WARNING: not utf8 compatible
func LevenshteinDistanceLowMem(s, t string) int {
	ls := len(s)
	lt := len(t)

	if ls > lt {
		s, t = t, s
		ls, lt = lt, ls
	}

	// s is the shorter string
	dist := make([]int, ls+1)
	for is := 1; is <= ls; is++ {
		dist[is] = is
	}

	for it := 0; it < lt; it++ {
		olddiag := dist[0]
		dist[0] = it + 1
		for is := 0; is < ls; is++ {
			newdist := olddiag
			if s[is] != t[it] {
				newdist = min3int(
					dist[is+1]+1, // above
					dist[is+0]+1, // left
					olddiag+1)    // diag
			}
			olddiag = dist[is+1]
			dist[is+1] = newdist
		}
	}
	return dist[ls]
}

// Generalized cost function for two runes. This can be used to model
// a 'distance' between letter. Physical distance on a keyboard, ignoring
// case differences or treating i,j,y as the same letter.
// The complexity of this callback will greatly impact the performance
// of the distance algorithm.
type RuneCost func(a, b rune) float32

// Compares two unicode runes in the most trivial way.
func EqualRune(a, b rune) float32 {
	if a == b {
		return 0.0
	}
	return 1.0
}

// Iterates through all versions (casings etc.) of a rune and compares to the
// other rune. A generalized case insensitive compare.
func EqualRuneFold(a, b rune) float32 {
	if a == b {
		return 0.0
	}
	for c := unicode.SimpleFold(a); c != a; c = unicode.SimpleFold(c) {
		//fmt.Println("Compare", c, "and", b)
		if c == b {
			return 0.0
		}
	}
	return 1.0
}

// Generalization of the LowMem algorithm. Allows custom values for
// insertions, deletions and substitutions. Usually insertions and deletions
// are weighted identical to retains symmetry for Dist(a,b) and Dist(b, a)
// This version is also compatible with utf8 encoded input strings.
// Performance will be half or worse compared to the LowMem variant, highly
// depending on which cost method is used.
func LevenshteinDistanceCustom(s, t string, ins, del, sub float32,
	cost RuneCost) float32 {

	ls := utf8.RuneCountInString(s)
	lt := utf8.RuneCountInString(t)

	if ls > lt {
		// string s is shorter, swap strings, length and costs of ins/del.
		s, t = t, s
		ls, lt = lt, ls
		ins, del = del, ins
	}

	// s is now always the shorter string
	dist := make([]float32, ls+1)
	for is := 1; is <= ls; is++ {
		dist[is] = float32(is) * del
	}

	// it and is need to be counted outside of range since 'range <string>'
	// does increment them by the byte size of the rune, not by 1.
	// rs and rt are the runes of string s and t.
	it := 0
	for _, rt := range t {
		olddiag := dist[0]
		dist[0] = float32(it+1) * ins
		is := 0
		for _, rs := range s {
			c := cost(rs, rt)
			newdist := min3float32(
				dist[is+1]+ins,
				dist[is+0]+del,
				olddiag+sub*c)
			olddiag = dist[is+1]
			dist[is+1] = newdist
			is++
		}

		it++
	}

	return dist[ls]
}
