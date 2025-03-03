package v1beta1

import (
    "reflect"
    "testing"
    "time"

    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
    runtime "k8s.io/apimachinery/pkg/runtime"
    istiov1beta1 "github.com/fluxcd/flagger/pkg/apis/istio/v1beta1"
    gatewayapiv1 "github.com/fluxcd/flagger/pkg/apis/gatewayapi/v1"
    gatewayapiv1beta1 "github.com/fluxcd/flagger/pkg/apis/gatewayapi/v1beta1"
)

// TestDeepCopyFunctions verifies that all DeepCopy and DeepCopyObject methods create deep copies.
func TestDeepCopyFunctions(t *testing.T) {
    // Prepare a common timestamp for testing.
    now := time.Now()

    // Test AlertProvider DeepCopy
    alertProvider := &AlertProvider{
    TypeMeta: metav1.TypeMeta{
    Kind:       "AlertProvider",
    APIVersion: "flagger.app/v1beta1",
    },
    ObjectMeta: metav1.ObjectMeta{
    Name:      "test-alertprovider",
    Namespace: "default",
    },
    Spec: AlertProviderSpec{
    SecretRef: &corev1.LocalObjectReference{
    Name: "secret1",
    },
    },
    Status: AlertProviderStatus{
    Conditions: []AlertProviderCondition{
    {
        LastUpdateTime:     metav1.NewTime(now),
        LastTransitionTime: metav1.NewTime(now),
    },
    },
    },
    }

    copiedAlertProvider := alertProvider.DeepCopy()
    if copiedAlertProvider == alertProvider {
    t.Error("DeepCopy did not create a new object for AlertProvider")
    }
    if !reflect.DeepEqual(alertProvider, copiedAlertProvider) {
    t.Error("Deep copied AlertProvider is not equal to the original")
    }

    // Modify the copy and verify that it does not affect the original.
    copiedAlertProvider.Spec.SecretRef.Name = "modified-secret"
    if alertProvider.Spec.SecretRef.Name == copiedAlertProvider.Spec.SecretRef.Name {
    t.Error("Deep copy is not deep: modification affected the original AlertProvider")
    }

    // Verify DeepCopy on a nil AlertProvider returns nil.
    var nilAlertProvider *AlertProvider
    if nilAlertProvider.DeepCopy() != nil {
    t.Error("DeepCopy on nil AlertProvider should return nil")
    }

    // Test Canary DeepCopy and DeepCopyObject
    canary := &Canary{
    TypeMeta: metav1.TypeMeta{
    Kind:       "Canary",
    APIVersion: "flagger.app/v1beta1",
    },
    ObjectMeta: metav1.ObjectMeta{
    Name:      "test-canary",
    Namespace: "default",
    },
    Spec: CanarySpec{
    TargetRef: LocalObjectReference{
    Name: "target-svc",
    },
    Service: CanaryService{
    TargetPort: intstr.FromInt(8080),
    Gateways:   []string{"gateway1"},
    },
    },
    Status: CanaryStatus{
    TrackedConfigs: new(map[string]string),
    Conditions: []CanaryCondition{
    {
        LastUpdateTime:     metav1.NewTime(now),
        LastTransitionTime: metav1.NewTime(now),
    },
    },
    },
    }
    *canary.Status.TrackedConfigs = map[string]string{"config": "value"}

    copiedCanary := canary.DeepCopy()
    if copiedCanary == canary {
    t.Error("DeepCopy did not create a new object for Canary")
    }
    if !reflect.DeepEqual(canary, copiedCanary) {
    t.Error("Deep copied Canary is not equal to the original")
    }

    // Test the runtime.Object interface through DeepCopyObject.
    var obj runtime.Object = canary
    copiedObj := obj.DeepCopyObject()
    if copiedObj == nil {
    t.Error("DeepCopyObject returned nil for Canary")
    }
    canaryFromObj, ok := copiedObj.(*Canary)
    if !ok {
    t.Error("DeepCopyObject did not return a *Canary type")
    }
    if !reflect.DeepEqual(canary, canaryFromObj) {
    t.Error("Deep copied Canary from DeepCopyObject is not equal to the original")
    }

    // Test MetricTemplate DeepCopy
    metricTemplate := &MetricTemplate{
    TypeMeta: metav1.TypeMeta{
    Kind:       "MetricTemplate",
    APIVersion: "flagger.app/v1beta1",
    },
    ObjectMeta: metav1.ObjectMeta{
    Name:      "test-metrictemplate",
    Namespace: "default",
    },
    Spec: MetricTemplateSpec{
    Provider: MetricTemplateProvider{
    SecretRef: &corev1.LocalObjectReference{Name: "secret-metric"},
    },
    },
    Status: MetricTemplateStatus{
    Conditions: []MetricTemplateCondition{
    {
        LastUpdateTime:     metav1.NewTime(now),
        LastTransitionTime: metav1.NewTime(now),
    },
    },
    },
    }

    copiedMetricTemplate := metricTemplate.DeepCopy()
    if copiedMetricTemplate == metricTemplate {
    t.Error("DeepCopy did not create a new object for MetricTemplate")
    }
    if !reflect.DeepEqual(metricTemplate, copiedMetricTemplate) {
    t.Error("Deep copied MetricTemplate is not equal to the original")
    }

    // Verify DeepCopy on a nil MetricTemplate returns nil.
    var nilMetricTemplate *MetricTemplate
    if nilMetricTemplate.DeepCopy() != nil {
    t.Error("DeepCopy on nil MetricTemplate should return nil")
    }
    // Test DeepCopy for AlertProviderCondition: verify deep copy independence of embedded metav1.Time fields.
    t.Run("DeepCopy AlertProviderCondition", func(t *testing.T) {
        condition := AlertProviderCondition{
            LastUpdateTime:     metav1.NewTime(now),
            LastTransitionTime: metav1.NewTime(now),
        }
        copiedCondition := condition.DeepCopy()
        if copiedCondition == &condition {
            t.Error("DeepCopy did not create a new object for AlertProviderCondition")
        }
        if !reflect.DeepEqual(condition, *copiedCondition) {
            t.Error("Deep copied AlertProviderCondition is not equal to the original")
        }
        // Modify the copy and verify that it does not affect the original.
        copiedCondition.LastUpdateTime.Time = copiedCondition.LastUpdateTime.Add(1 * time.Hour)
        if condition.LastUpdateTime.Time.Equal(copiedCondition.LastUpdateTime.Time) {
            t.Error("Deep copy is not deep: modification affected the original AlertProviderCondition")
        }
    })
    // Test DeepCopy for AutoscalerRefernce: verify map and sub‐object deep copy.
    t.Run("DeepCopy AutoscalerRefernce", func(t *testing.T) {
        replicas := &ScalerReplicas{
            MinReplicas: new(int32),
            MaxReplicas: new(int32),
        }
        *replicas.MinReplicas = 1
        *replicas.MaxReplicas = 5
        original := &AutoscalerRefernce{
            PrimaryScalerQueries: map[string]string{"query": "value"},
            PrimaryScalerReplicas: replicas,
        }
        copied := original.DeepCopy()
        if copied == original {
            t.Error("DeepCopy did not create a new object for AutoscalerRefernce")
        }
        if !reflect.DeepEqual(original, copied) {
            t.Error("Deep copied AutoscalerRefernce is not equal to the original")
        }
        // Modify the copy and check independence of the inner map.
        copied.PrimaryScalerQueries["query"] = "modified"
        if original.PrimaryScalerQueries["query"] == "modified" {
            t.Error("Deep copy is not deep: modification affected the original AutoscalerRefernce map")
        }
    })
    // Test DeepCopy for CanaryAnalysis: check deep copying of slices, pointers and embedded structs.
    t.Run("DeepCopy CanaryAnalysis", func(t *testing.T) {
        analysis := &CanaryAnalysis{
            StepWeights: []int{10, 20, 30},
            PrimaryReadyThreshold: new(int),
            CanaryReadyThreshold: new(int),
                    Alerts: []CanaryAlert{{ProviderRef: CrossNamespaceObjectReference{Name: "test-provider"}}},
            Metrics: []CanaryMetric{
                {
                    TemplateVariables: map[string]string{"key": "val"},
                },
            },
            Webhooks: []CanaryWebhook{
                {
                    Metadata: &map[string]string{"hook": "value"},
                },
            },
            Match: []istiov1beta1.HTTPMatchRequest{},
            SessionAffinity: &SessionAffinity{},
        }
        *analysis.PrimaryReadyThreshold = 5
        *analysis.CanaryReadyThreshold = 7
        copiedAnalysis := analysis.DeepCopy()
        if copiedAnalysis == analysis {
            t.Error("DeepCopy did not create a new object for CanaryAnalysis")
        }
        if !reflect.DeepEqual(analysis, copiedAnalysis) {
            t.Error("Deep copied CanaryAnalysis is not equal to the original")
        }
        // Modify a slice element in the copy and check independence.
        copiedAnalysis.StepWeights[0] = 99
        if analysis.StepWeights[0] == 99 {
            t.Error("Deep copy is not deep: modification affected the original CanaryAnalysis")
        }
    })
    // Test DeepCopy for CanaryWebhook and CanaryWebhookPayload: verify metadata maps are correctly deep copied.
    t.Run("DeepCopy CanaryWebhook and CanaryWebhookPayload", func(t *testing.T) {
        webhook := &CanaryWebhook{
            Metadata: &map[string]string{"hook": "initial"},
        }
        copiedWebhook := webhook.DeepCopy()
        if copiedWebhook == webhook {
            t.Error("DeepCopy did not create a new object for CanaryWebhook")
        }
        if !reflect.DeepEqual(webhook, copiedWebhook) {
            t.Error("Deep copied CanaryWebhook is not equal to the original")
        }
        (*copiedWebhook.Metadata)["hook"] = "changed"
        if (*webhook.Metadata)["hook"] == "changed" {
            t.Error("Deep copy is not deep: modification affected the original CanaryWebhook")
        }

        payload := &CanaryWebhookPayload{
            Metadata: map[string]string{"payload": "data"},
        }
        copiedPayload := payload.DeepCopy()
        if copiedPayload == payload {
            t.Error("DeepCopy did not create a new object for CanaryWebhookPayload")
        }
        if !reflect.DeepEqual(payload, copiedPayload) {
            t.Error("Deep copied CanaryWebhookPayload is not equal to the original")
        }
        copiedPayload.Metadata["payload"] = "modified"
        if payload.Metadata["payload"] == "modified" {
            t.Error("Deep copy is not deep: modification affected the original CanaryWebhookPayload")
        }
    })
    // Test DeepCopy for CustomBackend: verify that BackendObjectReference and Filters are deep copied.
    t.Run("DeepCopy CustomBackend", func(t *testing.T) {
        backend := &CustomBackend{
            BackendObjectReference: &gatewayapiv1.BackendObjectReference{
                Name: "backend1",
            },
            Filters: []gatewayapiv1.HTTPRouteFilter{
                {Type: "TestFilter"},
            },
        }
        copiedBackend := backend.DeepCopy()
        if copiedBackend == backend {
            t.Error("DeepCopy did not create a new object for CustomBackend")
        }
        if !reflect.DeepEqual(backend, copiedBackend) {
            t.Error("Deep copied CustomBackend is not equal to the original")
        }
        // Modify the copy and verify that changes do not affect the original.
        copiedBackend.BackendObjectReference.Name = "modified-backend"
        if backend.BackendObjectReference.Name == "modified-backend" {
            t.Error("Deep copy is not deep: modification affected the original CustomBackend")
        }
    })
    // Test DeepCopy for SessionAffinity: trivial deep copy.
    t.Run("DeepCopy SessionAffinity", func(t *testing.T) {
        affinity := &SessionAffinity{}
        copiedAffinity := affinity.DeepCopy()
        if copiedAffinity == affinity {
            t.Error("DeepCopy did not create a new object for SessionAffinity")
        }
        if !reflect.DeepEqual(affinity, copiedAffinity) {
            t.Error("Deep copied SessionAffinity is not equal to the original")
        }
    })
    // Test DeepCopy for CustomMetadata: verify maps are deep copied.
    t.Run("DeepCopy CustomMetadata", func(t *testing.T) {
        meta := &CustomMetadata{
            Labels: map[string]string{"app": "demo"},
            Annotations: map[string]string{"anno": "value"},
        }
        copiedMeta := meta.DeepCopy()
        if copiedMeta == meta {
            t.Error("DeepCopy did not create a new object for CustomMetadata")
        }
        if !reflect.DeepEqual(meta, copiedMeta) {
            t.Error("Deep copied CustomMetadata is not equal to the original")
        }
        // Modify the copy and check that the original remains unchanged.
        copiedMeta.Labels["app"] = "modified"
        if meta.Labels["app"] == "modified" {
            t.Error("Deep copy is not deep: modification affected the original CustomMetadata")
        }
    })
}
// TestDeepCopySimpleObjects tests deep copy methods for simple types that were not already covered.
func TestDeepCopySimpleObjects(t *testing.T) {
    // Test CrossNamespaceObjectReference DeepCopy
    cns := &CrossNamespaceObjectReference{
        Name: "cross-namespace",
    }
    copiedCNS := cns.DeepCopy()
    if copiedCNS == cns {
        t.Error("DeepCopy did not create a new object for CrossNamespaceObjectReference")
    }
    if !reflect.DeepEqual(cns, copiedCNS) {
        t.Error("Deep copied CrossNamespaceObjectReference is not equal to the original")
    }

    // Test HTTPRewrite DeepCopy
    rewrite := &HTTPRewrite{}
    copiedRewrite := rewrite.DeepCopy()
    if copiedRewrite == rewrite {
        t.Error("DeepCopy did not create a new object for HTTPRewrite")
    }
    if !reflect.DeepEqual(rewrite, copiedRewrite) {
        t.Error("Deep copied HTTPRewrite is not equal to the original")
    }

    // Test LocalObjectReference DeepCopy
    local := &LocalObjectReference{Name: "local-ref"}
    copiedLocal := local.DeepCopy()
    if copiedLocal == local {
        t.Error("DeepCopy did not create a new object for LocalObjectReference")
    }
    if !reflect.DeepEqual(local, copiedLocal) {
        t.Error("Deep copied LocalObjectReference is not equal to the original")
    }

    // Test MetricTemplateModel DeepCopy
    mtm := &MetricTemplateModel{
        Variables: map[string]string{"key": "value"},
    }
    copiedMTM := mtm.DeepCopy()
    if copiedMTM == mtm {
        t.Error("DeepCopy did not create a new object for MetricTemplateModel")
    }
    if !reflect.DeepEqual(mtm, copiedMTM) {
        t.Error("Deep copied MetricTemplateModel is not equal to the original")
    }
    copiedMTM.Variables["key"] = "changed"
    if mtm.Variables["key"] == "changed" {
        t.Error("Deep copy is not deep for MetricTemplateModel: modification affected the original")
    }

    // Test MetricTemplateProvider DeepCopy
    mtp := &MetricTemplateProvider{
        SecretRef: &corev1.LocalObjectReference{Name: "secret1"},
    }
    copiedMTP := mtp.DeepCopy()
    if copiedMTP == mtp {
        t.Error("DeepCopy did not create a new object for MetricTemplateProvider")
    }
    if !reflect.DeepEqual(mtp, copiedMTP) {
        t.Error("Deep copied MetricTemplateProvider is not equal to the original")
    }
    copiedMTP.SecretRef.Name = "secret2"
    if mtp.SecretRef.Name == "secret2" {
        t.Error("Deep copy is not deep for MetricTemplateProvider: modification affected the original")
    }

    // Test MetricTemplateSpec DeepCopy
    mts := &MetricTemplateSpec{
        Provider: MetricTemplateProvider{
            SecretRef: &corev1.LocalObjectReference{Name: "secret1"},
        },
    }
    copiedMTS := mts.DeepCopy()
    if copiedMTS == mts {
        t.Error("DeepCopy did not create a new object for MetricTemplateSpec")
    }
    if !reflect.DeepEqual(mts, copiedMTS) {
        t.Error("Deep copied MetricTemplateSpec is not equal to the original")
    }
    copiedMTS.Provider.SecretRef.Name = "changed"
    if mts.Provider.SecretRef.Name == "changed" {
        t.Error("Deep copy is not deep for MetricTemplateSpec: modification affected the original")
    }
}
// TestAdditionalDeepCopies increases test coverage by verifying DeepCopy functions for list and other types.
func TestAdditionalDeepCopies(t *testing.T) {
    now := time.Now()

    // Test AlertProviderList DeepCopy
    alert := AlertProvider{
        TypeMeta: metav1.TypeMeta{Kind: "AlertProvider", APIVersion: "flagger.app/v1beta1"},
        ObjectMeta: metav1.ObjectMeta{Name: "alert-1", Namespace: "default"},
        Spec: AlertProviderSpec{
            SecretRef: &corev1.LocalObjectReference{Name: "secret-alert"},
        },
        Status: AlertProviderStatus{
            Conditions: []AlertProviderCondition{
                {
                    LastUpdateTime:     metav1.NewTime(now),
                    LastTransitionTime: metav1.NewTime(now),
                },
            },
        },
    }
    apList := &AlertProviderList{
        Items: []AlertProvider{alert},
    }
    copiedAPList := apList.DeepCopy()
    if copiedAPList == apList {
        t.Error("DeepCopy did not create a new object for AlertProviderList")
    }
    if !reflect.DeepEqual(apList, copiedAPList) {
        t.Error("Deep copied AlertProviderList is not equal to the original")
    }
    // Modify the copy and ensure original remains unchanged.
    if len(copiedAPList.Items) > 0 {
        copiedAPList.Items[0].Spec.SecretRef.Name = "modified-secret"
        if apList.Items[0].Spec.SecretRef.Name == "modified-secret" {
            t.Error("Deep copy is not deep for AlertProviderList: modification affected the original")
        }
    }

    // Test CanaryList DeepCopy
    canary := Canary{
        TypeMeta: metav1.TypeMeta{Kind: "Canary", APIVersion: "flagger.app/v1beta1"},
        ObjectMeta: metav1.ObjectMeta{Name: "canary-1", Namespace: "default"},
        Spec: CanarySpec{
            TargetRef: LocalObjectReference{Name: "target-svc"},
            Service: CanaryService{
                TargetPort: intstr.FromInt(8080),
                Gateways:   []string{"gateway1"},
            },
        },
        Status: CanaryStatus{
            TrackedConfigs: new(map[string]string),
            Conditions: []CanaryCondition{
                {
                    LastUpdateTime:     metav1.NewTime(now),
                    LastTransitionTime: metav1.NewTime(now),
                },
            },
        },
    }
    *canary.Status.TrackedConfigs = map[string]string{"config": "value"}
    canaryList := &CanaryList{
        Items: []Canary{canary},
    }
    copiedCanaryList := canaryList.DeepCopy()
    if copiedCanaryList == canaryList {
        t.Error("DeepCopy did not create a new object for CanaryList")
    }
    if !reflect.DeepEqual(canaryList, copiedCanaryList) {
        t.Error("Deep copied CanaryList is not equal to the original")
    }

    // Test MetricTemplateList DeepCopy
    metricTemplate := MetricTemplate{
        TypeMeta: metav1.TypeMeta{Kind: "MetricTemplate", APIVersion: "flagger.app/v1beta1"},
        ObjectMeta: metav1.ObjectMeta{Name: "metric-1", Namespace: "default"},
        Spec: MetricTemplateSpec{
            Provider: MetricTemplateProvider{
                SecretRef: &corev1.LocalObjectReference{Name: "secret-metric"},
            },
        },
        Status: MetricTemplateStatus{
            Conditions: []MetricTemplateCondition{
                {
                    LastUpdateTime:     metav1.NewTime(now),
                    LastTransitionTime: metav1.NewTime(now),
                },
            },
        },
    }
    mtList := &MetricTemplateList{
        Items: []MetricTemplate{metricTemplate},
    }
    copiedMTList := mtList.DeepCopy()
    if copiedMTList == mtList {
        t.Error("DeepCopy did not create a new object for MetricTemplateList")
    }
    if !reflect.DeepEqual(mtList, copiedMTList) {
        t.Error("Deep copied MetricTemplateList is not equal to the original")
    }

    // Test CanaryAlert DeepCopy
    alertObj := CanaryAlert{
        ProviderRef: CrossNamespaceObjectReference{Name: "provider1"},
    }
    copiedAlertObj := alertObj.DeepCopy()
    if copiedAlertObj == &alertObj {
        t.Error("DeepCopy did not create a new object for CanaryAlert")
    }
    if !reflect.DeepEqual(alertObj, *copiedAlertObj) {
        t.Error("Deep copied CanaryAlert is not equal to the original")
    }
    // Modify the copy and check that the original is not affected.
    copiedAlertObj.ProviderRef.Name = "modified-provider"
    if alertObj.ProviderRef.Name == "modified-provider" {
        t.Error("Deep copy is not deep for CanaryAlert: modification affected the original")
    }

    // Test CanaryCondition DeepCopy
    condition := CanaryCondition{
        LastUpdateTime:     metav1.NewTime(now),
        LastTransitionTime: metav1.NewTime(now),
    }
    copiedCondition := condition.DeepCopy()
    if copiedCondition == &condition {
        t.Error("DeepCopy did not create a new object for CanaryCondition")
    }
    if !reflect.DeepEqual(condition, *copiedCondition) {
        t.Error("Deep copied CanaryCondition is not equal to the original")
    }
    copiedCondition.LastUpdateTime.Time = copiedCondition.LastUpdateTime.Add(2 * time.Hour)
    if condition.LastUpdateTime.Time.Equal(copiedCondition.LastUpdateTime.Time) {
        t.Error("Deep copy is not deep for CanaryCondition: modification affected the original")
    }

    // Test MetricTemplateCondition DeepCopy
    mtCondition := MetricTemplateCondition{
        LastUpdateTime:     metav1.NewTime(now),
        LastTransitionTime: metav1.NewTime(now),
    }
    copiedMTCondition := mtCondition.DeepCopy()
    if copiedMTCondition == &mtCondition {
        t.Error("DeepCopy did not create a new object for MetricTemplateCondition")
    }
    if !reflect.DeepEqual(mtCondition, *copiedMTCondition) {
        t.Error("Deep copied MetricTemplateCondition is not equal to the original")
    }
    copiedMTCondition.LastTransitionTime.Time = copiedMTCondition.LastTransitionTime.Add(30 * time.Minute)
    if mtCondition.LastTransitionTime.Time.Equal(copiedMTCondition.LastTransitionTime.Time) {
        t.Error("Deep copy is not deep for MetricTemplateCondition: modification affected the original")
    }

    // Test ScalerReplicas DeepCopy
    minVal := int32(2)
    maxVal := int32(10)
    replicas := &ScalerReplicas{
        MinReplicas: &minVal,
        MaxReplicas: &maxVal,
    }
    copiedReplicas := replicas.DeepCopy()
    if copiedReplicas == replicas {
        t.Error("DeepCopy did not create a new object for ScalerReplicas")
    }
    if !reflect.DeepEqual(replicas, copiedReplicas) {
        t.Error("Deep copied ScalerReplicas is not equal to the original")
    }
    *copiedReplicas.MinReplicas = 5
    if *replicas.MinReplicas == 5 {
        t.Error("Deep copy is not deep for ScalerReplicas: modification affected the original")
    }

    // Test CanaryService DeepCopy
    canaryService := CanaryService{
        TargetPort: intstr.FromInt(9090),
        Gateways:   []string{"gatewayA"},
        GatewayRefs: []gatewayapiv1beta1.ParentReference{
            {Name: "parent1"},
        },
        Hosts: []string{"hostA"},
        TrafficPolicy: &istiov1beta1.TrafficPolicy{},
        Match: []istiov1beta1.HTTPMatchRequest{},
        Rewrite:      &HTTPRewrite{},
        Retries:      &istiov1beta1.HTTPRetry{},
        Headers:      &istiov1beta1.Headers{},
        Mirror:       []gatewayapiv1beta1.HTTPRequestMirrorFilter{{}},
        CorsPolicy:   &istiov1beta1.CorsPolicy{},
        Backends:     []string{"backendA"},
        Apex:         &CustomMetadata{Labels: map[string]string{"key": "value"}},
        Primary:      &CustomMetadata{Annotations: map[string]string{"anno": "orig"}},
        Canary:       &CustomMetadata{Labels: map[string]string{"app": "canary"}},
        PrimaryBackend: &CustomBackend{
            BackendObjectReference: &gatewayapiv1.BackendObjectReference{Name: "backendRef"},
        },
        CanaryBackend: &CustomBackend{
            BackendObjectReference: &gatewayapiv1.BackendObjectReference{Name: "backendRef2"},
        },
    }
    copiedCanaryService := canaryService.DeepCopy()
    if copiedCanaryService == &canaryService {
        t.Error("DeepCopy did not create a new object for CanaryService")
    }
    if !reflect.DeepEqual(canaryService, *copiedCanaryService) {
        t.Error("Deep copied CanaryService is not equal to the original")
    }
    copiedCanaryService.Gateways[0] = "modified-gateway"
    if canaryService.Gateways[0] == "modified-gateway" {
        t.Error("Deep copy is not deep for CanaryService: modification affected the original")
    }
}
// TestDeepCopyCanaryMetric tests that the DeepCopy method for CanaryMetric
// creates a deep copy that is independent from the original.
func TestDeepCopyCanaryMetric(t *testing.T) {
    minVal := 0.5
    maxVal := 1.5
    orig := &CanaryMetric{
        TemplateVariables: map[string]string{"key": "initial"},
        TemplateRef: &CrossNamespaceObjectReference{Name: "original-ref"},
        ThresholdRange: &CanaryThresholdRange{
            Min: &minVal,
            Max: &maxVal,
        },
    }
    copied := orig.DeepCopy()
    if copied == orig {
        t.Error("DeepCopy did not create a new object for CanaryMetric")
    }
    if !reflect.DeepEqual(orig, copied) {
        t.Error("Deep copied CanaryMetric is not equal to the original")
    }
    // Modify the copy and make sure the original is not changed.
    copied.TemplateVariables["key"] = "changed"
    copied.TemplateRef.Name = "changed-ref"
    *copied.ThresholdRange.Min = 2.0
    if orig.TemplateVariables["key"] == "changed" {
        t.Error("Deep copy is not deep for CanaryMetric: modification affected original TemplateVariables")
    }
    if orig.TemplateRef.Name == "changed-ref" {
        t.Error("Deep copy is not deep for CanaryMetric: modification affected original TemplateRef")
    }
    if *orig.ThresholdRange.Min == 2.0 {
        t.Error("Deep copy is not deep for CanaryMetric: modification affected original ThresholdRange")
    }
}

// TestDeepCopyCanaryThresholdRange tests that the DeepCopy for CanaryThresholdRange
// creates a deep copy that is independent of the original pointer values.
func TestDeepCopyCanaryThresholdRange(t *testing.T) {
    minVal := 0.3
    maxVal := 1.2
    orig := &CanaryThresholdRange{
        Min: &minVal,
        Max: &maxVal,
    }
    copied := orig.DeepCopy()
    if copied == orig {
        t.Error("DeepCopy did not create a new object for CanaryThresholdRange")
    }
    if !reflect.DeepEqual(orig, copied) {
        t.Error("Deep copied CanaryThresholdRange is not equal to the original")
    }
    // Modify the copied values and expect the original remains unaffected.
    *copied.Min = 0.7
    *copied.Max = 2.3
    if *orig.Min == 0.7 || *orig.Max == 2.3 {
        t.Error("Deep copy is not deep for CanaryThresholdRange: modification affected original")
    }
}
// TestDeepCopyNilFields tests deep copies for objects with nil or missing optional fields.
func TestDeepCopyNilFields(t *testing.T) {
    // Test AlertProviderSpec: when SecretRef is nil the deep copy should have nil as well.
    spec := AlertProviderSpec{}
    specCopy := spec.DeepCopy()
    if specCopy.SecretRef != nil {
        t.Error("Expected nil SecretRef in AlertProviderSpec deep copy")
    }

    // Test AutoscalerRefernce: when PrimaryScalerQueries and PrimaryScalerReplicas are nil.
    autoRef := &AutoscalerRefernce{
        PrimaryScalerQueries: nil,
        PrimaryScalerReplicas: nil,
    }
    autoRefCopy := autoRef.DeepCopy()
    if autoRefCopy.PrimaryScalerQueries != nil {
        t.Error("Expected nil PrimaryScalerQueries in AutoscalerRefernce deep copy")
    }
    if autoRefCopy.PrimaryScalerReplicas != nil {
        t.Error("Expected nil PrimaryScalerReplicas in AutoscalerRefernce deep copy")
    }

    // Test CanarySpec: optional fields IngressRef, RouteRef, UpstreamRef should be nil.
    cs := &CanarySpec{
        TargetRef: LocalObjectReference{Name: "default-target"},
        Service:   CanaryService{}, // default empty service struct
    }
    csCopy := cs.DeepCopy()
    if csCopy.IngressRef != nil || csCopy.RouteRef != nil || csCopy.UpstreamRef != nil {
        t.Error("Expected nil IngressRef, RouteRef, and UpstreamRef in CanarySpec deep copy")
    }

    // Test CustomBackend: when BackendObjectReference and Filters are nil.
    cb := &CustomBackend{
        BackendObjectReference: nil,
        Filters:                nil,
    }
    cbCopy := cb.DeepCopy()
    if cbCopy.BackendObjectReference != nil {
        t.Error("Expected nil BackendObjectReference in CustomBackend deep copy")
    }
    if cbCopy.Filters != nil {
        t.Error("Expected nil Filters in CustomBackend deep copy")
    }
}
// TestChainedDeepCopies verifies that chained DeepCopy calls produce objects completely independent from the original.
func TestChainedDeepCopies(t *testing.T) {
    // Create an original Canary object with some nested values.
    original := &Canary{
        TypeMeta: metav1.TypeMeta{
            Kind:       "Canary",
            APIVersion: "flagger.app/v1beta1",
        },
        ObjectMeta: metav1.ObjectMeta{
            Name:      "chain-canary",
            Namespace: "default",
        },
        Spec: CanarySpec{
            TargetRef: LocalObjectReference{
                Name: "target-chain",
            },
            Service: CanaryService{
                TargetPort: intstr.FromInt(8081),
                Gateways:   []string{"chain-gateway"},
            },
        },
        Status: CanaryStatus{
            TrackedConfigs: new(map[string]string),
            Conditions: []CanaryCondition{
                {
                    LastUpdateTime:     metav1.Now(),
                    LastTransitionTime: metav1.Now(),
                },
            },
        },
    }
    *original.Status.TrackedConfigs = map[string]string{"chain": "original"}

    // Perform a first deep copy.
    copy1 := original.DeepCopy()

    // Chain a second deep copy from the first copy.
    copy2 := copy1.DeepCopy()

    // Modify fields in the second copy.
    copy2.ObjectMeta.Name = "chain-canary-copy2"
    copy2.Spec.Service.Gateways[0] = "modified-chain-gateway"
    (*copy2.Status.TrackedConfigs)["chain"] = "modified"

    // Verify that the original object remains unchanged.
    if original.ObjectMeta.Name == copy2.ObjectMeta.Name {
        t.Error("Chained deep copy failed: modification in copy2 affected the original ObjectMeta.Name")
    }
    if original.Spec.Service.Gateways[0] == copy2.Spec.Service.Gateways[0] {
        t.Error("Chained deep copy failed: modification in copy2 affected the original Service.Gateways")
    }
    if (*original.Status.TrackedConfigs)["chain"] == "modified" {
        t.Error("Chained deep copy failed: modification in copy2 affected the original TrackedConfigs")
    }

    // Also verify that the first copy remains unchanged by the changes in copy2.
    if copy1.ObjectMeta.Name == copy2.ObjectMeta.Name {
        t.Error("Chained deep copy failed: modification in copy2 affected copy1 ObjectMeta.Name")
    }
    if copy1.Spec.Service.Gateways[0] == copy2.Spec.Service.Gateways[0] {
        t.Error("Chained deep copy failed: modification in copy2 affected copy1 Service.Gateways")
    }
    if (*copy1.Status.TrackedConfigs)["chain"] == "modified" {
        t.Error("Chained deep copy failed: modification in copy2 affected copy1 TrackedConfigs")
    }
}
// TestDeepCopyEdgeCases tests that deep copy functions correctly handle nil and empty optional fields.
func TestDeepCopyEdgeCases(t *testing.T) {
    // Test AlertProviderSpec with nil SecretRef remains nil after deep copy.
    spec := AlertProviderSpec{SecretRef: nil}
    specCopy := spec.DeepCopy()
    if specCopy.SecretRef != nil {
        t.Error("Expected nil SecretRef in deep-copied AlertProviderSpec")
    }

    // Test AutoscalerRefernce with nil fields.
    autoRef := &AutoscalerRefernce{
        PrimaryScalerQueries: nil,
        PrimaryScalerReplicas: nil,
    }
    autoRefCopy := autoRef.DeepCopy()
    if autoRefCopy.PrimaryScalerQueries != nil {
        t.Error("Expected nil PrimaryScalerQueries in deep-copied AutoscalerRefernce")
    }
    if autoRefCopy.PrimaryScalerReplicas != nil {
        t.Error("Expected nil PrimaryScalerReplicas in deep-copied AutoscalerRefernce")
    }

    // Test CanarySpec where optional IngressRef, RouteRef, and UpstreamRef are nil.
    cs := &CanarySpec{
        TargetRef: LocalObjectReference{Name: "test-target"},
        Service:   CanaryService{},
    }
    csCopy := cs.DeepCopy()
    if csCopy.IngressRef != nil || csCopy.RouteRef != nil || csCopy.UpstreamRef != nil {
        t.Error("Expected nil IngressRef, RouteRef, and UpstreamRef in deep-copied CanarySpec")
    }

    // Test CustomBackend with nil BackendObjectReference and Filters.
    cb := &CustomBackend{
        BackendObjectReference: nil,
        Filters:                nil,
    }
    cbCopy := cb.DeepCopy()
    if cbCopy.BackendObjectReference != nil {
        t.Error("Expected nil BackendObjectReference in deep-copied CustomBackend")
    }
    if cbCopy.Filters != nil {
        t.Error("Expected nil Filters in deep-copied CustomBackend")
    }

    // Test MetricTemplate with empty Conditions and nil Spec.Provider.SecretRef.
    mt := &MetricTemplate{
        TypeMeta: metav1.TypeMeta{Kind: "MetricTemplate", APIVersion: "flagger.app/v1beta1"},
        ObjectMeta: metav1.ObjectMeta{Name: "empty-mt", Namespace: "default"},
        Spec: MetricTemplateSpec{
            Provider: MetricTemplateProvider{
                SecretRef: nil,
            },
        },
        Status: MetricTemplateStatus{
            Conditions: nil, // nil slice
        },
    }
    mtCopy := mt.DeepCopy()
    if mtCopy.Spec.Provider.SecretRef != nil {
        t.Error("Expected nil SecretRef in deep-copied MetricTemplate")
    }
    if mtCopy.Status.Conditions != nil {
        t.Error("Expected nil Conditions in deep-copied MetricTemplate")
    }
}
// TestDeepCopyRoundTrip creates a fully populated Canary, performs two rounds of deep copy, and verifies that modifications in the second copy do not affect the original.
func TestDeepCopyRoundTrip(t *testing.T) {
    now := time.Now()
    // Create a Canary with many non-nil fields populated.
    original := &Canary{
        TypeMeta: metav1.TypeMeta{Kind: "Canary", APIVersion: "flagger.app/v1beta1"},
        ObjectMeta: metav1.ObjectMeta{
            Name:      "round-trip",
            Namespace: "default",
        },
        Spec: CanarySpec{
            TargetRef: LocalObjectReference{Name: "svc"},
            IngressRef: &LocalObjectReference{Name: "ingress"},
            RouteRef:   &LocalObjectReference{Name: "route"},
            UpstreamRef: &CrossNamespaceObjectReference{Name: "upstream"},
            Service: CanaryService{
                TargetPort: intstr.FromInt(8080),
                Gateways:   []string{"gateway1"},
                GatewayRefs: []gatewayapiv1beta1.ParentReference{
                    {Name: "parent1"},
                },
            },
            Analysis: &CanaryAnalysis{
                StepWeights: []int{20, 40},
                Alerts:      []CanaryAlert{{ProviderRef: CrossNamespaceObjectReference{Name: "provider1"}}},
            },
        },
        Status: CanaryStatus{
            TrackedConfigs: new(map[string]string),
            Conditions: []CanaryCondition{
                {
                    LastUpdateTime:     metav1.NewTime(now),
                    LastTransitionTime: metav1.NewTime(now),
                },
            },
        },
    }
    *original.Status.TrackedConfigs = map[string]string{"key": "value"}

    // Perform first deep copy.
    copy1 := original.DeepCopy()
    // Chain a second deep copy from the first copy.
    copy2 := copy1.DeepCopy()

    // Modify the second copy's top-level ObjectMeta and nested fields.
    copy2.ObjectMeta.Name = "modified"
    copy2.Spec.Service.Gateways[0] = "modified-gateway"
    copy2.Spec.IngressRef.Name = "modified-ingress"

    // Verify that none of these modifications affected the original.
    if original.ObjectMeta.Name == "modified" {
        t.Error("Deep copy round trip: modification in copy2 affected original ObjectMeta.Name")
    }
    if original.Spec.Service.Gateways[0] == "modified-gateway" {
        t.Error("Deep copy round trip: modification in copy2 affected original Service.Gateways")
    }
    if original.Spec.IngressRef.Name == "modified-ingress" {
        t.Error("Deep copy round trip: modification in copy2 affected original IngressRef")
    }
}

// TestDeepCopyComplexNested creates a complex nested object (using Canary and AutoscalerRefernce) and ensures that changes in nested maps and pointer values do not propagate to the original.
func TestDeepCopyComplexNested(t *testing.T) {
    now := time.Now()
    replicas := &ScalerReplicas{
        MinReplicas: new(int32),
        MaxReplicas: new(int32),
    }
    *replicas.MinReplicas = 2
    *replicas.MaxReplicas = 8

    // Create a Canary that uses AutoscalerRefernce to hold deep nested map and pointer values.
    canary := &Canary{
        TypeMeta: metav1.TypeMeta{Kind: "Canary", APIVersion: "flagger.app/v1beta1"},
        ObjectMeta: metav1.ObjectMeta{
            Name:      "complex-canary",
            Namespace: "default",
        },
        Spec: CanarySpec{
            TargetRef: LocalObjectReference{Name: "svc"},
            AutoscalerRef: &AutoscalerRefernce{
                PrimaryScalerQueries: map[string]string{"query1": "value1"},
                PrimaryScalerReplicas: replicas,
            },
            Service: CanaryService{
                TargetPort:    intstr.FromInt(80),
                Gateways:      []string{"gw"},
                Hosts:         []string{"host"},
                TrafficPolicy: &istiov1beta1.TrafficPolicy{},
                Match:         []istiov1beta1.HTTPMatchRequest{{}},
                Rewrite:       &HTTPRewrite{},
                Retries:       &istiov1beta1.HTTPRetry{},
                Headers:       &istiov1beta1.Headers{},
                Mirror:        []gatewayapiv1beta1.HTTPRequestMirrorFilter{{}},
                CorsPolicy:    &istiov1beta1.CorsPolicy{},
                Backends:      []string{"backend"},
                Apex:          &CustomMetadata{Labels: map[string]string{"a": "1"}},
                Primary:       &CustomMetadata{Annotations: map[string]string{"p": "1"}},
                Canary:        &CustomMetadata{Labels: map[string]string{"c": "1"}},
                PrimaryBackend: &CustomBackend{
                    BackendObjectReference: &gatewayapiv1.BackendObjectReference{Name: "b1"},
                },
                CanaryBackend: &CustomBackend{
                    BackendObjectReference: &gatewayapiv1.BackendObjectReference{Name: "b2"},
                },
            },
        },
        Status: CanaryStatus{
            TrackedConfigs: new(map[string]string),
            Conditions: []CanaryCondition{
                {
                    LastUpdateTime:     metav1.NewTime(now),
                    LastTransitionTime: metav1.NewTime(now),
                },
            },
        },
    }
    *canary.Status.TrackedConfigs = map[string]string{"a": "b"}

    // Perform a deep copy of the complex object.
    copyCanary := canary.DeepCopy()

    // Modify nested fields in the copy.
    copyCanary.Spec.AutoscalerRef.PrimaryScalerQueries["query1"] = "modified"
    *copyCanary.Spec.AutoscalerRef.PrimaryScalerReplicas.MinReplicas = 10

    // Verify that modifications in nested maps and pointer values did not affect the original.
    if canary.Spec.AutoscalerRef.PrimaryScalerQueries["query1"] == "modified" {
        t.Error("Deep copy complex nested: modification affected original autoscaler queries")
    }
    if *canary.Spec.AutoscalerRef.PrimaryScalerReplicas.MinReplicas == 10 {
        t.Error("Deep copy complex nested: modification affected original scaler replicas")
    }
}
// TestDeepCopyCanaryAlertIndependence tests that the DeepCopy method for CanaryAlert produces an independent copy.
func TestDeepCopyCanaryAlertIndependence(t *testing.T) {
    orig := &CanaryAlert{
        ProviderRef: CrossNamespaceObjectReference{Name: "provider-original"},
    }
    copyAlert := orig.DeepCopy()
    if copyAlert == orig {
        t.Error("DeepCopy did not create a distinct CanaryAlert object")
    }
    if !reflect.DeepEqual(orig, copyAlert) {
        t.Error("Deep copied CanaryAlert is not equal to the original")
    }
    // Modify the copy and ensure the original is not affected.
    copyAlert.ProviderRef.Name = "provider-changed"
    if orig.ProviderRef.Name == "provider-changed" {
        t.Error("Modification on deep copied CanaryAlert affected the original")
    }
}

// TestDeepCopyNilReceiver verifies that calling DeepCopy on a nil receiver returns nil for various types.
func TestDeepCopyNilReceiver(t *testing.T) {
    var nilAlertProvider *AlertProvider
    if nilAlertProvider.DeepCopy() != nil {
        t.Error("Expected nil result for nil AlertProvider deep copy")
    }
    var nilCanary *Canary
    if nilCanary.DeepCopy() != nil {
        t.Error("Expected nil result for nil Canary deep copy")
    }
    var nilMetricTemplate *MetricTemplate
    if nilMetricTemplate.DeepCopy() != nil {
        t.Error("Expected nil result for nil MetricTemplate deep copy")
    }
    var nilAutoscalerRef *AutoscalerRefernce
    if nilAutoscalerRef.DeepCopy() != nil {
        t.Error("Expected nil result for nil AutoscalerRefernce deep copy")
    }
    var nilScalerReplicas *ScalerReplicas
    if nilScalerReplicas.DeepCopy() != nil {
        t.Error("Expected nil result for nil ScalerReplicas deep copy")
    }
}

// TestDeepCopyHTTPRewriteIndependence tests that a deep copy of HTTPRewrite (a trivial struct) returns an independent copy.
func TestDeepCopyHTTPRewriteIndependence(t *testing.T) {
    orig := &HTTPRewrite{}
    copyRewrite := orig.DeepCopy()
    if copyRewrite == orig {
        t.Error("DeepCopy did not create a distinct HTTPRewrite object")
    }
    // HTTPRewrite has no fields so we assert deep equality.
    if !reflect.DeepEqual(orig, copyRewrite) {
        t.Error("Deep copied HTTPRewrite is not equal to the original")
    }
}