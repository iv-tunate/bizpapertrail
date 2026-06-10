package utils

func Tenary[T any](condition bool, res1, res2 T) T{
	if condition{
		return  res1
	}
	return res2
}