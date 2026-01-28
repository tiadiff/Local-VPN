package tunnel

import "strings"

// BlockList is a simple list of domains to block.
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

	// Retargeting & Programmatic
	"criteo.com",
	"criteo.net",
	"teads.tv",
	"adroll.com",
	"rubiconproject.com",
	"pubmatic.com",
	"openx.net",
	"adnxs.com", // AppNexus/Xandr
	"smartadserver.com",
	"moatads.com",

	// Heatmaps & Session Recording
	"hotjar.com",
	"hotjar.io",
	"crazyegg.com",
	"luckyorange.com",
	"fullstory.com",
	"logrocket.io", // Often used for recording

	// Social & Viral (Trackers Only)
	"ads.tiktok.com",
	"analytics.tiktok.com",
	"business-api.tiktok.com",
	"ads.snapchat.com",
	"tr.snapchat.com",  // Snapchat Pixel
	"sc-static.net",    // Often ads
	"ct.pinterest.com", // Pinterest Pixel
	"ads.pinterest.com",
	"analytics.pinterest.com",

	// Email Tracking / ESPs
	"sendgrid.com",  // Often used for tracking pixels
	"mailchimp.com", // Tracking pixels
	"list-manage.com",
	"mandrillapp.com",

	// Mobile & Push Notification Tracking
	"onesignal.com",
	"urbanairship.com",
	"braze.com",
	"appboy.com", // Old Braze domain
	"leanplum.com",
	"airship.com",
	"kochava.com", // Attribution

	// Product Analytics & UX
	"amplitude.com",
	"heapanalytics.com",
	"pendo.io",
	"split.io",
	"optimizely.com",
	"inspectlet.com",
	"mouseflow.com",

	// Customer Feedback & Surveys
	"qualtrics.com",
	"medallia.com",
	"usabilla.com",
	"delighted.com",

	// Performance Monitoring (RUM)
	"nr-data.net", // New Relic
	"bam.nr-data.net",
	"browser-intake-datadoghq.com",
	"sentry-cdn.com",

	// Famous Ad Tech & Data Brokers
	"adform.com",
	"adform.net",
	"demdex.net",      // Adobe Audience Manager
	"everesttech.net", // Adobe
	"omtrdc.net",      // Adobe SiteCatalyst
	"2o7.net",         // Adobe
	"bluekai.com",     // Oracle Data Cloud
	"bkrtx.com",       // BlueKai
	"exelator.com",    // Nielsen
	"krxd.net",        // Salesforce DMP (Krux)
	"liadm.com",       // LiveIntent
	"rlcdn.com",       // Rapleaf
	"pippio.com",      // LiveRamp
	"acxiom.com",
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
