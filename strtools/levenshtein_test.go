package strtools

import (
	"strings"
	"testing"
)

func test_dist(s1, s2 string, dist int, t *testing.T) {
	dist1 := LevenshteinDistance(s1, s2)
	dist2 := LevenshteinDistanceLowMem(s1, s2)
	fdist := LevenshteinDistanceCustom(s1, s2, 1.0, 1.0, 1.0, EqualRuneFold)
	if dist1 != dist || dist2 != dist || fdist != float32(dist) {
		t.Errorf("Wrong distance for '%s' and '%s': Correct: %d, Returned: %d %d %f",
			s1, s2, dist, dist1, dist2, fdist)
	}
}

func test_custom(s1, s2 string, dist float32, t *testing.T) {
	ldist := LevenshteinDistanceCustom(s1, s2, 1.0, 2.0, 2.5, EqualRuneFold)
	if ldist != dist {
		t.Error("Wrong distance for", s1, s2, "is", ldist, "should be", dist)
	}
}

func TestLevenshteinDistance(t *testing.T) {
	// test empty strings
	test_dist("", "Default", 7, t)
	test_dist("Default", "", 7, t)
	// variations in abcd
	test_dist("abcd", "abcd", 0, t) // id
	test_dist("bcd", "abcd", 1, t)  // insertion at the beginning
	test_dist("abd", "abcd", 1, t)  // insertion in the middle
	test_dist("abc", "abcd", 1, t)  // insertion at the end
	test_dist("abcd", "bcd", 1, t)  // deletion at the beginning
	test_dist("abcd", "abd", 1, t)  // deletion in the middle
	test_dist("abcd", "abc", 1, t)  // deletion at the end
	test_dist("xbcd", "abcd", 1, t) // substitution at the beginning
	test_dist("abxd", "abcd", 1, t) // substitution in the middle
	test_dist("abcx", "abcd", 1, t) // substitution at the end
	// random samples
	test_dist("kitten", "sitting", 3, t)
	test_dist("Saturday", "Sunday", 3, t)

	// test the more advanced version with more challenging test cases
	// insertions cost 1, deletions 2 and subsitutions 2.5
	test_custom("Abcd", "abcd", 0, t)   // case insensitive id
	test_custom("Abcd", "bcd", 2, t)    // deletion
	test_custom("aBcd", "bcd", 2, t)    // deletion
	test_custom("bCd", "abcd", 1, t)    // insertion
	test_custom("abcd", "cbcd", 2.5, t) // substitution
	// Simplified Chinese, "Hello World" and "Bye World" (google translate)
	test_custom("你好世界", "再见世界", 5, t)
}

// always use strings of this length for the benchmark
const inputsize = 10

func BenchmarkLevenshteinRegular(b *testing.B) {
	b.StopTimer()
	s1 := strings.Repeat("a", inputsize)
	s2 := strings.Repeat("b", inputsize)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		LevenshteinDistance(s1, s2)
	}
}

func BenchmarkLevenshteinLowMem(b *testing.B) {
	b.StopTimer()
	s1 := strings.Repeat("a", inputsize)
	s2 := strings.Repeat("b", inputsize)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		LevenshteinDistanceLowMem(s1, s2)
	}
}

func BenchmarkLevenshteinCustom(b *testing.B) {
	b.StopTimer()
	s1 := strings.Repeat("a", inputsize)
	s2 := strings.Repeat("b", inputsize)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		LevenshteinDistanceCustom(s1, s2, 1.0, 1.0, 1.0, EqualRune)
	}
}

func BenchmarkLevenshteinUnicode(b *testing.B) {
	b.StopTimer()
	s1 := strings.Repeat("a", inputsize)
	s2 := strings.Repeat("b", inputsize)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		LevenshteinDistanceCustom(s1, s2, 1.0, 1.0, 1.0, EqualRuneFold)
	}
}
