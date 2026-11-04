package keqing

import (
	"errors"
	"os"
	"os/exec"
	"sync/atomic"
	"testing"
	"time"
)

// The clock eventually advances so a regressed wait loop fails without hanging.
func snowflakeClockSteps(steps ...int64) func() int64 {
	index := 0
	return func() int64 {
		if index < len(steps)-1 {
			value := steps[index]
			index++
			return value
		}
		return steps[len(steps)-1]
	}
}

func TestSnowflakeWaitRollbackPreservesExhaustedState(t *testing.T) {
	instance := &snowflake{
		workerId: 7, lastTimestamp: 1000, sequence: maxSequence,
		nowMillis: snowflakeClockSteps(1000, 999, 1001),
	}
	id, err := instance.nextID()
	if id != 0 || !errors.Is(err, errTimeBackwards) {
		t.Fatalf("wait rollback returned (%d, %v), want (0, %v)", id, err, errTimeBackwards)
	}
	if instance.sequence != maxSequence || instance.lastTimestamp != 1000 {
		t.Fatalf("failed generation changed exhausted state: timestamp=%d sequence=%d", instance.lastTimestamp, instance.sequence)
	}
}

func TestSnowflakeWaitRollbackRecoveryDoesNotReuseIDs(t *testing.T) {
	instance := &snowflake{workerId: 7, lastTimestamp: -1, nowMillis: func() int64 { return 1000 }}
	seen := make(map[int64]bool)
	for sequence := int64(0); sequence <= maxSequence; sequence++ {
		id, err := instance.nextID()
		if err != nil || seen[id] || id&sequenceMask != sequence {
			t.Fatalf("sequence %d: id=%d err=%v duplicate=%v", sequence, id, err, seen[id])
		}
		seen[id] = true
	}
	instance.nowMillis = snowflakeClockSteps(1000, 999, 1001)
	if _, err := instance.nextID(); !errors.Is(err, errTimeBackwards) {
		t.Fatalf("exhausted generator ignored wait rollback: %v", err)
	}
	// Recovering only to the exhausted millisecond must still wait for 1001.
	instance.nowMillis = snowflakeClockSteps(1000, 1000, 1001)
	id, err := instance.nextID()
	if err != nil || seen[id] || id>>timestampShift != 1001 || id&sequenceMask != 0 {
		t.Fatalf("recovery reused an ID or skipped next-ms sequence zero: id=%d err=%v", id, err)
	}
	seen[id] = true
	id, err = instance.nextID()
	if err != nil || seen[id] || id>>timestampShift != 1001 || id&sequenceMask != 1 {
		t.Fatalf("next recovered ID is invalid: id=%d err=%v", id, err)
	}
}

func TestSnowflakeWaitRollbackReleasesConcurrentCallers(t *testing.T) {
	var clockValue atomic.Int64
	clockValue.Store(999)
	var reads atomic.Int64
	waitEntered := make(chan struct{})
	resumeClock := make(chan struct{})
	instance := &snowflake{
		workerId: 7, lastTimestamp: 1000, sequence: maxSequence,
		nowMillis: func() int64 {
			switch reads.Add(1) {
			case 1:
				return 1000
			case 2:
				close(waitEntered)
				<-resumeClock
			}
			return clockValue.Load()
		},
	}
	const count = 32
	results := make(chan error, count)
	generate := func() { _, err := instance.nextID(); results <- err }
	go generate()
	<-waitEntered
	for i := 1; i < count; i++ {
		go generate()
	}
	close(resumeClock)
	deadline := time.NewTimer(2 * time.Second)
	defer deadline.Stop()
	for i := 0; i < count; i++ {
		select {
		case err := <-results:
			if !errors.Is(err, errTimeBackwards) {
				t.Errorf("concurrent caller returned %v instead of rollback", err)
			}
		case <-deadline.C:
			// Let the old busy-wait implementation finish before reporting failure.
			clockValue.Store(1001)
			for ; i < count; i++ {
				<-results
			}
			t.Fatal("wait rollback kept the generator lock until the wall clock caught up")
		}
	}
}

func TestRandomIDWaitRollbackPublicAPIs(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^TestRandomIDWaitRollbackProcess$", "-test.count=1")
	cmd.Env = append(os.Environ(), "KEQING_WAIT_ROLLBACK_PROCESS=1")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("public API subprocess: %v\n%s", err, output)
	}
}

func TestRandomIDWaitRollbackProcess(t *testing.T) {
	if os.Getenv("KEQING_WAIT_ROLLBACK_PROCESS") != "1" {
		return
	}
	// This subprocess runs only this test, so its process-global generator is isolated.
	instance := &snowflake{workerId: 7, lastTimestamp: 1000, sequence: maxSequence,
		nowMillis: snowflakeClockSteps(1000, 999, 1001)}
	sfMu.Lock()
	sf, sfStarted = instance, false
	sfMu.Unlock()
	if id, err := RandomIdWithError(); id != 0 || !errors.Is(err, errTimeBackwards) {
		t.Fatalf("WithError returned (%d, %v)", id, err)
	}
	instance.nowMillis = snowflakeClockSteps(1000, 999, 1001)
	defer func() {
		got, ok := recover().(error)
		if !ok || !errors.Is(got, errTimeBackwards) {
			t.Errorf("RandomId panic = %v, want rollback error", got)
		}
	}()
	RandomId()
}
