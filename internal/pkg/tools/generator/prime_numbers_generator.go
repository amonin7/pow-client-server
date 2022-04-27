package generator

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

func GeneratePrime() int {
	f, err := os.OpenFile("internal/pkg/tools/generator/primes.txt", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return 1_000_000_007
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("failed to close file - %s\n", err)
		}
	}(f)

	primeIndex := rand.Intn(650_000)
	sc := bufio.NewScanner(f)
	i := 0
	for sc.Scan() {
		if i == primeIndex {
			primeStr := sc.Text() // GET the line string
			prime, err := strconv.Atoi(primeStr)
			if err != nil {
				log.Fatalf("failed to parse int: %v", err)
				return 1_000_000_007
			} else {
				return prime
			}
		} else {
			i++
		}
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return 1_000_000_007
	}
	return 1_000_000_007
}
