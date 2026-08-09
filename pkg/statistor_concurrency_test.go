package pkg

import (
	"sync"
	"testing"

	"github.com/chainreactors/utils/parsers"
)

func TestStatistorRecordsConcurrentResults(t *testing.T) {
	const count = 1000
	stat := &Statistor{}
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			stat.RecordResult(NewResult(&Task{ZombieResult: &parsers.ZombieResult{}}, nil))
		}()
	}
	wg.Wait()

	if stat.Success != count {
		t.Fatalf("success = %d, want %d", stat.Success, count)
	}
}
