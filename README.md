# VPN Prototype Walkthrough

This document explains how to use the Go-based encrypted tunnel system.

<img width="568" height="368" alt="Schermata 2026-01-28 alle 22 40 56" src="https://github.com/user-attachments/assets/00680175-4b8c-4c7d-b22c-ae0734dcef65" />

## Prerequisites
- **Go 1.18+** (Installed via `brew install go`)
- **Sudo privileges** (Required for creating network interfaces)

## Building
If not already built:
```bash
go build -o vpn
```

## Quick Start (AUTOMATIC Script)
To start the VPN and automatically configure your Mac's proxy settings:
1.  Navigate to the project folder.
2.  Double-click `run_vpn.command`.
3.  Enter your password when prompted (required to change network settings).
4.  Press P from yourkeyboard to pause/resume the Tunneling.
5.  **To stop**: Press CTRL+C in the terminal window or close all terminal windows with CMD+Q.. in both cases, the proxy will be automatically disabled.

Or..

## SOCKS5 Proxy Mode (MANUAL)
This mode creates a local proxy (port 1080) that tunnels traffic to the server safely without sudo.

1.  **Start Server** (no sudo needed for this mode, but sudo enables fallback to TUN):
    ```bash
    ./vpn -mode server -port 3000 -secret mysecret
    ```

2.  **Start Client** (Proxy Mode):
    ```bash
    ./vpn -mode socks -server 127.0.0.1 -port 3000 -socks 1080 -secret mysecret
    ```

3.  **Configure Safari / System**:
    - Go to **System Preferences** > **Network** > **Advanced** > **Proxies**.
    - Check **SOCKS Proxy**.
    - Server: `127.0.0.1`, Port: `1080`.
    - Click OK > Apply.

4.  **Verify**:
## Auto-Start at Login 🚀
To run the VPN automatically when you log in:
1.  Open **System Settings** -> **General** -> **Login Items**.
2.  Click the `+` button in the "Open at Login" section.
3.  Navigate to your project folder and select `run_vpn.command`.
4.  **Note**: Upon login, a terminal window will open and ask for your password to activate the proxy.

## Privacy & Security Features 🛡️

### 1. Mutual TLS (mTLS) 🔐
Authentication is strictly enforced using mTLS.
- **Client Side**: Must present a valid `client.crt` signed by your private CA.
- **Server Side**: Verifies the client certificate before accepting any connection.
- **Benefit**: Immune to password brute-forcing and unauthorized scanning.

### 2. DNS over HTTPS (DoH) 🕵️
The Server automatically resolves domain names using **Cloudflare (1.1.1.1)** via encrypted HTTPS.
- **Benefit**: Your ISP cannot see your DNS queries (hides which sites you visit).
- **Log Message**: `[SEC] Resolving example.com via DoH...`

### 3. Tracker Blocker 🚫
The VPN Client includes a built-in **Tracker Blocker** that actively filters specific tracking domains.
- **Mechanism**: Checks requests against a local blocklist before they leave your machine.
- **Benefit**: Blocks ads and trackers at the source, saving bandwidth and protecting privacy.
- **Log Message**: `[BLOCK] Connection denied (Client-side): doubleclick.net`
The server listens for incoming connections and acts as the tunnel endpoint.
```bash
# Listen on port 3000 (default) with a shared secret
./vpn -mode server -port 3000 -secret mysecret
```
*Note: The server uses a self-signed certificate generated on startup.*

## Running the Client
The client connects to the server and creates a TUN interface (e.g., `utun3`).
```bash
# Connect to server
sudo ./vpn -mode client -server 127.0.0.1 -port 3000 -secret mysecret
```
*Note: Sudo is required to create the TUN interface.*

## Networking Configuration (Manual Steps)
Since this is a prototype, routing rules are not securely applied automatically to prevent locking you out of the system.

Once the client is running, you will see a message like:
`Interface utun3 created`

You need to assign an IP and set up routes in a **separate terminal**:

1.  **Assign IP to Tunnel**:
    ```bash
    # Replace utun3 with your actual interface name
    sudo ifconfig utun3 10.0.0.2 10.0.0.1 up
    ```

2.  **Route Traffic**:
    To route specific traffic (e.g., to a specific IP) through the tunnel:
    ```bash
    sudo route add <DESTINATION_IP> 10.0.0.1
    ```
    
    To route **ALL Internet Traffic** (VPN):
    > [!WARNING]
    > Doing this incorrectly can disconnect your internet.
    
    1.  Add specific route to VPN Server IP via your Gateway (to avoid loops).
    2.  Delete default route.
    3.  Add default route via `10.0.0.1`.

## Verification

### Basic Test
On Client:
```bash
curl -v --socks5 127.0.0.1:1080 https://www.google.com
```
If you get a response and see `[INFO] SOCKS Request...` in your valid VPN terminal, it works!

### Pro Tip: Verify Encryption Yourself 🕵️‍♂️
Want to see the matrix? You can prove the traffic is encrypted by sniffing your own loopback interface:

```bash
sudo tcpdump -i lo0 -X port 3000
```

**What to look for:**
- **Gibberish**: The data columns (right side) should look like random characters (`...T...F.E>.?...`). This is good! It means the payload is encrypted.
- **TLS Handshake**: Look for `16 03 01` at the start of packets. This is the "Client Hello" signature of a TLS secured connection.
- **No Plaintext**: You should NOT see any website names (like `google.com`) or HTML content in the data dump.
