package v1beta1

import (
    "reflect"
    "testing"
    "time"

    istiov1beta1 "github.com/fluxcd/flagger/pkg/apis/istio/v1beta1"
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    gatewayapiv1 "github.com/fluxcd/flagger/pkg/apis/gatewayapi/v1"
    gatewayapiv1beta1 "github.com/fluxcd/flagger/pkg/apis/gatewayapi/v1beta1"
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
// TestDeepCopy_AutoscalerRefernce tests deep copy functions for AutoscalerRefernce.
func TestDeepCopy_AutoscalerRefernce(t *testing.T) {
    orig := &AutoscalerRefernce{
        PrimaryScalerQueries: map[string]string{"query": "value"},
        PrimaryScalerReplicas: &ScalerReplicas{
            MinReplicas: func(i int32) *int32 { x := i; return &x }(1),
            MaxReplicas: func(i int32) *int32 { x := i; return &x }(5),
        },
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for AutoscalerRefernce")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for AutoscalerRefernce did not produce an equal object")
    }
    // Modify original to ensure a deep copy.
    orig.PrimaryScalerQueries["query"] = "modified"
    *orig.PrimaryScalerReplicas.MinReplicas = 10
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for AutoscalerRefernce is not independent from the original")
    }
}

// TestDeepCopy_CustomBackend tests deep copy on CustomBackend.
func TestDeepCopy_CustomBackend(t *testing.T) {
    orig := &CustomBackend{
        BackendObjectReference: &gatewayapiv1.BackendObjectReference{},
        Filters: []gatewayapiv1.HTTPRouteFilter{
            {},
        },
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for CustomBackend")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for CustomBackend did not produce an equal object")
    }
    // Change original and check that copy stays unchanged.
    grp := "grp"
    kindVal := "kind"
    nameVal := "name"
    orig.BackendObjectReference = &gatewayapiv1.BackendObjectReference{
        Group: &grp,
        Kind:  &kindVal,
        Name:  &nameVal,
    }
    orig.Filters = append(orig.Filters, gatewayapiv1.HTTPRouteFilter{})
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for CustomBackend is not independent from the original")
    }
}

// TestDeepCopy_CustomMetadata tests deep copy on CustomMetadata.
func TestDeepCopy_CustomMetadata(t *testing.T) {
    orig := &CustomMetadata{
        Labels:      map[string]string{"key1": "value1"},
        Annotations: map[string]string{"anno1": "valueA"},
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for CustomMetadata")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for CustomMetadata did not produce an equal object")
    }
    orig.Labels["key1"] = "modified"
    orig.Annotations["anno1"] = "modified"
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for CustomMetadata is not independent from the original")
    }
}

// TestDeepCopy_LocalObjectReference tests deep copy on LocalObjectReference.
func TestDeepCopy_LocalObjectReference(t *testing.T) {
    orig := &LocalObjectReference{Name: "test"}
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for LocalObjectReference")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for LocalObjectReference did not produce an equal object")
    }
    orig.Name = "modified"
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for LocalObjectReference is not independent from the original")
    }
}

// TestDeepCopy_HTTPRewrite tests deep copy on HTTPRewrite.
func TestDeepCopy_HTTPRewrite(t *testing.T) {
    orig := &HTTPRewrite{}
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for HTTPRewrite")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for HTTPRewrite did not produce an equal object")
    }
}

// TestDeepCopy_MetricTemplateList tests deep copy on MetricTemplateList.
func TestDeepCopy_MetricTemplateList(t *testing.T) {
    now := metav1.NewTime(time.Now())
    orig := &MetricTemplateList{
        TypeMeta: metav1.TypeMeta{Kind: "MetricTemplate", APIVersion: "flagger.app/v1beta1"},
        ListMeta: metav1.ListMeta{ResourceVersion: "v1"},
        Items: []MetricTemplate{
            {
                ObjectMeta: metav1.ObjectMeta{Name: "metric1", CreationTimestamp: now},
                Spec: MetricTemplateSpec{
                    Provider: MetricTemplateProvider{
                        SecretRef: &v1.LocalObjectReference{Name: "secret1"},
                    },
                },
                Status: MetricTemplateStatus{},
            },
        },
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for MetricTemplateList")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for MetricTemplateList did not produce an equal object")
    }
    orig.Items[0].ObjectMeta.Name = "modified-metric1"
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for MetricTemplateList is not independent from the original")
    }
}

// TestDeepCopy_MetricTemplateModel tests deep copy on MetricTemplateModel.
func TestDeepCopy_MetricTemplateModel(t *testing.T) {
    orig := &MetricTemplateModel{
        Variables: map[string]string{"var": "val"},
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for MetricTemplateModel")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for MetricTemplateModel did not produce an equal object")
    }
    orig.Variables["var"] = "modified"
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for MetricTemplateModel is not independent from the original")
    }
}

// TestDeepCopy_MetricTemplateProvider tests deep copy on MetricTemplateProvider.
func TestDeepCopy_MetricTemplateProvider(t *testing.T) {
    orig := &MetricTemplateProvider{
        SecretRef: &v1.LocalObjectReference{Name: "provider"},
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for MetricTemplateProvider")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for MetricTemplateProvider did not produce an equal object")
    }
    orig.SecretRef.Name = "modified-provider"
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for MetricTemplateProvider is not independent from the original")
    }
}

// TestDeepCopy_MetricTemplateStatus tests deep copy on MetricTemplateStatus.
func TestDeepCopy_MetricTemplateStatus(t *testing.T) {
    now := metav1.NewTime(time.Now())
    orig := &MetricTemplateStatus{
        Conditions: []MetricTemplateCondition{
            {
                LastUpdateTime:     now,
                LastTransitionTime: now,
            },
        },
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for MetricTemplateStatus")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for MetricTemplateStatus did not produce an equal object")
    }
    orig.Conditions[0].LastUpdateTime = metav1.NewTime(time.Now().Add(1 * time.Hour))
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for MetricTemplateStatus is not independent from the original")
    }
}

// TestDeepCopy_SessionAffinity tests deep copy on SessionAffinity.
func TestDeepCopy_SessionAffinity(t *testing.T) {
    orig := &SessionAffinity{}
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for SessionAffinity")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for SessionAffinity did not produce an equal object")
    }
}

// TestDeepCopy_CanaryCondition tests deep copy on CanaryCondition.
func TestDeepCopy_CanaryCondition(t *testing.T) {
    now := metav1.NewTime(time.Now())
    orig := &CanaryCondition{
        Type:               "TestCondition",
        Status:             "True",
        LastUpdateTime:     now,
        LastTransitionTime: now,
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for CanaryCondition")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for CanaryCondition did not produce an equal object")
    }
    orig.Status = "False"
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for CanaryCondition is not independent from the original")
    }
}

// TestDeepCopy_AlertProviderCondition tests deep copy on AlertProviderCondition.
func TestDeepCopy_AlertProviderCondition(t *testing.T) {
    now := metav1.NewTime(time.Now())
    orig := &AlertProviderCondition{
        LastUpdateTime:     now,
        LastTransitionTime: now,
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for AlertProviderCondition")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for AlertProviderCondition did not produce an equal object")
    }
    cp.LastUpdateTime = metav1.NewTime(time.Now().Add(1 * time.Hour))
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for AlertProviderCondition is not independent from the original")
    }
}

// TestDeepCopy_CanaryWebhook tests deep copy on CanaryWebhook.
func TestDeepCopy_CanaryWebhook(t *testing.T) {
    orig := &CanaryWebhook{
        Metadata: func() *map[string]string {
            m := map[string]string{"hook": "orig"}
            return &m
        }(),
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for CanaryWebhook")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for CanaryWebhook did not produce an equal object")
    }
    // Modify original metadata.
    (*orig.Metadata)["hook"] = "modified"
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for CanaryWebhook is not independent from the original")
    }
}

// TestDeepCopy_CanaryWebhookPayload tests deep copy on CanaryWebhookPayload.
func TestDeepCopy_CanaryWebhookPayload(t *testing.T) {
    orig := &CanaryWebhookPayload{
        Metadata: map[string]string{"key": "value"},
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for CanaryWebhookPayload")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for CanaryWebhookPayload did not produce an equal object")
    }
    orig.Metadata["key"] = "modified"
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for CanaryWebhookPayload is not independent from the original")
    }
}
// TestDeepCopy_BackendObjectReference tests deep copy for BackendObjectReference in gatewayapi v1beta1 package
func TestDeepCopy_BackendObjectReference(t *testing.T) {
    // Setup original BackendObjectReference with pointer fields
    group := gatewayapiv1beta1.Group("group1")
    kind := gatewayapiv1beta1.Kind("kind1")
    namespace := gatewayapiv1beta1.Namespace("ns1")
    port := gatewayapiv1beta1.PortNumber(8080)
    orig := &gatewayapiv1beta1.BackendObjectReference{
        Group:     &group,
        Kind:      &kind,
        Namespace: &namespace,
        Port:      &port,
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for BackendObjectReference")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy did not produce an equal BackendObjectReference")
    }
    // Mutate original and check independence
    newGroup := gatewayapiv1beta1.Group("modified-group")
    orig.Group = &newGroup
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for BackendObjectReference is not independent from the original")
    }
}

// TestDeepCopy_HTTPRouteFilter tests deep copy for HTTPRouteFilter in gatewayapi v1beta1 package
func TestDeepCopy_HTTPRouteFilter(t *testing.T) {
    // Setup original HTTPRouteFilter with various nested fields
    orig := &gatewayapiv1beta1.HTTPRouteFilter{
        RequestHeaderModifier: &gatewayapiv1beta1.HTTPHeaderFilter{
            Set:    []gatewayapiv1beta1.HTTPHeader{{Name: "X-Req", Value: "val1"}},
            Add:    []gatewayapiv1beta1.HTTPHeader{{Name: "X-Add", Value: "val2"}},
            Remove: []string{"X-Remove"},
        },
        ResponseHeaderModifier: &gatewayapiv1beta1.HTTPHeaderFilter{
            Set: []gatewayapiv1beta1.HTTPHeader{{Name: "X-Resp", Value: "val3"}},
        },
        RequestMirror: &gatewayapiv1beta1.HTTPRequestMirrorFilter{
            BackendRef: gatewayapiv1beta1.BackendRef{
                BackendObjectReference: gatewayapiv1beta1.BackendObjectReference{},
            },
        },
        RequestRedirect: &gatewayapiv1beta1.HTTPRequestRedirectFilter{
            Scheme:   func() *string { s := "http"; return &s }(),
            Hostname: func() *gatewayapiv1beta1.PreciseHostname { s := gatewayapiv1beta1.PreciseHostname("host1"); return &s }(),
        },
        URLRewrite: &gatewayapiv1beta1.HTTPURLRewriteFilter{
            Hostname: func() *gatewayapiv1beta1.PreciseHostname { s := gatewayapiv1beta1.PreciseHostname("rewrite-host"); return &s }(),
            Path: &gatewayapiv1beta1.HTTPPathModifier{
                ReplaceFullPath: func() *string { s := "/full"; return &s }(),
            },
        },
        ExtensionRef: &gatewayapiv1beta1.LocalObjectReference{Name: "ext-ref"},
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for HTTPRouteFilter")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy did not produce an equal HTTPRouteFilter")
    }
    // Mutate a nested field in the original
    orig.RequestHeaderModifier.Set[0].Value = "modified-val"
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for HTTPRouteFilter is not independent from the original")
    }
}

// TestDeepCopy_HTTPPathMatch tests deep copy for HTTPPathMatch in gatewayapi v1beta1 package
func TestDeepCopy_HTTPPathMatch(t *testing.T) {
    // Setup original HTTPPathMatch with Type and Value
    matchType := gatewayapiv1beta1.PathMatchType("Exact")
    val := "/home"
    orig := &gatewayapiv1beta1.HTTPPathMatch{
        Type:  &matchType,
        Value: &val,
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for HTTPPathMatch")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy did not produce an equal HTTPPathMatch")
    }
    // Mutate original and ensure the copy remains unchanged
    newVal := "/modified"
    orig.Value = &newVal
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for HTTPPathMatch is not independent from the original")
    }
}
// TestDeepCopy_BackendRef tests deep copy functions on BackendRef.
func TestDeepCopy_BackendRef(t *testing.T) {
    // Create a sample BackendRef with a pointer field Weight.
    weight := int32(100)
    orig := &gatewayapiv1beta1.BackendRef{
        BackendObjectReference: gatewayapiv1beta1.BackendObjectReference{
            Name: "backend",
        },
        Weight: &weight,
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for BackendRef")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for BackendRef did not produce an equal object")
    }
    // Modify the original and check that the copy stays unchanged.
    *orig.Weight = 200
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for BackendRef is not independent from the original")
    }
}

// TestDeepCopy_HTTPHeaderMatch tests deep copy functions on HTTPHeaderMatch.
func TestDeepCopy_HTTPHeaderMatch(t *testing.T) {
    matchType := gatewayapiv1beta1.HeaderMatchType("Exact")
    orig := &gatewayapiv1beta1.HTTPHeaderMatch{
        Name:  "X-Test",
        Value: "value1",
        Type:  &matchType,
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for HTTPHeaderMatch")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for HTTPHeaderMatch did not produce an equal object")
    }
    // Change the original and ensure the copy remains unchanged.
    orig.Value = "modified"
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for HTTPHeaderMatch is not independent from the original")
    }
}

// TestDeepCopy_HTTPPathModifier tests deep copy functions on HTTPPathModifier.
func TestDeepCopy_HTTPPathModifier(t *testing.T) {
    orig := &gatewayapiv1beta1.HTTPPathModifier{
        ReplaceFullPath:   func() *string { s := "/full"; return &s }(),
        ReplacePrefixMatch: func() *string { s := "/prefix"; return &s }(),
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for HTTPPathModifier")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for HTTPPathModifier did not produce an equal object")
    }
    // Modify the original.
    orig.ReplaceFullPath = func() *string { s := "/modified"; return &s }()
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for HTTPPathModifier is not independent from the original")
    }
}

// TestDeepCopy_HTTPQueryParamMatch tests deep copy functions on HTTPQueryParamMatch.
func TestDeepCopy_HTTPQueryParamMatch(t *testing.T) {
    matchType := gatewayapiv1beta1.QueryParamMatchType("Exact")
    orig := &gatewayapiv1beta1.HTTPQueryParamMatch{
        Name:  "param",
        Value: "1",
        Type:  &matchType,
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for HTTPQueryParamMatch")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for HTTPQueryParamMatch did not produce an equal object")
    }
    // Modify the original.
    orig.Value = "2"
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for HTTPQueryParamMatch is not independent from the original")
    }
}

// TestDeepCopy_HTTPRequestMirrorFilter tests deep copy functions on HTTPRequestMirrorFilter.
func TestDeepCopy_HTTPRequestMirrorFilter(t *testing.T) {
    orig := &gatewayapiv1beta1.HTTPRequestMirrorFilter{
        BackendRef: gatewayapiv1beta1.BackendRef{
            BackendObjectReference: gatewayapiv1beta1.BackendObjectReference{
                Name: "mirror-backend",
            },
        },
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for HTTPRequestMirrorFilter")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for HTTPRequestMirrorFilter did not produce an equal object")
    }
    // Modify the original.
    orig.BackendRef.Name = "modified-backend"
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for HTTPRequestMirrorFilter is not independent from the original")
    }
}

// TestDeepCopy_HTTPRequestRedirectFilter tests deep copy functions on HTTPRequestRedirectFilter.
func TestDeepCopy_HTTPRequestRedirectFilter(t *testing.T) {
    scheme := "http"
    hostname := gatewayapiv1beta1.PreciseHostname("host")
    port := gatewayapiv1beta1.PortNumber(80)
    orig := &gatewayapiv1beta1.HTTPRequestRedirectFilter{
        Scheme:   &scheme,
        Hostname: &hostname,
        Port:     &port,
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for HTTPRequestRedirectFilter")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for HTTPRequestRedirectFilter did not produce an equal object")
    }
    // Modify the original.
    newScheme := "https"
    orig.Scheme = &newScheme
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for HTTPRequestRedirectFilter is not independent from the original")
    }
}

// TestDeepCopy_ParentReference tests deep copy functions on ParentReference.
func TestDeepCopy_ParentReference(t *testing.T) {
    sectionName := "section"
    port := gatewayapiv1beta1.PortNumber(8080)
    orig := &gatewayapiv1beta1.ParentReference{
        Group:       func() *gatewayapiv1beta1.Group { g := gatewayapiv1beta1.Group("group"); return &g }(),
        Kind:        func() *gatewayapiv1beta1.Kind { k := gatewayapiv1beta1.Kind("kind"); return &k }(),
        Namespace:   func() *gatewayapiv1beta1.Namespace { n := gatewayapiv1beta1.Namespace("default"); return &n }(),
        SectionName: &sectionName,
        Port:        &port,
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for ParentReference")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for ParentReference did not produce an equal object")
    }
    // Modify the original.
    newSection := "modified-section"
    orig.SectionName = &newSection
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for ParentReference is not independent from the original")
    }
}

// TestDeepCopy_ReferenceGrant tests deep copy functions on ReferenceGrant.
func TestDeepCopy_ReferenceGrant(t *testing.T) {
    orig := &gatewayapiv1beta1.ReferenceGrant{
        ObjectMeta: metav1.ObjectMeta{
            Name: "refgrant",
        },
        Spec: gatewayapiv1beta1.ReferenceGrantSpec{
            From: []gatewayapiv1beta1.ReferenceGrantFrom{
                {},
            },
            To: []gatewayapiv1beta1.ReferenceGrantTo{
                {
                    Name: func() *gatewayapiv1beta1.ObjectName { n := gatewayapiv1beta1.ObjectName("obj"); return &n }(),
                },
            },
        },
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for ReferenceGrant")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for ReferenceGrant did not produce an equal object")
    }
    // Modify the original.
    orig.Spec.To[0].Name = func() *gatewayapiv1beta1.ObjectName { n := gatewayapiv1beta1.ObjectName("modified"); return &n }()
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for ReferenceGrant is not independent from the original")
    }
}

// TestDeepCopy_RouteParentStatus tests deep copy functions on RouteParentStatus.
func TestDeepCopy_RouteParentStatus(t *testing.T) {
    now := metav1.NewTime(time.Now())
    orig := &gatewayapiv1beta1.RouteParentStatus{
        ParentRef: gatewayapiv1beta1.ParentReference{
            Group: func() *gatewayapiv1beta1.Group { g := gatewayapiv1beta1.Group("group"); return &g }(),
        },
        Conditions: []metav1.Condition{
            {
                Type:   "Ready",
                Status: "True",
            },
        },
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for RouteParentStatus")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for RouteParentStatus did not produce an equal object")
    }
    // Modify the original.
    orig.Conditions[0].Status = "False"
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for RouteParentStatus is not independent from the original")
    }
}

// TestDeepCopy_SecretObjectReference tests deep copy functions on SecretObjectReference.
func TestDeepCopy_SecretObjectReference(t *testing.T) {
    orig := &gatewayapiv1beta1.SecretObjectReference{
        Name: "secret",
        Group: func() *gatewayapiv1beta1.Group {
            g := gatewayapiv1beta1.Group("grp")
            return &g
        }(),
        Kind: func() *gatewayapiv1beta1.Kind {
            k := gatewayapiv1beta1.Kind("kind")
            return &k
        }(),
        Namespace: func() *gatewayapiv1beta1.Namespace {
            n := gatewayapiv1beta1.Namespace("ns")
            return &n
        }(),
    }
    cp := orig.DeepCopy()
    if cp == nil {
        t.Fatalf("DeepCopy returned nil for SecretObjectReference")
    }
    if !reflect.DeepEqual(orig, cp) {
        t.Error("DeepCopy for SecretObjectReference did not produce an equal object")
    }
    // Modify the original.
    newGroup := gatewayapiv1beta1.Group("modified-grp")
    orig.Group = &newGroup
    if reflect.DeepEqual(orig, cp) {
        t.Error("Deep copy for SecretObjectReference is not independent from the original")
    }
}