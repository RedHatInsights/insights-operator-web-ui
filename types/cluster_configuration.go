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

// ClusterConfiguration represents cluster configuration record in the controller service.
//     ID: unique key
//     Cluster: cluster ID (not name)
//     Configuration: a JSON structure stored in a string
//     ChangeAt: timestamp of the last configuration change
//     ChangeBy: username of admin that created or updated the configuration
//     Active: flag indicating whether the configuration is active or not
//     Reason: a string with any comment(s) about the cluster configuration
type ClusterConfiguration struct {
	ID            int    `json:"id"`
	Cluster       string `json:"cluster"`
	Configuration string `json:"configuration"`
	ChangedAt     string `json:"changed_at"`
	ChangedBy     string `json:"changed_by"`
	Active        string `json:"active"`
	Reason        string `json:"reason"`
}
