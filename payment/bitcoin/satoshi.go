// Copyright (c) 2014-2017 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package bitcoin

// convert a string to a Satoshi value
//
// i.e. "0.00000001" will convert to uint64(1)
//
// Note: Invalid characters are simply ignored and the conversion
//       simply stops after 8 decimal places have been processed.
//       Extra decimal points will also be ignored.
func convertToSatoshi(btc []byte) uint64 {

	s := uint64(0)
	point := false
	decimals := 0
	for _, b := range btc {
		if b >= '0' && b <= '9' {
			s *= 10
			s += uint64(b - '0')
			if point {
				decimals += 1
				if decimals >= 8 {
					break
				}
			}
		} else if '.' == b {
			point = true
		}
	}
	for decimals < 8 {
		s *= 10
		decimals += 1
	}

	return s
}
