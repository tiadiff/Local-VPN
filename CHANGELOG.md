# Changelog

All notable changes to this project will be documented in this file.

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
