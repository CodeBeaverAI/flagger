package v1beta1

import (
    "reflect"
    "testing"

    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/runtime/schema"
)

// TestResourceFunction tests that the Resource function returns a properly
// qualified GroupResource with the correct group and resource name.
func TestResourceFunction(t *testing.T) {
    resourceName := "testroute"
    gr := Resource(resourceName)

    // Verify that the group part is correctly set.
    expectedGroup := SchemeGroupVersion.Group
    if gr.Group != expectedGroup {
    t.Errorf("expected group %s, got %s", expectedGroup, gr.Group)
    }

    // Verify that the resource part is correctly set.
    if gr.Resource != resourceName {
    t.Errorf("expected resource %s, got %s", resourceName, gr.Resource)
    }
}

// TestAddKnownTypes verifies that AddToScheme properly registers the expected types in the scheme.
func TestAddKnownTypes(t *testing.T) {
    scheme := runtime.NewScheme()

    // AddToScheme should register the types without error.
    if err := AddToScheme(scheme); err != nil {
    t.Errorf("AddToScheme returned an unexpected error: %v", err)
    }

    // Retrieve the known types for our SchemeGroupVersion.
    knownTypes := scheme.KnownTypes(SchemeGroupVersion)

    // Expected types that should be registered.
    expectedTypes := []string{"HTTPRoute", "HTTPRouteList", "ReferenceGrant", "ReferenceGrantList"}

    // Check that each expected type is registered.
    for _, typeName := range expectedTypes {
    if _, exists := knownTypes[typeName]; !exists {
    t.Errorf("expected type %s to be registered", typeName)
    }
    }

    // Verify the group version added by metav1.
    gv := schema.GroupVersion{Group: SchemeGroupVersion.Group, Version: SchemeGroupVersion.Version}
    if !reflect.DeepEqual(gv, SchemeGroupVersion) {
    t.Errorf("group version mismatch: expected %+v, got %+v", SchemeGroupVersion, gv)
    }
}
// TestResourceFunctionEmpty tests that Resource() returns the expected GroupResource when given an empty resource string.
func TestResourceFunctionEmpty(t *testing.T) {
    resourceName := ""
    gr := Resource(resourceName)

    // Verify group part
    expectedGroup := SchemeGroupVersion.Group
    if gr.Group != expectedGroup {
        t.Errorf("expected group %s, got %s", expectedGroup, gr.Group)
    }

    // Verify resource part is empty.
    if gr.Resource != resourceName {
        t.Errorf("expected empty resource, got %s", gr.Resource)
    }
}

// TestResourceFunctionSpecialChars tests that Resource() returns the expected GroupResource
// when given a resource string containing special characters.
func TestResourceFunctionSpecialChars(t *testing.T) {
    resourceName := "route-test_special"
    gr := Resource(resourceName)

    // Verify that the group part is correct.
    expectedGroup := SchemeGroupVersion.Group
    if gr.Group != expectedGroup {
        t.Errorf("expected group %s, got %s", expectedGroup, gr.Group)
    }

    // Verify that the resource part is correctly set with special characters.
    if gr.Resource != resourceName {
        t.Errorf("expected resource %s, got %s", resourceName, gr.Resource)
    }
}

// TestAddKnownTypesNilScheme tests that addKnownTypes panics if passed a nil scheme.
func TestAddKnownTypesNilScheme(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Errorf("expected panic when passing nil scheme to addKnownTypes, but no panic occurred")
        }
    }()

    // Calling addKnownTypes with nil should cause a panic.
    addKnownTypes(nil)
}
// TestResourceFunctionWithSpaces verifies that Resource()
func TestResourceFunctionWithSpaces(t *testing.T) {
    resourceName := "test route with spaces"
    gr := Resource(resourceName)

    // Verify that the group part is correctly set.
    if gr.Group != SchemeGroupVersion.Group {
        t.Errorf("expected group %s, got %s", SchemeGroupVersion.Group, gr.Group)
    }
    // Verify that the resource part is exactly the provided string.
    if gr.Resource != resourceName {
        t.Errorf("expected resource %q, got %q", resourceName, gr.Resource)
    }
}

// TestAddToSchemeMultiple verifies that calling AddToScheme multiple times
// is idempotent and does not cause errors or duplicate registrations.
func TestAddToSchemeMultiple(t *testing.T) {
    scheme := runtime.NewScheme()
    // Call AddToScheme three times.
    for i := 0; i < 3; i++ {
        if err := AddToScheme(scheme); err != nil {
            t.Errorf("call %d: AddToScheme returned an unexpected error: %v", i, err)
        }
    }

    // Retrieve the known types for our SchemeGroupVersion.
    knownTypes := scheme.KnownTypes(SchemeGroupVersion)
    // Expected types that should be registered.
    expectedTypes := []string{"HTTPRoute", "HTTPRouteList", "ReferenceGrant", "ReferenceGrantList"}

    // Verify each type is registered.
    for _, typeName := range expectedTypes {
        if _, exists := knownTypes[typeName]; !exists {
            t.Errorf("expected type %s to be registered", typeName)
        }
    }
}

// TestSchemeGroupVersionString verifies that the String method on SchemeGroupVersion
// returns the expected group/version string.
func TestSchemeGroupVersionString(t *testing.T) {
    expected := SchemeGroupVersion.Group + "/" + SchemeGroupVersion.Version
    if SchemeGroupVersion.String() != expected {
        t.Errorf("expected SchemeGroupVersion string %q, got %q", expected, SchemeGroupVersion.String())
    }
}
// TestResourceFunctionUnicode tests that Resource() correctly handles Unicode resource names.
func TestResourceFunctionUnicode(t *testing.T) {
    resourceName := "ルート" // Japanese for "route"
    gr := Resource(resourceName)

    // Verify that the group part is correctly set.
    if gr.Group != SchemeGroupVersion.Group {
        t.Errorf("expected group %s, got %s", SchemeGroupVersion.Group, gr.Group)
    }
    // Verify that the resource part matches the Unicode input.
    if gr.Resource != resourceName {
        t.Errorf("expected resource %s, got %s", resourceName, gr.Resource)
    }
}

// TestResourceFunctionNumeric tests that Resource() correctly handles numeric resource names.
func TestResourceFunctionNumeric(t *testing.T) {
    resourceName := "1234"
    gr := Resource(resourceName)

    // Verify that the group part is correctly set.
    if gr.Group != SchemeGroupVersion.Group {
        t.Errorf("expected group %s, got %s", SchemeGroupVersion.Group, gr.Group)
    }
    // Verify that the resource part matches the numeric string.
    if gr.Resource != resourceName {
        t.Errorf("expected resource %s, got %s", resourceName, gr.Resource)
    }
}

// TestResourceFunctionNewlineTab tests that Resource() correctly handles resource names containing newline and tab characters.
func TestResourceFunctionNewlineTab(t *testing.T) {
    resourceName := "route\n\t test"
    gr := Resource(resourceName)

    // Verify that the group part is correctly set.
    if gr.Group != SchemeGroupVersion.Group {
        t.Errorf("expected group %s, got %s", SchemeGroupVersion.Group, gr.Group)
    }
    // Verify that the resource part retains the newline and tab characters.
    if gr.Resource != resourceName {
        t.Errorf("expected resource %q, got %q", resourceName, gr.Resource)
    }
}