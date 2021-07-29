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
	"context"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
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

func generateRoleReference(kind, name string) rbacv1.RoleRef {
	const rbacAPIGroup = "rbac.authorization.k8s.io"
	return rbacv1.RoleRef{
		Kind:     kind,
		Name:     name,
		APIGroup: rbacAPIGroup,
	}
}

func generateSubject(kind, name, namespace string) rbacv1.Subject {
	return rbacv1.Subject{
		Kind:      kind,
		Name:      name,
		Namespace: namespace,
	}
}

func ApplyRbac(client k8sclient.Client) error {

	const serviceAccountKind = "ServiceAccount"
	const roleKind = "Role"
	const clusterRoleKind = "ClusterRole"

	roles := []rbacv1.Role{
		{
			TypeMeta:   getRoleTypeMeta(),
			ObjectMeta: generateObjectMeta("leader-election-clusterRoleBinding", "monitoring-system", map[string]string{}, map[string]string{}),
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
			TypeMeta:   getRoleBindingTypeMeta(),
			ObjectMeta: generateObjectMeta("leader-election-rolebinding", "monitoring-system", map[string]string{}, map[string]string{}),
			RoleRef:    generateRoleReference(roleKind, "leader-election-clusterRoleBinding"),
			Subjects: []rbacv1.Subject{
				generateSubject(serviceAccountKind, "vm-operator", "monitoring-system"),
				generateSubject(serviceAccountKind, "vm-operator2", "monitoring-system"),
			},
		},
	}

	clusterRoles := []rbacv1.ClusterRole{
		{
			TypeMeta:   getClusterRoleTypeMeta(),
			ObjectMeta: generateObjectMeta("vm-operator-psp-clusterRoleBinding", "", map[string]string{}, map[string]string{}),
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

	clusterRoleBindings := []rbacv1.ClusterRoleBinding{
		{
			TypeMeta:   getClusterRoleBindingTypeMeta(),
			ObjectMeta: generateObjectMeta("manager-rolebinding", "", map[string]string{}, map[string]string{}),
			RoleRef:    generateRoleReference(clusterRoleKind, "manager-clusterRoleBinding"),
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					APIGroup:  "",
					Name:      "vm-operator",
					Namespace: "monitoring-system",
				},
			},
		},
	}

	for _, role := range roles {
		if createErr := client.Create(context.TODO(), &role); createErr != nil {
			return createErr
		}
	}

	for _, clusterRole := range clusterRoles {
		if createErr := client.Create(context.TODO(), &clusterRole); createErr != nil {
			return createErr
		}
	}

	for _, roleBinding := range rolebindings {
		if createErr := client.Create(context.TODO(), &roleBinding); createErr != nil {
			return createErr
		}
	}

	for _, clusterRoleBinding := range clusterRoleBindings {
		if createErr := client.Create(context.TODO(), &clusterRoleBinding); createErr != nil {
			return createErr
		}
	}

	return nil
}
