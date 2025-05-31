/**
 * Author:  Nyxvectar Yan
 * Repo:    CipherHive
 * Created: 05/31/2025
 */

package main

import "fmt"

func main() {
	fmt.Printf("Hello, gopher!")
}

const (
	batchSize = 10000
	maxRange  = 1 << 60
)

var (
	targetMD5   [16]byte
	foundChan   = make(chan uint64, 1)
	globalIndex uint64
)
