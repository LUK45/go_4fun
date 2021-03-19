/*
### Example

Input

```text
7 3
1 2 3 4 5 6 7
8 9 1 2 3 4 5
6 7 8 9 1 2 3
4 5 6 7 8 9 1
2 3 4 5 6 7 8
9 1 2 3 4 5 6
7 8 9 1 2 3 4
1 1 7 7
3 3 3 3
4 2 6 2
```

Output

```text
235
8
9
```
*/

package main

import (
	"bufio"
	"fmt"
	"strings"
	"strconv"
	"os"
)

func main() {
	input_data := [][]int{}
	input := bufio.NewScanner(os.Stdin)
	line := 0
	for input.Scan() {
		input_data = append(input_data, []int{})
		if input.Text() == "" {
		    break
		}
		strs := strings.Split(input.Text()," ")
		ari := make([]int, len(strs))
		for i :=  range ari {
			ari[i], _ = strconv.Atoi(strs[i])
		}
		input_data[line] = append(input_data[line], ari...)
		line++
	}
	N := input_data[0][0]
	Q := input_data[0][1]
	queries_index_start := N + 1
	for q := 0; q < Q; q++ {
		query := input_data[queries_index_start + q]
		sum := 0
		for row := query[0]; row <= query[2]; row++ {
			for col := query[1]-1; col <= query[3]-1; col++ {
				sum += input_data[row][col]
			}
		}
		fmt.Println(sum)
	}
}
