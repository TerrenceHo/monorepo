package watcher

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
)

// An atomic counter
type counter struct {
	val int32
}

func (c *counter) increment() {
	atomic.AddInt32(&c.val, 1)
}

func (c *counter) value() int32 {
	return atomic.LoadInt32(&c.val)
}

func (c *counter) reset() {
	atomic.StoreInt32(&c.val, 0)
}

func renameFile(file, rename string) error {
	switch runtime.GOOS {
	case "windows", "plan9":
		return os.Rename(file, rename)
	default:
		cmd := exec.Command("mv", file, rename)
		return cmd.Run()
	}
}

func tempFile(t *testing.T, dir string) string {
	file, err := ioutil.TempFile(dir, "tempfile")
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer file.Close()
	return file.Name()
}

func tempDir(t *testing.T, dir string) string {
	directory, err := ioutil.TempDir(dir, "tempdir")
	if err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}
	return directory
}

func newRecursiveWatcher(t *testing.T) *RecursiveWatcher {
	rWatcher, err := New()
	if err != nil {
		t.Fatalf("newCursiveWatcher error: %v", err)
	}
	return rWatcher
}

func addRecursiveWatch(t *testing.T, watcher *RecursiveWatcher, dir string) {
	if err := watcher.AddRecursive(dir); err != nil {
		t.Fatalf("addRecursiveWatch err: %v", err)
	}
}

func removeRecursiveWatch(t *testing.T, watcher *RecursiveWatcher, dir string) {
	if err := watcher.RemoveRecursive(dir); err != nil {
		t.Fatalf("addRecursiveWatch err: %v", err)
	}
}

func TestRecursive(t *testing.T) {
	assert := assert.New(t)
	rWatcher := newRecursiveWatcher(t)
	defer rWatcher.Close()

	go func() {
		for err := range rWatcher.Errors {
			t.Fatalf("error received: %v", err)
		}
	}()

	rootDir := t.TempDir()
	addRecursiveWatch(t, rWatcher, rootDir)

	var createCounter, writeCounter, removeCounter, renameCounter counter
	done := make(chan bool)
	go func() {
		for event := range rWatcher.Events {
			t.Logf("event received: %s", event)
			if event.Op&fsnotify.Create == fsnotify.Create {
				createCounter.increment()
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				writeCounter.increment()
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				removeCounter.increment()
			}
			if event.Op&fsnotify.Rename == fsnotify.Rename {
				renameCounter.increment()
			}
		}
		t.Log("event loop done")
		done <- true
	}()

	// Create subdirectory
	subdir := tempDir(t, rootDir)
	t.Logf("rootdir: %s", rootDir)
	t.Logf("subdir: %s", subdir)
	testFile := filepath.Join(subdir, "TestRecursiveWatcher.testfile")
	renamedFile := filepath.Join(subdir, "TestRecursiveWatcher.renamefile")

	time.Sleep(50 * time.Millisecond)

	// Create a file in subdirectory
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Creating test file failed: %s", testFile)
	}
	f.Sync()
	time.Sleep(50 * time.Millisecond)

	// Write to file in subdirectory
	f.WriteString("test data")
	f.Sync()
	time.Sleep(50 * time.Millisecond)

	// Write to file a second time in subdirectory
	f.WriteString("second write")
	f.Sync()
	f.Close()
	time.Sleep(50 * time.Millisecond)

	// Rename file in subdirectory
	if err := os.Rename(testFile, renamedFile); err != nil {
		t.Fatalf("rename failed: %s", renamedFile)
	}

	time.Sleep(50 * time.Millisecond)
	// Write to renamed file
	f, err = os.OpenFile(renamedFile, os.O_WRONLY, 0666)
	if err != nil {
		t.Fatalf("opening renamed file failed: %s", renamedFile)
	}
	if _, err := f.WriteString("renamed file writes"); err != nil {
		t.Fatalf("writing to renamed file %q failed: %v", renamedFile, err)
	}
	f.Sync()
	f.Close()
	time.Sleep(50 * time.Millisecond)

	// Delete renamed file
	err = os.Remove(renamedFile)
	if err != nil {
		t.Fatalf("error removing renamed file: %s", renamedFile)
	}

	time.Sleep(500 * time.Millisecond) // wait for all events to catch up

	assert.EqualValuesf(
		3, createCounter.value(),
		"got %d create events, wanted %d", createCounter.value(), 3,
	)
	assert.EqualValuesf(
		3, writeCounter.value(),
		"got %d write events, wanted %d", writeCounter.value(), 3,
	)
	assert.EqualValuesf(
		1, removeCounter.value(),
		"got %d remove events, wanted %d", removeCounter.value(), 1,
	)
	assert.EqualValuesf(
		1, renameCounter.value(),
		"got %d rename events, wanted %d", renameCounter.value(), 1,
	)

	t.Log("calling Close()")
	rWatcher.Close()
	t.Log("waiting for the event channel to become closed...")
	select {
	case <-done:
		t.Log("Events channel closed")
	case <-time.After(2 * time.Second):
		t.Fatal("event stream was not closed after 2 seconds")
	}
}

func TestRecursiveRemove(t *testing.T) {
	assert := assert.New(t)
	rWatcher := newRecursiveWatcher(t)
	defer rWatcher.Close()

	go func() {
		for err := range rWatcher.Errors {
			t.Fatalf("error received: %v", err)
		}
	}()

	rootDir := t.TempDir()
	addRecursiveWatch(t, rWatcher, rootDir)

	var createCounter, writeCounter, removeCounter, renameCounter counter
	done := make(chan bool)
	go func() {
		for event := range rWatcher.Events {
			t.Logf("event received: %s", event)
			if event.Op&fsnotify.Create == fsnotify.Create {
				createCounter.increment()
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				writeCounter.increment()
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				removeCounter.increment()
			}
			if event.Op&fsnotify.Rename == fsnotify.Rename {
				renameCounter.increment()
			}
		}
		t.Log("event loop done")
		done <- true
	}()

	// Create subdirectory
	subdir := tempDir(t, rootDir)
	t.Logf("rootdir: %s", rootDir)
	t.Logf("subdir: %s", subdir)
	testFile1 := filepath.Join(subdir, "TestRemoveRecursiveWatcher.testfile")
	testFile2 := filepath.Join(subdir, "TestRemoveRecursiveWatcher.fake")

	time.Sleep(50 * time.Millisecond)

	// Create a file in subdirectory
	f, err := os.Create(testFile1)
	if err != nil {
		t.Fatalf("Creating test file failed: %s", testFile1)
	}
	f.Sync()
	time.Sleep(50 * time.Millisecond)

	// Stop watching recursively
	removeRecursiveWatch(t, rWatcher, rootDir)

	// Create a file in subdirectory, should not be watched
	f, err = os.Create(testFile2)
	if err != nil {
		t.Fatalf("Creating test file failed: %s", testFile1)
	}
	f.Sync()
	time.Sleep(50 * time.Millisecond)

	assert.EqualValuesf(
		2, createCounter.value(),
		"got %d create events, wanted %d", createCounter.value(), 2,
	)
	assert.EqualValuesf(
		0, writeCounter.value(),
		"got %d write events, wanted %d", writeCounter.value(), 0,
	)
	assert.EqualValuesf(
		0, removeCounter.value(),
		"got %d remove events, wanted %d", removeCounter.value(), 0,
	)
	assert.EqualValuesf(
		0, renameCounter.value(),
		"got %d rename events, wanted %d", renameCounter.value(), 0,
	)

	t.Log("calling Close()")
	rWatcher.Close()
	t.Log("waiting for the event channel to become closed...")
	select {
	case <-done:
		t.Log("Events channel closed")
	case <-time.After(2 * time.Second):
		t.Fatal("event stream was not closed after 2 seconds")
	}
}

func TestIsClosed(t *testing.T) {
	assert := assert.New(t)
	rWatcher := newRecursiveWatcher(t)
	if err := rWatcher.Close(); err != nil {
		t.Fatal("Close did not succeed")
	}

	err := rWatcher.Add("/fake/path")
	assert.EqualErrorf(
		err, "Add failed: RecursiveWatcher has already been closed",
		"%s != %s",
		err.Error(), "Add failed: RecursiveWatcher has already been closed",
	)
	err = rWatcher.AddRecursive("/fake/path")
	assert.EqualErrorf(
		err, "AddRecursive failed: RecursiveWatcher has already been closed",
		"%s != %s",
		err.Error(), "AddRecursive failed: RecursiveWatcher has already been closed",
	)
	err = rWatcher.Remove("/fake/path")
	assert.EqualErrorf(
		err, "Remove failed: RecursiveWatcher has already been closed",
		"%s != %s",
		err.Error(), "Remove failed: RecursiveWatcher has already been closed",
	)
	err = rWatcher.RemoveRecursive("/fake/path")
	assert.EqualErrorf(
		err, "RemoveRecursive failed: RecursiveWatcher has already been closed",
		"%s != %s",
		err.Error(), "RemoveRecursive failed: RecursiveWatcher has already been closed",
	)
}
