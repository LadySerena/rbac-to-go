/*
 * Copyright © 2021 Serena Tiede
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package parser

import (
	"bytes"
	"encoding/json"
	"os"

	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

func Run() {

}

func ParseRole() (*v1.ClusterRole, error) {
	content, readErr := os.ReadFile("test-resources/victoria-metrics-role.yaml")
	if readErr != nil {
		return nil, readErr
	}
	documentList := bytes.Split(content, []byte("---\n"))
	for _, document := range documentList {
		if len(document) < 2 {
			continue // skip things without content to deal with spliting yaml documents within 1 file
		}
		raw := unstructured.Unstructured{}
		unmarshalErr := yaml.Unmarshal(document, &raw)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
		metadata := metav1.ObjectMeta{
			Name:                       raw.GetName(),
			GenerateName:               raw.GetGenerateName(),
			Namespace:                  raw.GetNamespace(),
			UID:                        raw.GetUID(),
			ResourceVersion:            raw.GetResourceVersion(),
			Generation:                 raw.GetGeneration(),
			CreationTimestamp:          raw.GetCreationTimestamp(),
			DeletionTimestamp:          raw.GetDeletionTimestamp(),
			DeletionGracePeriodSeconds: raw.GetDeletionGracePeriodSeconds(),
			Labels:                     raw.GetLabels(),
			Annotations:                raw.GetAnnotations(),
			OwnerReferences:            raw.GetOwnerReferences(),
			Finalizers:                 raw.GetFinalizers(),
			ClusterName:                raw.GetClusterName(),
			ManagedFields:              raw.GetManagedFields(),
		}

		typeMeta := metav1.TypeMeta{
			Kind:       raw.GetKind(),
			APIVersion: raw.GetAPIVersion(),
		}

		var rbacRule []v1.PolicyRule

		if typeMeta.Kind == "ClusterRole" || typeMeta.Kind == "Role" {
			rawRules := raw.Object["rules"]
			ruleInterfaceList := rawRules.([]interface{}) // todo safely assert type
			for _, ruleInterface := range ruleInterfaceList {
				testRule := ruleInterface.(map[string]interface{}) // todo safely assert type
				ruleBytes, marshalErr := json.Marshal(testRule)
				if marshalErr != nil {
					return nil, marshalErr
				}
				parsedRule := v1.PolicyRule{}
				parseErr := json.Unmarshal(ruleBytes, &parsedRule)
				if parseErr != nil {
					return nil, parseErr
				}
				rbacRule = append(rbacRule, parsedRule)
			}
		}
		return &v1.ClusterRole{
			TypeMeta:        typeMeta,
			ObjectMeta:      metadata,
			Rules:           rbacRule,
			AggregationRule: nil,
		}, nil

	}
	return nil, nil
}
