// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Preload hack to allow GLFW with software rendering on macOS.
// Allows running the binary inside a VM for testing on CI.

#import "macos_gl_allow_software.h"

#import <Foundation/Foundation.h>
#import <objc/runtime.h>
#import <stdio.h>

typedef id (*initWithAttributesPtr) (id, SEL, const uint32_t *);

static initWithAttributesPtr origInitWithAttributes = NULL;

@implementation GLAllowSoftware

+ (void) load {
	fprintf(stderr, "Injecting initWithAttributes...\n");

	Class origClass = NSClassFromString(@"NSOpenGLPixelFormat");
	Method origMeth = class_getInstanceMethod(origClass, @selector(initWithAttributes:));
	origInitWithAttributes = (initWithAttributesPtr) method_getImplementation(origMeth);

	Class replClass = NSClassFromString(@"GLAllowSoftware");
	Method replMeth = class_getInstanceMethod(replClass, @selector(initWithAttributes:));
	IMP replInitWithAttributes = method_getImplementation(replMeth);

	method_setImplementation(origMeth, replInitWithAttributes);
}

- (id) initWithAttributes: (const uint32_t *) attribs {
	// Just allow ANY pixel format.
	// Yes, the app may misbehave with this, but it'll _run_.
	static uint32_t replAttribs[] = {0};
	id ret = origInitWithAttributes(self, @selector(initWithAttributes:), replAttribs);
	return ret;
}

@end
