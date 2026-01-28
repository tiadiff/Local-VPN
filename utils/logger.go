package utils

import (
	"fmt"
	"os/exec"
	"sync"
	"time"
)

// ANSI Colors
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
	Bold   = "\033[1m"
)

var (
	lastNotification time.Time
	notifyMutex      sync.Mutex
)

func timestamp() string {
	return time.Now().Format("15:04:05")
}

func NotifyBlocked(target string) {
	notifyMutex.Lock()
	defer notifyMutex.Unlock()

	// Cooldown: 10 seconds
	if time.Since(lastNotification) < 10*time.Second {
		return
	}

	lastNotification = time.Now()

	// Send macOS Notification
	msg := fmt.Sprintf("display notification \"Blocked: %s\" with title \"Secure Tunnel\" subtitle \"Tracker Neutralized 🛡️\" sound name \"Pop\"", target)
	exec.Command("osascript", "-e", msg).Start()
}

func Info(format string, args ...interface{}) {
	fmt.Printf("%s[INFO]  %s %s\n", Cyan, fmt.Sprintf(format, args...), Reset)
}

func Success(format string, args ...interface{}) {
	fmt.Printf("%s[OK]    %s %s\n", Green, fmt.Sprintf(format, args...), Reset)
}

func Warn(format string, args ...interface{}) {
	fmt.Printf("%s[WARN]  %s %s\n", Yellow, fmt.Sprintf(format, args...), Reset)
}

func Error(format string, args ...interface{}) {
	fmt.Printf("%s[ERR]   %s %s\n", Red, fmt.Sprintf(format, args...), Reset)
}

func Secure(format string, args ...interface{}) {
	fmt.Printf("%s[SEC]   %s %s\n", Purple, fmt.Sprintf(format, args...), Reset)
}

func Block(format string, args ...interface{}) {
	fmt.Printf("%s[BLOCK] %s %s\n", Red, fmt.Sprintf(format, args...), Reset)
}

func Debug(format string, args ...interface{}) {
	// Optional: Only if verbose
	// fmt.Printf("%s[DEBUG] %s %s\n", Gray, fmt.Sprintf(format, args...), Reset)
}
