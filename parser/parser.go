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
	"fmt"
	"os"

	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

type ParsingError interface {
	getKind() string
}

type ShortDocumentError struct {
}

func (s ShortDocumentError) Error() string {
	return fmt.Sprint("document was too short to be parsed")
}

func (s ShortDocumentError) getKind() string {
	return "short-document"
}

type DownStreamError struct {
}

func Parse() ([]v1.ClusterRole, []v1.ClusterRoleBinding, []v1.Role, []v1.RoleBinding, error) {
	var clusterRoleList []v1.ClusterRole
	var clusterRoleBindingList []v1.ClusterRoleBinding
	var roleList []v1.Role
	var roleBindingList []v1.RoleBinding

	content, readErr := os.ReadFile("release/operator/rbac.yaml")
	if readErr != nil {
		return nil, nil, nil, nil, readErr
	}
	documentList := bytes.Split(content, []byte("---\n"))
	for _, document := range documentList {
		if len(document) < 2 {
			continue // skip things without content to deal with spliting yaml documents within 1 file
		}
		raw := unstructured.Unstructured{}
		unmarshalErr := yaml.Unmarshal(document, &raw)
		if unmarshalErr != nil {
			return nil, nil, nil, nil, unmarshalErr
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
					return nil, nil, nil, nil, marshalErr
				}
				parsedRule := v1.PolicyRule{}
				parseErr := json.Unmarshal(ruleBytes, &parsedRule)
				if parseErr != nil {
					return nil, nil, nil, nil, parseErr
				}
				rbacRule = append(rbacRule, parsedRule)
			}
		}

	}
	return nil, nil, nil, nil, nil
}

func FirstRound(document []byte) (*metav1.TypeMeta, *metav1.ObjectMeta, *unstructured.Unstructured, ParsingError) {
	if len(document) < 2 {
		return nil, nil, nil, nil // skip things without content to deal with spliting yaml documents within 1 file
	}
	raw := unstructured.Unstructured{}
	unmarshalErr := yaml.Unmarshal(document, &raw)
	if unmarshalErr != nil {
		return nil, nil, nil,
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
}
