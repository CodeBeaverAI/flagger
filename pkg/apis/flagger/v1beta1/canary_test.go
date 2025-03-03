package v1beta1

import (
    "testing"
    "time"
    "encoding/json"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
    gatewayv1beta1 "github.com/fluxcd/flagger/pkg/apis/gatewayapi/v1beta1"
)

// TestHTTPRewriteGetType tests the GetType method of HTTPRewrite.
func TestHTTPRewriteGetType(t *testing.T) {
    // Test when r.Type is set to PrefixMatchHTTPPathModifier.
    rewrite := &HTTPRewrite{
    Type: string(gatewayv1beta1.PrefixMatchHTTPPathModifier),
    }
    if got := rewrite.GetType(); got != string(gatewayv1beta1.PrefixMatchHTTPPathModifier) {
    t.Errorf("HTTPRewrite.GetType() = %v, want %v", got, string(gatewayv1beta1.PrefixMatchHTTPPathModifier))
    }

    // Test when r.Type is set to a non-prefix value. Expect default to FullPathHTTPPathModifier.
    rewrite.Type = "non-prefix"
    if got := rewrite.GetType(); got != string(gatewayv1beta1.FullPathHTTPPathModifier) {
    t.Errorf("HTTPRewrite.GetType() = %v, want %v", got, string(gatewayv1beta1.FullPathHTTPPathModifier))
    }
}

// TestCanaryServiceGetIstioRewrite tests the CanaryService.GetIstioRewrite method.
func TestCanaryServiceGetIstioRewrite(t *testing.T) {
    // When Rewrite is nil, expect nil result.
    service := CanaryService{}
    if res := service.GetIstioRewrite(); res != nil {
    t.Error("Expected nil IstioRewrite when Rewrite is nil")
    }

    // When Rewrite is provided, expect matching authority and uri values.
    service.Rewrite = &HTTPRewrite{
    Authority: "example.com",
    Uri:       "/new",
    }
    istioRewrite := service.GetIstioRewrite()
    if istioRewrite == nil {
    t.Fatal("Expected non-nil IstioRewrite")
    }
    if istioRewrite.Authority != "example.com" || istioRewrite.Uri != "/new" {
    t.Errorf("Got IstioRewrite = %+v, want Authority: example.com, Uri: /new", istioRewrite)
    }
}

// TestSessionAffinityGetMaxAge tests the GetMaxAge method of SessionAffinity.
func TestSessionAffinityGetMaxAge(t *testing.T) {
    // When MaxAge is 0, default should be 86400 (24 hours).
    affinity := &SessionAffinity{}
    if got := affinity.GetMaxAge(); got != 86400 {
    t.Errorf("SessionAffinity.GetMaxAge() = %d, want %d", got, 86400)
    }

    // When MaxAge is set, it should return that value.
    affinity.MaxAge = 3600
    if got := affinity.GetMaxAge(); got != 3600 {
    t.Errorf("SessionAffinity.GetMaxAge() = %d, want %d", got, 3600)
    }
}

// TestGetServiceNames tests Canary.GetServiceNames method.
func TestGetServiceNames(t *testing.T) {
    canary := &Canary{
    Spec: CanarySpec{
    TargetRef: LocalObjectReference{
    Name: "foo",
    },
    },
    }
    // When Service.Name is not set, apex should be TargetRef.Name.
    apex, primary, canaryName := canary.GetServiceNames()
    if apex != "foo" {
    t.Errorf("Expected apex name 'foo', got '%s'", apex)
    }
    if primary != "foo-primary" {
    t.Errorf("Expected primary name 'foo-primary', got '%s'", primary)
    }
    if canaryName != "foo-canary" {
    t.Errorf("Expected canary name 'foo-canary', got '%s'", canaryName)
    }

    // When Service.Name is provided, it should override the apex name.
    canary.Spec.Service.Name = "bar"
    apex, primary, canaryName = canary.GetServiceNames()
    if apex != "bar" {
    t.Errorf("Expected apex name 'bar', got '%s'", apex)
    }
    if primary != "bar-primary" {
    t.Errorf("Expected primary name 'bar-primary', got '%s'", primary)
    }
    if canaryName != "bar-canary" {
    t.Errorf("Expected canary name 'bar-canary', got '%s'", canaryName)
    }
}

// TestGetProgressDeadlineSeconds tests Canary.GetProgressDeadlineSeconds.
func TestGetProgressDeadlineSeconds(t *testing.T) {
    // When ProgressDeadlineSeconds is set, it should be returned.
    seconds := int32(300)
    canary := &Canary{
    Spec: CanarySpec{
    ProgressDeadlineSeconds: &seconds,
    },
    }
    if got := canary.GetProgressDeadlineSeconds(); got != 300 {
    t.Errorf("GetProgressDeadlineSeconds() = %d, want %d", got, 300)
    }

    // When not set, it should return the default ProgressDeadlineSeconds.
    canary.Spec.ProgressDeadlineSeconds = nil
    if got := canary.GetProgressDeadlineSeconds(); got != ProgressDeadlineSeconds {
    t.Errorf("GetProgressDeadlineSeconds() = %d, want %d", got, ProgressDeadlineSeconds)
    }
}

// TestGetAnalysis tests the GetAnalysis method of Canary.
func TestGetAnalysis(t *testing.T) {
    // When Analysis is not nil, it should return Analysis.
    analysis := &CanaryAnalysis{
    Interval: "2m",
    }
    canary := &Canary{
    Spec: CanarySpec{
    Analysis: analysis,
    },
    }
    if got := canary.GetAnalysis(); got != analysis {
    t.Errorf("GetAnalysis() = %v, want %v", got, analysis)
    }

    // When Analysis is nil but deprecated CanaryAnalysis is provided.
    deprecated := &CanaryAnalysis{
    Interval: "3m",
    }
    canary.Spec.Analysis = nil
    canary.Spec.CanaryAnalysis = deprecated
    if got := canary.GetAnalysis(); got != deprecated {
    t.Errorf("GetAnalysis() = %v, want %v", got, deprecated)
    }
}

// TestGetAnalysisInterval tests the GetAnalysisInterval method.
func TestGetAnalysisInterval(t *testing.T) {
    // Case 1: Empty interval should return default AnalysisInterval (60s)
    analysis := &CanaryAnalysis{
    Interval: "",
    }
    canary := &Canary{
    Spec: CanarySpec{
    Analysis: analysis,
    },
    }
    if d := canary.GetAnalysisInterval(); d != AnalysisInterval {
    t.Errorf("GetAnalysisInterval() = %v, want %v", d, AnalysisInterval)
    }

    // Case 2: An invalid duration string should return default AnalysisInterval.
    analysis.Interval = "invalid"
    if d := canary.GetAnalysisInterval(); d != AnalysisInterval {
    t.Errorf("GetAnalysisInterval() with invalid value = %v, want %v", d, AnalysisInterval)
    }

    // Case 3: A duration less than 10 seconds should be clamped to 10 seconds.
    analysis.Interval = "5s"
    if d := canary.GetAnalysisInterval(); d != 10*time.Second {
    t.Errorf("GetAnalysisInterval() with low value = %v, want %v", d, 10*time.Second)
    }

    // Case 4: A valid duration greater or equal to 10 seconds should be returned.
    analysis.Interval = "15s"
    if d := canary.GetAnalysisInterval(); d != 15*time.Second {
    t.Errorf("GetAnalysisInterval() = %v, want %v", d, 15*time.Second)
    }
}

// TestGetAnalysisThreshold tests the GetAnalysisThreshold method.
func TestGetAnalysisThreshold(t *testing.T) {
    analysis := &CanaryAnalysis{
    Threshold: 5,
    }
    canary := &Canary{
    Spec: CanarySpec{
    Analysis: analysis,
    },
    }
    if got := canary.GetAnalysisThreshold(); got != 5 {
    t.Errorf("GetAnalysisThreshold() = %d, want %d", got, 5)
    }

    // When Threshold is 0, should default to 1.
    analysis.Threshold = 0
    if got := canary.GetAnalysisThreshold(); got != 1 {
    t.Errorf("GetAnalysisThreshold() = %d, want %d", got, 1)
    }
}

// TestGetAnalysisPrimaryReadyThreshold tests the primary ready threshold getter.
func TestGetAnalysisPrimaryReadyThreshold(t *testing.T) {
    // When set, should return the provided value.
    threshold := 80
    analysis := &CanaryAnalysis{
    PrimaryReadyThreshold: &threshold,
    }
    canary := &Canary{
    Spec: CanarySpec{
    Analysis: analysis,
    },
    }
    if got := canary.GetAnalysisPrimaryReadyThreshold(); got != 80 {
    t.Errorf("GetAnalysisPrimaryReadyThreshold() = %d, want %d", got, 80)
    }

    // When not set, should return the constant PrimaryReadyThreshold.
    analysis.PrimaryReadyThreshold = nil
    if got := canary.GetAnalysisPrimaryReadyThreshold(); got != PrimaryReadyThreshold {
    t.Errorf("GetAnalysisPrimaryReadyThreshold() = %d, want %d", got, PrimaryReadyThreshold)
    }
}

// TestGetAnalysisCanaryReadyThreshold tests the canary ready threshold getter.
func TestGetAnalysisCanaryReadyThreshold(t *testing.T) {
    // When set, should return the provided value.
    threshold := 75
    analysis := &CanaryAnalysis{
    CanaryReadyThreshold: &threshold,
    }
    canary := &Canary{
    Spec: CanarySpec{
    Analysis: analysis,
    },
    }
    if got := canary.GetAnalysisCanaryReadyThreshold(); got != 75 {
    t.Errorf("GetAnalysisCanaryReadyThreshold() = %d, want %d", got, 75)
    }

    // When not set, should return the constant CanaryReadyThreshold.
    analysis.CanaryReadyThreshold = nil
    if got := canary.GetAnalysisCanaryReadyThreshold(); got != CanaryReadyThreshold {
    t.Errorf("GetAnalysisCanaryReadyThreshold() = %d, want %d", got, CanaryReadyThreshold)
    }
}

// TestGetMetricInterval tests that the metric interval is always "1m".
func TestGetMetricInterval(t *testing.T) {
    canary := &Canary{}
    if got := canary.GetMetricInterval(); got != MetricInterval {
    t.Errorf("GetMetricInterval() = %s, want %s", got, MetricInterval)
    }
}

// TestSkipAnalysis tests the SkipAnalysis method.
func TestSkipAnalysis(t *testing.T) {
    // Case 1: When both Analysis and CanaryAnalysis are nil, should skip analysis.
    canary := &Canary{
    Spec: CanarySpec{},
    }
    if !canary.SkipAnalysis() {
    t.Error("SkipAnalysis() = false, want true when both Analysis and CanaryAnalysis are nil")
    }

    // Case 2: When either Analysis or CanaryAnalysis is provided without SkipAnalysis true, should not skip.
    analysis := &CanaryAnalysis{
    Interval: "1m",
    }
    canary.Spec.Analysis = analysis
    canary.Spec.SkipAnalysis = false
    if canary.SkipAnalysis() {
    t.Error("SkipAnalysis() = true, want false when Analysis is provided and SkipAnalysis is false")
    }

    // Case 3: When SkipAnalysis is explicitly set to true.
    canary.Spec.SkipAnalysis = true
    if !canary.SkipAnalysis() {
    t.Error("SkipAnalysis() = false, want true when SkipAnalysis is true")
    }
}
// TestHTTPRewriteGetTypeEmpty tests the GetType method of HTTPRewrite when the type is empty.
func TestHTTPRewriteGetTypeEmpty(t *testing.T) {
    rewrite := &HTTPRewrite{
        Type: "",
    }
    want := string(gatewayv1beta1.FullPathHTTPPathModifier)
    if got := rewrite.GetType(); got != want {
        t.Errorf("HTTPRewrite.GetType() with empty type = %v, want %v", got, want)
    }
}

// TestGetAnalysisWithBothProvided tests that when both Analysis and deprecated CanaryAnalysis are provided,
// GetAnalysis returns the Analysis field (i.e. the non-deprecated one).
func TestGetAnalysisWithBothProvided(t *testing.T) {
    analysis := &CanaryAnalysis{
        Interval: "5m",
    }
    deprecated := &CanaryAnalysis{
        Interval: "4m",
    }
    canary := &Canary{
        Spec: CanarySpec{
            Analysis:       analysis,
            CanaryAnalysis: deprecated,
        },
    }
    if got := canary.GetAnalysis(); got != analysis {
        t.Errorf("GetAnalysis() with both fields provided = %v, want %v", got, analysis)
    }
}
// TestCanaryJSONMarshalling tests JSON marshalling and unmarshalling of the Canary struct,
func TestCanaryJSONMarshalling(t *testing.T) {
    // Create a Canary object with various fields set, including nested structures.
    canary := Canary{
        TypeMeta: metav1.TypeMeta{
            Kind:       "Canary",
            APIVersion: "flagger.app/v1beta1",
        },
        ObjectMeta: metav1.ObjectMeta{
            Name:      "test-canary",
            Namespace: "default",
        },
        Spec: CanarySpec{
            Provider: "istio",
            TargetRef: LocalObjectReference{
                Name: "foo",
            },
            Service: CanaryService{
                Name:      "bar",
                Port:      8080,
                PortName:  "http",
                TargetPort: intstr.FromInt(8080),
            },
            Analysis: &CanaryAnalysis{
                Interval:  "30s",
                Threshold: 3,
                Metrics: []CanaryMetric{
                    {
                        Name:      "error-rate",
                        Interval:  "30s",
                        Threshold: 0.1,
                    },
                },
            },
        },
    }

    // Marshal the Canary object to JSON.
    data, err := json.Marshal(canary)
    if err != nil {
        t.Errorf("Failed to marshal Canary: %v", err)
    }

    // Unmarshal the JSON back into an object.
    var newCanary Canary
    err = json.Unmarshal(data, &newCanary)
    if err != nil {
        t.Errorf("Failed to unmarshal Canary: %v", err)
    }

    // Verify that key fields are unmarshalled correctly.
    if newCanary.ObjectMeta.Name != "test-canary" {
        t.Errorf("Expected name 'test-canary', got '%s'", newCanary.ObjectMeta.Name)
    }
    if newCanary.Spec.Provider != "istio" {
        t.Errorf("Expected provider 'istio', got '%s'", newCanary.Spec.Provider)
    }
    if newCanary.Spec.Service.Port != 8080 {
        t.Errorf("Expected service port 8080, got %d", newCanary.Spec.Service.Port)
    }
    if newCanary.Spec.Analysis == nil || newCanary.Spec.Analysis.Threshold != 3 {
        t.Errorf("Expected analysis threshold 3, got %v", newCanary.Spec.Analysis)
    }
}
// TestGetAnalysisNil tests that GetAnalysis returns nil when both Analysis and CanaryAnalysis are nil.
func TestGetAnalysisNil(t *testing.T) {
    canary := &Canary{
        Spec: CanarySpec{},
    }
    if analysis := canary.GetAnalysis(); analysis != nil {
        t.Errorf("Expected GetAnalysis() to return nil when both Analysis and CanaryAnalysis are nil")
    }
}

// TestSessionAffinityNegativeMaxAge tests that a negative MaxAge is returned as set.
func TestSessionAffinityNegativeMaxAge(t *testing.T) {
    affinity := &SessionAffinity{MaxAge: -5}
    if got := affinity.GetMaxAge(); got != -5 {
        t.Errorf("Expected GetMaxAge() to return -5, got %d", got)
    }
}

// int32Ptr is a helper function that returns a pointer to an int32 value.
func int32Ptr(i int32) *int32 {
    return &i
}

// TestAutoscalerReferenceJSONMarshalling tests JSON marshalling and unmarshalling of AutoscalerRefernce.
func TestAutoscalerReferenceJSONMarshalling(t *testing.T) {
    autoscaler := AutoscalerRefernce{
        APIVersion: "v1",
        Kind:       "HorizontalPodAutoscaler",
        Name:       "hpa-test",
        PrimaryScalerQueries: map[string]string{
            "q1": "query1",
        },
        PrimaryScalerReplicas: &ScalerReplicas{
            MinReplicas: int32Ptr(1),
            MaxReplicas: int32Ptr(5),
        },
    }
    data, err := json.Marshal(autoscaler)
    if err != nil {
        t.Fatalf("Failed to marshal autoscaler: %v", err)
    }

    var newAutoscaler AutoscalerRefernce
    err = json.Unmarshal(data, &newAutoscaler)
    if err != nil {
        t.Fatalf("Failed to unmarshal autoscaler: %v", err)
    }

    if newAutoscaler.Name != "hpa-test" || newAutoscaler.APIVersion != "v1" {
        t.Errorf("Unmarshalled autoscaler does not match expected values")
    }
}