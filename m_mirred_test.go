package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMirred(t *testing.T) {
	tests := map[string]struct {
		val  Mirred
		err1 error
		err2 error
	}{
		"simple":          {val: Mirred{Parms: &MirredParam{Index: 42, Action: 1}}},
		"invalidArgument": {val: Mirred{Tm: &Tcft{Install: 1}}, err1: ErrNoArgAlter},
	}
	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalMirred(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			newData, tm := injectTcft(t, data, tcaMirredTm)
			newData = injectAttribute(t, newData, []byte{}, tcaMirredPad)
			val := Mirred{}
			err2 := unmarshalMirred(newData, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			testcase.val.Tm = tm
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Mirred missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalMirred(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
