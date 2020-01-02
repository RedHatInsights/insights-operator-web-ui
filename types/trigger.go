/*
Copyright © 2019, 2020 Red Hat, Inc.

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

// Trigger represents trigger record in the controller service
//     ID: unique key
//     Type: ID of trigger type
//     Cluster: cluster ID (not name)
//     Reason: a string with any comment(s) about the trigger
//     Link: link to any document with customer ACK with the trigger
//     TriggeredAt: timestamp of the last configuration change
//     TriggeredBy: username of admin that created or updated the trigger
//     AckedAt: timestamp where the insights operator acked the trigger
//     Parameters: parameters that needs to be pass to trigger code
//     Active: flag indicating whether the trigger is still active or not
type Trigger struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Cluster     string `json:"cluster"`
	Reason      string `json:"reason"`
	Link        string `json:"link"`
	TriggeredAt string `json:"triggered_at"`
	TriggeredBy string `json:"triggered_by"`
	AckedAt     string `json:"acked_at"`
	Parameters  string `json:"parameters"`
	Active      int    `json:"active"`
}
