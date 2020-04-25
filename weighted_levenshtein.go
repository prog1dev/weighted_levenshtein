package weighted_levenshtein

import (
  "bufio"
  "fmt"
  "log"
  "os"
  "strconv"
  "strings"
  "time"
  "unicode/utf8"
)

// func main() {
//   costs := make(map[rune]map[rune]float64)
//   // readFile(costs, "/Users/filenkoivan/go/src/analyze_es_picks/symbol_statistics_moscow_77")
//   readFile(costs, "/Users/filenkoivan/go/src/analyze_es_picks/new_symbol_statistics")
//   // log.Printf("costs len: %v", len(costs))

//   dist := distance("аэрапор внук", "аэропорт внук", costs)
//   log.Printf("distance: %v", dist)
//   // log.Printf("P(orig|fix): %v", math.Pow(1234, (-dist)))
// }

func Distance(a string, b string, costs map[rune]map[rune]float64) float64 {
  defer timeTrack(time.Now(), "distance")

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

  // initialize v0 (the previous row of distances)
  // this row is A[0][i]: edit distance for an empty s1
  // the distance is the cost of inserting each character of s2
  // we start from 1 because index 0 is already 0.
  for i := 1; i < len(v0); i++ {
    v0[i] = v0[i-1] + insertionCost(s2[i-1], costs)
  }

  // make a dummy bounds check to prevent the 2 bounds check down below.
  // The one inside the loop is particularly costly.
  _ = v0[lenS2]

  for i := 0; i < lenS1; i++ {
    s1i := s1[i]
    deletion_cost := deletionCost(s1i, costs)

    // calculate v1 (current row distances) from the previous row v0
    // first element of v1 is A[i+1][0]
    // Edit distance is the cost of deleting characters from s1
    // to match empty t.
    v1[0] = v0[0] + deletion_cost

    minv1 := v1[0]

    // use formula to fill in the rest of the row
    for j := 0; j < lenS2; j++ {
      s2j := s2[j]
      substitution_cost := float64(0)
      if s1i != s2j {
        substitution_cost = substitutionCost(s1i, s2j, costs)
      }
      insertion_cost := insertionCost(s2j, costs)

      v1[j+1] = minFloat64(
        v1[j]+insertion_cost, // Cost of insertion
        minFloat64(
          v0[j+1]+deletion_cost,    // Cost of deletion
          v0[j]+substitution_cost)) // Cost of substitution

      minv1 = minFloat64(minv1, v1[j+1])
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

func minFloat64(a, b float64) float64 {
  if a < b {
    return a
  }
  return b
}

func insertionCost(c rune, costs map[rune]map[rune]float64) float64 {
  return costs[' '][c]
}

func deletionCost(c rune, costs map[rune]map[rune]float64) float64 {
  return costs[c][' ']
}

func substitutionCost(c1 rune, c2 rune, costs map[rune]map[rune]float64) float64 {
  return costs[c1][c2]
}

func timeTrack(start time.Time, name string) {
  elapsed := time.Since(start)
  log.Printf("%s took %s", name, elapsed)
}

func readFile(costs map[rune]map[rune]float64, filepath string) {
  if !FileExists(filepath) {
    errorMessage := fmt.Sprintf("[SPELL CHECKER]: does not exist %v", filepath)
    log.Println(errorMessage)
  }

  file, _ := os.Open(filepath) //nolint

  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    parts := strings.Split(scanner.Text(), "\t")

    if len(parts) == 2 {
      s1 := ' '
      s2 := []rune(parts[0])[0]
      frequency, _ := strconv.ParseFloat(parts[1], 64)

      if costs[s1] == nil {
        costs[s1] = make(map[rune]float64)
      }

      costs[s1][s2] = frequency
    } else if len(parts) == 3 {
      s1 := []rune(parts[0])[0]
      s2 := []rune(parts[1])[0]
      frequency, _ := strconv.ParseFloat(parts[2], 64)

      if costs[s1] == nil {
        costs[s1] = make(map[rune]float64)
      }

      costs[s1][s2] = frequency
    }
  }

  file.Close()
}

func FileExists(name string) bool {
  _, err := os.Stat(name)
  if os.IsNotExist(err) {
    return false
  }
  return err == nil
}
