package tunnel

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"vpn_proto/config"
	"vpn_proto/crypto"
	"vpn_proto/utils"
)

func StartSocksClient(cfg *config.Config) error {
	serverEndpoint := fmt.Sprintf("%s:%d", cfg.ServerAddr, cfg.Port)
	log.Printf("Connecting to server at %s...", serverEndpoint)

	// 1. Connect to Tunnel Server
    // In SOCKS mode, we don't open a persistent connection for control.
    // Instead we listen and open a NEW connection for every SOCKS request.
    // So distinct from 'StartClient' which does a persistent tunnel.
    
    // We just start the listener here.
    
    log.Printf("Starting SOCKS server on :%d", cfg.SocksPort)

	// 3. Start Local SOCKS Server
    listenAddr := fmt.Sprintf("127.0.0.1:%d", cfg.SocksPort)
    ln, err := net.Listen("tcp", listenAddr)
    if err != nil {
        return err
    }
    
    for {
        localConn, err := ln.Accept()
        if err != nil {
            log.Printf("Socks accept error: %v", err)
            continue
        }
        
        log.Printf("Accepted SOCKS connection from %s", localConn.RemoteAddr())
        go handleSocksConnection(localConn, cfg)
    }
}

func handleSocksConnection(localConn net.Conn, cfg *config.Config) {
    defer localConn.Close()
    
    // 1. Dial Server for THIS specific connection
	serverEndpoint := fmt.Sprintf("%s:%d", cfg.ServerAddr, cfg.Port)
    tlsConfig := crypto.ClientTLSConfig()
    remoteConn, err := tls.Dial("tcp", serverEndpoint, tlsConfig)
    if err != nil {
        log.Printf("Failed to dial upstream: %v", err)
        return
    }
    defer remoteConn.Close()
    
    // 2. Auth with "SOCKS" prefix so server knows it's a proxy request?
    // Actually, if we just use the same protocol, the Server will try to write it to TUN.
    // We need to tell the server "I want to CONNECT to target:port".
    
    // Let's modify the authentication packet to include mode.
    // "SECRET|MODE|TARGET"
    
    // Wait, standard SOCKS5 handshake is:
    // Client -> SOCKS Server (Us) -> Handshake -> Request (Connect google.com:80)
    // We need to READ the request from localConn, extract target, send it to Remote Server.
    
    // Using go-socks5 is complex if we want to channel it.
    // Let's implement a very stupid SOCKS5 server helper.
    
    // SOCKS5 Init
    buf := make([]byte, 256)
    // Read version
    n, err := localConn.Read(buf)
    if err != nil || n < 2 || buf[0] != 0x05 {
        return // Not SOCKS5
    }
    // No auth required, rely
    localConn.Write([]byte{0x05, 0x00}) 
    
    // Read Request
    n, err = localConn.Read(buf)
    if err != nil || n < 4 {
        return
    }
    
    cmd := buf[1]
    if cmd != 0x01 { // CONNECT
        return // We only support CONNECT
    }
    
    // Parse Address
    var target string
    addrType := buf[3]
    switch addrType {
    case 0x01: // IPv4
        if n < 10 { return }
        ip := net.IP(buf[4:8])
        port := int(buf[8])<<8 | int(buf[9])
        target = fmt.Sprintf("%s:%d", ip, port)
    case 0x03: // Domain
        domLen := int(buf[4])
        if n < 5+domLen+2 { return }
        domain := string(buf[5 : 5+domLen])
        port := int(buf[5+domLen])<<8 | int(buf[5+domLen+1])
        target = fmt.Sprintf("%s:%d", domain, port)
    case 0x04: // IPv6
        return // Skip for prototype
    }
    
    // Now we have the target.
    // Connect to VPN Server
    // Send Auth Header: "SECRET|CONNECT|target"
    authMsg := fmt.Sprintf("%s|CONNECT|%s", cfg.Secret, target)
    err = utils.WritePacket(remoteConn, []byte(authMsg))
    if err != nil {
        log.Printf("Upstream handshake failed: %v", err)
        return
    }
    
    // Read response from VPN server (OK/Fail)
    resp, err := utils.ReadPacket(remoteConn)
    if err != nil || string(resp) != "OK" {
        localConn.Write([]byte{0x05, 0x01, 0x00, 0x01, 0,0,0,0, 0,0}) // Fail
        return
    }
    
    // Success SOCKS reply
    localConn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0,0,0,0, 0,0}) // OK
    
    // Bi-directional Copy
    go io.Copy(localConn, remoteConn)
    io.Copy(remoteConn, localConn)
}
