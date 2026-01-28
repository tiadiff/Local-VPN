package tunnel

import "strings"

// BlockList is a simple list of domains to block.
// In a real app, this would be loaded from a file or external source.
var BlockList = []string{
	"doubleclick.net",
	"googleadservices.com",
	"googlesyndication.com",
	"adservice.google.com",
	"facebook.net",
	"connect.facebook.net",
	"analytics.google.com",
	"c.bing.com",
	"pagead2.googlesyndication.com",
	"tpc.googlesyndication.com",
	"www.googletagservices.com",
    "pixel.facebook.com",
    "ads.twitter.com",
    "static.ads-twitter.com",
}

// IsBlocked checks if the target host contains any of the blocked domains.
func IsBlocked(target string) bool {
    // target usually comes as "hostname:port"
    host := target
    if idx := strings.Index(host, ":"); idx != -1 {
        host = host[:idx]
    }
    
    for _, blocked := range BlockList {
        // Simple suffix match covers "subdomain.tracker.com" and "tracker.com"
        if strings.HasSuffix(host, blocked) {
            return true
        }
    }
    return false
}
