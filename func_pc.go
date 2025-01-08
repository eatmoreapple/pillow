/*
Copyright 2025 eatmoreapple

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package juice

import (
	"runtime"
	"strings"
	"sync"
)

// replacer defines the replacer of function name
var replacer = strings.NewReplacer("/", ".", `\`, ".", "*", "", "(", "", ")", "")

// runtimeFuncName returns the function name of runtime
func runtimeFuncName(addr uintptr) string {
	// one id from function name
	name := runtime.FuncForPC(addr).Name()
	name = replacer.Replace(name)
	return strings.TrimSuffix(name, "-fm")
}

// _cachedRuntimeFuncName initializes a cached version of runtimeFuncName
// The function is unexported to avoid polluting the package namespace,
// while its returned function is assigned to an exported variable.
// This pattern keeps the initialization logic private while exposing
// only the necessary functionality.
func _cachedRuntimeFuncName() func(addr uintptr) string {
	var cache sync.Map
	return func(addr uintptr) string {
		if name, ok := cache.Load(addr); ok {
			return name.(string)
		}
		name := runtimeFuncName(addr)
		cache.Store(addr, name)
		return name
	}
}

// cachedRuntimeFuncName is a cached version of runtimeFuncName
// It stores function names in memory to avoid repeated processing
var cachedRuntimeFuncName = _cachedRuntimeFuncName()