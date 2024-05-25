package utils

import "errors"

const maxInt = int(^uint(0) >> 1)

var ErrTooLarge = errors.New("buf too large")

// GrowByteSlice a byte slice's capacity or allocate more space(len) by n. It returns
// the previous length of b.
func GrowByteSlice(b *[]byte, n int) (int, error) {
	c := cap(*b)
	l := len(*b)
	a := c - l // available
	if n <= a {
		*b = (*b)[:l+n]
		return l, nil
	}

	// above n <= a condition will handle this as well
	// if l+n <= c {
	// 	// if needed length is lte c
	// 	return l, nil
	// }

	if c > maxInt-c-n {
		// too large
		return l, ErrTooLarge
	}

	newBuf := make([]byte, c*2+n)
	copy(newBuf, (*b)[0:])
	*b = newBuf[:l+n]
	return l, nil
}

// GrowByteSliceCap grows the capacity of byte slice by n
func GrowByteSliceCap(b *[]byte, n int) error {
	c := cap(*b)
	l := len(*b)
	a := c - l // available
	if n <= a {
		return nil
	}

	// if l+n <= c {
	// 	// if needed length is lte c
	// 	return nil
	// }

	if c > maxInt-c-n {
		// too large
		return ErrTooLarge
	}

	newBuf := make([]byte, c+n)
	copy(newBuf, (*b)[0:])
	*b = newBuf[:l]
	return nil
}
