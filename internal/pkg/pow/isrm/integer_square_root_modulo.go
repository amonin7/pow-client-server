package isrm

import (
	"fmt"
	"math"
	"strconv"
)

type IntegerSquareRootModulo struct {
	Proof       int
	QuadResidue int
	Modulo      int
}

func (isrm IntegerSquareRootModulo) IsProofCorrect() bool {
	return FindQuadResidue(isrm.Proof, isrm.Modulo) == isrm.QuadResidue
}

func (isrm IntegerSquareRootModulo) FindSolution() error {
	for i := 0; i < math.MaxInt; i++ {
		if isrm.IsProofCorrect() {
			println(isrm.Proof)
			return nil
		}
	}
	return fmt.Errorf("failed to find solution for %s", isrm.Serialize())
}

func (isrm IntegerSquareRootModulo) Serialize() string {
	return "ISRM{" +
		"quadResidue=" + strconv.Itoa(isrm.QuadResidue) +
		"modulo=" + strconv.Itoa(isrm.Modulo) +
		"}"
}

func FindQuadResidue(number int, modulo int) int {
	return (number * number) % modulo
}

func NewISRM(residue int, modulo int) *IntegerSquareRootModulo {
	return &IntegerSquareRootModulo{
		Proof:       2,
		QuadResidue: residue,
		Modulo:      modulo,
	}
}
