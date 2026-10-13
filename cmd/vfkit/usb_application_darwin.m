#import <AppKit/AppKit.h>
#include <signal.h>
#include <unistd.h>
#include "usb_application.h"

extern void usbApplicationReady(uintptr_t context);

@interface VFUSBApplicationDelegate : NSObject <NSApplicationDelegate>
@property(nonatomic) uintptr_t context;
@end

@implementation VFUSBApplicationDelegate
- (void)applicationDidFinishLaunching:(NSNotification *)notification
{
    usbApplicationReady(self.context);
}

- (NSApplicationTerminateReply)applicationShouldTerminate:(NSApplication *)sender
{
    // Use the same guest shutdown and timeout as Ctrl-C or SIGTERM.
    kill(getpid(), SIGTERM);
    return NSTerminateCancel;
}
@end

bool prepareUSBApplication(void)
{
    return [[NSApplication sharedApplication] setActivationPolicy:NSApplicationActivationPolicyRegular];
}

void runUSBApplication(uintptr_t context)
{
    @autoreleasepool {
        VFUSBApplicationDelegate *delegate = [VFUSBApplicationDelegate new];
        delegate.context = context;
        [NSApp setDelegate:delegate];
        NSMenu *menu = [[NSMenu alloc] initWithTitle:@"vfkit"];
        NSMenuItem *item = [[NSMenuItem alloc] initWithTitle:@"vfkit" action:nil keyEquivalent:@""];
        NSMenu *applicationMenu = [[NSMenu alloc] initWithTitle:@"vfkit"];
        [applicationMenu addItemWithTitle:@"Quit vfkit" action:@selector(terminate:) keyEquivalent:@"q"];
        [item setSubmenu:applicationMenu];
        [menu addItem:item];
        [NSApp setMainMenu:menu];
        [applicationMenu release];
        [item release];
        [menu release];
        [NSApp run];
        [NSApp setDelegate:nil];
        [delegate release];
    }
}

void stopUSBApplication(void)
{
    dispatch_async(dispatch_get_main_queue(), ^{
        [NSApp stop:nil];
        // Wake nextEventMatchingMask so the run loop observes stop immediately.
        [NSApp postEvent:[NSEvent otherEventWithType:NSEventTypeApplicationDefined
                                          location:NSZeroPoint modifierFlags:0 timestamp:0
                                      windowNumber:0 context:nil subtype:0 data1:0 data2:0]
                  atStart:YES];
    });
}
