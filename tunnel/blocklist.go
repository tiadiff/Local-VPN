package tunnel

import "strings"

// BlockList is a simple list of domains to block.
// In a real app, this would be loaded from a file or external source.
var BlockList = []string{
	// Google / DoubleClick
	"doubleclick.net",
	"googleadservices.com",
	"googlesyndication.com",
	"adservice.google.com",
	"analytics.google.com",
	"p.googletagservices.com",
	"www.googletagservices.com",
	"pagead2.googlesyndication.com",
	"tpc.googlesyndication.com",
	"google-analytics.com",
	"tagmanager.google.com",

	// Facebook / Meta
	"facebook.net",
	"connect.facebook.net",
	"pixel.facebook.com",
	"graph.facebook.com",
	"atlas.com",
	"an.facebook.com",

	// Microsoft / LinkedIn
	"ads.linkedin.com",
	"analytics.linkedin.com",
	"browser.events.data.microsoft.com", // Telemetry
	"vortex.data.microsoft.com",         // Telemetry
	"settings-win.data.microsoft.com",   // Telemetry
	"c.bing.com",

	// Twitter / X
	"ads.twitter.com",
	"static.ads-twitter.com",
	"analytics.twitter.com",
	"p.twitter.com",

	// Amazon
	"amazon-adsystem.com",
	"aan.amazon.com",
	"aax.amazon-adsystem.com",
	"device-metrics-us.amazon.com",

	// Generic Ad Networks & Trackers
	"ads.yahoo.com",
	"analytics.yahoo.com",
	"flurry.com",
	"quantserve.com",
	"quantcount.com",
	"scorecardresearch.com",
	"taboola.com",
	"outbrain.com",
	"zedo.com",
	"adcolony.com",
	"applovin.com",
	"chartboost.com",
	"appsflyer.com",
	"adjust.com",
	"kochava.com",
	"branch.io",
	"mixpanel.com",
	"segment.io",
	"bugsnag.com",
	"sentry.io",
	"crashlytics.com",
	"firebase-logging.googleapis.com",
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
