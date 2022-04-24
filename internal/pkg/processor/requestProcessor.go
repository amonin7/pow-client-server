package processor

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"pow-client-server/internal/pkg/model"
	"pow-client-server/internal/pkg/pow/isrm"
)

func Process(input string, clientInfo string) (*model.Message, error) {
	msg, err := model.DeserializeMessage(input)
	if err != nil {
		return nil, err
	}
	switch msg.Header {
	case model.ChallengeReq:
		prime := generatePrime()
		residue := rand.Intn(prime)
		integerSquareRootModulo := isrm.NewISRM(residue, prime)
		isrmJson, err := json.Marshal(integerSquareRootModulo)
		if err != nil {
			return nil, fmt.Errorf("failed to generato json from isrm: %v", err)
		}
		msg := model.Message{
			Header:  model.ChallengeRes,
			Payload: string(isrmJson),
		}
		return &msg, nil
	case model.ResourceReq:
		fmt.Printf("client %s requests resource with payload %s\n", clientInfo, msg.Payload)
		// parse client's solution
		var integerSquareRootModulo isrm.IntegerSquareRootModulo
		err := json.Unmarshal([]byte(msg.Payload), &integerSquareRootModulo)
		if err != nil {
			return nil, fmt.Errorf("err unmarshal hashcash: %w", err)
		}

		isCorrect := integerSquareRootModulo.IsProofCorrect()
		if !isCorrect {
			return nil, fmt.Errorf("invalid proof for %s", integerSquareRootModulo.Serialize())
		}
		fmt.Printf("client %s succesfully computed proof %s\n", clientInfo, msg.Payload)
		msg := model.Message{
			Header:  model.ResourceRes,
			Payload: getResource(),
		}
		return &msg, nil
	default:
		return nil, fmt.Errorf("unknown header")
	}
}

func generatePrime() int {
	return 13
}

func getResource() string {
	return "this is the resource wanted to be received"
}
