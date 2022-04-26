package processor

import (
	"fmt"
	"math"
	"math/rand"
	"pow-client-server/internal/pkg/model"
	"pow-client-server/internal/pkg/pow/isrm"
	"pow-client-server/internal/pkg/tools/generator"
	"pow-client-server/internal/pkg/tools/manager"
	"strconv"
)

// Process - this function is designed to process request from the client on the server side
func Process(input string, clientUrl string) (*model.Message, error) {
	message, err := model.DeserializeMessage(input)
	if err != nil {
		return nil, err
	}
	switch message.Header {
	case model.ChallengeReq: // in this case we need to generate the challenge for the client
		return processChallengeRequest(clientUrl)
	case model.ResourceReq: // in this case we need to check, whether clients solution is correct. And if so - provide him with resource
		return processResourceRequest(message, clientUrl)
	default:
		return nil, fmt.Errorf("received unsupported message type: %d", message.Header)
	}
}

// processChallengeRequest - processes particular request for challenge, received from client
func processChallengeRequest(clientUrl string) (*model.Message, error) {
	fmt.Printf("received request for challenge from %s\n", clientUrl)
	prime := generator.GeneratePrime()
	integerSquareRootModulo := isrm.NewISRM(rand.Intn(math.MaxInt16)%prime, prime)
	isrmJson, err := integerSquareRootModulo.Serialized()
	if err != nil {
		return nil, fmt.Errorf("failed to generate json from isrm: %w", err)
	}
	msg := model.Message{
		Header:  model.ChallengeRes,
		Payload: string(isrmJson),
	}
	return &msg, nil
}

// processResourceRequest - processes request for resource, received from client
// firstly, checks, whether the proof, found by client is correct
// 	- if yes - sends the resource to client
// 	- if no - returns the corresponding error
func processResourceRequest(message *model.Message, clientUrl string) (*model.Message, error) {
	fmt.Printf("received request for resource from %s, checking if proof is correct...\n", clientUrl)
	integerSquareRootModulo, err := isrm.DeserializeIsrm(message.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize IntegerSquareRootModulo: %w", err)
	}

	if !integerSquareRootModulo.IsProofCorrect() {
		return nil, fmt.Errorf(
			"invalid proof %s for %s",
			strconv.Itoa(integerSquareRootModulo.Proof),
			integerSquareRootModulo.ToShortString())
	}

	fmt.Printf("proof, received from client is correct %s\n", clientUrl)
	msg := model.Message{
		Header:  model.ResourceRes,
		Payload: manager.GetResource(),
	}
	return &msg, nil
}
