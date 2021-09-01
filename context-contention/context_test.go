package context

import (
	"context"
	"sync"
	"testing"
	"time"
)

const (
	waiters = 2500
)

func TestSingleContextErr(t *testing.T) {
	testSingleContextErr(t)
}

func BenchmarkSingleContextErr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testSingleContextErr(b)
	}
}

func testSingleContextErr(t testing.TB) {
	ctx, cancel := makeTestContext()
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(waiters)
	for i := 0; i < waiters; i++ {
		go func() {
			defer wg.Done()
			if err := ctx.Err(); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
}

func TestMultipleContextsErr(t *testing.T) {
	testMultipleContextsErr(t)
}

func BenchmarkMultipleContextsErr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testMultipleContextsErr(b)
	}
}

func testMultipleContextsErr(t testing.TB) {
	ctx, cancel := makeTestContext()
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(waiters)
	for i := 0; i < waiters; i++ {
		go func() {
			defer wg.Done()
			cctx, cancel := context.WithCancel(ctx)
			defer cancel()
			if err := cctx.Err(); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
}

func TestSingleContextDone(t *testing.T) {
	testSingleContextDone(t)
}

func BenchmarkSingleContextDone(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testSingleContextDone(b)
	}
}

func testSingleContextDone(t testing.TB) {
	ctx, cancel := makeTestContext()
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(waiters)
	for i := 0; i < waiters; i++ {
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				t.Error(ctx.Err())
			default:
			}
		}()
	}
	wg.Wait()
}

func makeTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second)
}
