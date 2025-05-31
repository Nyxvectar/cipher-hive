/**
 * Author:  Nyxvectar Yan
 * Repo:    CipherHive
 * Created: 05/31/2025
 */

package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("32-BIT MD5 HASH : ")
	input, _ := reader.ReadString('\n')
	targetHash := strings.ToLower(strings.TrimSpace(input))
	if len(targetHash) != 32 {
		fmt.Println("MD5_ERR")
		os.Exit(1)
	}

	targetBytes, err := hex.DecodeString(targetHash)
	if err != nil || len(targetBytes) != 16 {
		fmt.Println("INVALID HASH")
		os.Exit(1)
	}
	copy(targetMD5[:], targetBytes)

	start := time.Now()
	numWorkers := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			worker()
		}()
	}

	go func() {
		wg.Wait()
		close(foundChan)
	}()

	if result, ok := <-foundChan; ok {
		fmt.Printf("MATCH UID FOUND : %d\n", result)
	} else {
		fmt.Println("NO MATCH FOUND")
	}
	elapsed := time.Since(start)
	fmt.Printf("TOTAL USED TIME : %.3f\n", elapsed.Seconds())
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

func uitoa(n uint64, buf []byte) []byte {
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

func worker() {
	var localBuf [20]byte
	for {
		start := atomic.LoadUint64(&globalIndex)
		if start+batchSize > maxRange {
			return
		}
		end := atomic.AddUint64(&globalIndex, batchSize)

		for n := start + 1; n <= end; n++ {
			digits := uitoa(n, localBuf[:])
			sum := md5.Sum(digits)
			if sum == targetMD5 {
				select {
				case foundChan <- n:
				default:
				}
				return
			}
		}
	}
}
