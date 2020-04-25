package weighted_levenshtein

import (
  "unicode/utf8"
)

var defaultWeight = float64(1)

func Distance(a string, b string, weights map[rune]map[rune]float64) float64 {
  if len(a) == 0 {
    return float64(utf8.RuneCountInString(b))
  }

  if len(b) == 0 {
    return float64(utf8.RuneCountInString(a))
  }

  if a == b {
    return 0
  }

  // We need to convert to []rune if the strings are non-ASCII.
  // This could be avoided by using utf8.RuneCountInString
  // and then doing some juggling with rune indices,
  // but leads to far more bounds checks. It is a reasonable trade-off.
  s1 := []rune(a)
  s2 := []rune(b)

  lenS1 := len(s1)
  lenS2 := len(s2)

  // create two work vectors of floating point (i.e. weighted) distances
  v0 := make([]float64, lenS2+1)
  v1 := make([]float64, lenS2+1)
  vtemp := make([]float64, lenS2+1)

  // we start from 1 because index 0 is already 0.
  for i := 1; i < len(v0); i++ {
    v0[i] = v0[i-1] + insertionWeight(s2[i-1], weights)
  }

  // make a dummy bounds check to prevent the 2 bounds check down below.
  // The one inside the loop is particularly costly.
  _ = v0[lenS2]

  for i := 0; i < lenS1; i++ {
    s1i := s1[i]
    deletion_weight := deletionWeight(s1i, weights)

    // calculate v1 (current row distances) from the previous row v0
    // first element of v1 is A[i+1][0]
    // Edit distance is the cost of deleting characters from s1
    // to match empty t.
    v1[0] = v0[0] + deletion_weight

    minv1 := v1[0]

    // use formula to fill in the rest of the row
    for j := 0; j < lenS2; j++ {
      s2j := s2[j]
      substitution_weight := float64(0)
      if s1i != s2j {
        substitution_weight = substitutionWeight(s1i, s2j, weights)
      }
      insertion_weight := insertionWeight(s2j, weights)

      v1[j+1] = min(
        v1[j]+insertion_weight, // Weight of insertion
        min(
          v0[j+1]+deletion_weight,    // Weight of deletion
          v0[j]+substitution_weight)) // Weight of substitution

      minv1 = min(minv1, v1[j+1])
    }

    // Copy v1 (current row) to v0 (previous row) for next iteration
    // Flip references to current and previous row
    vtemp = v0
    v0 = v1
    v1 = vtemp
  }

  return v0[lenS2]
}

func min(a, b float64) float64 {
  if a < b {
    return a
  }
  return b
}

func insertionWeight(c rune, weights map[rune]map[rune]float64) float64 {
  weight := weights[' '][c]

  if weight == 0 {
    weight = defaultWeight
  }

  return weight
}

func deletionWeight(c rune, weights map[rune]map[rune]float64) float64 {
  weight := weights[c][' ']

  if weight == 0 {
    weight = defaultWeight
  }

  return weight
}

func substitutionWeight(c1 rune, c2 rune, weights map[rune]map[rune]float64) float64 {
  weight := weights[c1][c2]

  if weight == 0 {
    weight = defaultWeight
  }

  return weight
}
