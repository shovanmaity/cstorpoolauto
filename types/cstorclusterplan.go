/*
Copyright 2019 The MayaData Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// CStorClusterPlan is a kubernetes custom resource that plans
// the resources especially nodes to form the CStorPoolCluster
type CStorClusterPlan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   CStorClusterPlanSpec   `json:"spec"`
	Status CStorClusterPlanStatus `json:"status"`
}

// CStorClusterPlanSpec has the plan details required to form
// CStorPoolCluster
type CStorClusterPlanSpec struct {
	Nodes []CStorClusterPlanNode `json:"nodes"`
}

// CStorClusterPlanNode has the node details that is used to
// form CStorPoolCluster
type CStorClusterPlanNode struct {
	Name string    `json:"name"`
	UID  types.UID `json:"uid"`
}

// CStorClusterConfigReference refers to CStorClusterConfig
// which is the trigger to CStorClusterPlan
type CStorClusterConfigReference struct {
	UID types.UID `json:"uid"`
}

// CStorClusterPlanStatus represents the current state of
// CStorClusterPlan
type CStorClusterPlanStatus struct {
	Phase      CStorClusterPlanStatusPhase       `json:"phase"`
	Conditions []CStorClusterPlanStatusCondition `json:"conditions"`
}

// CStorClusterPlanStatusPhase reports the current phase of
// CStorClusterPlan
type CStorClusterPlanStatusPhase string

const (
	// CStorClusterPlanStatusPhaseError indicates error in
	// CStorClusterPlan
	CStorClusterPlanStatusPhaseError CStorClusterPlanStatusPhase = "Error"

	// CStorClusterPlanStatusPhaseOnline indicates
	// CStorClusterPlan in Online state i.e. no error or warning
	CStorClusterPlanStatusPhaseOnline CStorClusterPlanStatusPhase = "Online"
)

// CStorClusterPlanStatusCondition represents a condition
// that represents the current state of CStorClusterPlan
type CStorClusterPlanStatusCondition struct {
	Type             ConditionType  `json:"type"`
	Status           ConditionState `json:"status"`
	Reason           string         `json:"reason,omitempty"`
	LastObservedTime string         `json:"lastObservedTime"`
}

// MakeListMapOfPlanNodes returns a slice of maps from
// the given slice of CStorClusterPlanNode
func MakeListMapOfPlanNodes(given []CStorClusterPlanNode) []interface{} {
	var listMap []interface{}
	for _, node := range given {
		listMap = append(listMap,
			map[string]interface{}{
				"name": node.Name,
				"uid":  string(node.UID),
			},
		)
	}
	return listMap
}

// CStorClusterPlanNodeList is a typed representation of a list
// of CStorClusterPlanNode(s)
type CStorClusterPlanNodeList []CStorClusterPlanNode

// CStorClusterPlanNodeNil represents a nil CStorClusterPlanNode
var CStorClusterPlanNodeNil = CStorClusterPlanNode{}

// FindByNameAndUID returns the node instance based on the
// given name & uid
func (l CStorClusterPlanNodeList) FindByNameAndUID(name string, uid types.UID) CStorClusterPlanNode {
	for _, node := range l {
		if node.Name == name && node.UID == uid {
			return node
		}
	}
	return CStorClusterPlanNodeNil
}

// Contains returns true if the given node name & uid
// is available in this list
func (l CStorClusterPlanNodeList) Contains(name string, uid types.UID) bool {
	return l.FindByNameAndUID(name, uid) != CStorClusterPlanNodeNil
}

// ContainsAll returns true if the given CStorClusterPlanNode
// list is available in this list
func (l CStorClusterPlanNodeList) ContainsAll(given []CStorClusterPlanNode) bool {
	if len(l) == 0 && len(l) == len(given) {
		return true
	}
	if len(l) != len(given) {
		return false
	}
	for _, planNode := range given {
		if !l.Contains(planNode.Name, planNode.UID) {
			return false
		}
	}
	return true
}
