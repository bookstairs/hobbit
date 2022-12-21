package logger

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/syhily/hobbit/config"
)

func Test_Logger(t *testing.T) {
	logger1 := GetLogger("PKG", "Log")
	RunningAtomicLevel.SetLevel(zapcore.DebugLevel)

	fmt.Println(White.Add("white"))
	logger1.Warn("warn for test", String("count", "1"), Reflect("v1", map[string]string{"a": "1"}))
	logger1.Info("info for test", Uint16("value", 1), Int32("v1", 2),
		Int64("v2", 2), Any("v3", 3))
	logger1.Debug("debug for test", Uint32("value", 2))
	logger1.Error("error for test", Error(fmt.Errorf("error")))

	assert.NotNil(t, defaultLogger)

	logger3 := GetLogger("PKG", "")
	logger3.Error("error test")
}

func Test_Access_logger(t *testing.T) {
	assert.Nil(t, InitLogger(config.Logging{Level: "debug"}, "access.log"))
	logger1 := GetLogger(HTTPModule, "Access")
	logger1.Info("access log")
	isTerminal = true
	defer func() {
		isTerminal = false
	}()
	logger1.Info("access log")
}

func Test_Level_String(t *testing.T) {
	isTerminal = true
	defer func() {
		isTerminal = false
	}()
	assert.Equal(t, "[35mDEBUG[0m", LevelString(zapcore.DebugLevel))   //nolint:stylecheck
	assert.Equal(t, "[31mDPANIC[0m", LevelString(zapcore.DPanicLevel)) //nolint:stylecheck
	assert.Equal(t, "[32mINFO[0m", LevelString(zapcore.InfoLevel))     //nolint:stylecheck
	assert.Equal(t, "[33mWARN[0m", LevelString(zapcore.WarnLevel))     //nolint:stylecheck
	assert.Equal(t, "[31mERROR[0m", LevelString(zapcore.ErrorLevel))   //nolint:stylecheck
	isTerminal = false
	assert.Equal(t, "ERROR", LevelString(zapcore.ErrorLevel))
}

func Test_Logger_Stack(t *testing.T) {
	panicFunc := func() {
		defer func() {
			if r := recover(); r != nil {
				GetLogger("PKG", "Logger").
					getInitializedOrDefaultLogger().Panic("panic stack", Stack())
			}
		}()
		panic("test-panic")
	}
	assert.Panics(t, panicFunc)
}

func Test_IsTerminal(t *testing.T) {
	assert.False(t, IsTerminal(os.Stdout))
}

func Test_IsDebug(t *testing.T) {
	RunningAtomicLevel.SetLevel(zapcore.InfoLevel)
	assert.False(t, IsDebug())
	RunningAtomicLevel.SetLevel(zapcore.DebugLevel)
	assert.True(t, IsDebug())
}

func Test_InitLogger(t *testing.T) {
	assert.NotNil(t, GetLogger("PKG", "Log").getInitializedOrDefaultLogger())

	cfg1 := config.Logging{Level: "LLL"}
	assert.NotNil(t, InitLogger(cfg1, "test.log"))

	cfg2 := config.NewDefaultLogging()
	assert.Nil(t, InitLogger(*cfg2, "test.log"))
	thisLogger := GetLogger("PKG", "Log")
	assert.NotNil(t, thisLogger.getInitializedOrDefaultLogger())
	assert.NotNil(t, thisLogger.getInitializedOrDefaultLogger())

	cfg3 := config.Logging{Level: "info"}
	assert.Nil(t, InitLogger(cfg3, "test.log"))

	cfg4 := config.Logging{Level: "debug"}
	assert.Nil(t, InitLogger(cfg4, "test.log"))

	isTerminal = true
	defer func() {
		isTerminal = false
	}()
	assert.Nil(t, InitLogger(cfg4, "test.log"))
}
