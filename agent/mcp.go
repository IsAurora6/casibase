package agent

import (
	"github.com/mark3labs/mcp-go/client"
)

type McpAgentClient struct {
	Client *client.Client
}

func NewMcpClient(text string) (*McpAgentClient, error) {
	var c *client.Client
	var err error
	isSSE, command, envVars, args, url, err := ParseMCPConfig(text)
	if !isSSE {
		c, err = client.NewStdioMCPClient(
			command,
			envVars,
			args...,
		)
		if err != nil {
			return nil, err
		}
	} else {
		c, err = client.NewSSEMCPClient(url)
		if err != nil {
			return nil, err
		}
	}
	// Create a new MCP client
	return &McpAgentClient{
		Client: c,
	}, nil
}

func (client McpAgentClient) Query(blockId string, data string) (string, error) {
	return "", nil
	//// simulate the situation that error occurs
	//if strings.HasSuffix(data["id"], "0") {
	//	return "", fmt.Errorf("some error occurred in the ChainTencentChainmakerClient::Commit operation")
	//}
	//
	//// Query the data from the blockchain
	//// Write some code... (if error occurred, handle it as above)
	//
	//// assume the chain data are retrieved from the blockchain, here we just generate it statically
	//chainData := map[string]string{"organization": "casbin"}
	//
	//// Check if the data are matched with the chain data
	//res := "Matched"
	//if chainData["organization"] != data["organization"] {
	//	res = "Mismatched"
	//}
	//
	//// simulate the situation that mismatch occurs
	//if strings.HasSuffix(blockId, "2") || strings.HasSuffix(blockId, "4") || strings.HasSuffix(blockId, "6") || strings.HasSuffix(blockId, "8") || strings.HasSuffix(blockId, "0") {
	//	res = "Mismatched"
	//}
	//
	//return fmt.Sprintf("The query result for block [%s] is: %s", blockId, res), nil
}
