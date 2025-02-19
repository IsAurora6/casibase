// Copyright 2025 The Casibase Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sashabaranov/go-openai"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	hunyuan "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/hunyuan/v20230901"
)

type TencentCloudClient struct {
	credential  *common.Credential
	endpoint    string
	subType     string
	apiKey      string
	temperature float32
	topP        float32
}

func NewTencentCloudProvider(secretID, secretKey, endpoint, subType string, temperature, topP float32) (*TencentCloudClient, error) {
	if strings.TrimSpace(endpoint) == "" || endpoint == "hunyuan.tencentcloudapi.com" {
		if strings.TrimSpace(secretID) == "" || strings.TrimSpace(secretKey) == "" {
			return nil, fmt.Errorf("invalid credentials: secretID and secretKey cannot be empty")
		}
		endpoint = "hunyuan.tencentcloudapi.com"
		return &TencentCloudClient{
			credential: common.NewCredential(secretID, secretKey),
			endpoint:   endpoint,
			subType:    subType,
		}, nil
	} else {
		return &TencentCloudClient{
			apiKey:      secretKey,
			endpoint:    endpoint,
			subType:     subType,
			temperature: temperature,
			topP:        topP,
		}, nil
	}
}

func (c *TencentCloudClient) GetPricing() string {
	return `Pricing information for Tencent Cloud models is not yet available.`
}

func (c *TencentCloudClient) QueryText(question string, writer io.Writer, history []*RawMessage, prompt string, knowledgeMessages []*RawMessage) (*ModelResult, error) {
	if c.credential == nil {
		return c.QueryTextByOpenAI(question, writer, history, prompt, knowledgeMessages)
	} else {
		return c.QueryTextByHunyuan(question, writer, history, prompt, knowledgeMessages)
	}
}

func (c *TencentCloudClient) QueryTextByHunyuan(question string, writer io.Writer, history []*RawMessage, prompt string, knowledgeMessages []*RawMessage) (*ModelResult, error) {
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = c.endpoint
	client, err := hunyuan.NewClient(c.credential, "", clientProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to create Hunyuan client: %v", err)
	}

	request := hunyuan.NewChatCompletionsRequest()
	request.Model = common.StringPtr(c.subType)

	request.Messages = []*hunyuan.Message{
		{
			Role:    common.StringPtr("user"),
			Content: common.StringPtr(question),
		},
	}

	response, err := client.ChatCompletions(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("TencentCloud SDK error: %s", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}

	if response.Response == nil || len(response.Response.Choices) == 0 {
		return nil, fmt.Errorf("unexpected empty response")
	}
	respText := strings.TrimSpace(*response.Response.Choices[0].Message.Content)

	_, err = fmt.Fprint(writer, respText)
	if err != nil {
		return nil, fmt.Errorf("failed to write response: %v", err)
	}

	modelResult, err := getDefaultModelResult(c.subType, question, respText)
	if err != nil {
		return nil, err
	}

	return modelResult, nil
}

func (c *TencentCloudClient) QueryTextByOpenAI(question string, writer io.Writer, history []*RawMessage, prompt string, knowledgeMessages []*RawMessage) (*ModelResult, error) {
	ctx := context.Background()
	flusher, ok := writer.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("writer does not implement http.Flusher")
	}

	config := openai.DefaultConfig(c.apiKey)
	config.BaseURL = c.endpoint
	client := openai.NewClientWithConfig(config)

	// set request params
	messages := []openai.ChatCompletionMessage{
		{
			Role:    "user",
			Content: question,
		},
	}

	modelSplit := strings.Split(c.endpoint, "/")
	model := modelSplit[len(modelSplit)-1]

	request := openai.ChatCompletionRequest{
		Model:       model,
		Messages:    messages,
		Temperature: c.temperature,
		TopP:        c.topP,
		Stream:      true,
	}

	flushData := func(data string) error {
		if _, err := fmt.Fprintf(writer, "event: message\ndata: %s\n\n", data); err != nil {
			return err
		}
		flusher.Flush()
		return nil
	}
	modelResult := &ModelResult{}

	promptTokenCount, err := OpenaiNumTokensFromMessages(messages, "gpt-4") // calculate token
	if err != nil {
		return nil, err
	}
	modelResult.PromptTokenCount = promptTokenCount
	modelResult.TotalTokenCount = modelResult.PromptTokenCount + modelResult.ResponseTokenCount

	if strings.HasPrefix(question, "$CasibaseDryRun$") {
		if GetOpenAiMaxTokens(c.subType) > modelResult.TotalTokenCount {
			return modelResult, nil
		} else {
			return nil, fmt.Errorf("exceed max tokens")
		}
	}

	stream, err := client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		return nil, err
	}
	defer stream.Close()
	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if len(response.Choices) == 0 {
			continue
		}

		data := response.Choices[0].Delta.Content
		err = flushData(data)
		if err != nil {
			return nil, err
		}

		responseTokenCount, err := GetTokenSize("gpt-4", data)
		if err != nil {
			return nil, err
		}
		modelResult.ResponseTokenCount += responseTokenCount
		modelResult.TotalTokenCount = modelResult.PromptTokenCount + modelResult.ResponseTokenCount
	}
	return modelResult, nil
}
