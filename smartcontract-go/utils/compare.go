/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package utils

import (
	"strings"
)

func StringCompare(expected, actual string) int {
	return strings.Compare(expected, actual)
}

func IntCompare(expected, actual int) int {
	return expected - actual
}
