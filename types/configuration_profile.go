/*
Copyright © 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

// ConfigurationProfile represents configuration profile record in the controller service.
//     ID: unique key
//     Configuration: a JSON structure stored in a string
//     ChangeAt: username of admin that created or updated the configuration
//     ChangeBy: timestamp of the last configuration change
//     Description: a string with any comment(s) about the configuration
type ConfigurationProfile struct {
	ID            int    `json:"id"`
	Configuration string `json:"configuration"`
	ChangedAt     string `json:"changed_at"`
	ChangedBy     string `json:"changed_by"`
	Description   string `json:"description"`
}
