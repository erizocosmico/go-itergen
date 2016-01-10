package examples

import (
	"fmt"
	"math/rand"
	"time"
)

//go:generate go-itergen -t "float64" --pkg="examples" --map="int" --filter --all --some --foreach --concat --find --reverse --splice --reduce="int"
//go:generate go-itergen -t "chan float64" --pkg="examples" --map="int" --filter --foreach --concat --reduce="int" --array

func produce(ch chan float64) {
	var n int
	for n < 1000 {
		ch <- rand.Float64() * 100.
		<-time.After(1 * time.Second)
		n++
	}
	close(ch)
}

func higherThan50(f float64) bool {
	return f > 50.
}

func round(i int, f float64) interface{} {
	return int(f)
}

func main() {
	in := make(chan float64)
	go produce(in)

	out, _ := Float64ChanIter(in).Filter(higherThan50).Map(round).ToInt()
	for v := range out {
		fmt.Println(v)
	}
}
