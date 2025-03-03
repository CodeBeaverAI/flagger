package v1beta1

import (
    "testing"

    "k8s.io/apimachinery/pkg/runtime"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

    "github.com/fluxcd/flagger/pkg/apis/gatewayapi"
)

// TestResource verifies that Resource returns the expected GroupResource.
func TestResource(t *testing.T) {
    resourceName := "dummy"
    grpRes := Resource(resourceName)
    if grpRes.Group != SchemeGroupVersion.Group {
    t.Errorf("expected group %s, got %s", SchemeGroupVersion.Group, grpRes.Group)
    }
    if grpRes.Resource != resourceName {
    t.Errorf("expected resource %s, got %s", resourceName, grpRes.Resource)
    }
}

// TestAddKnownTypes checks that the known types are properly added to the scheme.
func TestAddKnownTypes(t *testing.T) {
    scheme := runtime.NewScheme()
    err := addKnownTypes(scheme)
    if err != nil {
    t.Fatalf("unexpected error from addKnownTypes: %v", err)
    }

    types := scheme.KnownTypes(SchemeGroupVersion)
    expectedTypes := []string{"HTTPRoute", "HTTPRouteList", "ReferenceGrant", "ReferenceGrantList"}
    for _, et := range expectedTypes {
    if _, found := types[et]; !found {
    t.Errorf("expected type %s to be registered", et)
    }
    }

    // Verify that the group version is reflected in the scheme.
    gvMeta := metav1.GroupVersion{Group: SchemeGroupVersion.Group, Version: SchemeGroupVersion.Version}
    if gvMeta.String() != SchemeGroupVersion.String() {
        t.Errorf("expected GroupVersion string %s, got %s", SchemeGroupVersion.String(), gvMeta.String())
    }
}

// TestAddToScheme runs AddToScheme to ensure that all known types are added.
func TestAddToScheme(t *testing.T) {
    scheme := runtime.NewScheme()
    err := AddToScheme(scheme)
    if err != nil {
    t.Fatalf("unexpected error from AddToScheme: %v", err)
    }

    types := scheme.KnownTypes(SchemeGroupVersion)
    expectedTypes := []string{"HTTPRoute", "HTTPRouteList", "ReferenceGrant", "ReferenceGrantList"}
    for _, et := range expectedTypes {
    if _, found := types[et]; !found {
    t.Errorf("expected type %s to be registered", et)
    }
    }
}

// TestSchemeGroupVersion ensures that SchemeGroupVersion is set as expected.
func TestSchemeGroupVersion(t *testing.T) {
    expectedGroup := gatewayapi.GroupName
    if SchemeGroupVersion.Group != expectedGroup {
    t.Errorf("expected group %s, got %s", expectedGroup, SchemeGroupVersion.Group)
    }
    if SchemeGroupVersion.Version != "v1beta1" {
    t.Errorf("expected version v1beta1, got %s", SchemeGroupVersion.Version)
    }
}
// TestResourceEmpty verifies that Resource correctly handles an empty resource string.
func TestResourceEmpty(t *testing.T) {
    resourceName := ""
    grpRes := Resource(resourceName)
    if grpRes.Group != SchemeGroupVersion.Group {
        t.Errorf("expected group %s, got %s", SchemeGroupVersion.Group, grpRes.Group)
    }
    if grpRes.Resource != resourceName {
        t.Errorf("expected empty resource, got %q", grpRes.Resource)
    }
}

// TestAddKnownTypesIdempotency ensures that calling addKnownTypes repeatedly does not cause errors
// and that the known types are registered consistently.
func TestAddKnownTypesIdempotency(t *testing.T) {
    scheme := runtime.NewScheme()
    // Call addKnownTypes twice to test idempotency.
    for i := 0; i < 2; i++ {
        if err := addKnownTypes(scheme); err != nil {
            t.Fatalf("unexpected error calling addKnownTypes at iteration %d: %v", i, err)
        }
    }

    types := scheme.KnownTypes(SchemeGroupVersion)
    expectedTypes := []string{"HTTPRoute", "HTTPRouteList", "ReferenceGrant", "ReferenceGrantList"}
    for _, et := range expectedTypes {
        if _, found := types[et]; !found {
            t.Errorf("expected type %s to be registered after idempotent calls", et)
        }
    }
}

// TestAddKnownTypesNilScheme verifies that calling addKnownTypes with a nil scheme causes a panic.
func TestAddKnownTypesNilScheme(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Errorf("expected panic when calling addKnownTypes with nil scheme, but no panic occurred")
        }
    }()
    // This call is expected to panic since the scheme is nil.
    _ = addKnownTypes(nil)
    t.Errorf("should not reach here")
}