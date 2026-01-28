#!/bin/bash

# Configuration
SERVICE="Wi-Fi"
PROXY_HOST="127.0.0.1"
PROXY_PORT="1080"
SECRET="mysecret"
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Ensure sudo
if [ "$EUID" -ne 0 ]; then
  echo "Richiesto privilegio di amministratore per modificare le impostazioni di rete."
  sudo "$0"
  exit
fi

cleanup() {
    echo ""
    echo "🛑 Arresto in corso..."
    
    # Kill background jobs
    kill $(jobs -p) 2>/dev/null
    
    # Disable Proxy
    echo "🔌 Disattivazione Proxy SOCKS su $SERVICE..."
    networksetup -setsocksfirewallproxystate "$SERVICE" off
    
    echo "✅ VPN spenta e impostazioni ripristinate."
    exit
}

# Trap Ctrl+C (SIGINT), Exit (SIGTERM), and Window Close (SIGHUP)
trap cleanup SIGINT SIGTERM SIGHUP EXIT

echo "** VPN Automation Wrapper **"
echo "📂 Directory: $DIR"

# 1. Enable Proxy
echo "🔌 Attivazione Proxy SOCKS su $SERVICE ($PROXY_HOST:$PROXY_PORT)..."
networksetup -setsocksfirewallproxy "$SERVICE" "$PROXY_HOST" "$PROXY_PORT"
networksetup -setsocksfirewallproxystate "$SERVICE" on

# 2. Start Server (Background)
echo "🚀 Avvio Server VPN..."
"$DIR/vpn" -mode server -port 3000 -secret "$SECRET" > /dev/null 2>&1 &
SERVER_PID=$!

sleep 1

# 3. Start Client (Background) -> Foreground Log
echo "🚀 Avvio Client VPN..."
"$DIR/vpn" -mode socks -server 127.0.0.1 -port 3000 -socks "$PROXY_PORT" -secret "$SECRET" &
CLIENT_PID=$!

echo "✅ VPN Attiva e Protetta!"
echo "🌐 Naviga con Safari. Premi CTRL+C per spegnere."
echo "------------------------------------------------"

# Wait for client
wait $CLIENT_PID
