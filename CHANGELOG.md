# Changelog

All notable changes to this project will be documented in this file.

## [v1.0.5] - macOS App Edition 🚀
### Added
- **Native macOS App Bundle**: Created `SecureTunnel.app` for a native experience.
- **GUI Launcher**: New Swift-based launcher that handles administrator privileges via standard macOS dialogs.
- **Status Bar Integration**: Native 🔒 icon with live status (● Green for Active, ● Red for Paused).
- **Log Visibility**: Added "Show Logs" menu item to monitor VPN activity in real-time.
- **Automatic Cleanup**: The app now automatically detects and terminates old VPN processes on startup to prevent port conflicts (Port 1080).
- **Smart Shutdown**: Closing the app via "Quit" or `Cmd+Q` automatically stops all VPN services and restores network settings.

### Changed
- Improved the build system to generate icons from emojis and package the full `.app` structure.

## [v1.0.4] - Smart Notifications 🧠
### Added
- **Smart Notification Deduplication**: Notifications are now only triggered for *new* unique trackers.
  - A persistent history file (`seen_trackers.txt`) is created in the project folder to remember previously blocked domains.
  - Prevents repetitive notifications from the same tracker while keeping the user informed of new threats.
- **Dynamic Path Resolution**: The application now correctly locates its data files relative to the executable path, ensuring reliability regardless of the launch directory.

### Changed
- Improved notification logic: New trackers are saved immediately, while the visual notification adheres to a 30-second cooldown to avoid spam.

## [v1.0.3] - Privacy Shield Edition 🛡️
### Added
- **Global Tracker Blocklist Expansion**: Added hundreds of new domains including:
  - **Mobile Push**: OneSignal, Braze, Urban Airship.
  - **Product Analytics**: Amplitude, Heap, Pendo, Mixpanel.
  - **Data Brokers**: Adobe Demdex, Oracle BlueKai, Acxiom.
  - **Enterprise Monitoring**: New Relic, Datadog.
  - **Social Tracking**: TikTok Ads, Snapchat Pixel, Pinterest Pixel.
- **Native macOS Notifications**: You now receive a notification when a tracker is blocked.
  - Includes a **30-second cooldown** to prevent spam on heavy sites.
  - Sound effect ("Pop") for auditory feedback.

### Changed
- Refined blocklist to remove root domains (e.g., `tiktok.com`) that caused app breakage, targeting only ad/tracking subdomains.

## [v1.0.2] - Interactive Control ⏯️
### Added
- **Interactive Pause/Resume**: Press `p` in the terminal window to toggle the VPN on the fly.
- **Visual Feedback**: Clear "PAUSED" (Yellow) and "RESUMED" (Green) status messages.

### Fixed
- Fixed a bug where the interactive loop would consume CPU or misread input.
- improved network status detection logic.

## [v1.0.1] - Stability & Logging 🪵
### Fixed
- Re-enabled critical log messages for SOCKS requests and connection status.
- Moved blocklist check to client-side for better visibility.

## [v1.0.0] - Initial Release 🚀
- **Core VPN**: SOCKS5 over mTLS tunnel.
- **Security**: Mutual TLS authentication, DNS over HTTPS (DoH).
- **Automation**: `run_vpn.command` for one-click startup.
- **Privacy**: Basic ad-blocker included.
