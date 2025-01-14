package util

import (
	"os"
	"testing"
	"time"
)

const testTimeout = 100 * time.Millisecond

func TestMain(t *testing.M) {
	waitInitialTimeout = 2 * testTimeout
	os.Exit(t.Run())
}

func TestWaiterInitialUpdateInTime(t *testing.T) {
	for _, timeout := range []time.Duration{0, testTimeout} {
		w := NewWaiter(timeout, func() {})

		go func() {
			time.Sleep(testTimeout / 2)
			w.Lock()
			w.Update()
			w.Unlock()
		}()

		w.Lock()
		if elapsed := w.Overdue(); elapsed != 0 {
			t.Errorf("expected %v, got %v", 0, elapsed)
		}
		w.Unlock()
	}
}

func TestWaiterInitialUpdateNotReceived(t *testing.T) {
	for _, timeout := range []time.Duration{0, testTimeout} {
		w := NewWaiter(timeout, func() {})

		w.Lock()
		if elapsed := w.Overdue(); elapsed != waitInitialTimeout {
			t.Errorf("expected %v, got %v", waitInitialTimeout, elapsed)
		}
		w.Unlock()
	}
}

func TestWaiterUpdateInTime(t *testing.T) {
	w := NewWaiter(testTimeout, func() {})
	w.Lock()
	w.Update()
	w.Unlock()

	go func() {
		time.Sleep(testTimeout / 2)
		w.Lock()
		w.Update()
		w.Unlock()
	}()

	w.Lock()
	defer w.Unlock()

	if elapsed := w.Overdue(); elapsed != 0 {
		t.Errorf("expected %v, got %v", 0, elapsed)
	}
}

func TestWaiterUpdateNotReceived(t *testing.T) {
	w := NewWaiter(testTimeout, func() {})
	w.Lock()
	w.Update()
	w.Unlock()

	time.Sleep(2 * testTimeout)

	w.Lock()
	defer w.Unlock()

	if elapsed := w.Overdue(); elapsed < 2*testTimeout {
		t.Errorf("expected >%v, got %v", 2*testTimeout, elapsed)
	}
}
