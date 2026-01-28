#!/bin/bash

# Configuration (must match run_vpn.command)
SERVICE="Wi-Fi"

# Get current SOCKS status (Enabled: Yes / No)
STATUS=$(networksetup -getsocksfirewallproxy "$SERVICE" | grep "Enabled:" | awk '{print $2}')

if [ "$STATUS" == "Yes" ]; then
    networksetup -setsocksfirewallproxystate "$SERVICE" off
    osascript -e 'display notification "VPN Disattivata. Traffico NON protetto." with title "Secure Tunnel" subtitle "Paused ⏸️" sound name "Basso"'
    echo "VPN Paused"
else
    networksetup -setsocksfirewallproxystate "$SERVICE" on
    osascript -e 'display notification "VPN Attiva. Traffico Protetto." with title "Secure Tunnel" subtitle "Resumed ▶️" sound name "Glass"'
    echo "VPN Resumed"
fi
