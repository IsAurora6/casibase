// Copyright 2023 The Casibase Authors. All Rights Reserved.
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
	"io"
)

type OllamaModelProvider struct {
	endpoint    string
	subType     string
	apiKey      string
	temperature float32
	topP        float32
}

func NewOllamaModelProvider(secretKey, endpoint, subType string, temperature, topP float32) (*OllamaModelProvider, error) {
	return &OllamaModelProvider{
		apiKey:      secretKey,
		endpoint:    endpoint,
		subType:     subType,
		temperature: temperature,
		topP:        topP,
	}, nil
}

func (p *OllamaModelProvider) GetPricing() string {
	return `URL:
https://ollama.com/search

Not charged
`
}

func (p *OllamaModelProvider) QueryText(question string, writer io.Writer, history []*RawMessage, prompt string, knowledgeMessages []*RawMessage) (*ModelResult, error) {
	baseUrl := p.endpoint
	localProvider, err := NewLocalModelProvider("Custom", "custom-model", p.apiKey, p.temperature, p.topP, 0, 0, baseUrl, p.subType)
	if err != nil {
		return nil, err
	}

	modelResult, err := localProvider.QueryText(question, writer, history, prompt, knowledgeMessages)
	if err != nil {
		return nil, err
	}
	return modelResult, nil
}
