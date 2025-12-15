package client

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// IsUnrecoverableError checks if an error is unrecoverable and should not be retried.
// Returns true for errors like Forbidden, Unauthorized, Invalid, MethodNotSupported, and NotAcceptable.
func IsUnrecoverableError(err error) bool {
	if apierrors.IsForbidden(err) ||
		apierrors.IsUnauthorized(err) ||
		apierrors.IsInvalid(err) ||
		apierrors.IsMethodNotSupported(err) ||
		apierrors.IsNotAcceptable(err) {
		return true
	}

	return false
}
