#!/bin/bash
# Advanced Build script for SecureTunnel.app

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"

APP_NAME="SecureTunnel"
BUNDLE_ID="com.vpnproto.securetunnel"
APP_DIR="$DIR/$APP_NAME.app"
CONTENTS="$APP_DIR/Contents"
MACOS="$CONTENTS/MacOS"
RESOURCES="$CONTENTS/Resources"

echo "🧹 Cleaning old builds..."
rm -rf "$APP_DIR"
mkdir -p "$MACOS" "$RESOURCES"

echo "🔨 Building Go VPN binary..."
go build -o "$RESOURCES/vpn" main.go

echo "🔨 Building Swift Menubar utility..."
swiftc -o "$MACOS/vpn_menubar" "$DIR/menubar_ext.swift" -framework Cocoa

echo "🔨 Building Swift Launcher..."
swiftc -o "$MACOS/$APP_NAME" "$DIR/launcher.swift" -framework Cocoa -framework AppKit

echo "🎨 Generating App Icon (Emoji -> .icns)..."
ICONSET="$DIR/AppIcon.iconset"
mkdir -p "$ICONSET"

# Function to create an image from emoji using a temporary swift file
generate_emoji_png() {
    cat > icon_gen.swift <<EOF
import Cocoa
let emoji = "🔒"
let size = NSSize(width: 512, height: 512)
let img = NSImage(size: size)
img.lockFocus()
let rect = NSRect(origin: .zero, size: size)
(emoji as NSString).draw(in: rect, withAttributes: [.font: NSFont.systemFont(ofSize: 400)])
img.unlockFocus()
if let data = img.tiffRepresentation, let rep = NSBitmapImageRep(data: data) {
    if let pngData = rep.representation(using: .png, properties: [:]) {
        try? pngData.write(to: URL(fileURLWithPath: "emoji_temp.png"))
    }
}
EOF
    swift icon_gen.swift
    rm icon_gen.swift
}

generate_emoji_png
sips -z 16 16   emoji_temp.png --out "$ICONSET/icon_16x16.png" > /dev/null
sips -z 32 32   emoji_temp.png --out "$ICONSET/icon_16x16@2x.png" > /dev/null
sips -z 32 32   emoji_temp.png --out "$ICONSET/icon_32x32.png" > /dev/null
sips -z 64 64   emoji_temp.png --out "$ICONSET/icon_32x32@2x.png" > /dev/null
sips -z 128 128 emoji_temp.png --out "$ICONSET/icon_128x128.png" > /dev/null
sips -z 256 256 emoji_temp.png --out "$ICONSET/icon_128x128@2x.png" > /dev/null
sips -z 256 256 emoji_temp.png --out "$ICONSET/icon_256x256.png" > /dev/null
sips -z 512 512 emoji_temp.png --out "$ICONSET/icon_256x256@2x.png" > /dev/null
sips -z 512 512 emoji_temp.png --out "$ICONSET/icon_512x512.png" > /dev/null
sips -z 1024 1024 emoji_temp.png --out "$ICONSET/icon_512x512@2x.png" > /dev/null

iconutil -c icns "$ICONSET" -o "$RESOURCES/AppIcon.icns"
rm -rf "$ICONSET" emoji_temp.png

echo "📝 Creating Info.plist..."
cat > "$CONTENTS/Info.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>$APP_NAME</string>
    <key>CFBundleIconFile</key>
    <string>AppIcon</string>
    <key>CFBundleIdentifier</key>
    <string>$BUNDLE_ID</string>
    <key>CFBundleName</key>
    <string>$APP_NAME</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.13</string>
    <key>NSPrincipalClass</key>
    <string>NSApplication</string>
    <key>LSUIElement</key>
    <string>YES</string>
</dict>
</plist>
EOF

echo "🔑 Copying certificates and supplemental files..."
cp "$DIR"/*.crt "$DIR"/*.key "$RESOURCES/" 2>/dev/null

echo "✅ $APP_NAME.app build complete."
echo "👉 You can now move $APP_NAME.app to your Applications folder."
