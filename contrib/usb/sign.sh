#!/bin/sh
set -eu

app=$1
identity=$2
profile=$3
plistbuddy=/usr/libexec/PlistBuddy
profile_plist=$(mktemp)
entitlements=$(mktemp)
trap 'rm -f "$profile_plist" "$entitlements"' EXIT HUP INT TERM

security cms -D -i "$profile" > "$profile_plist"
if [ "$("$plistbuddy" -c 'Print :Entitlements:com.apple.developer.accessory-access.usb' "$profile_plist")" != true ]; then
    echo "The provisioning profile does not authorize Claim USB Accessory" >&2
    exit 1
fi
app_id=$("$plistbuddy" -c 'Print :Entitlements:com.apple.application-identifier' "$profile_plist")
team_id=$("$plistbuddy" -c 'Print :Entitlements:com.apple.developer.team-identifier' "$profile_plist")
case "$app_id" in
    *\**|*\?*|*\'*) echo "Use an explicit app ID provisioning profile" >&2; exit 1 ;;
esac
bundle_id=${app_id#*.}

cp contrib/usb/vf-usb.entitlements "$entitlements"
"$plistbuddy" -c "Add :com.apple.application-identifier string $app_id" "$entitlements"
"$plistbuddy" -c "Add :com.apple.developer.team-identifier string $team_id" "$entitlements"
"$plistbuddy" -c "Set :CFBundleIdentifier $bundle_id" "$app/Contents/Info.plist"
cp "$profile" "$app/Contents/embedded.provisionprofile"
codesign -f --entitlements "$entitlements" -s "$identity" "$app"
codesign --verify --strict "$app"
