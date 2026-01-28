#!/bin/bash

# Configuration
SERVICE="Wi-Fi"
PROXY_HOST="127.0.0.1"
PROXY_PORT="1080"
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Ensure sudo
if [ "$EUID" -ne 0 ]; then
  printf "${YELLOW}[sudo] Richiesto privilegio di amministratore per modificare le impostazioni di rete.${NC}\n"
  sudo "$0"
  exit
fi

# ** Check for Certificates **
if [ ! -f "$DIR/ca.key" ]; then
    printf "${YELLOW}[INIT] Certificati mTLS mancanti. Generazione in corso...${NC}\n"
    "$DIR/vpn" -gen-certs
    chmod 600 "$DIR"/*.key
fi

cleanup() {
    echo ""
    printf "${YELLOW}[STOP] Arresto in corso...${NC}\n"
    
    # Kill background jobs
    kill $(jobs -p) 2>/dev/null
    
    # Disable Proxy
    # echo "🔌 Disattivazione Proxy SOCKS su $SERVICE..."
    networksetup -setsocksfirewallproxystate "$SERVICE" off
    
    printf "${GREEN}[OK]   VPN disattivata.${NC}\n"
    exit
}

# Trap Ctrl+C (SIGINT), Exit (SIGTERM), and Window Close (SIGHUP)
trap cleanup SIGINT SIGTERM SIGHUP EXIT

clear
printf "${CYAN}VPN AUTOMATION WRAPPER${NC}\n"
printf "${CYAN}Directory: $DIR${NC}\n"
echo ""

# 1. Enable Proxy
printf "${BLUE}[NET]  Attivazione Proxy SOCKS su $SERVICE ($PROXY_HOST:$PROXY_PORT)...${NC}\n"
networksetup -setsocksfirewallproxy "$SERVICE" "$PROXY_HOST" "$PROXY_PORT"
networksetup -setsocksfirewallproxystate "$SERVICE" on

# 2. Start Server (Background)
# Loads CA, server.crt, server.key. Enforces client Auth.
"$DIR/vpn" -mode server -port 3000 -cert "$DIR/server.crt" -key "$DIR/server.key" -ca "$DIR/ca.crt" > /dev/null 2>&1 &
SERVER_PID=$!

sleep 1

# 3. Start Client (Background) -> Foreground Log
# Presents client.crt, client.key to server.
# We let vpn process output its own cool logs
"$DIR/vpn" -mode socks -server 127.0.0.1 -port 3000 -socks "$PROXY_PORT" \
    -cert "$DIR/client.crt" -key "$DIR/client.key" -ca "$DIR/ca.crt" &
CLIENT_PID=$!

printf "${GREEN}[RUN]  Sistema pronto. \nPremi 'p' per Pausa/Ripresa, CTRL+C per terminare.${NC}\n"
echo "------------------------------------------------------------"

while kill -0 $CLIENT_PID 2>/dev/null; do
    key=""
    # Silent read of 1 char with 1s timeout
    read -t 1 -n 1 -s key
    
    # DEBUG: Show what is happening
    # printf "DEBUG: Key read: '%s'\n" "$key"
    
    # Check if key is actually empty (timeout) or not
    if [[ -n "$key" ]]; then
        # printf "DEBUG: Key detected: '%s'\n" "$key"
        
        if [[ "$key" == "p" || "$key" == "P" ]]; then
             # DEBUG: Print exact status
             RAW_STATUS=$(networksetup -getsocksfirewallproxy "$SERVICE" | grep "^Enabled:")
             STATUS=$(echo "$RAW_STATUS" | awk '{print $2}')
             
             # Trim potential whitespace
             STATUS=$(echo "$STATUS" | xargs)

             # printf "DEBUG: RawStatus=['%s'] ParsedStatus=['%s']\n" "$RAW_STATUS" "$STATUS"

             if [ "$STATUS" == "Yes" ]; then
                 networksetup -setsocksfirewallproxystate "$SERVICE" off
                 printf "\n${YELLOW}[PAUSED] VPN disattivata temporaneamente (Traffico Diretto).${NC}\n"
                 osascript -e 'display notification "VPN Disattivata" with title "Secure Tunnel" subtitle "Paused ⏸️"'
             else
                 networksetup -setsocksfirewallproxystate "$SERVICE" on
                 printf "\n${GREEN}[RESUMED] VPN riattivata (Traffico Protetto).${NC}\n"
                 osascript -e 'display notification "VPN Attiva" with title "Secure Tunnel" subtitle "Resumed ▶️"'
             fi
        fi
    fi
done
