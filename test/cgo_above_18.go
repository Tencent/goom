//go:build go1.20
// +build go1.20

package test

/*
#include <stdio.h>

void printint(int v) {
	printf("printint: %d\n", v);
}
*/
import "C"

func cgoFuncAny() {
	v := 42
	C.printint(C.int(v))
}
