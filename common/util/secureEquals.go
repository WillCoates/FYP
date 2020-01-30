package util

// SecureEqualsStr compares two strings in constant time.
func SecureEqualsStr(a, b string) bool {
	lengthA := len(a)
	lengthB := len(b)
	check := lengthA ^ lengthB

	for i := 0; i < lengthB; i++ {
		check |= int(a[i%lengthA] ^ b[i])
	}

	return check == 0
}

// SecureEquals compares two byte slices in constant time.
func SecureEquals(a, b []byte) bool {
	lengthA := len(a)
	lengthB := len(b)
	check := lengthA ^ lengthB

	for i := 0; i < lengthB; i++ {
		check |= int(a[i%lengthA] ^ b[i])
	}

	return check == 0
}
