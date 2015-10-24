package goselect

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

type fder interface {
	Fd() uintptr
}

func TestReadWriteSync(t *testing.T) {
	const count = 500
	rrs := []io.Reader{}
	wws := []io.Writer{}
	rFDSet := &FDSet{}
	for i := 0; i < count; i++ {
		rr, ww, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		rrs = append(rrs, rr)
		wws = append(wws, ww)
	}

	go func() {
		time.Sleep(time.Second)
		for i := 0; i < count; i++ {
			fmt.Fprintf(wws[i], "hello %d", i)
			time.Sleep(time.Millisecond)
		}
	}()

	buf := make([]byte, 1024)
	for i := 0; i < count; i++ {
		rFDSet.Zero()
		for i := 0; i < count; i++ {
			rFDSet.Set(rrs[i].(fder).Fd())
		}

		if err := RetrySelect(1024, rFDSet, nil, nil, -1, 10, 10*time.Millisecond); err != nil {
			t.Fatalf("select call failed: %s", err)
		}
		for j := 0; j < count; j++ {
			if rFDSet.IsSet(rrs[j].(fder).Fd()) {
				//				println(i, j)
				if i != j {
					t.Fatalf("unexpected fd ready: %d,expected: %d", j, i)
				}
				_, err := rrs[j].Read(buf)
				if err != nil {
					t.Fatalf("read call failed: %s", err)
				}
			}
		}
	}
}
