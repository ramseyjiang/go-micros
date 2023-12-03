package strhelpers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ramseyjiang/go-micros/shared/apierror"
)

func AnyToBool(src interface{}) (bool, error) {
	if src == nil {
		return false, errors.New("src is nil")
	}
	if b, castOK := src.(int); castOK {
		return (b == 1), nil
	}
	if b, castOK := src.(uint8); castOK {
		return (b == 1), nil
	}
	if b, castOK := src.(int8); castOK {
		return (b == 1), nil
	}
	if b, castOK := src.(uint32); castOK {
		return (b == 1), nil
	}
	if b, castOK := src.(int32); castOK {
		return (b == 1), nil
	}
	if b, castOK := src.([]byte); castOK {
		return StringToBool(string(b)), nil
	}
	if b, castOK := src.([]uint8); castOK {
		return StringToBool(string(b)), nil
	}
	if b, castOK := src.(string); castOK {
		return StringToBool(b), nil
	}
	if b, castOK := src.(bool); castOK {
		return (b), nil
	}
	strVersion := fmt.Sprintf("%v", src)
	return StringToBool(strVersion), nil
}

func StringToBool(str string) bool {
	value, err := StringToBoolOrErr(str)
	if err != nil {
		return false
	}
	return value
}

func BoolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func StringToBoolOrErr(str string) (bool, error) {
	str = strings.ToLower(strings.TrimSpace(str))
	if strings.HasPrefix(str, "\"") && strings.HasSuffix(str, "\"") {
		str = strings.TrimPrefix(strings.TrimSuffix(str, "\""), "\"")
	}
	switch str {
	case "":
		return false, apierror.NewAPIDebug(nil, 400, "StringToBoolOrErr", "Empty Boolean value")
	case "on":
		return true, nil
	case "off":
		return false, nil
	}
	switch str[0] {
	case 't', 'y', 'e', '1': // t, true, y, yes, enable, 1
		return true, nil
	case 'f', 'n', 'd', '0': // f, false, n, no, disabled, 0
		return false, nil
	}

	return false, apierror.NewAPIError(nil, 400, "StringToBoolOrErr", "Invalid Boolean value (%s)", str)
}

func StripStringQuotes(in string) string {
	if len(in) >= 2 && in[0] == byte('"') && in[len(in)-1] == byte('"') {
		return in[1 : len(in)-1]
	}
	return in
}

func StringIsNil(input string) bool {
	if input == "" || input == "null" || input == "nil" || input == "<nil>" {
		return true
	}
	return false
}

func NullStringToInt(input string) (int, error) {
	if StringIsNil(input) {
		return 0, nil
	}
	return strconv.Atoi(input)
}

func CombineStringSlices(slice1 []string, slice2 []string) []string {
	for i := range slice2 {
		if StringInSlice(slice2[i], slice1) {
			continue
		}
		slice1 = append(slice1, slice2[i])
	}
	return slice1
}

func StringInSlice(str string, strslice []string) bool {
	for _, v := range strslice {
		if v == str {
			return true
		}
	}
	return false
}

func StringInSliceCaseInsensitive(casedstr string, strslice []string) bool {
	for _, v := range strslice {
		if strings.EqualFold(casedstr, v) {
			return true
		}
	}
	return false
}

// NumStringInSlice returns the numbers of instances of str in strslice
func NumStringInSlice(str string, strslice []string) int {
	count := 0
	for _, v := range strslice {
		if v == str {
			count++
		}
	}
	return count
}

func SliceInSlice(isSlice []string, inSlice []string) bool {
	for _, v := range isSlice {
		if NumStringInSlice(v, isSlice) > NumStringInSlice(v, inSlice) {
			return false
		}
	}
	return true
}

func SliceEqualsSlice(isSlice []string, inSlice []string) bool {
	if len(isSlice) != len(inSlice) {
		return false
	}
	return SliceInSlice(isSlice, inSlice)
}

func Int32InSlice(i int32, int32slice []int32) bool {
	for _, v := range int32slice {
		if v == i {
			return true
		}
	}
	return false
}

func Int32SliceInInt32Slice(isSlice []int32, inSlice []int32) bool {
	for _, v := range isSlice {
		if Int32InSlice(v, inSlice) {
			return true
		}
	}
	return false
}

func CSVtoSlice(in string) []string {
	if len(in) == 0 {
		return make([]string, 0)
	}
	return strings.Split(in, ",")
}

func StringReverse(s string) string {
	strLen := len(s)
	resp := make([]byte, strLen)
	for i := 0; i < strLen; {
		thisrune, n := utf8.DecodeRuneInString(s[i:])
		i += n
		utf8.EncodeRune(resp[strLen-i:], thisrune)
	}
	return string(resp)
}

func ConvertToNonMacrons(s string) string {
	r := strings.NewReplacer(
		"ā", "a",
		"ē", "e",
		"ī", "i",
		"ō", "o",
		"ū", "u",
		"Ā", "A",
		"Ē", "E",
		"Ī", "I",
		"Ō", "O",
		"Ū", "U",
	)
	return r.Replace(s)
}
