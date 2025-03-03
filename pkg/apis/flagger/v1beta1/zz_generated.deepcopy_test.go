package v1beta1

import (
    "reflect"
    "testing"
    "time"

    istiov1beta1 "github.com/fluxcd/flagger/pkg/apis/istio/v1beta1"
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestDeepCopy_Canary tests the DeepCopy, DeepCopyInto, and DeepCopyObject functions on Canary.
func TestDeepCopy_Canary(t *testing.T) {
    // Create a sample Canary instance
    now := metav1.NewTime(time.Now())
    orig := &Canary{
    TypeMeta: metav1.TypeMeta{Kind: "Canary", APIVersion: "flagger.app/v1beta1"},
    ObjectMeta: metav1.ObjectMeta{
    Name:              "test-canary",
    Namespace:         "default",
    CreationTimestamp: now,
    },
    Spec: CanarySpec{
    TargetRef: LocalObjectReference{Name: "test-target"},
    Service: CanaryService{
    Gateways: []string{"gw1", "gw2"},
    Hosts:    []string{"example.com"},
    },
    },
    Status: CanaryStatus{
    TrackedConfigs: &map[string]string{"config1": "v1"},
    Conditions: []CanaryCondition{
    {
        Type:               "Ready",
        Status:             "True",
        LastTransitionTime: now,
    },
    },
    },
    }

    // Make a deep copy of the Canary
    cp := orig.DeepCopy()
    if cp == nil {
    t.Fatalf("DeepCopy returned nil")
    }

    // Check DeepCopyObject
    var objCopy interface{} = orig.DeepCopyObject()
    if objCopy == nil {
    t.Fatalf("DeepCopyObject returned nil")
    }

    // Verify that the copy is equal
    if !reflect.DeepEqual(orig, cp) {
    t.Error("DeepCopy did not produce an equal object")
    }

    // Mutate the original and make sure the copy does not change to ensure deep copy
    orig.Spec.TargetRef.Name = "modified-target"
    (*orig.Status.TrackedConfigs)["config1"] = "modified"
    orig.Spec.Service.Gateways[0] = "modified-gw"

    // The copy should remain unchanged
    if reflect.DeepEqual(orig, cp) {
    t.Error("Deep copy is not independent from the original")
    }
}

// TestDeepCopy_AlertProvider tests the deep copy functions on AlertProvider.
func TestDeepCopy_AlertProvider(t *testing.T) {
    now := metav1.NewTime(time.Now())
    orig := &AlertProvider{
    TypeMeta: metav1.TypeMeta{Kind: "AlertProvider", APIVersion: "flagger.app/v1beta1"},
    ObjectMeta: metav1.ObjectMeta{
    Name:              "test-alert",
    Namespace:         "default",
    CreationTimestamp: now,
    },
    Spec: AlertProviderSpec{
    SecretRef: &v1.LocalObjectReference{Name: "secret-ref"},
    },
    Status: AlertProviderStatus{
    Conditions: []AlertProviderCondition{
    {
        LastUpdateTime: now,
        LastTransitionTime: now,
    },
    },
    },
    }

    cp := orig.DeepCopy()
    if cp == nil {
    t.Fatalf("DeepCopy returned nil for AlertProvider")
    }

    if !reflect.DeepEqual(orig, cp) {
    t.Error("DeepCopy for AlertProvider did not produce an equal object")
    }

    // Mutate original
    orig.Spec.SecretRef.Name = "modified-secret"
    orig.Status.Conditions[0].LastUpdateTime = metav1.NewTime(time.Now().Add(2 * time.Hour))

    if reflect.DeepEqual(orig, cp) {
    t.Error("Deep copy for AlertProvider is not independent from the original")
    }
}

// TestDeepCopy_MetricTemplate tests the deep copy functions on MetricTemplate.
func TestDeepCopy_MetricTemplate(t *testing.T) {
    now := metav1.NewTime(time.Now())
    orig := &MetricTemplate{
    TypeMeta: metav1.TypeMeta{Kind: "MetricTemplate", APIVersion: "flagger.app/v1beta1"},
    ObjectMeta: metav1.ObjectMeta{
    Name:              "test-metric",
    Namespace:         "default",
    CreationTimestamp: now,
    },
    Spec: MetricTemplateSpec{
    Provider: MetricTemplateProvider{
    SecretRef: &v1.LocalObjectReference{Name: "provider-secret"},
    },
    },
    Status: MetricTemplateStatus{
    Conditions: []MetricTemplateCondition{
    {
        LastUpdateTime:     now,
        LastTransitionTime: now,
    },
    },
    },
    }

    cp := orig.DeepCopy()
    if cp == nil {
    t.Fatalf("DeepCopy returned nil for MetricTemplate")
    }

    if !reflect.DeepEqual(orig, cp) {
    t.Error("DeepCopy for MetricTemplate did not produce an equal object")
    }

    // Change some fields on the original to ensure the copy is deep
    orig.Spec.Provider.SecretRef.Name = "modified-provider-secret"
    orig.Status.Conditions[0].LastTransitionTime = metav1.NewTime(time.Now().Add(3 * time.Hour))

    if reflect.DeepEqual(orig, cp) {
    t.Error("Deep copy for MetricTemplate is not independent from the original")
    }
}

// TestDeepCopy_ScalerReplicas tests the deep copy functions on ScalerReplicas.
func TestDeepCopy_ScalerReplicas(t *testing.T) {
    orig := &ScalerReplicas{
    MinReplicas: new(int32),
    MaxReplicas: new(int32),
    }
    *orig.MinReplicas = 1
    *orig.MaxReplicas = 10

    cp := orig.DeepCopy()
    if cp == nil {
    t.Fatalf("DeepCopy returned nil for ScalerReplicas")
    }

    if !reflect.DeepEqual(orig, cp) {
    t.Error("DeepCopy for ScalerReplicas did not produce an equal object")
    }

    // Mutate original
    *orig.MinReplicas = 2
    *orig.MaxReplicas = 20

    if reflect.DeepEqual(orig, cp) {
    t.Error("Deep copy for ScalerReplicas is not independent from the original")
    }
}
// TestDeepCopy_CanaryAlert tests the deep copy functions on CanaryAlert.
func TestDeepCopy_CanaryAlert(t *testing.T) {
    // Create a sample CanaryAlert instance
    orig := &CanaryAlert{
        ProviderRef: CrossNamespaceObjectReference{Name: "alert-provider"},
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for CanaryAlert")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for CanaryAlert did not produce an equal object")
    }
    // Mutate the original to check that the copy is independent
    orig.ProviderRef.Name = "modified-alert-provider"
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for CanaryAlert is not independent from the original")
    }
}

// TestDeepCopy_CanaryAnalysis tests the deep copy functions on CanaryAnalysis.
func TestDeepCopy_CanaryAnalysis(t *testing.T) {
    // Prepare pointer values for thresholds
    prt := new(int)
    *prt = 5
    crt := new(int)
    *crt = 7

    // Create a sample CanaryAnalysis with non-nil slices and pointer fields.
    orig := &CanaryAnalysis{
        StepWeights:           []int{10, 20},
        PrimaryReadyThreshold: prt,
        CanaryReadyThreshold:  crt,
        Alerts: []CanaryAlert{
            {ProviderRef: CrossNamespaceObjectReference{Name: "alert"}},
        },
        Metrics: []CanaryMetric{
            {
                ThresholdRange: &CanaryThresholdRange{
                    Min: func() *float64 { f := 0.1; return &f }(),
                    Max: func() *float64 { f := 1.0; return &f }(),
                },
                TemplateRef: &CrossNamespaceObjectReference{Name: "ref", Namespace: "ns"},
                TemplateVariables: map[string]string{
                    "key": "val",
                },
            },
        },
        Webhooks: []CanaryWebhook{
            {Metadata: func() *map[string]string { m := map[string]string{"hook": "value"}; return &m }()},
        },
        Match: []istiov1beta1.HTTPMatchRequest{
            {},
        },
        SessionAffinity: &SessionAffinity{},
    }

    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for CanaryAnalysis")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for CanaryAnalysis did not produce an equal object")
    }

    // Mutate the original to test deep copy isolation.
    orig.StepWeights[0] = 99
    *orig.PrimaryReadyThreshold = 50
    orig.Alerts[0].ProviderRef.Name = "modified-alert"
    orig.Metrics[0].TemplateVariables["key"] = "modified-val"
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for CanaryAnalysis is not independent from the original")
    }
}

// TestDeepCopy_CanaryList tests the deep copy functions on CanaryList.
func TestDeepCopy_CanaryList(t *testing.T) {
    // Create a sample CanaryList with two Canary items.
    now := metav1.NewTime(time.Now())
    orig := &CanaryList{
        TypeMeta: metav1.TypeMeta{Kind: "Canary", APIVersion: "flagger.app/v1beta1"},
        ListMeta: metav1.ListMeta{ResourceVersion: "v1"},
        Items: []Canary{
            {
                ObjectMeta: metav1.ObjectMeta{Name: "canary-1", CreationTimestamp: now},
                Spec:       CanarySpec{TargetRef: LocalObjectReference{Name: "target-1"}},
            },
            {
                ObjectMeta: metav1.ObjectMeta{Name: "canary-2", CreationTimestamp: now},
                Spec:       CanarySpec{TargetRef: LocalObjectReference{Name: "target-2"}},
            },
        },
    }

    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for CanaryList")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for CanaryList did not produce an equal object")
    }

    // Mutate original and verify change does not affect the copy.
    orig.Items[0].ObjectMeta.Name = "modified-canary-1"
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for CanaryList is not independent from the original")
    }
}

// TestDeepCopy_CanaryMetric tests the deep copy functions on CanaryMetric.
func TestDeepCopy_CanaryMetric(t *testing.T) {
    // Create a sample CanaryMetric instance
    orig := &CanaryMetric{
        ThresholdRange: &CanaryThresholdRange{
            Min: func() *float64 { f := 0.5; return &f }(),
            Max: func() *float64 { f := 2.5; return &f }(),
        },
        TemplateRef: &CrossNamespaceObjectReference{Name: "metric-ref", Namespace: "default"},
        TemplateVariables: map[string]string{
            "var1": "value1",
        },
    }

    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for CanaryMetric")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for CanaryMetric did not produce an equal object")
    }

    // Mutate original to ensure independency.
    orig.ThresholdRange.Min = func() *float64 { f := 1.0; return &f }()
    orig.TemplateVariables["var1"] = "modified-value1"
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for CanaryMetric is not independent from the original")
    }
}