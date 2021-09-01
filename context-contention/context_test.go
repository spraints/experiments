package context

import (
	"context"
	"sync"
	"testing"
	"time"
)

const (
	waiters = 25
	waits   = 1000
)

func TestSingleContextErr(t *testing.T) {
	f := setupSingleContextErrTest()
	defer f.stop()
	f.un(t)
}

func BenchmarkSingleContextErr(b *testing.B) {
	f := setupSingleContextErrTest()
	defer f.stop()
	for i := 0; i < b.N; i++ {
		f.un(b)
	}
}

func setupSingleContextErrTest() *singleContextErrTest {
	ctx, cancel := makeTestContext()
	return &singleContextErrTest{
		ctx:    ctx,
		cancel: cancel,
	}
}

type singleContextErrTest struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *singleContextErrTest) stop() {
	s.cancel()
}

func (s *singleContextErrTest) un(t testing.TB) {
	runWaiters(t, func(_ int) error {
		return s.ctx.Err()
	})
}

func TestMultipleContextsErr(t *testing.T) {
	f := setupMultipleContextsErrTest()
	defer f.stop()
	f.un(t)
}

func BenchmarkMultipleContextsErr(b *testing.B) {
	f := setupMultipleContextsErrTest()
	defer f.stop()
	for i := 0; i < b.N; i++ {
		f.un(b)
	}
}

func setupMultipleContextsErrTest() *multipleContextsErrTest {
	ctx, cancel := makeTestContext()
	ctxs := make([]context.Context, 0, waiters)
	cancels := make([]context.CancelFunc, 0, waiters+1)
	for i := 0; i < waiters; i++ {
		cctx, ccancel := context.WithCancel(ctx)
		ctxs = append(ctxs, cctx)
		cancels = append(cancels, ccancel)
	}
	cancels = append(cancels, cancel)
	return &multipleContextsErrTest{
		ctxs:    ctxs,
		cancels: cancels,
	}
}

type multipleContextsErrTest struct {
	ctxs    []context.Context
	cancels []context.CancelFunc
}

func (m *multipleContextsErrTest) stop() {
	for _, cancel := range m.cancels {
		cancel()
	}
}

func (m *multipleContextsErrTest) un(t testing.TB) {
	runWaiters(t, func(i int) error {
		return m.ctxs[i].Err()
	})
}

func TestSingleContextDone(t *testing.T) {
	f := setupSingleContextDoneTest()
	defer f.stop()
	f.un(t)
}

func BenchmarkSingleContextDone(b *testing.B) {
	f := setupSingleContextDoneTest()
	defer f.stop()
	for i := 0; i < b.N; i++ {
		f.un(b)
	}
}

func setupSingleContextDoneTest() *singleContextDoneTest {
	ctx, cancel := makeTestContext()
	return &singleContextDoneTest{
		ctx:    ctx,
		cancel: cancel,
	}
}

type singleContextDoneTest struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *singleContextDoneTest) stop() {
	s.cancel()
}

func (s *singleContextDoneTest) un(t testing.TB) {
	runWaiters(t, func(_ int) error {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		default:
			return nil
		}
	})
}

func makeTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

func runWaiters(t testing.TB, fn func(int) error) {
	var wg sync.WaitGroup
	wg.Add(waiters)
	for i := 0; i < waiters; i++ {
		i := i
		go func() {
			defer wg.Done()
			for n := 0; n < waits; n++ {
				if err := fn(i); err != nil {
					t.Fatal(err)
				}
			}
		}()
	}
	wg.Wait()
}
