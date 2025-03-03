package v1beta1

import (
    gwv1beta1 "github.com/fluxcd/flagger/pkg/apis/gatewayapi/v1beta1"
    v1 "github.com/fluxcd/flagger/pkg/apis/gatewayapi/v1"
    "testing"
    "time"
    "encoding/json"
)

// TestHTTPRewrite_GetType tests the GetType method of HTTPRewrite.
func TestHTTPRewrite_GetType(t *testing.T) {
    // when the type equals PrefixMatchHTTPPathModifier, it is returned unchanged.
    r1 := &HTTPRewrite{Type: string(gwv1beta1.PrefixMatchHTTPPathModifier)}
    if got := r1.GetType(); got != string(gwv1beta1.PrefixMatchHTTPPathModifier) {
        t.Errorf("Expected %s but got %s", gwv1beta1.PrefixMatchHTTPPathModifier, got)
    }

    // when any other type is set, it returns FullPathHTTPPathModifier.
    r2 := &HTTPRewrite{Type: "other"}
    if got := r2.GetType(); got != string(gwv1beta1.FullPathHTTPPathModifier) {
        t.Errorf("Expected %s but got %s", gwv1beta1.FullPathHTTPPathModifier, got)
    }
}

// TestCanaryService_GetIstioRewrite tests the GetIstioRewrite method of CanaryService.
func TestCanaryService_GetIstioRewrite(t *testing.T) {
    // Test with nil Rewrite, we expect nil.
    service := &CanaryService{}
    if rw := service.GetIstioRewrite(); rw != nil {
    t.Errorf("Expected nil rewrite, got %+v", rw)
    }

    // Test with a non-nil Rewrite.
    service.Rewrite = &HTTPRewrite{
    Uri:       "/test",
    Authority: "example.com",
    }
    rw := service.GetIstioRewrite()
    if rw == nil {
    t.Errorf("Expected non nil rewrite, got nil")
    }
    if rw.Uri != "/test" || rw.Authority != "example.com" {
    t.Errorf("Unexpected rewrite values: %+v", rw)
    }
}

// TestSessionAffinity_GetMaxAge tests the GetMaxAge method of SessionAffinity.
func TestSessionAffinity_GetMaxAge(t *testing.T) {
    sa := &SessionAffinity{}
    // when MaxAge is zero, it should return the default value of 86400 seconds.
    if age := sa.GetMaxAge(); age != 86400 {
    t.Errorf("Expected default max age 86400, got %d", age)
    }
    sa.MaxAge = 500
    if age := sa.GetMaxAge(); age != 500 {
    t.Errorf("Expected max age 500, got %d", age)
    }
}

// TestCanary_GetServiceNames tests the GetServiceNames method of Canary.
func TestCanary_GetServiceNames(t *testing.T) {
    c := &Canary{
    Spec: CanarySpec{
    TargetRef: LocalObjectReference{Name: "myservice"},
    Service: CanaryService{
    Name: "",
    },
    },
    }
    apex, primary, canaryName := c.GetServiceNames()
    if apex != "myservice" {
    t.Errorf("Expected apex service name 'myservice', got %s", apex)
    }
    if primary != "myservice-primary" {
    t.Errorf("Expected primary service name 'myservice-primary', got %s", primary)
    }
    if canaryName != "myservice-canary" {
    t.Errorf("Expected canary service name 'myservice-canary', got %s", canaryName)
    }

    // When Service.Name is provided, it should override the target name.
    c.Spec.Service.Name = "custom"
    apex, primary, canaryName = c.GetServiceNames()
    if apex != "custom" {
    t.Errorf("Expected apex service name 'custom', got %s", apex)
    }
    if primary != "custom-primary" || canaryName != "custom-canary" {
    t.Errorf("Unexpected primary or canary names: %s, %s", primary, canaryName)
    }
}

// TestCanary_GetProgressDeadlineSeconds tests the GetProgressDeadlineSeconds method of Canary.
func TestCanary_GetProgressDeadlineSeconds(t *testing.T) {
    c := &Canary{
    Spec: CanarySpec{},
    }
    // When ProgressDeadlineSeconds is nil, should return the default (600 seconds).
    if secs := c.GetProgressDeadlineSeconds(); secs != 600 {
    t.Errorf("Expected progress deadline 600, got %d", secs)
    }
    val := int32(900)
    c.Spec.ProgressDeadlineSeconds = &val
    if secs := c.GetProgressDeadlineSeconds(); secs != 900 {
    t.Errorf("Expected progress deadline 900, got %d", secs)
    }
}

// TestCanary_GetAnalysis tests the GetAnalysis method of Canary.
func TestCanary_GetAnalysis(t *testing.T) {
    c := &Canary{
    Spec: CanarySpec{
    Analysis: &CanaryAnalysis{
    Interval: "10s",
    },
    },
    }
    a := c.GetAnalysis()
    if a.Interval != "10s" {
    t.Errorf("Expected analysis interval '10s', got %s", a.Interval)
    }

    // Test fallback to deprecated CanaryAnalysis when Analysis is nil.
    c = &Canary{
    Spec: CanarySpec{
    Analysis:       nil,
    CanaryAnalysis: &CanaryAnalysis{Interval: "20s"},
    },
    }
    a = c.GetAnalysis()
    if a.Interval != "20s" {
    t.Errorf("Expected analysis interval '20s', got %s", a.Interval)
    }
}

// TestCanary_GetAnalysisInterval tests the GetAnalysisInterval method.
func TestCanary_GetAnalysisInterval(t *testing.T) {
    // With an empty interval, should use the default AnalysisInterval (60s).
    c := &Canary{
    Spec: CanarySpec{
    Analysis: &CanaryAnalysis{Interval: ""},
    },
    }
    if d := c.GetAnalysisInterval(); d != 60*time.Second {
    t.Errorf("Expected default analysis interval 60s, got %v", d)
    }

    // With an invalid duration, should fallback to the default.
    c.Spec.Analysis.Interval = "invalid"
    if d := c.GetAnalysisInterval(); d != 60*time.Second {
    t.Errorf("Expected default analysis interval 60s on invalid input, got %v", d)
    }

    // With a valid duration less than 10 seconds, should use a minimum of 10s.
    c.Spec.Analysis.Interval = "5s"
    if d := c.GetAnalysisInterval(); d != 10*time.Second {
    t.Errorf("Expected minimum analysis interval 10s, got %v", d)
    }

    // With a valid duration of 15 seconds.
    c.Spec.Analysis.Interval = "15s"
    if d := c.GetAnalysisInterval(); d != 15*time.Second {
    t.Errorf("Expected analysis interval 15s, got %v", d)
    }
}

// TestCanary_GetAnalysisThreshold tests the GetAnalysisThreshold method.
func TestCanary_GetAnalysisThreshold(t *testing.T) {
    c := &Canary{
    Spec: CanarySpec{
    Analysis: &CanaryAnalysis{Threshold: 5},
    },
    }
    if thresh := c.GetAnalysisThreshold(); thresh != 5 {
    t.Errorf("Expected threshold 5, got %d", thresh)
    }

    // If the threshold is zero, it should return default value of 1.
    c.Spec.Analysis.Threshold = 0
    if thresh := c.GetAnalysisThreshold(); thresh != 1 {
    t.Errorf("Expected threshold 1 when value is zero, got %d", thresh)
    }
}

// TestCanary_GetAnalysisPrimaryReadyThreshold tests the GetAnalysisPrimaryReadyThreshold method.
func TestCanary_GetAnalysisPrimaryReadyThreshold(t *testing.T) {
    ready := 80
    c := &Canary{
    Spec: CanarySpec{
    Analysis: &CanaryAnalysis{PrimaryReadyThreshold: &ready},
    },
    }
    if pct := c.GetAnalysisPrimaryReadyThreshold(); pct != 80 {
    t.Errorf("Expected primary ready threshold 80, got %d", pct)
    }

    // If nil, it should return the default (100).
    c.Spec.Analysis.PrimaryReadyThreshold = nil
    if pct := c.GetAnalysisPrimaryReadyThreshold(); pct != 100 {
    t.Errorf("Expected primary ready threshold 100 when nil, got %d", pct)
    }
}

// TestCanary_GetAnalysisCanaryReadyThreshold tests the GetAnalysisCanaryReadyThreshold method.
func TestCanary_GetAnalysisCanaryReadyThreshold(t *testing.T) {
    ready := 70
    c := &Canary{
    Spec: CanarySpec{
    Analysis: &CanaryAnalysis{CanaryReadyThreshold: &ready},
    },
    }
    if pct := c.GetAnalysisCanaryReadyThreshold(); pct != 70 {
    t.Errorf("Expected canary ready threshold 70, got %d", pct)
    }

    // If nil, it should return the default (100).
    c.Spec.Analysis.CanaryReadyThreshold = nil
    if pct := c.GetAnalysisCanaryReadyThreshold(); pct != 100 {
    t.Errorf("Expected canary ready threshold 100 when nil, got %d", pct)
    }
}

// TestCanary_GetMetricInterval tests the GetMetricInterval method.
func TestCanary_GetMetricInterval(t *testing.T) {
    c := &Canary{}
    if mi := c.GetMetricInterval(); mi != "1m" {
    t.Errorf("Expected metric interval '1m', got %s", mi)
    }
}

// TestCanary_SkipAnalysis tests the SkipAnalysis method of Canary.
func TestCanary_SkipAnalysis(t *testing.T) {
    c := &Canary{
    Spec: CanarySpec{
    Analysis:       nil,
    CanaryAnalysis: nil,
    SkipAnalysis:   false,
    },
    }
    // When both Analysis and CanaryAnalysis are nil, SkipAnalysis should return true.
    if !c.SkipAnalysis() {
    t.Errorf("Expected SkipAnalysis to return true when analysis fields are nil")
    }

    // When Analysis is provided, even with SkipAnalysis false, it should not skip.
    c.Spec.Analysis = &CanaryAnalysis{}
    if c.SkipAnalysis() {
    t.Errorf("Expected SkipAnalysis to return false when analysis is provided")
    }

    // When deprecated CanaryAnalysis is provided.
    c.Spec.Analysis = nil
    c.Spec.CanaryAnalysis = &CanaryAnalysis{}
    if c.SkipAnalysis() {
    t.Errorf("Expected SkipAnalysis to return false when deprecated analysis is provided")
    }

    // When SkipAnalysis flag is explicitly set to true.
    c.Spec.SkipAnalysis = true
    if !c.SkipAnalysis() {
    t.Errorf("Expected SkipAnalysis to return true when skipAnalysis flag is set")
    }
}

// TestHTTPRewrite_GetType_Default verifies that GetType returns the default FullPath modifier when Type is empty.
func TestHTTPRewrite_GetType_Default(t *testing.T) {
    r := &HTTPRewrite{}
    if got := r.GetType(); got != string(gwv1beta1.FullPathHTTPPathModifier) {
        t.Errorf("Expected %s but got %s", gwv1beta1.FullPathHTTPPathModifier, got)
    }
}

// TestCanary_GetAnalysisNil verifies that GetAnalysis returns nil when both Analysis and CanaryAnalysis are nil.
func TestCanary_GetAnalysisNil(t *testing.T) {
    c := &Canary{
        Spec: CanarySpec{
            Analysis:       nil,
            CanaryAnalysis: nil,
        },
    }
    if analysis := c.GetAnalysis(); analysis != nil {
        t.Errorf("Expected nil analysis, got %+v", analysis)
    }
}

// TestCanary_GetAnalysisThreshold_Negative verifies that a negative threshold returns the default value of 1.
func TestCanary_GetAnalysisThreshold_Negative(t *testing.T) {
    c := &Canary{
        Spec: CanarySpec{
            Analysis: &CanaryAnalysis{Threshold: -5},
        },
    }
    if thresh := c.GetAnalysisThreshold(); thresh != 1 {
        t.Errorf("Expected threshold 1 for negative value, got %d", thresh)
    }
}

// TestCanary_GetAnalysisInterval_Boundary verifies that an interval of exactly 10s is accepted.
func TestCanary_GetAnalysisInterval_Boundary(t *testing.T) {
    c := &Canary{
        Spec: CanarySpec{
            Analysis: &CanaryAnalysis{Interval: "10s"},
        },
    }
    if d := c.GetAnalysisInterval(); d != 10*time.Second {
        t.Errorf("Expected analysis interval 10s, got %v", d)
    }
}
// End of test file.
// TestCanary_DoubleAnalysis verifies that when both Analysis and CanaryAnalysis are provided,
func TestCanary_DoubleAnalysis(t *testing.T) {
    analysis := &CanaryAnalysis{Interval: "30s", Threshold: 3}
    deprecated := &CanaryAnalysis{Interval: "40s", Threshold: 5}
    c := &Canary{
        Spec: CanarySpec{
            Analysis:       analysis,
            CanaryAnalysis: deprecated,
        },
    }
    a := c.GetAnalysis()
    if a.Interval != "30s" || a.Threshold != 3 {
        t.Errorf("Expected Analysis to take precedence, got interval %s, threshold %d", a.Interval, a.Threshold)
    }
}

// TestCanary_GetAnalysisInterval_Nil verifies that GetAnalysisInterval panics when both Analysis and CanaryAnalysis are nil.
func TestCanary_GetAnalysisInterval_Nil(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Errorf("Expected panic when calling GetAnalysisInterval with nil analysis, but did not panic")
        }
    }()
    c := &Canary{
        Spec: CanarySpec{
            Analysis:       nil,
            CanaryAnalysis: nil,
        },
    }
    _ = c.GetAnalysisInterval() // should panic
}

// TestCanary_GetAnalysisThreshold_Nil verifies that GetAnalysisThreshold panics when both analysis fields are nil.
func TestCanary_GetAnalysisThreshold_Nil(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Errorf("Expected panic when calling GetAnalysisThreshold with nil analysis, but did not panic")
        }
    }()
    c := &Canary{
        Spec: CanarySpec{
            Analysis:       nil,
            CanaryAnalysis: nil,
        },
    }
    _ = c.GetAnalysisThreshold() // should panic
}

// TestCanary_GetAnalysisPrimaryReadyThreshold_Nil verifies that GetAnalysisPrimaryReadyThreshold panics when both analysis fields are nil.
func TestCanary_GetAnalysisPrimaryReadyThreshold_Nil(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Errorf("Expected panic when calling GetAnalysisPrimaryReadyThreshold with nil analysis, but did not panic")
        }
    }()
    c := &Canary{
        Spec: CanarySpec{
            Analysis:       nil,
            CanaryAnalysis: nil,
        },
    }
    _ = c.GetAnalysisPrimaryReadyThreshold() // should panic
}

// TestCanary_GetAnalysisCanaryReadyThreshold_Nil verifies that GetAnalysisCanaryReadyThreshold panics when both analysis fields are nil.
func TestCanary_GetAnalysisCanaryReadyThreshold_Nil(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Errorf("Expected panic when calling GetAnalysisCanaryReadyThreshold with nil analysis, but did not panic")
        }
    }()
    c := &Canary{
        Spec: CanarySpec{
            Analysis:       nil,
            CanaryAnalysis: nil,
        },
    }
    _ = c.GetAnalysisCanaryReadyThreshold() // should panic
}
// TestCanaryWebhookPayload_JSON tests JSON marshalling and unmarshalling of CanaryWebhookPayload.
func TestCanaryWebhookPayload_JSON(t *testing.T) {
    payload := &CanaryWebhookPayload{
        Name:      "test-canary",
        Namespace: "default",
        Phase:     "Progressing",
        Checksum:  "abc123",
        Metadata:  map[string]string{"key": "value"},
    }
    data, err := json.Marshal(payload)
    if err != nil {
        t.Fatalf("unexpected error during marshalling: %v", err)
    }
    var newPayload CanaryWebhookPayload
    if err := json.Unmarshal(data, &newPayload); err != nil {
        t.Fatalf("unexpected error during unmarshalling: %v", err)
    }
    if newPayload.Name != payload.Name || newPayload.Namespace != payload.Namespace ||
        newPayload.Phase != payload.Phase || newPayload.Checksum != payload.Checksum ||
        newPayload.Metadata["key"] != "value" {
        t.Errorf("unmarshalled payload does not match original")
    }
}

// TestCrossNamespaceObjectReference_JSON tests JSON marshalling of CrossNamespaceObjectReference.
func TestCrossNamespaceObjectReference_JSON(t *testing.T) {
    ref := CrossNamespaceObjectReference{
        APIVersion: "v1",
        Kind:       "Service",
        Name:       "my-service",
        Namespace:  "default",
    }
    data, err := json.Marshal(ref)
    if err != nil {
        t.Fatalf("unexpected error during marshalling: %v", err)
    }
    var newRef CrossNamespaceObjectReference
    if err := json.Unmarshal(data, &newRef); err != nil {
        t.Fatalf("unexpected error during unmarshalling: %v", err)
    }
    if newRef != ref {
        t.Errorf("expected %+v but got %+v", ref, newRef)
    }
}

// TestAutoscalerRefernce_JSON tests JSON marshalling and unmarshalling of AutoscalerRefernce.
func TestAutoscalerRefernce_JSON(t *testing.T) {
    replicas := ScalerReplicas{
        MinReplicas: func(i int32) *int32 { return &i }(1),
        MaxReplicas: func(i int32) *int32 { return &i }(5),
    }
    aRef := AutoscalerRefernce{
        APIVersion:           "v1",
        Kind:                 "HorizontalPodAutoscaler",
        Name:                 "autoscaler",
        PrimaryScalerQueries: map[string]string{"query": "value"},
        PrimaryScalerReplicas: &replicas,
    }
    data, err := json.Marshal(aRef)
    if err != nil {
        t.Fatalf("unexpected error during marshalling: %v", err)
    }
    var newARef AutoscalerRefernce
    if err := json.Unmarshal(data, &newARef); err != nil {
        t.Fatalf("unexpected error during unmarshalling: %v", err)
    }
    if newARef.Name != aRef.Name || newARef.APIVersion != aRef.APIVersion || newARef.Kind != aRef.Kind {
        t.Errorf("unmarshalled autoscaler reference does not match original")
    }
}

// TestCustomBackend_JSON tests JSON marshalling and unmarshalling of CustomBackend.
func TestCustomBackend_JSON(t *testing.T) {
    backend := CustomBackend{
        BackendObjectReference: &v1.BackendObjectReference{
            Name: "backend-service",
        },
        Filters: []v1.HTTPRouteFilter{
            {Type: "RequestHeaderModifier"},
        },
    }
    data, err := json.Marshal(backend)
    if err != nil {
        t.Fatalf("unexpected error during marshalling: %v", err)
    }
    var newBackend CustomBackend
    if err := json.Unmarshal(data, &newBackend); err != nil {
        t.Fatalf("unexpected error during unmarshalling: %v", err)
    }
    if newBackend.BackendObjectReference == nil || newBackend.BackendObjectReference.Name != "backend-service" {
        t.Errorf("unmarshalled custom backend does not match original")
    }
    if len(newBackend.Filters) != 1 || newBackend.Filters[0].Type != "RequestHeaderModifier" {
        t.Errorf("unmarshalled filters do not match original")
    }
}