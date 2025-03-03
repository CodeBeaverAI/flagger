package v1beta1

import (
    "testing"
    "time"
    gwv1beta1 "github.com/fluxcd/flagger/pkg/apis/gatewayapi/v1beta1"
)

// TestHTTPRewriteGetType tests the GetType method of HTTPRewrite.
func TestHTTPRewriteGetType(t *testing.T) {
    // When Type matches the constant for prefix, it should return that constant.
    rw := &HTTPRewrite{Type: string(gwv1beta1.PrefixMatchHTTPPathModifier)}
    if rw.GetType() != string(gwv1beta1.PrefixMatchHTTPPathModifier) {
        t.Errorf("Expected type '%s', got %s", string(gwv1beta1.PrefixMatchHTTPPathModifier), rw.GetType())
    }

    // When Type is empty, GetType should return the default full path modifier.
    rw.Type = ""
    if rw.GetType() != string(gwv1beta1.FullPathHTTPPathModifier) {
        t.Errorf("Expected type '%s', got %s", string(gwv1beta1.FullPathHTTPPathModifier), rw.GetType())
    }

}

// TestCanaryServiceGetIstioRewrite tests the GetIstioRewrite method of CanaryService.
func TestCanaryServiceGetIstioRewrite(t *testing.T) {
    cs := &CanaryService{}
    // With nil Rewrite, the method should return nil.
    if cs.GetIstioRewrite() != nil {
    t.Error("Expected nil IstioRewrite when Rewrite is nil")
    }

    // With a proper Rewrite, the method should return an istio rewrite object.
    cs.Rewrite = &HTTPRewrite{
    Authority: "example.com",
    Uri:       "/test",
    }
    istioRewrite := cs.GetIstioRewrite()
    if istioRewrite == nil {
    t.Error("Expected valid IstioRewrite object")
    }
    if istioRewrite.Authority != "example.com" || istioRewrite.Uri != "/test" {
    t.Errorf("IstioRewrite mismatch: got %+v", istioRewrite)
    }
}

// TestSessionAffinityGetMaxAge tests the GetMaxAge method of SessionAffinity.
func TestSessionAffinityGetMaxAge(t *testing.T) {
    sa := &SessionAffinity{}
    // Expect the default value of 86400 seconds if MaxAge is 0.
    if sa.GetMaxAge() != 86400 {
    t.Errorf("Expected default max age 86400, got %d", sa.GetMaxAge())
    }
    sa.MaxAge = 100
    if sa.GetMaxAge() != 100 {
    t.Errorf("Expected max age 100, got %d", sa.GetMaxAge())
    }
}

// TestGetServiceNames tests the GetServiceNames method of Canary.
func TestGetServiceNames(t *testing.T) {
    c := &Canary{
    Spec: CanarySpec{
    TargetRef: LocalObjectReference{
    Name: "myapp",
    },
    Service: CanaryService{
    Name: "custom-service",
    },
    },
    }
    apex, primary, canarySvc := c.GetServiceNames()
    if apex != "custom-service" {
    t.Errorf("Expected apex name 'custom-service', got %s", apex)
    }
    if primary != "custom-service-primary" {
    t.Errorf("Expected primary name 'custom-service-primary', got %s", primary)
    }
    if canarySvc != "custom-service-canary" {
    t.Errorf("Expected canary name 'custom-service-canary', got %s", canarySvc)
    }
}

// TestProgressDeadline tests the GetProgressDeadlineSeconds method of Canary.
func TestProgressDeadline(t *testing.T) {
    // Test default deadline.
    c := &Canary{
    Spec: CanarySpec{},
    }
    if c.GetProgressDeadlineSeconds() != ProgressDeadlineSeconds {
    t.Errorf("Expected default progress deadline %d, got %d", ProgressDeadlineSeconds, c.GetProgressDeadlineSeconds())
    }
    // Test a custom deadline.
    deadline := int32(300)
    c.Spec.ProgressDeadlineSeconds = &deadline
    if c.GetProgressDeadlineSeconds() != 300 {
    t.Errorf("Expected custom progress deadline 300, got %d", c.GetProgressDeadlineSeconds())
    }
}

// TestGetAnalysis tests the GetAnalysis method of Canary.
func TestGetAnalysis(t *testing.T) {
    c := &Canary{
    Spec: CanarySpec{
    Analysis: &CanaryAnalysis{Interval: "30s"},
    },
    }
    analysis := c.GetAnalysis()
    if analysis == nil || analysis.Interval != "30s" {
    t.Errorf("Expected analysis interval '30s', got %+v", analysis)
    }

    // When Analysis is nil, but CanaryAnalysis is provided.
    c.Spec.Analysis = nil
    c.Spec.CanaryAnalysis = &CanaryAnalysis{Interval: "45s"}
    analysis = c.GetAnalysis()
    if analysis == nil || analysis.Interval != "45s" {
    t.Errorf("Expected canaryAnalysis interval '45s', got %+v", analysis)
    }
}

// TestGetAnalysisInterval tests the GetAnalysisInterval method of Canary.
func TestGetAnalysisInterval(t *testing.T) {
    c := &Canary{
    Spec: CanarySpec{
    Analysis: &CanaryAnalysis{Interval: "15s"},
    },
    }
    if c.GetAnalysisInterval() != 15*time.Second {
    t.Errorf("Expected analysis interval 15s, got %v", c.GetAnalysisInterval())
    }

    // Test with an invalid duration: should fall back to default.
    c.Spec.Analysis.Interval = "invalid"
    if c.GetAnalysisInterval() != AnalysisInterval {
    t.Errorf("Expected fallback analysis interval %v, got %v", AnalysisInterval, c.GetAnalysisInterval())
    }

    // Test enforcing a minimum duration (if less than 10s, return 10s).
    c.Spec.Analysis.Interval = "5s"
    if c.GetAnalysisInterval() != 10*time.Second {
    t.Errorf("Expected enforced minimum analysis interval 10s, got %v", c.GetAnalysisInterval())
    }
}

// TestGetAnalysisThreshold tests the GetAnalysisThreshold method of Canary.
func TestGetAnalysisThreshold(t *testing.T) {
    c := &Canary{
    Spec: CanarySpec{
    Analysis: &CanaryAnalysis{Threshold: 3},
    },
    }
    if c.GetAnalysisThreshold() != 3 {
    t.Errorf("Expected analysis threshold 3, got %d", c.GetAnalysisThreshold())
    }

    // With threshold 0, the default value (1) should be returned.
    c.Spec.Analysis.Threshold = 0
    if c.GetAnalysisThreshold() != 1 {
    t.Errorf("Expected analysis threshold default 1, got %d", c.GetAnalysisThreshold())
    }
}

// TestReadyThresholds tests GetAnalysisPrimaryReadyThreshold and GetAnalysisCanaryReadyThreshold.
func TestReadyThresholds(t *testing.T) {
    // Test defaults when ready thresholds are not set.
    c := &Canary{
    Spec: CanarySpec{
    Analysis: &CanaryAnalysis{},
    },
    }
    if c.GetAnalysisPrimaryReadyThreshold() != PrimaryReadyThreshold {
    t.Errorf("Expected primary ready threshold %d, got %d", PrimaryReadyThreshold, c.GetAnalysisPrimaryReadyThreshold())
    }
    if c.GetAnalysisCanaryReadyThreshold() != CanaryReadyThreshold {
    t.Errorf("Expected canary ready threshold %d, got %d", CanaryReadyThreshold, c.GetAnalysisCanaryReadyThreshold())
    }

    // Test custom ready threshold values.
    customThreshold := 80
    c.Spec.Analysis.PrimaryReadyThreshold = &customThreshold
    c.Spec.Analysis.CanaryReadyThreshold = &customThreshold
    if c.GetAnalysisPrimaryReadyThreshold() != 80 {
        t.Errorf("Expected primary ready threshold 80, got %d", c.GetAnalysisPrimaryReadyThreshold())
    }
    if c.GetAnalysisCanaryReadyThreshold() != 80 {
        t.Errorf("Expected canary ready threshold 80, got %d", c.GetAnalysisCanaryReadyThreshold())
    }
    c.Spec.Analysis.PrimaryReadyThreshold = &customThreshold
    c.Spec.Analysis.CanaryReadyThreshold = &customThreshold
    if c.GetAnalysisCanaryReadyThreshold() != 80 {
    t.Errorf("Expected canary ready threshold 80, got %d", c.GetAnalysisCanaryReadyThreshold())
    }
}

// TestGetMetricInterval tests that GetMetricInterval returns the default metric interval.
func TestGetMetricInterval(t *testing.T) {
    c := &Canary{}
    if c.GetMetricInterval() != MetricInterval {
    t.Errorf("Expected metric interval %s, got %s", MetricInterval, c.GetMetricInterval())
    }
}

// TestSkipAnalysis tests the SkipAnalysis method of Canary.
func TestSkipAnalysis(t *testing.T) {
    // When both Analysis and CanaryAnalysis are nil, SkipAnalysis should return true.
    c := &Canary{
    Spec: CanarySpec{
    Analysis:       nil,
    CanaryAnalysis: nil,
    SkipAnalysis:   false,
    },
    }
    if !c.SkipAnalysis() {
    t.Error("Expected SkipAnalysis to return true when analysis fields are nil")
    }

    // When an analysis is present, the SkipAnalysis flag should be effective.
    c.Spec.Analysis = &CanaryAnalysis{}
    c.Spec.SkipAnalysis = true
    if !c.SkipAnalysis() {
    t.Errorf("Expected SkipAnalysis to return true when SkipAnalysis flag is true")
    }
    c.Spec.SkipAnalysis = false
    if c.SkipAnalysis() {
    t.Errorf("Expected SkipAnalysis to return false when analysis is present and flag is false")
    }
}
// TestHTTPRewriteGetType_NonPrefix tests that if HTTPRewrite type is not the prefix match,
func TestHTTPRewriteGetType_NonPrefix(t *testing.T) {
    rw := &HTTPRewrite{Type: "custom"}
    expected := string(gwv1beta1.FullPathHTTPPathModifier)
    if got := rw.GetType(); got != expected {
        t.Errorf("Expected type '%s', got '%s'", expected, got)
    }
}

// TestGetAnalysisNil tests that GetAnalysis returns nil when both Analysis and CanaryAnalysis are nil.
func TestGetAnalysisNil(t *testing.T) {
    c := &Canary{
        Spec: CanarySpec{},
    }
    if analysis := c.GetAnalysis(); analysis != nil {
        t.Errorf("Expected nil analysis when both Analysis and CanaryAnalysis are nil, got %+v", analysis)
    }
}

// TestGetAnalysis_BothSet tests that when both Analysis and CanaryAnalysis are provided,
func TestGetAnalysis_BothSet(t *testing.T) {
    c := &Canary{
        Spec: CanarySpec{
            Analysis:       &CanaryAnalysis{Interval: "20s"},
            CanaryAnalysis: &CanaryAnalysis{Interval: "40s"},
        },
    }
    analysis := c.GetAnalysis()
    if analysis == nil || analysis.Interval != "20s" {
        t.Errorf("Expected analysis interval '20s', got %+v", analysis)
    }
}

// TestGetAnalysisInterval_NoAnalysis tests that GetAnalysisInterval panics when no Analysis (or CanaryAnalysis) is set.
func TestGetAnalysisInterval_NoAnalysis(t *testing.T) {
    c := &Canary{
        Spec: CanarySpec{},
    }
    defer func() {
        if r := recover(); r == nil {
            t.Error("Expected panic when GetAnalysisInterval is called with nil analysis")
        }
    }()
    _ = c.GetAnalysisInterval()
}
// TestGetServiceNames_Default tests that when Service.Name is empty, the apex name falls back to TargetRef.Name.
func TestGetServiceNames_Default(t *testing.T) {
    c := &Canary{
        Spec: CanarySpec{
            TargetRef: LocalObjectReference{Name: "app-default"},
            Service: CanaryService{
                Port: 80,
            },
        },
    }
    apex, primary, canarySvc := c.GetServiceNames()
    if apex != "app-default" {
        t.Errorf("Expected apex name 'app-default', got %s", apex)
    }
    if primary != "app-default-primary" {
        t.Errorf("Expected primary name 'app-default-primary', got %s", primary)
    }
    if canarySvc != "app-default-canary" {
        t.Errorf("Expected canary name 'app-default-canary', got %s", canarySvc)
    }
}

// TestCanaryServiceGetIstioRewrite_Empty tests that GetIstioRewrite returns an object (even with empty Rewrite values).
func TestCanaryServiceGetIstioRewrite_Empty(t *testing.T) {
    cs := &CanaryService{
        Rewrite: &HTTPRewrite{},
    }
    istioRewrite := cs.GetIstioRewrite()
    if istioRewrite == nil {
        t.Error("Expected IstioRewrite object, got nil")
    }
    if istioRewrite.Authority != "" || istioRewrite.Uri != "" {
        t.Errorf("Expected empty IstioRewrite fields, got %+v", istioRewrite)
    }
}

// TestGetAnalysisInterval_MinimumEnforcement tests that extremely short intervals (e.g., "1s") are bumped to the minimum (10s).
func TestGetAnalysisInterval_MinimumEnforcement(t *testing.T) {
    c := &Canary{
        Spec: CanarySpec{
            Analysis: &CanaryAnalysis{Interval: "1s"},
        },
    }
    interval := c.GetAnalysisInterval()
    if interval != 10*time.Second {
        t.Errorf("Expected enforced minimum analysis interval 10s, got %v", interval)
    }
}

// TestFullCanaryGetterDefaults tests a fully default Canary configuration and verifies the getters return expected default values.
func TestFullCanaryGetterDefaults(t *testing.T) {
    deadline := int32(400)
    analysis := &CanaryAnalysis{
        // Threshold not set, should default to 1 in getter.
    }

    c := &Canary{
        Spec: CanarySpec{
            TargetRef: LocalObjectReference{Name: "default-app"},
            Service: CanaryService{
                Port: 8080,
            },
            Analysis: analysis,
            ProgressDeadlineSeconds: &deadline,
            SkipAnalysis: false,
        },
    }

    // Check GetServiceNames
    apex, primary, canarySvc := c.GetServiceNames()
    if apex != "default-app" {
        t.Errorf("Expected apex name 'default-app', got %s", apex)
    }
    if primary != "default-app-primary" {
        t.Errorf("Expected primary name 'default-app-primary', got %s", primary)
    }
    if canarySvc != "default-app-canary" {
        t.Errorf("Expected canary name 'default-app-canary', got %s", canarySvc)
    }

    // Check GetProgressDeadlineSeconds should return 400
    if c.GetProgressDeadlineSeconds() != 400 {
        t.Errorf("Expected progress deadline 400, got %d", c.GetProgressDeadlineSeconds())
    }

    // GetAnalysisInterval: Analysis.Interval is empty so should fall back to default (60s)
    if c.GetAnalysisInterval() != AnalysisInterval {
        t.Errorf("Expected analysis interval %v, got %v", AnalysisInterval, c.GetAnalysisInterval())
    }

    // GetAnalysisThreshold should return 1 (default) as threshold is unset
    if c.GetAnalysisThreshold() != 1 {
        t.Errorf("Expected analysis threshold default 1, got %d", c.GetAnalysisThreshold())
    }

    // GetAnalysisPrimaryReadyThreshold should return default PrimaryReadyThreshold (100)
    if c.GetAnalysisPrimaryReadyThreshold() != PrimaryReadyThreshold {
        t.Errorf("Expected primary ready threshold %d, got %d", PrimaryReadyThreshold, c.GetAnalysisPrimaryReadyThreshold())
    }

    // GetAnalysisCanaryReadyThreshold should return default CanaryReadyThreshold (100)
    if c.GetAnalysisCanaryReadyThreshold() != CanaryReadyThreshold {
        t.Errorf("Expected canary ready threshold %d, got %d", CanaryReadyThreshold, c.GetAnalysisCanaryReadyThreshold())
    }

    // GetMetricInterval should return the default metric interval "1m"
    if c.GetMetricInterval() != MetricInterval {
        t.Errorf("Expected metric interval %s, got %s", MetricInterval, c.GetMetricInterval())
    }

    // Since Analysis is set, SkipAnalysis should return the flag (false)
    if c.SkipAnalysis() != false {
        t.Errorf("Expected SkipAnalysis false, got true")
    }
}
// TestGetAnalysisInterval_ZeroDuration tests that when the Analysis interval is "0s", 
// the enforced minimum duration of 10 seconds is returned.
func TestGetAnalysisInterval_ZeroDuration(t *testing.T) {
    c := &Canary{
        Spec: CanarySpec{
            Analysis: &CanaryAnalysis{Interval: "0s"},
        },
    }
    expected := 10 * time.Second
    if dur := c.GetAnalysisInterval(); dur != expected {
        t.Errorf("Expected analysis interval %v for '0s', got %v", expected, dur)
    }
}

// TestAutoscalerReference tests that the AutoscalerReference fields are correctly assigned.
func TestAutoscalerReference(t *testing.T) {
    asr := &AutoscalerRefernce{
        APIVersion: "autoscaling/v1",
        Kind:       "HorizontalPodAutoscaler",
        Name:       "hpa-test",
        PrimaryScalerQueries: map[string]string{
            "query1": "value1",
        },
        PrimaryScalerReplicas: &ScalerReplicas{
            MinReplicas: int32Ptr(1),
            MaxReplicas: int32Ptr(10),
        },
    }

    c := &Canary{
        Spec: CanarySpec{
            AutoscalerRef: asr,
        },
    }

    if c.Spec.AutoscalerRef.APIVersion != "autoscaling/v1" {
        t.Errorf("Expected AutoscalerRef.APIVersion 'autoscaling/v1', got %s", c.Spec.AutoscalerRef.APIVersion)
    }
    if c.Spec.AutoscalerRef.Kind != "HorizontalPodAutoscaler" {
        t.Errorf("Expected AutoscalerRef.Kind 'HorizontalPodAutoscaler', got %s", c.Spec.AutoscalerRef.Kind)
    }
    if c.Spec.AutoscalerRef.Name != "hpa-test" {
        t.Errorf("Expected AutoscalerRef.Name 'hpa-test', got %s", c.Spec.AutoscalerRef.Name)
    }
    if q, ok := c.Spec.AutoscalerRef.PrimaryScalerQueries["query1"]; !ok || q != "value1" {
        t.Errorf("Expected PrimaryScalerQueries[\"query1\"] 'value1', got %s", q)
    }
    if c.Spec.AutoscalerRef.PrimaryScalerReplicas == nil {
        t.Error("Expected PrimaryScalerReplicas to be non-nil")
    } else {
        if c.Spec.AutoscalerRef.PrimaryScalerReplicas.MinReplicas == nil || *c.Spec.AutoscalerRef.PrimaryScalerReplicas.MinReplicas != 1 {
            t.Errorf("Expected MinReplicas to be 1, got %v", c.Spec.AutoscalerRef.PrimaryScalerReplicas.MinReplicas)
        }
        if c.Spec.AutoscalerRef.PrimaryScalerReplicas.MaxReplicas == nil || *c.Spec.AutoscalerRef.PrimaryScalerReplicas.MaxReplicas != 10 {
            t.Errorf("Expected MaxReplicas to be 10, got %v", c.Spec.AutoscalerRef.PrimaryScalerReplicas.MaxReplicas)
        }
    }
}

// int32Ptr is a helper function that returns a pointer to the given int32 value.
func int32Ptr(i int32) *int32 {
    return &i
}