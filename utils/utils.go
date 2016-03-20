package utils

import "strconv"

func StrToFlt32(str string) (res float32, err error) {
	var flt64 float64
	flt64, err = strconv.ParseFloat(str, 32)
	if err == nil {
		res = float32(flt64)
	}
	return
}

func Flt32ToStr(flt float32) string {
	return strconv.FormatFloat(float64(flt), 'f', 6, 32)
}
