package visca

import (
	"context"
	"testing"
	"time"
)

func TestMemoryRecall(t *testing.T) {
	cam := New("10.13.34.8:52381", WithLogger(testLogger{t}))

	for i := 0; i < 3; i++ {
		if i != 0 {
			time.Sleep(5000 * time.Millisecond)
		}

		t.Logf("Memory recall %v\n", i)

		err := cam.MemoryRecall(context.Background(), byte(i))
		if err != nil {
			t.Fatalf("failed to memory recall: %s\n", err)
		}
	}
}
