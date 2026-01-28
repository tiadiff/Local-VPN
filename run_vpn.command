#!/bin/bash

# Configuration
SERVICE="Wi-Fi"
PROXY_HOST="127.0.0.1"
PROXY_PORT="1080"
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Ensure sudo
if [ "$EUID" -ne 0 ]; then
  echo "Richiesto privilegio di amministratore per modificare le impostazioni di rete."
  sudo "$0"
  exit
fi

# ** Check for Certificates **
if [ ! -f "$DIR/ca.key" ]; then
    echo "⚠️ Certificati mTLS mancanti. Generazione in corso..."
    "$DIR/vpn" -gen-certs
    echo "✅ Certificati generati (ca.crt, server.crt, client.crt)."
    chmod 600 "$DIR"/*.key
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

echo "** VPN Automation Wrapper (mTLS Secured) **"
echo "📂 Directory: $DIR"

# 1. Enable Proxy
echo "🔌 Attivazione Proxy SOCKS su $SERVICE ($PROXY_HOST:$PROXY_PORT)..."
networksetup -setsocksfirewallproxy "$SERVICE" "$PROXY_HOST" "$PROXY_PORT"
networksetup -setsocksfirewallproxystate "$SERVICE" on

# 2. Start Server (Background)
# Loads CA, server.crt, server.key. Enforces client Auth.
echo "🚀 Avvio Server VPN (mTLS)..."
"$DIR/vpn" -mode server -port 3000 -cert "$DIR/server.crt" -key "$DIR/server.key" -ca "$DIR/ca.crt" > /dev/null 2>&1 &
SERVER_PID=$!

sleep 1

# 3. Start Client (Background) -> Foreground Log
# Presents client.crt, client.key to server.
echo "🚀 Avvio Client VPN (mTLS)..."
"$DIR/vpn" -mode socks -server 127.0.0.1 -port 3000 -socks "$PROXY_PORT" \
    -cert "$DIR/client.crt" -key "$DIR/client.key" -ca "$DIR/ca.crt" &
CLIENT_PID=$!

echo "✅ VPN Attiva e Protetta con Mutual TLS!"
echo "🌐 Naviga con Safari. Premi CTRL+C per spegnere."
echo "------------------------------------------------"

# Wait for client
wait $CLIENT_PID
