package pow

type ProofOfWork interface {
	IsProofCorrect() bool
	FindSolution()
}
