package observe

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var encoderCfg = zapcore.EncoderConfig{
	MessageKey:   "msg",
	NameKey:      "name",
	LevelKey:     "level",
	EncodeLevel:  zapcore.LowercaseLevelEncoder,
	CallerKey:    "caller",
	EncodeCaller: zapcore.ShortCallerEncoder,
}

func Example_zapJSON() {
	zl := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			zapcore.DebugLevel,
		),
		zap.AddCaller(),
		zap.Fields(
			zap.String("version", runtime.Version()),
		),
	)

	defer func() {
		_ = zl.Sync()
	}()

	example := zl.Named("example")
	example.Debug("test debug message")
	example.Info("test info message")

	// Output:
	// {"level":"debug","name":"example","caller":"observe/zap_test.go:38","msg":"test debug message","version":"go1.25.3"}
	// {"level":"info","name":"example","caller":"observe/zap_test.go:39","msg":"test info message","version":"go1.25.3"}
}

func Example_zapConsole() {
	zl := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			zapcore.InfoLevel,
		),
	)

	defer func() {
		_ = zl.Sync()
	}()

	console := zl.Named("[console]")
	console.Info("this is logged by the logger")
	console.Debug("this is below the logger's threshold and won't log")
	console.Error("this is also logged by the logger")

	// Output:
	// info	[console]	this is logged by the logger
	// error	[console]	this is also logged by the logger

}

func Example_zapSampling() {
	zl := zap.New(
		zapcore.NewSamplerWithOptions(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderCfg),
				zapcore.Lock(os.Stdout),
				zapcore.DebugLevel,
			),
			time.Second, 1, 3,
		),
	)

	defer func() {
		_ = zl.Sync()
	}()

	for i := 0; i < 10; i++ {
		if i == 5 {
			time.Sleep(time.Second)
		}
		zl.Debug(fmt.Sprintf("%d", i))
		zl.Debug("debug message")
	}

	// Output:
	// {"level":"debug","msg":"0"}
	// {"level":"debug","msg":"debug message"}
	// {"level":"debug","msg":"1"}
	// {"level":"debug","msg":"2"}
	// {"level":"debug","msg":"3"}
	// {"level":"debug","msg":"debug message"}
	// {"level":"debug","msg":"4"}
	// {"level":"debug","msg":"5"}
	// {"level":"debug","msg":"debug message"}
	// {"level":"debug","msg":"6"}
	// {"level":"debug","msg":"7"}
	// {"level":"debug","msg":"8"}
	// {"level":"debug","msg":"debug message"}
	// {"level":"debug","msg":"9"}
}

func Example_zapDynamicDebugging() {
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	debugLevelFile := filepath.Join(tempDir, "level.debug")
	atomicLevel := zap.NewAtomicLevel()

	zl := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			atomicLevel,
		),
	)

	defer func() {
		_ = zl.Sync()
	}()

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = watcher.Close()
	}()

	err = watcher.Add(tempDir)
	if err != nil {
		log.Fatal(err)
	}

	ready := make(chan struct{})

	go func() {
		defer close(ready)
		originalLevel := atomicLevel.Level()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Name == debugLevelFile {
					switch {
					case event.Op&fsnotify.Create == fsnotify.Create:
						atomicLevel.SetLevel(zapcore.DebugLevel)
						ready <- struct{}{}
					case event.Op&fsnotify.Remove == fsnotify.Remove:
						atomicLevel.SetLevel(originalLevel)
						ready <- struct{}{}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				zl.Error(err.Error())
			}
		}
	}()

	zl.Debug("debug log 는 찍히지 않습니다.")
	df, err := os.Create(debugLevelFile)
	if err != nil {
		log.Fatal(err)
	}

	err = df.Close()
	if err != nil {
		log.Fatal(err)
	}
	<-ready

	zl.Debug("debug log 가 찍힙니다.")

	err = os.Remove(debugLevelFile)
	if err != nil {
		log.Fatal(err)
	}
	<-ready

	zl.Debug("debug log 는 찍히지 않습니다.")
	zl.Info("info log 가 찍힙니다.")

	// Output:
	// {"level":"debug","msg":"debug log 가 찍힙니다."}
	// {"level":"info","msg":"info log 가 찍힙니다."}

}
