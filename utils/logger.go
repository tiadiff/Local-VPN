package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	seenTrackers     map[string]bool
	trackingFile     string // Dynamic path
)

func init() {
	// Determine absolute path to the binary to keep seen_trackers.txt in the same folder
	exeProxy, err := os.Executable()
	if err != nil {
		// Fallback to CWD if we can't find executable path
		trackingFile = "seen_trackers.txt"
	} else {
		// e.g. /Applications/MAMP/htdocs/vpn_proto/seen_trackers.txt
		trackingFile = filepath.Join(filepath.Dir(exeProxy), "seen_trackers.txt")
	}

	seenTrackers = make(map[string]bool)
	loadSeenTrackers()
}

func loadSeenTrackers() {
	// fmt.Printf("DEBUG: Loading trackers from: %s\n", trackingFile)

	file, err := os.Open(trackingFile)
	if err != nil {
		// fmt.Printf("DEBUG: Could not open tracker file (might be new): %v\n", err)
		return // File might not exist yet
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			seenTrackers[line] = true
		}
	}
}

func saveSeenTracker(target string) {
	file, err := os.OpenFile(trackingFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Error saving tracker db: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(target + "\n"); err != nil {
		fmt.Printf("Error writing tracker db: %v\n", err)
	}
}

func timestamp() string {
	return time.Now().Format("15:04:05")
}

func NotifyBlocked(target string) {
	notifyMutex.Lock()
	defer notifyMutex.Unlock()

	// 1. Check if already seen (Persistence)
	if seenTrackers[target] {
		return // Already recorded.
	}

	// 2. It's a NEW tracker! Save it immediately.
	seenTrackers[target] = true
	saveSeenTracker(target)

	// 3. Cooldown check (Rate Limit Notifications Only)
	if time.Since(lastNotification) < 30*time.Second {
		// New tracker recorded in file, but we skip the "Pop" to avoid spam.
		return
	}

	lastNotification = time.Now()

	// 4. Send macOS Notification
	msg := fmt.Sprintf("display notification \"Blocked: %s\" with title \"Secure Tunnel\" subtitle \"New Tracker Detected 🛡️\" sound name \"Pop\"", target)
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
