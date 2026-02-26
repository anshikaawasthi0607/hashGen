package hashgen

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

const hashLength = 10
//no->string
func encodeBase62(num uint64) string {
	if num == 0 {
		return "0"
	}

	// creating a slice 
	//pre-allocating m/m for performance
	buf := make([]byte, 0, 11)
	for num > 0 {
		buf = append(buf, base62Chars[num%62])
		num /= 62
	}
	//swapping
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}

	return string(buf)
}

func padOrTruncate(s string, length int) string {
	switch {
	case len(s) == length:
		return s
	case len(s) > length:
		return s[:length]
	default:
		padded := make([]byte, length)
		for i := range padded {
			padded[i] = '0'
		}
		copy(padded[length-len(s):], s)
		return string(padded)
	}
}