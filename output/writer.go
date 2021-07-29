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

package output

import "bytes"

const (
	name = "test.go"
)

func WritePreamble() error {
	var buffer bytes.Buffer
	writePackageName("foo", &buffer)
	writeLine("", &buffer)
	writeImports(&buffer)
	return nil
}

func writePackageName(name string, buffer *bytes.Buffer) error {
	return writeLine("package "+name, buffer)
}

func writeImports(buffer *bytes.Buffer) error {
	writeLine("imports (", buffer)
	writeLine("\trbacv1 \"k8s.io/api/rbac/v1\"", buffer)
	writeLine(")", buffer)
	return nil
}

func write(content string, buffer *bytes.Buffer) error {
	_, writeErr := buffer.WriteString(content)
	return writeErr
}

func writeLine(content string, buffer *bytes.Buffer) error {
	_, writeErr := buffer.WriteString(content + "\n")
	return writeErr
}
