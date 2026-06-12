package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// dailyWriter 按天分割的日志写入器
type dailyWriter struct {
	mu      sync.Mutex
	dir     string
	maxAge  int      // 保留天数
	curFile *os.File // 当前日志文件
	curDate string   // 当前日期（YYYY-MM-DD）
	curPath string   // 当前文件路径
}

func newDailyWriter(dir string, maxAge int) (*dailyWriter, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	w := &dailyWriter{dir: dir, maxAge: maxAge}
	if err := w.rotate(time.Now()); err != nil {
		return nil, err
	}
	// 启动清理协程
	go w.cleanupLoop()
	return w, nil
}

// rotate 切换到指定日期的日志文件
func (w *dailyWriter) rotate(t time.Time) error {
	dateStr := t.Format("2006-01-02")
	if w.curDate == dateStr && w.curFile != nil {
		return nil
	}

	newPath := path.Join(w.dir, fmt.Sprintf("app-%s.log", dateStr))
	f, err := os.OpenFile(newPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// 关闭旧文件
	if w.curFile != nil {
		w.curFile.Close()
	}

	w.curFile = f
	w.curDate = dateStr
	w.curPath = newPath
	return nil
}

// Write 实现io.Writer接口
func (w *dailyWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	today := now.Format("2006-01-02")

	// 日期变更时切换文件
	if w.curDate != today {
		if err := w.rotate(now); err != nil {
			return 0, err
		}
	}

	return w.curFile.Write(p)
}

// cleanupLoop 定期清理过期日志
func (w *dailyWriter) cleanupLoop() {
	// 启动时先清理一次
	w.cleanOldLogs()

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		w.cleanOldLogs()
	}
}

// cleanOldLogs 删除超过maxAge天的日志文件
func (w *dailyWriter) cleanOldLogs() {
	w.mu.Lock()
	defer w.mu.Unlock()

	cutoff := time.Now().AddDate(0, 0, -w.maxAge)

	entries, err := os.ReadDir(w.dir)
	if err != nil {
		return
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasPrefix(name, "app-") || !strings.HasSuffix(name, ".log") {
			continue
		}
		files = append(files, name)
	}

	// 按文件名排序（文件名包含日期，排序后方便处理）
	sort.Strings(files)

	for _, name := range files {
		// 从文件名提取日期：app-2024-01-15.log -> 2024-01-15
		dateStr := strings.TrimPrefix(name, "app-")
		dateStr = strings.TrimSuffix(dateStr, ".log")

		fileDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		if fileDate.Before(cutoff) {
			fullPath := filepath.Join(w.dir, name)
			os.Remove(fullPath)
		}
	}
}

// Sync 实现zapcore.WriteSyncer接口
func (w *dailyWriter) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.curFile != nil {
		return w.curFile.Sync()
	}
	return nil
}

func Init(logDir string) error {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 按天分割日志文件，保留最新20天
	writer, err := newDailyWriter(logDir, 20)
	if err != nil {
		return err
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 文件使用JSON格式
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		zapcore.DebugLevel,
	)

	// 控制台使用友好格式
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	// 合并输出到文件和控制台
	core := zapcore.NewTee(fileCore, consoleCore)
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// 确保dailyWriter实现了io.Writer和zapcore.WriteSyncer接口
var _ io.Writer = (*dailyWriter)(nil)
var _ zapcore.WriteSyncer = (*dailyWriter)(nil)
