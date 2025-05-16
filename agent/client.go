package agent

import "fmt"

type AgentClientInterface interface {
	Query(txId string, data string) (string, error)
}

func NewAgentClient(providerType, text string) (AgentClientInterface, error) {
	var res AgentClientInterface
	var err error
	if providerType == "mcp" {
		res, err = NewMcpClient(text)
	} else {
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}
