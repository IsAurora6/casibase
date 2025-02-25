// Copyright 2024 The Casibase Authors. All Rights Reserved.
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

package embedding

import (
	"context"
)

type OllamaEmbeddingProvider struct {
	subType     string
	secretKey   string
	providerUrl string
}

func NewOllamaEmbeddingProvider(subType string, secretKey string, providerUrl string) (*OllamaEmbeddingProvider, error) {
	return &OllamaEmbeddingProvider{
		subType:     subType,
		secretKey:   secretKey,
		providerUrl: providerUrl,
	}, nil
}

func (p *OllamaEmbeddingProvider) GetPricing() string {
	return `URL:
https://ollama.com/search

Not charged
`
}

func (p *OllamaEmbeddingProvider) calculatePrice(res *EmbeddingResult) error {
	return nil
}

func (p *OllamaEmbeddingProvider) QueryVector(text string, ctx context.Context) ([]float32, *EmbeddingResult, error) {
	localEmbeddingProvider, err := NewLocalEmbeddingProvider("Custom", p.subType, p.secretKey, p.providerUrl)
	if err != nil {
		return nil, nil, err
	}
	vector, embeddingResult, err := localEmbeddingProvider.QueryVector(text, ctx)
	if err != nil {
		return nil, nil, err
	}
	return vector, embeddingResult, nil
}
