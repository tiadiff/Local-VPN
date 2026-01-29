import Cocoa

class StatusMenuController: NSObject, NSUserNotificationCenterDelegate {
    let statusItem = NSStatusBar.system.statusItem(withLength: NSStatusItem.variableLength)
    let menu = NSMenu()
    
    var isPaused = false
    
    override init() {
        super.init()
        
        setStatusIcon(active: true)
        
        let pauseItem = NSMenuItem(title: "Pausa VPN", action: #selector(togglePause), keyEquivalent: "p")
        pauseItem.target = self
        menu.addItem(pauseItem)
        
        let logsItem = NSMenuItem(title: "Visualizza Log", action: #selector(showLogs), keyEquivalent: "l")
        logsItem.target = self
        menu.addItem(logsItem)
        
        menu.addItem(NSMenuItem.separator())
        
        let quitItem = NSMenuItem(title: "Esci", action: #selector(quit), keyEquivalent: "q")
        quitItem.target = self
        menu.addItem(quitItem)
        
        statusItem.menu = menu
        
        NSUserNotificationCenter.default.delegate = self
    }

    func setStatusIcon(active: Bool) {
        guard let button = statusItem.button else { return }
        
        let lockEmoji = "🔒"
        let dotChar = " ●" // Unicode circle
        let dotColor = active ? NSColor.systemGreen : NSColor.systemRed
        
        let attrString = NSMutableAttributedString(string: lockEmoji + dotChar)
        
        // Lock Font (Standard)
        attrString.addAttribute(.font, value: NSFont.systemFont(ofSize: 14), range: NSRange(location: 0, length: (lockEmoji as NSString).length))
        
        // Dot Font (Smaller) and Color
        let dotRange = NSRange(location: (lockEmoji as NSString).length, length: (dotChar as NSString).length)
        attrString.addAttribute(.font, value: NSFont.boldSystemFont(ofSize: 10), range: dotRange)
        attrString.addAttribute(.foregroundColor, value: dotColor, range: dotRange)
        
        button.attributedTitle = attrString
    }
    
    func userNotificationCenter(_ center: NSUserNotificationCenter, shouldPresent notification: NSUserNotification) -> Bool {
        return true
    }
    
    func notify(title: String, subtitle: String) {
        let notification = NSUserNotification()
        notification.title = title
        notification.subtitle = subtitle
        notification.soundName = NSUserNotificationDefaultSoundName
        NSUserNotificationCenter.default.deliver(notification)
    }
    
    @objc func togglePause() {
        isPaused = !isPaused
        if isPaused {
            setStatusIcon(active: false)
            menu.item(at: 0)?.title = "Riprendi VPN"
            shell("networksetup -setsocksfirewallproxystate Wi-Fi off")
            notify(title: "VPN In Pausa", subtitle: "Traffico diretto attivato ⏸️")
        } else {
            setStatusIcon(active: true)
            menu.item(at: 0)?.title = "Pausa VPN"
            shell("networksetup -setsocksfirewallproxystate Wi-Fi on")
            notify(title: "VPN Attiva", subtitle: "Traffico protetto attivato ▶️")
        }
    }
    
    @objc func showLogs() {
        // Open client.log using the default text editor or Console
        shell("open client.log")
    }
    
    @objc func quit() {
        // Find parent process (the shell script) and send SIGINT
        let parentPID = getppid()
        if parentPID > 1 {
            kill(parentPID, SIGINT)
        }
        NSApplication.shared.terminate(self)
    }
    
    @discardableResult
    func shell(_ command: String) -> String {
        let task = Process()
        let pipe = Pipe()
        
        task.standardOutput = pipe
        task.standardError = pipe
        task.arguments = ["-c", command]
        task.launchPath = "/bin/bash"
        task.launch()
        
        let data = pipe.fileHandleForReading.readDataToEndOfFile()
        let output = String(data: data, encoding: .utf8) ?? ""
        return output
    }
}

let app = NSApplication.shared
app.setActivationPolicy(.accessory) 
let controller = StatusMenuController()
app.run()
