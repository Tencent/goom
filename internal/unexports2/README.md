Subvert
=======

Package subvert provides functions to subvert go's runtime system, allowing you to:

* Get addresses of stack-allocated or otherwise protected values
* Access unexported values
* Call unexported functions
* Apply patches to memory (even if it's read-only)
* Make aliases to functions

![Now I know what it feels like to be God!](power.gif)

This is not a power to be taken lightly! It's expected that you're fully
versed in how the go type system works, and why there are protections and
restrictions in the first place. Using this package incorrectly will quickly
lead to undefined behavior and bizarre crashes, even segfaults or nuclear
missile launches.

**YOU HAVE BEEN WARNED!**



Example
-------

```golang
import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/kstenerud/go-subvert"
)

type SubvertTester struct {
	A int
	a int
	int
}

const constString = "testing"

func Demonstrate() {
	v := SubvertTester{1, 2, 3}

	rv := reflect.ValueOf(v)
	rv_A := rv.FieldByName("A")
	rv_a := rv.FieldByName("a")
	rv_int := rv.FieldByName("int")

	fmt.Printf("Interface of A: %v\n", rv_A.Interface())

	// MakeWritable

	// rv_a.Interface() // This would panic
	if err := subvert.MakeWritable(&rv_a); err != nil {
		// TODO: Handle this
	}
	fmt.Printf("Interface of a: %v\n", rv_a.Interface())

	// rv_int.Interface() // This would panic
	if err := subvert.MakeWritable(&rv_int); err != nil {
		// TODO: Handle this
	}
	fmt.Printf("Interface of int: %v\n", rv_int.Interface())

	// MakeAddressable

	// rv.Addr() // This would panic
	if err := subvert.MakeAddressable(&rv); err != nil {
		// TODO: Handle this
	}
	fmt.Printf("Pointer to v: %v\n", rv.Addr())

	// ExposeFunction

	exposed, err := subvert.ExposeFunction("reflect.methodName", (func() string)(nil))
	if err != nil {
		// TODO: Handle this
	}
	f := exposed.(func() string)
	fmt.Printf("Result of reflect.methodName: %v\n", f())

	// PatchMemory

	rv = reflect.ValueOf(constString)
	if err := subvert.MakeAddressable(&rv); err != nil {
		// TODO: Handle this
	}
	strAddr := rv.Addr().Pointer()
	strBytes := *((*unsafe.Pointer)(unsafe.Pointer(strAddr)))
	if oldMem, err := subvert.PatchMemory(uintptr(strBytes), []byte("XX")); err != nil {
		// TODO: Handle this
	}
	fmt.Printf("constString is now: %v, Oldmem = %v\n", constString, string(oldMem))
}
```

**Output:**

```
Interface of A: 1
Interface of a: 2
Interface of int: 3
Pointer to v: &{1 2 3}
Result of reflect.methodName: github.com/kstenerud/go-subvert.TestDemonstrate
constString is now: XXsting, Oldmem = te
```



License
-------

MIT License:

Copyright 2020 Karl Stenerud

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
