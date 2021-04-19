package watcher

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"
	"github.com/fsnotify/fsnotify"
)

type RecursiveWatcherClosedError struct{}

func (e RecursiveWatcherClosedError) Error() string {
	return "RecursiveWatcher has already been closed"
}

type RecursiveWatcher struct {
	watcher  *fsnotify.Watcher
	Events   chan fsnotify.Event
	Errors   chan error
	stop     chan bool
	isClosed bool
}

func New() (*RecursiveWatcher, error) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		stackerrors.Wrap(err, "fsnotify watcher could not be created")
	}

	rWatcher := &RecursiveWatcher{
		watcher:  watcher,
		Events:   make(chan fsnotify.Event),
		Errors:   make(chan error),
		stop:     make(chan bool),
		isClosed: false,
	}
	go rWatcher.run()
	return rWatcher, nil
}

func (rw *RecursiveWatcher) Add(path string) error {
	if rw.isClosed {
		return stackerrors.Wrap(
			RecursiveWatcherClosedError{}, "Add failed",
		)
	}
	return rw.watcher.Add(path)
}

func (rw *RecursiveWatcher) AddRecursive(path string) error {
	if rw.isClosed {
		return stackerrors.Wrap(
			RecursiveWatcherClosedError{}, "AddRecursive failed",
		)
	}
	if err := rw.walkdir(path, true); err != nil {
		return stackerrors.Wrap(err, "Add Recursive walkdir failed")
	}
	return nil
}

func (rw *RecursiveWatcher) Remove(path string) error {
	if rw.isClosed {
		return stackerrors.Wrap(
			RecursiveWatcherClosedError{}, "Remove failed",
		)
	}
	return rw.watcher.Remove(path)
}

func (rw *RecursiveWatcher) RemoveRecursive(path string) error {
	if rw.isClosed {
		return stackerrors.Wrap(
			RecursiveWatcherClosedError{}, "RemoveRecursive failed",
		)
	}
	if err := rw.walkdir(path, false); err != nil {
		return stackerrors.Wrap(err, "RemoveRecursive walkdir failed")
	}
	return nil
}

func (rw *RecursiveWatcher) Close() error {
	if rw.isClosed {
		return nil
	}
	close(rw.stop)
	rw.isClosed = true
	return nil
}

func (rw *RecursiveWatcher) run() {
	for {
		select {
		case event, ok := <-rw.watcher.Events:
			if !ok {
				return
			}
			stats, err := os.Stat(event.Name)
			if err == nil && stats != nil && stats.IsDir() &&
				event.Op&fsnotify.Create != 0 {
				rw.walkdir(event.Name, true)
			}

			if event.Op&fsnotify.Remove != 0 {
				rw.watcher.Remove(event.Name)
			}

			rw.Events <- event
		case err, ok := <-rw.watcher.Errors:
			if !ok {
				return
			}
			rw.Errors <- err
		case <-rw.stop:
			rw.watcher.Close()
			close(rw.Events)
			close(rw.Errors)
			return
		}
	}
}

func (rw *RecursiveWatcher) walkdir(rootpath string, isAdd bool) error {
	walkfunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return stackerrors.Wrap(err, "walkdir-walkfunc failed")
		}
		if d.IsDir() {
			if isAdd {
				if err := rw.watcher.Add(path); err != nil {
					return stackerrors.Wrapf(
						err, "walkdir-walkfunc Add path %s failed", path,
					)
				}
			} else {
				if err := rw.watcher.Remove(path); err != nil {
					return stackerrors.Wrapf(
						err, "walkdir-walkfunc Remove path %s failed", path,
					)
				}
			}
		}
		return nil
	}

	if err := filepath.WalkDir(rootpath, walkfunc); err != nil {
		return stackerrors.Wrap(err, "walkdir failed")
	}
	return nil
}
