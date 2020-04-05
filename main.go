package main

import (
  "log"
  "time"
  "unicode/utf8"
)

func main() {
  log.Printf("distance: %v", distance("BANANAS", "BANANAQ", float32(5)))
}

func distance(a string, b string, limit float32) float32 {
  defer timeTrack(time.Now(), "distance")

  if len(a) == 0 {
    return float32(utf8.RuneCountInString(b))
  }

  if len(b) == 0 {
    return float32(utf8.RuneCountInString(a))
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

  // init the row
  // x := make([]uint16, lenS1+1)
  // we start from 1 because index 0 is already 0.
  // for i := 1; i < len(x); i++ {
  //  x[i] = uint16(i)
  // }

  // create two work vectors of floating point (i.e. weighted) distances
  v0 := make([]float32, lenS2+1)
  v1 := make([]float32, lenS2+1)
  vtemp := make([]float32, lenS2+1)

  // initialize v0 (the previous row of distances)
  // this row is A[0][i]: edit distance for an empty s1
  // the distance is the cost of inserting each character of s2
  // for (int i = 1; i < v0.length; i++) {
  //     v0[i] = v0[i - 1] + insertionCost(s2.charAt(i - 1));
  // }
  // we start from 1 because index 0 is already 0.
  for i := 1; i < len(v0); i++ {
    // x[i] = uint16(i)
    v0[i] = v0[i-1] + insertionCost(s2[i-1])
  }

  // double[] v0 = new double[s2.length() + 1];
  // double[] v1 = new double[s2.length() + 1];

  // make a dummy bounds check to prevent the 2 bounds check down below.
  // The one inside the loop is particularly costly.
  _ = v0[lenS2]

  for i := 0; i < lenS1; i++ {
    // for (int i = 0; i < s1.length(); i++) {
    s1i := s1[i]
    deletion_cost := deletionCost(s1i)

    // calculate v1 (current row distances) from the previous row v0
    // first element of v1 is A[i+1][0]
    // Edit distance is the cost of deleting characters from s1
    // to match empty t.
    v1[0] = v0[0] + deletion_cost

    minv1 := v1[0]

    // use formula to fill in the rest of the row
    // for (int j = 0; j < s2.length(); j++) {
    for j := 0; j < lenS2; j++ {
      s2j := s2[j]
      substitution_cost := float32(0)
      if s1i != s2j {
        // substitution_cost = charsub.cost(s1i, s2j)
        substitution_cost = substitutionCost(s1i, s2j)
      }
      // substitution_cost := float32(1)
      insertion_cost := insertionCost(s2j)

      // current = minUint16(minUint16(x[j-1]+1, prev+1), x[j]+1)
      // current = minUint16(minUint16(x[j-1]+1, prev+1), x[j]+substitution_cost)

      v1[j+1] = minFloat32(
        v1[j]+insertion_cost, // Cost of insertion
        minFloat32(
          v0[j+1]+deletion_cost,    // Cost of deletion
          v0[j]+substitution_cost)) // Cost of substitution

      minv1 = minFloat32(minv1, v1[j+1])
    }

    if minv1 >= limit {
      return limit
    }

    // copy v1 (current row) to v0 (previous row) for next iteration
    //System.arraycopy(v1, 0, v0, 0, v0.length);
    // Flip references to current and previous row
    vtemp = v0
    v0 = v1
    v1 = vtemp
  }

  return v0[lenS2]
}

func minUint16(a, b uint16) uint16 {
  if a < b {
    return a
  }
  return b
}

func minFloat32(a, b float32) float32 {
  if a < b {
    return a
  }
  return b
}

func insertionCost(c rune) float32 {
  // if (charchange == null) {
  return 2.0
  // } else {
  //     return charchange.insertionCost(c);
  // }
}

func deletionCost(c rune) float32 {
  // if (charchange == null) {
  // if c == 'S' {
  //  return 0.5
  // } else {
  return 1.0
  // }
  // } else {
  //     return charchange.deletionCost(c);
  // }
}

func substitutionCost(c1 rune, c2 rune) float32 {
  // if (charchange == null) {
  return 2.5
  // } else {
  //     return charchange.substitutionCost(c);
  // }
}

func timeTrack(start time.Time, name string) {
  elapsed := time.Since(start)
  log.Printf("%s took %s", name, elapsed)
}
