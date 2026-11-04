package keqing_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
)

// This program mirrors the README's minimal example.
func Example() {
	list := keqing.AsList([]int{3, 1}, nil)
	list.Add(2)
	list.Set(0, 4)
	copied := list.DeepCopy(func(v int) int { return v })
	values := copied.ToArray()
	keqing.Shuffle(values)
	keqing.SortArrayDesc(values)
	keqing.SortArray(values)
	fmt.Println(values)

	set := keqing.AsSet([]int{1, 1, 2}, nil)
	fmt.Println(set.DeepCopy(func(v int) int { return v }).ToArray())
	fmt.Println(keqing.ToArray[int]([]int{1, 2}))
	fmt.Println(keqing.ToObject[string]("keqing"))
	date := keqing.ParseDate("2024-02-29", "2006-01-02", nil)
	fmt.Println(date.Format("2006-01-02 MST"))
	// Output:
	// [1 2 4]
	// [1 2]
	// [1 2]
	// keqing
	// 2024-02-29 UTC
}

func ExampleCopyProperties() {
	type Person struct{ Name string }
	var dst Person
	if err := keqing.CopyProperties(&dst, Person{Name: "keqing"}); err != nil {
		panic(err)
	}
	fmt.Println(dst.Name)
	// Output: keqing
}

func ExampleFileSize() {
	dir, err := os.MkdirTemp("", "keqing-example-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "example.txt")
	if err := os.WriteFile(path, []byte("abc"), 0600); err != nil {
		panic(err)
	}
	size := keqing.FileSize(path)
	exists, err := keqing.FileExistsWithError(path)
	if err != nil {
		panic(err)
	}
	fmt.Println(size, exists)
	// Output: 3 true
}

// Run at application startup before any goroutine generates an ID. Each live
// instance needs a different workerID assigned by the deployment configuration.
func ExampleConfigureRandomID() {
	if err := keqing.ConfigureRandomID(7); err != nil {
		panic(err)
	}
	id, err := keqing.RandomIdWithError()
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
}
