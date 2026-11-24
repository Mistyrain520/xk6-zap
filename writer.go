package zaplogger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type writerConfig struct {
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

type managedWriter struct {
	mu sync.Mutex

	filename        string
	config          writerConfig
	writer          *lumberjack.Logger
	logger          *zap.SugaredLogger
	firstWriteError error
	reportedConfig  bool
}

var sharedWriters = struct {
	sync.Mutex
	entries map[string]*managedWriter
}{
	entries: make(map[string]*managedWriter),
}

func getSharedLogger(path string, overwrite bool, config writerConfig) *zap.SugaredLogger {
	key, filename, err := canonicalLoggerPath(path)
	if err != nil {
		panic(fmt.Sprintf("zaplogger: invalid log path %q: %v", path, err))
	}

	entry := getOrCreateWriter(key, filename, config)
	entry.mu.Lock()
	defer entry.mu.Unlock()

	if entry.config != config && !entry.reportedConfig {
		entry.reportedConfig = true
		fmt.Fprintf(os.Stderr,
			"zaplogger: ignoring conflicting configuration for %s; keeping MaxSize=%d, MaxBackups=%d, MaxAge=%d\n",
			filename, entry.config.MaxSize, entry.config.MaxBackups, entry.config.MaxAge)
	}

	if overwrite {
		if err := entry.resetLocked(); err != nil {
			panic(fmt.Sprintf("zaplogger: cannot overwrite %q: %v", filename, err))
		}
	}
	return entry.logger
}

func getOrCreateWriter(key, filename string, config writerConfig) *managedWriter {
	sharedWriters.Lock()
	if entry := sharedWriters.entries[key]; entry != nil {
		sharedWriters.Unlock()
		return entry
	}

	entry := &managedWriter{
		filename: filename,
		config:   config,
	}
	entry.logger = zap.New(zapcore.NewCore(
		getEncoder(),
		zapcore.AddSync(entry),
		zapcore.DebugLevel,
	)).Sugar()
	sharedWriters.entries[key] = entry
	sharedWriters.Unlock()
	return entry
}

func (w *managedWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.writer == nil {
		w.writer = &lumberjack.Logger{
			Filename:   w.filename,
			MaxSize:    w.config.MaxSize,
			MaxBackups: w.config.MaxBackups,
			MaxAge:     w.config.MaxAge,
		}
	}

	n, err := w.writer.Write(p)
	if err != nil && w.firstWriteError == nil {
		w.firstWriteError = err
	}
	return n, err
}

func (w *managedWriter) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	err := w.firstWriteError
	w.firstWriteError = nil
	return err
}

func (w *managedWriter) closeLocked() error {
	if w.writer == nil {
		return nil
	}
	err := w.writer.Close()
	w.writer = nil
	return err
}

func (w *managedWriter) resetLocked() error {
	if err := w.closeLocked(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(w.filename), 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(w.filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	w.firstWriteError = nil
	return nil
}

func canonicalLoggerPath(path string) (string, string, error) {
	if strings.TrimSpace(path) == "" {
		return "", "", fmt.Errorf("path is empty")
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", "", err
	}
	filename := filepath.Clean(absolute)

	if resolved, err := filepath.EvalSymlinks(filename); err == nil {
		filename = resolved
	} else {
		parent := filepath.Dir(filename)
		if resolvedParent, err := filepath.EvalSymlinks(parent); err == nil {
			filename = filepath.Join(resolvedParent, filepath.Base(filename))
		}
	}

	key := filename
	if runtime.GOOS == "windows" {
		key = strings.ToLower(key)
	}
	return key, filename, nil
}

func closeSharedLogger(path string) error {
	key, _, err := canonicalLoggerPath(path)
	if err != nil {
		return err
	}

	sharedWriters.Lock()
	entry := sharedWriters.entries[key]
	delete(sharedWriters.entries, key)
	sharedWriters.Unlock()
	if entry == nil {
		return nil
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()
	return entry.closeLocked()
}

func closeAllSharedLoggers() int {
	sharedWriters.Lock()
	entries := sharedWriters.entries
	sharedWriters.entries = make(map[string]*managedWriter)
	sharedWriters.Unlock()

	for _, entry := range entries {
		entry.mu.Lock()
		_ = entry.closeLocked()
		entry.mu.Unlock()
	}
	return len(entries)
}
