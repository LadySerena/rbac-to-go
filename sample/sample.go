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

package sample

import (
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getRoleTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       "Role",
		APIVersion: rbacv1.SchemeGroupVersion.String(),
	}
}

func getClusterRoleTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       "ClusterRole",
		APIVersion: rbacv1.SchemeGroupVersion.String(),
	}
}

func getRoleBindingTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       "RoleBinding",
		APIVersion: rbacv1.SchemeGroupVersion.String(),
	}
}

func getClusterRoleBindingTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       "ClusterRoleBinding",
		APIVersion: rbacv1.SchemeGroupVersion.String(),
	}
}

func generateObjectMeta(name string, namespace string, labels map[string]string, annotations map[string]string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:        name,
		Namespace:   namespace,
		Labels:      labels,
		Annotations: annotations,
	}
}

func ApplyRbac(client client.Client) error {
	roles := []rbacv1.Role{
		{
			TypeMeta:   getRoleTypeMeta(),
			ObjectMeta: generateObjectMeta("leader-election-role", "monitoring-system", map[string]string{}, map[string]string{}),
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{
						"",
					},
					Resources: []string{
						"configmaps",
					},
					Verbs: []string{
						"get",
						"list",
						"watch",
						"create",
						"update",
						"patch",
						"delete",
					},
				},
			},
		},
	}
	rolebindings := []rbacv1.RoleBinding{
		{
			TypeMeta: getRoleBindingTypeMeta(),
			ObjectMeta: metav1.ObjectMeta{
				Name:        "",
				Namespace:   "",
				Labels:      nil,
				Annotations: nil,
			},
		},
	}

	clusterRoles := []rbacv1.ClusterRole{
		{
			TypeMeta: getClusterRoleTypeMeta(),
			ObjectMeta: metav1.ObjectMeta{
				Name:                       "",
				GenerateName:               "",
				Namespace:                  "",
				SelfLink:                   "",
				UID:                        "",
				ResourceVersion:            "",
				Generation:                 0,
				CreationTimestamp:          metav1.Time{},
				DeletionTimestamp:          nil,
				DeletionGracePeriodSeconds: nil,
				Labels:                     nil,
				Annotations:                nil,
				OwnerReferences:            nil,
				Finalizers:                 nil,
				ClusterName:                "",
				ManagedFields:              nil,
			},
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{
						"",
					},
					Resources: []string{
						"configmaps",
					},
					Verbs: []string{
						"get",
						"list",
						"watch",
						"create",
						"update",
						"patch",
						"delete",
					},
				},
			},
		},
	}
	fmt.Printf("%v", roles)
	fmt.Printf("%v", clusterRoles)
	fmt.Printf("%v", rolebindings)
	return nil
}
