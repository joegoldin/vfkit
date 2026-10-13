// xcrun clang -framework AppKit cmd/vfkit/testdata/usb_application_policy.m -o /tmp/vfkit-usb-policy-test
#import "../usb_application_darwin.m"
#include <assert.h>

void usbApplicationReady(uintptr_t context) {}

int main(void)
{
    @autoreleasepool {
        assert(prepareUSBApplication());
        // A bundled application may already have the regular activation policy.
        assert(prepareUSBApplication());
    }
    return 0;
}
