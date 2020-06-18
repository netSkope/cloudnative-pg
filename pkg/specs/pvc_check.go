/*
This file is part of Cloud Native PostgreSQL.

Copyright (C) 2019-2020 2ndQuadrant Italia SRL. Exclusively licensed to 2ndQuadrant Limited.
*/

package specs

import (
	corev1 "k8s.io/api/core/v1"
)

// IsPvcAvailable Check if a PVC with a certain key exists and is available or not
func IsPvcAvailable(pvc corev1.PersistentVolumeClaim) bool {
	// If the PVC is not usable is now available
	_, unusable := pvc.Annotations[PvcUnusableAnnotation]
	return !unusable
}