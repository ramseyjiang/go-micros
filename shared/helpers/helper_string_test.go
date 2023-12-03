package helpers

import "testing"

func TestSliceEqualsSlice(t *testing.T) {

	t1 := []string{"abc", "def"}
	t1a := []string{"abc"}
	if SliceEqualsSlice(t1a, t1) {
		t.Errorf("(%v) should NOT be in (%v)", t1a, t1)
	}

	t1b := []string{"def"}
	if SliceEqualsSlice(t1b, t1) {
		t.Errorf("(%v) should NOT be in (%v)", t1b, t1)
	}

	t1c := []string{""}
	if SliceEqualsSlice(t1c, t1) {
		t.Errorf("(%v) should NOT be in (%v)", t1c, t1)
	}

	t1d := []string{"abc", "def"}
	if !SliceEqualsSlice(t1d, t1) {
		t.Errorf("(%v) should be in (%v)", t1d, t1)
	}
}

func TestSliceInSlice(t *testing.T) {

	t1 := []string{"abc", "def"}

	t1a := []string{"abc"}
	if !SliceInSlice(t1a, t1) {
		t.Errorf("(%v) should be in (%v)", t1a, t1)
	}

	t1b := []string{"def"}
	if !SliceInSlice(t1b, t1) {
		t.Errorf("(%v) should be in (%v)", t1b, t1)
	}

	t1c := []string{"xyz"}
	if SliceInSlice(t1c, t1) {
		t.Errorf("(%v) should NOT be in (%v)", t1c, t1)
	}

	t1d := []string{""}
	if SliceInSlice(t1d, t1) {
		t.Errorf("(%v) should NOT be in (%v)", t1d, t1)
	}

	t1e := []string{"abc", "def", "xyz"}
	if SliceInSlice(t1e, t1) {
		t.Errorf("(%v) should NOT be in (%v)", t1e, t1)
	}

	t1f := []string{"abc", "abc", "def"}
	if SliceInSlice(t1f, t1) {
		t.Errorf("(%v) should NOT be in (%v)", t1f, t1)
	}
}

func TestAnyToBool(t *testing.T) {

	testset := make([]string, 0)
	testset = append(testset, "t")
	// testset = append(testset, "\"t\"")
	testset = append(testset, "true")
	// testset = append(testset, "\"true\"")
	testset = append(testset, "1")
	// testset = append(testset, "\"1\"")
	testset = append(testset, "y")
	// testset = append(testset, "\"y\"")
	testset = append(testset, "Y")
	// testset = append(testset, "\"Y\"")
	testset = append(testset, "yes")
	// testset = append(testset, "\"yes\"")
	testset = append(testset, "Yes")
	// testset = append(testset, "\"Yes\"")
	testset = append(testset, "YES")
	// testset = append(testset, "\"YES\"")
	testset = append(testset, "On")
	// testset = append(testset, "\"On\"")
	testset = append(testset, "on")
	// testset = append(testset, "\"on\"")
	testset = append(testset, "\"on\"")
	// testset = append(testset, "\"on\"")
	testset = append(testset, "\"true\"")
	// testset = append(testset, "\"true\"")

	for i := range testset {
		if val, err := AnyToBool(testset[i]); err != nil || val == false {
			if err != nil {
				t.Errorf("(%s) should be true. ERROR: %v", testset[i], err)
			} else {
				t.Errorf("(%s) should be true", testset[i])
			}
		}

		if val, err := AnyToBool([]byte(testset[i])); err != nil || val == false {
			if err != nil {
				t.Errorf("[]byte(%s) should be true. ERROR: %v", testset[i], err)
			} else {
				t.Errorf("[]byte(%s) should be true", testset[i])
			}
		}

		if val, err := AnyToBool([]uint8(testset[i])); err != nil || val == false {
			if err != nil {
				t.Errorf("[]uint8(%s) should be true. ERROR: %v", testset[i], err)
			} else {
				t.Errorf("[]uint8(%s) should be true", testset[i])
			}
		}
	}

	testset = make([]string, 0)
	testset = append(testset, "f")
	// testset = append(testset, "\"f\"")
	testset = append(testset, "false")
	// testset = append(testset, "\"false\"")
	testset = append(testset, "0")
	// testset = append(testset, "\"0\"")
	testset = append(testset, "n")
	// testset = append(testset, "\"n\"")
	testset = append(testset, "N")
	// testset = append(testset, "\"N\"")
	testset = append(testset, "no")
	// testset = append(testset, "\"no\"")
	testset = append(testset, "No")
	// testset = append(testset, "\"No\"")
	testset = append(testset, "NO")
	// testset = append(testset, "\"NO\"")
	testset = append(testset, "Off")
	// testset = append(testset, "\"Off\"")
	testset = append(testset, "off")
	// testset = append(testset, "\"off\"")
	testset = append(testset, "\"off\"")
	// testset = append(testset, "\"off\"")

	for i := range testset {
		if val, err := AnyToBool(testset[i]); err != nil || val == true {
			if err != nil {
				t.Errorf("(%s) should be false. ERROR: %v", testset[i], err)
			} else {
				t.Errorf("(%s) should be false", testset[i])
			}
		}

		if val, err := AnyToBool([]byte(testset[i])); err != nil || val == true {
			if err != nil {
				t.Errorf("[]byte(%s) should be false. ERROR: %v", testset[i], err)
			} else {
				t.Errorf("[]byte(%s) should be false", testset[i])
			}
		}

		if val, err := AnyToBool([]uint8(testset[i])); err != nil || val == true {
			if err != nil {
				t.Errorf("[]uint8(%s) should be false. ERROR: %v", testset[i], err)
			} else {
				t.Errorf("[]uint8(%s) should be false", testset[i])
			}
		}
	}

	// int
	if val, err := AnyToBool(int(1)); err != nil || val == false {
		t.Errorf("(%d) should be true. ERR: %v", int(1), err)
	}
	if val, err := AnyToBool(int(0)); err != nil || val == true {
		t.Errorf("(%d) should be false. ERR: %v", int(0), err)
	}

	// uint8
	if val, err := AnyToBool(uint8(1)); err != nil || val == false {
		t.Errorf("(%d) should be true. ERR: %v", uint8(1), err)
	}
	if val, err := AnyToBool(uint8(0)); err != nil || val == true {
		t.Errorf("(%d) should be false. ERR: %v", uint8(0), err)
	}

	// int8
	if val, err := AnyToBool(int8(1)); err != nil || val == false {
		t.Errorf("(%d) should be true. ERR: %v", int8(1), err)
	}
	if val, err := AnyToBool(int8(0)); err != nil || val == true {
		t.Errorf("(%d) should be false. ERR: %v", int8(0), err)
	}

	// int32
	if val, err := AnyToBool(int32(1)); err != nil || val == false {
		t.Errorf("(%d) should be true. ERR: %v", int32(1), err)
	}
	if val, err := AnyToBool(int32(0)); err != nil || val == true {
		t.Errorf("(%d) should be false. ERR: %v", int32(0), err)
	}

	// uint32
	if val, err := AnyToBool(uint32(1)); err != nil || val == false {
		t.Errorf("(%d) should be true. ERR: %v", uint32(1), err)
	}
	if val, err := AnyToBool(uint32(0)); err != nil || val == true {
		t.Errorf("(%d) should be false. ERR: %v", uint32(0), err)
	}

	// bool
	if val, err := AnyToBool(bool(true)); err != nil || val == false {
		t.Errorf("(%v) should be true. ERR: %v", bool(true), err)
	}
	if val, err := AnyToBool(bool(false)); err != nil || val == true {
		t.Errorf("(%v) should be false. ERR: %v", bool(false), err)
	}

	// bool
	if v, err := AnyToBool(nil); err == nil {
		t.Errorf("(nil) should be error. Value: %v", v)
	}

	// empty sting
	if val, err := AnyToBool(""); err != nil || val == true {
		t.Errorf("(%s) should be false. ERR: %v", "", err)
	}

}

func TestStripStringQuotes(t *testing.T) {

	if StripStringQuotes("") != "" {
		t.Errorf("() should be empty")
	}

	if StripStringQuotes("\"") != "\"" {
		t.Errorf("() should be 1 quote")
	}
	if StripStringQuotes("\"\"") != "" {
		t.Errorf("() should be empty")
	}
}

func TestStringIsNil(t *testing.T) {
	if StringIsNil("") != true {
		t.Errorf("() should be true")
	}
	if StringIsNil("null") != true {
		t.Errorf("null should be true")
	}
	if StringIsNil("nil") != true {
		t.Errorf("nil should be true")
	}
	if StringIsNil("<nil>") != true {
		t.Errorf("<nil> should be true")
	}
	if StringIsNil("<n>") != false {
		t.Errorf("<n> should be false")
	}
}

func TestConvertToNonMacrons(t *testing.T) {
	macrons := "āēīōū ĀĒĪŌŪ"
	nomacrons := "aeiou AEIOU"

	result := ConvertToNonMacrons(macrons)

	if result != nomacrons {
		t.Errorf("%s != %s", result, nomacrons)
	}
}
