.PHONY: all build clean test test-unit test-integration

GIT_VERSION ?= $(shell git describe --always --dirty)
CGO_CFLAGS=-mmacosx-version-min=13.0
VERSION_LDFLAGS=-X github.com/crc-org/vfkit/pkg/cmdline.gitVersion=$(GIT_VERSION)

all: build

TOOLS_DIR := tools
include tools/tools.mk

build: out/vfkit

test: test-unit test-integration

test-unit:
	@go test -v ./pkg/... ./cmd/vfkit/...

test-integration: build
	@go test -v -timeout 20m ./test

clean:
	rm -rf out

out/vfkit-amd64 out/vfkit-arm64: out/vfkit-%: force-build
	@mkdir -p $(@D)
	CGO_ENABLED=1 CGO_CFLAGS=$(CGO_CFLAGS) GOOS=darwin GOARCH=$* go build -ldflags "$(VERSION_LDFLAGS)" -o $@ ./cmd/vfkit
	codesign -f --entitlements vf.entitlements -s - $@

out/vfkit: out/vfkit-amd64 out/vfkit-arm64
	cd $(@D) && lipo -create $(^F) -output $(@F)

# the go compiler is doing a good job at not rebuilding unchanged files
# this phony target ensures out/vfkit-* are always considered out of date
# and rebuilt. If the code was unchanged, go won't rebuild anything so that's
# fast. Forcing the rebuild ensure we rebuild when needed, ie when the source code
# changed, without adding explicit dependencies to the go files/go.mod
.PHONY: force-build
force-build:

.PHONY: lint
lint: $(TOOLS_BINDIR)/golangci-lint
	"$(TOOLS_BINDIR)"/golangci-lint run

# A bundle gives Accessory Access a stable app identity for the approval menu.
USB_SIGN_IDENTITY ?=
USB_PROVISION_PROFILE ?=
.PHONY: usb-app
usb-app:
	@test "$$(uname -s)" = Darwin || (echo "usb-app requires macOS" >&2; exit 1)
	@test "$$(xcrun --sdk macosx --show-sdk-version | cut -d. -f1)" -ge 27 || (echo "usb-app requires the macOS 27 SDK" >&2; exit 1)
	@test -n "$(USB_SIGN_IDENTITY)" -a "$(USB_SIGN_IDENTITY)" != "-" || (echo "usb-app requires USB_SIGN_IDENTITY (ad-hoc signing cannot use Accessory Access)" >&2; exit 1)
	@test -f "$(USB_PROVISION_PROFILE)" || (echo "usb-app requires USB_PROVISION_PROFILE with the Claim USB Accessory capability" >&2; exit 1)
	@mkdir -p out/vfkit.app/Contents/MacOS
	cp contrib/usb/Info.plist out/vfkit.app/Contents/Info.plist
	CGO_ENABLED=1 CGO_CFLAGS=$(CGO_CFLAGS) go build -ldflags "$(VERSION_LDFLAGS)" -o out/vfkit.app/Contents/MacOS/vfkit ./cmd/vfkit
	sh contrib/usb/sign.sh out/vfkit.app "$(USB_SIGN_IDENTITY)" "$(USB_PROVISION_PROFILE)"
