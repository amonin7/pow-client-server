package pow

// ProofOfWork interface for any Proof Of Work algorithm implementation
type ProofOfWork interface {
	IsProofCorrect() bool
	FindSolution() error
}
