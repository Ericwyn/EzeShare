package utils

func Check(b bool, trueRes any, falseRes any) any {
	if b {
		return trueRes
	} else {
		return falseRes
	}
}
