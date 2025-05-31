/**
 * Author:  Nyxvectar Yan
 * Repo:    CipherHive
 * Created: 05/31/2025
 */

package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
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
