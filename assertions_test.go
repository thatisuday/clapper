package clapper

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func fail(t *testing.T) {
	_, f, l, _ := runtime.Caller(2)
	fmt.Printf("%s # %d\n", f, l)
	t.FailNow()
}

func toErrorS(msg ...interface{}) (s string) {
	args := msg[0].([]interface{})
	if len(args) > 0 {
		if fs, ok := args[0].(string); ok {
			if len(args) > 1 {
				s = fmt.Sprintf(": %s", fmt.Sprintf(fs, args[1:]...))
			} else {
				s = fs
			}
		}
	}
	return s
}

func assertError(t *testing.T, v error, msg ...interface{}) {
	if v == nil {
		t.Errorf("expected error%s", toErrorS(msg))
		fail(t)
	}
}

func assertNoError(t *testing.T, v error, msg ...interface{}) {
	if v != nil {
		t.Errorf("expected no error%s", toErrorS(msg))
		fail(t)
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}, msg ...interface{}) {
	if a != b {
		t.Errorf("expected %v, got %v%s", a, b, toErrorS(msg))
		fail(t)
	}
}

func assertNil(t *testing.T, a interface{}, msg ...interface{}) {
	if a != nil {
		t.Errorf("expected nil, got %v%s", a, toErrorS(msg))
		fail(t)
	}
}

func assertNotNil(t *testing.T, a interface{}, msg ...interface{}) {
	if a == nil || (reflect.ValueOf(a).Kind() == reflect.Ptr && reflect.ValueOf(a).IsNil()) {
		t.Errorf("expected not nil, got %v%s", a, toErrorS(msg))
		fail(t)
	}
}
