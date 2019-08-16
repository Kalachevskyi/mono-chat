// Copyright © 2019 Volodymyr Kalachevskyi <v.kalachevskyi@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package entities is an implementation of core application objects
package entities

// Transaction - represents monobank's transaction entity
type Transaction struct {
	ID              string `json:"id"`
	Time            int    `json:"time"`
	Description     string `json:"description"`
	Mcc             int    `json:"mcc"`
	Hold            bool   `json:"hold"`
	Amount          int    `json:"amount"`
	OperationAmount int    `json:"operationAmount"`
	CurrencyCode    int    `json:"currencyCode"`
	CommissionRate  int    `json:"commissionRate"`
	CashbackAmount  int    `json:"cashbackAmount"`
	Balance         int    `json:"balance"`
}

// MonoAPIError - represents monobank's error entity
type MonoAPIError struct {
	ErrorDescription string `json:"errorDescription"`
}
