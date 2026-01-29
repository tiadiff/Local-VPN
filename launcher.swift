import Foundation
import AppKit

// This launcher will run inside the .app bundle
// It needs to:
// 1. Get the path to the Resources folder
// 2. Launch vpn_menubar (as user)
// 3. Launch the main vpn process (via sudo)

let bundlePath = Bundle.main.bundlePath
let resourcesPath = "\(bundlePath)/Contents/Resources"
let binPath = "\(bundlePath)/Contents/MacOS"

// 0. Cleanup old instances
let cleanupScript = """
do shell script "pkill vpn; pkill vpn_menubar; networksetup -setsocksfirewallproxystate Wi-Fi off" with administrator privileges
"""
if let script = NSAppleScript(source: cleanupScript) {
    var error: NSDictionary?
    script.executeAndReturnError(&error)
}

// 1. Launch Menubar (As current user)
let menubarProcess = Process()
menubarProcess.executableURL = URL(fileURLWithPath: "\(binPath)/vpn_menubar")
menubarProcess.currentDirectoryURL = URL(fileURLWithPath: resourcesPath)

do {
    try menubarProcess.run()
} catch {
    print("Failed to launch menubar: \(error)")
}

// 2. Launch VPN via AppleScript (to get GUI Sudo prompt)
// We use osascript to run the shell command with administrator privileges
let scriptSource = """
do shell script "cd '\(resourcesPath)' && ./vpn -mode server -port 3000 -cert server.crt -key server.key -ca ca.crt > server.log 2>&1 & cd '\(resourcesPath)' && ./vpn -mode socks -server 127.0.0.1 -port 3000 -socks 1080 -cert client.crt -key client.key -ca ca.crt > client.log 2>&1 & networksetup -setsocksfirewallproxy Wi-Fi 127.0.0.1 1080 && networksetup -setsocksfirewallproxystate Wi-Fi on" with administrator privileges
"""

if let script = NSAppleScript(source: scriptSource) {
    var error: NSDictionary?
    script.executeAndReturnError(&error)
    if let err = error {
        print("AppleScript Error: \(err)")
        // If user cancels sudo, we should probably exit
        exit(1)
    }
}

// Keep the launcher alive so it can act as a parent if needed, 
// though the background processes will keep running.
// The menubar app handles termination of the VPN when "Quit" is clicked.
RunLoop.main.run()
