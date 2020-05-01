## Weighed Levenshtein

This implementation of Levenshtein allows to define different weights for character addition, deletion and substitution.

This algorithm is usually used for keyboard typing auto-correction and optical character recognition (OCR) applications.

For human typo correction, cost of substituting 'E' and 'R' is lower because these are located next to each other on an AZERTY or QWERTY keyboard. So the probability that the user mistyped the characters is higher.

If you are doing OCR correction, maybe substituting '0' for 'O' should have a smaller cost than substituting 'X' for 'O'.

## Installation

```go get -u github.com/prog1dev/weighed_levenshtein```

## Usage Example

```
package main

import (
  "fmt"
  levenshtein "github.com/prog1dev/weighed_levenshtein"
)

func main() {
  weights := make(map[rune]map[rune]float64)
  weights['s'] = make(map[rune]float64)
  weights[' '] = make(map[rune]float64)

  weights['s'][' '] = float64(0.5) // weight of addition 's' char
  weights[' ']['d'] = float64(0.6) // weight of deletion 'd' char
  weights['s']['a'] = float64(0.3) // weight of substitution 's' to 'a'

  s1 := "bananas"
  s2 := "banana"

  fmt.Printf("The distance between %v and %v is %v\n", s1, s2, levenshtein.Distance(s1, s2, weights))

  s1 = "bananas"
  s2 = "bananasd"

  fmt.Printf("The distance between %v and %v is %v\n", s1, s2, levenshtein.Distance(s1, s2, weights))

  s1 = "bananas"
  s2 = "bananaa"

  fmt.Printf("The distance between %v and %v is %v\n", s1, s2, levenshtein.Distance(s1, s2, weights))
}
```
### TODO

- add support for transposition edit
- add tests
- add docs

## LICENSE

MIT © [Ivan Filenko](https://github.com/prog1dev)
