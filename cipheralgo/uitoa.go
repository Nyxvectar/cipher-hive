/**
 * Author:  Nyxvectar Yan
 * Repo:    CipherHive
 * Created: 06/01/2025
 */

package cipheralgo

func Uitoa(n uint64, buf []byte) []byte {
	i := len(buf)
	for n >= 10 {
		i--
		q := n / 10
		buf[i] = byte(n%10) + '0'
		n = q
	}
	i--
	buf[i] = byte(n) + '0'
	return buf[i:]
}
