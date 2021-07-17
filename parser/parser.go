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
	"errors"
	"fmt"
	"os"

	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

type ParsingError interface {
	getKind() string
	Error() string
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
	err error
}

func (d DownStreamError) getKind() string {
	return "downstream-error"
}

func (d DownStreamError) Error() string {
	return d.err.Error()
}

func Parse() ([]v1.ClusterRole, []v1.ClusterRoleBinding, []v1.Role, []v1.RoleBinding, error) {
	var clusterRoleList []v1.ClusterRole
	var clusterRoleBindingList []v1.ClusterRoleBinding
	var roleList []v1.Role
	var roleBindingList []v1.RoleBinding
	const clusterRoleKind = "ClusterRole"
	const clusterRoleBindingKind = "ClusterRoleBinding"
	const roleKind = "Role"
	const roleBindingKind = "RoleBinding"

	content, readErr := os.ReadFile("test-resources/release/operator/rbac.yaml")
	if readErr != nil {
		return nil, nil, nil, nil, readErr
	}
	documentList := bytes.Split(content, []byte("---\n"))
	for _, document := range documentList {
		typeMeta, objectMeta, raw, firstRoundErr := FirstRound(document)
		if firstRoundErr != nil {
			return nil, nil, nil, nil, firstRoundErr
		}

		switch typeMeta.Kind {
		case roleKind:
			rules, ruleErr := ExtractRules(raw)
			if ruleErr != nil {
				return nil, nil, nil, nil, ruleErr
			}
			role := v1.Role{
				TypeMeta: *typeMeta,
				ObjectMeta: *objectMeta,
				Rules:      rules,
			}
			roleList = append(roleList, role)
		case clusterRoleKind:
			rules, ruleErr := ExtractRules(raw)
			if ruleErr != nil {
				return nil, nil, nil, nil, ruleErr
			}
			role := v1.ClusterRole{
				TypeMeta: *typeMeta,
				ObjectMeta: *objectMeta,
				Rules:      rules,
				AggregationRule: nil,
			}
			clusterRoleList = append(clusterRoleList, role)
		case roleBindingKind:
			roleRef, subjects, extractErr := ExtractRoleRefAndSubjects(raw)
			if extractErr != nil {
				return nil, nil, nil, nil, extractErr
			}
			roleBinding := v1.RoleBinding{
				TypeMeta: *typeMeta,
				ObjectMeta: *objectMeta,
				Subjects:   subjects,
				RoleRef: *roleRef,
			}
			roleBindingList = append(roleBindingList, roleBinding)

		case clusterRoleBindingKind:
			roleRef, subjects, extractErr := ExtractRoleRefAndSubjects(raw)
			if extractErr != nil {
				return nil, nil, nil, nil, extractErr
			}
			roleBinding := v1.ClusterRoleBinding{
				TypeMeta: *typeMeta,
				ObjectMeta: *objectMeta,
				Subjects:   subjects,
				RoleRef: *roleRef,
			}
			clusterRoleBindingList = append(clusterRoleBindingList, roleBinding)
		}
	}
	return clusterRoleList, clusterRoleBindingList, roleList, roleBindingList, nil
}

func FirstRound(document []byte) (*metav1.TypeMeta, *metav1.ObjectMeta, *unstructured.Unstructured, ParsingError) {
	if len(document) < 2 {
		return nil, nil, nil, nil // skip things without content to deal with spliting yaml documents within 1 file
	}
	raw := &unstructured.Unstructured{}
	unmarshalErr := yaml.Unmarshal(document, raw)
	if unmarshalErr != nil {
		return nil, nil, nil, DownStreamError{err: unmarshalErr}
	}
	metadata := &metav1.ObjectMeta{
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

	typeMeta := &metav1.TypeMeta{
		Kind:       raw.GetKind(),
		APIVersion: raw.GetAPIVersion(),
	}
	return typeMeta, metadata, raw, nil
}

func ExtractRules(raw *unstructured.Unstructured) ([]v1.PolicyRule, ParsingError) {
	var rbacRule []v1.PolicyRule
	ruleInterfaceList := raw.Object["rules"].([]interface{})
	ruleBytes, marshalErr := json.Marshal(ruleInterfaceList)
	if marshalErr != nil {
		return nil, DownStreamError{err: marshalErr}
	}
	parseErr := json.Unmarshal(ruleBytes, &rbacRule)
	if parseErr != nil {
		return nil, DownStreamError{err: parseErr}
	}
	return rbacRule, nil
}

func ExtractRoleRefAndSubjects(raw *unstructured.Unstructured) (*v1.RoleRef, []v1.Subject, error) {
	roleRef := &v1.RoleRef{}
	var subjects []v1.Subject
	rawRoleRef, exists := raw.Object["roleRef"]
	if !exists {
		return nil, nil, errors.New("required key is missing")
	}
	testRoleRef := rawRoleRef.(map[string]interface{})
	refBytes, marshalErr := json.Marshal(testRoleRef)
	if marshalErr != nil {
		return nil, nil, marshalErr
	}
	parsingErr := json.Unmarshal(refBytes, roleRef)
	if parsingErr != nil {
		return nil, nil, parsingErr
	}
	rawSubjects := raw.Object["subjects"].([]interface{})
	subjectBytes, subjectMarshalErr := json.Marshal(rawSubjects)
	if subjectMarshalErr != nil {
		return nil, nil, subjectMarshalErr
	}
	subjectParseErr := json.Unmarshal(subjectBytes, &subjects)
	if subjectParseErr != nil {
		return nil, nil, subjectParseErr
	}
	return roleRef, subjects, nil

}

