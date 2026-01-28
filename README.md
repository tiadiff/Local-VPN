# VPN Prototype Walkthrough

This document explains how to use the Go-based encrypted tunnel system.

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
4.  **To Stop**: Press `CTRL+C` in the terminal window. The proxy will be disabled automatically.

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

## Privacy Features
### Tracker Blocker 🛡️
The VPN Server now includes a built-in **Tracker Blocker**.
It automatically rejects connections to known ad and tracking domains (e.g., `doubleclick.net`, `facebook.net`, `googleadservices.com`).
- **Log Message**: `[BLOCKED] Connection to tracker denied: ...`
- **To Customize**: Edit `tunnel/blocklist.go` and rebuild (`go build -o vpn`).

## Security Note
This setup currently runs both Client and Server on your local machine (`localhost`).
**To actually hide your IP from your ISP:**
1.  Copy the `vpn` binary to a remote server (e.g., AWS, DigitalOcean, or a friend's PC).
2.  Run `./vpn -mode server ...` on that remote machine.
3.  Run `./vpn -mode socks -server <REMOTE_IP> ...` on your Mac.

Only then will your traffic truly emerge from a different location.
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
On Client:
```bash
ping 10.0.0.1
```
If you get a reply, the tunnel is up and encrypted!
