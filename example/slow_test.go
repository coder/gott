package example

import (
	"testing"
	"time"
)

const basis = time.Millisecond*10

func TestA(t *testing.T) {
	time.Sleep(basis)

	t.Run("A", func(t *testing.T) {
		t.Parallel()
		time.Sleep(2 * basis)
	})
	t.Run("B", func(t *testing.T) {
		time.Sleep(basis)
	})

	time.Sleep(basis)
}

func TestB(t *testing.T) {
	time.Sleep(basis)
	t.Run("BigFail", func(t *testing.T){
		time.Sleep(basis*4)
	    t.Fail()
	})
}
