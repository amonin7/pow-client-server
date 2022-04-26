package isrm

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// IntegerSquareRootModulo - SquareRootModulo implementation of PoW interface
// this means that to proof the work the client needs to calculate such number n
// 	that n * n = QuadResidue mod Modulo
type IntegerSquareRootModulo struct {
	Proof       int
	QuadResidue int
	Modulo      int
}

// NewISRM - generates the IntegerSquareRootModulo from the N (initial number) and modulo
func NewISRM(number int, modulo int) *IntegerSquareRootModulo {
	return &IntegerSquareRootModulo{
		Proof:       2,
		QuadResidue: findQuadResidue(number, modulo),
		Modulo:      modulo,
	}
}

// IsProofCorrect - checks, whether the found Proof is correct
func (isrm *IntegerSquareRootModulo) IsProofCorrect() bool {
	return findQuadResidue(isrm.Proof, isrm.Modulo) == isrm.QuadResidue
}

// FindSolution - finds the particular value of Proof, which is correct
func (isrm *IntegerSquareRootModulo) FindSolution() error {
	for ; isrm.Proof < math.MaxInt; isrm.Proof++ {
		if isrm.IsProofCorrect() {
			return nil
		}
	}
	return fmt.Errorf("failed to find solution for %s", isrm.ToShortString())
}

// ToShortString - generates the readable view of IntegerSquareRootModulo
func (isrm *IntegerSquareRootModulo) ToShortString() string {
	return "ISRM{" +
		"quadResidue=" + strconv.Itoa(isrm.QuadResidue) +
		", modulo=" + strconv.Itoa(isrm.Modulo) +
		"}"
}

// findQuadResidue - internal function to find quadratic residue modulo
func findQuadResidue(number int, modulo int) int {
	return ((number % modulo) * (number % modulo)) % modulo
}

// DeserializeIsrm - deserializes IntegerSquareRootModulo from the array of bytes
func DeserializeIsrm(serializedMessage string) (*IntegerSquareRootModulo, error) {
	var isrm IntegerSquareRootModulo
	err := json.Unmarshal([]byte(serializedMessage), &isrm)
	if err != nil {
		return nil, fmt.Errorf("cannot parse message")
	}
	return &isrm, nil
}

// Serialized - serializes IntegerSquareRootModulo to the array of bytes
func (isrm *IntegerSquareRootModulo) Serialized() ([]byte, error) {
	str, err := json.Marshal(isrm)
	if err != nil {
		return nil, fmt.Errorf("cannot serialize isrm - %w", err)
	}
	return str, nil
}
