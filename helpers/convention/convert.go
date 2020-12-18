package convention

import "strconv"

func StringToInt(text string) int {
	intType, _ := strconv.Atoi(text)
	return intType
}

func StringToInt8(text string) int8 {
	int8Type, _ := strconv.Atoi(text)
	return int8(int8Type)
}

func StringToUint64(text string) uint64 {
	int64Obj, _ := strconv.ParseUint(text, 10, 64)
	return int64Obj
}

func Uint64ToString(ui64 uint64) string {
	s := strconv.FormatUint(ui64, 10)
	return s
}
