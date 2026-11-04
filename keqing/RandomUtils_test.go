package keqing

import (
	"sync"
	"testing"
)

func TestRandomConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; n < 100; n++ {
				if v := Random(); v < 0 || v >= 1 {
					t.Errorf("Random out of range: %v", v)
				}
				if v := RandomNextInt(17); v < 0 || v >= 17 {
					t.Errorf("RandomNextInt out of range: %d", v)
				}
				if v := RandomNextInt(-7, 11); v < -7 || v >= 11 {
					t.Errorf("RandomNextInt interval: %d", v)
				}
				WinningRate(50)
			}
		}()
	}
	wg.Wait()
	if WinningRate(0) || !WinningRate(100) {
		t.Fatal("probability boundaries changed")
	}
}
