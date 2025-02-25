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

package object

import (
	"fmt"
	"strings"
)

type Factor struct {
	Name     string    `xorm:"varchar(100)" json:"name"`
	Category string    `xorm:"varchar(100)" json:"category"`
	Color    string    `xorm:"varchar(100)" json:"color"`
	Data     []float64 `xorm:"varchar(1000)" json:"data"`
}

func (factor *Factor) GetDataKey() string {
	sData := []string{}
	for _, f := range factor.Data {
		sData = append(sData, fmt.Sprintf("%f", f))
	}
	return strings.Join(sData, "|")
}
