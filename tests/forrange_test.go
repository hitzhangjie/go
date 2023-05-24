package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Note: forrange <expr>, here <expr> only evals once before forloop beginning
func Test_ForRange(t *testing.T) {
	t.Run("forrange <expr>: expr is slice", func(t *testing.T) {
		s := []int{1, 2, 3}
		// will forrange loop forever, ..., No!
		//
		// 's' is evaled and stored in a temporary variable '_tmp',
		// '_tmp' is an SliceHeader (value rather than pointer).
		//
		// when appending 10 to 's', '_tmp' won't be changed.
		for range s {
			s = append(s, 10)
		}
		assert.Len(t, s, 6)
	})

	t.Run("forrange <expr>: expr is chan", func(t *testing.T) {
		ch1 := make(chan int, 3)
		go func() {
			ch1 <- 1
			ch1 <- 2
			ch1 <- 3
			close(ch1)
		}()

		ch2 := make(chan int, 3)
		go func() {
			ch2 <- 4
			ch2 <- 5
			ch2 <- 6
			close(ch2)
		}()

		vals := []int{}
		ch := ch1
		// will forrange firstly read 1 element from ch1, then read 3 elements from ch2? ... No!
		// forrange <expr>, here expr is 'ch', and is evaled and stored in a temporary variable '_tmp',
		// even if we run 'ch = ch2' later, '_tmp' won't be changed.
		//
		// Actually, forrange here only traverse channel ch1.
		for v := range ch {
			vals = append(vals, v)
			ch = ch2
		}
		//assert.Equal(t, []int{1, 4, 5, 6}, vals)
		assert.Equal(t, []int{1, 2, 3}, vals)
	})

	t.Run("forrange <expr>: expr is array", func(t *testing.T) {
		// array is value type, rather than reference type
		nums := [...]int{0, 1, 2}
		vals := []int{}
		// when nums evaled, it is stored into a temporary variable '_tmp'. Note array is value type.
		// And when assigning an array1 to array2, all underlying buffer will be copied.
		// So even if nums[i] changed, it won't affect '_tmp'.
		for i, v := range nums {
			vals = append(vals, v)
			nums[i] = 0
		}
		//assert.Equal(t, []int{0, 0, 0}, vals)
		assert.Equal(t, []int{0, 1, 2}, vals)
	})

	t.Run("for i, v := range <expr>: variable v isn't per-variable-per-loop", func(t *testing.T) {
		type customer struct {
			id   int
			name string
		}
		// slice of customer (rather than *customer)
		customers := []customer{
			{1, "zhang"},
			{2, "wang"},
			{3, "li"},
			{4, "zhao"},
		}
		// supposing struct of customer is large, so we want to use pointer instead.
		maps := map[int]*customer{}
		// we want to reference the pointer to struct, so maps[v.id] = &v.
		// Will this work? ...No!
		//
		// the variable 'v' is declared once, all map keys will be mapped to the same *customer.
		for _, v := range customers {
			maps[v.id] = &v
		}
	})
}
