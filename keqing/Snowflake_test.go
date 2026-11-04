package keqing

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestRandomIDConcurrentFirstUse(t *testing.T) {
	const count = 1000
	start := make(chan struct{})
	ids := make(chan int64, count)
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); <-start; ids <- RandomId() }()
	}
	close(start)
	wg.Wait()
	close(ids)
	seen := make(map[int64]bool)
	for id := range ids {
		if id == 0 || seen[id] {
			t.Errorf("zero or duplicate ID: %d", id)
		}
		seen[id] = true
	}
}

func TestRandomIDConfiguredProcess(t *testing.T) {
	workerText := os.Getenv("KEQING_TEST_WORKER")
	if workerText == "" {
		return
	}
	worker, err := strconv.ParseInt(workerText, 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	if ConfigureRandomID(-1) == nil || ConfigureRandomID(1024) == nil {
		t.Fatal("invalid worker accepted")
	}
	if err := ConfigureRandomID(worker); err != nil {
		t.Fatal(err)
	}
	start, err := strconv.ParseInt(os.Getenv("KEQING_TEST_START"), 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	for time.Now().UnixMilli() < start {
		time.Sleep(time.Millisecond)
	}
	ids := make([]int64, 10000)
	for i := range ids {
		ids[i], err = RandomIdWithError()
		if err != nil {
			t.Fatal(err)
		}
		if (ids[i]>>workerShift)&maxWorkerId != worker {
			t.Fatal("wrong worker bits")
		}
	}
	if ConfigureRandomID(worker) == nil {
		t.Fatal("reconfiguration after use succeeded")
	}
	// 只修改此隔离子进程的实例，验证回拨错误，绝不修改系统时钟。
	sf.lastTimestamp = time.Now().Add(time.Hour).UnixMilli()
	if _, err := RandomIdWithError(); err == nil {
		t.Fatal("clock rollback was ignored")
	}
	requirePanic(t, func() { RandomId() })
	data, err := json.Marshal(ids)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(os.Getenv("KEQING_TEST_IDS"), data, 0600); err != nil {
		t.Fatal(err)
	}
}

func TestRandomIDMultipleProcesses(t *testing.T) {
	dir := t.TempDir()
	start := strconv.FormatInt(time.Now().Add(time.Second).UnixMilli(), 10)
	var commands []*exec.Cmd
	var paths []string
	for _, worker := range []string{"0", "1023"} {
		path := filepath.Join(dir, worker+".json")
		cmd := exec.Command(os.Args[0], "-test.run=^TestRandomIDConfiguredProcess$", "-test.count=1")
		cmd.Env = append(os.Environ(), "KEQING_TEST_WORKER="+worker, "KEQING_TEST_START="+start, "KEQING_TEST_IDS="+path)
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		commands, paths = append(commands, cmd), append(paths, path)
	}
	for _, cmd := range commands {
		if err := cmd.Wait(); err != nil {
			t.Fatal(err)
		}
	}
	seen := make(map[int64]bool)
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var ids []int64
		if err := json.Unmarshal(data, &ids); err != nil {
			t.Fatal(err)
		}
		for _, id := range ids {
			if id == 0 || seen[id] {
				t.Fatalf("duplicate/zero ID %d", id)
			}
			seen[id] = true
		}
	}
}
