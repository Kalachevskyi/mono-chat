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

// Package config implements the application configuration
package config

//Regexp patterns
const (
	DateShortRegexpPattern = `^\d{2}-\d{2}` //range of date for the current month
	DateRegexpPattern      = `\d{2}\.\d{2}\.\d{4}-\d{2}\.\d{2}\.\d{4}`
	DateTimeRegexpPattern  = `\d{2}\.\d{2}\.\d{4}T\d{2}\.\d{2}-\d{2}\.\d{2}\.\d{4}T\d{2}\.\d{2}`
)

// TimeLocation - application location
const TimeLocation = "Europe/Kiev"
