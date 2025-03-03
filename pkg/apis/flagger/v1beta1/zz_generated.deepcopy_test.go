package v1beta1

import (
    "reflect"
    "testing"
    "time"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    v1 "k8s.io/api/core/v1"
    gatewayapiv1 "github.com/fluxcd/flagger/pkg/apis/gatewayapi/v1"
    istiov1beta1 "github.com/fluxcd/flagger/pkg/apis/istio/v1beta1"
)

// TestDeepCopyAlertProvider tests the DeepCopy and DeepCopyObject functions for AlertProvider.
func TestDeepCopyAlertProvider(t *testing.T) {
    orig := &AlertProvider{
        TypeMeta:   metav1.TypeMeta{Kind: "AlertProvider", APIVersion: "v1beta1"},
        ObjectMeta: metav1.ObjectMeta{Name: "alert1", Namespace: "default"},
        Spec: AlertProviderSpec{
            SecretRef: &v1.LocalObjectReference{Name: "secret1"},
        },
        Status: AlertProviderStatus{
            Conditions: []AlertProviderCondition{
                {
                    LastUpdateTime:     metav1.Time{Time: time.Now()},
                    LastTransitionTime: metav1.Time{Time: time.Now()},
                },
            },
        },
    }

    copyObj := orig.DeepCopy()
    if copyObj == orig {
        t.Errorf("DeepCopy returned the same pointer")
    }

    if !reflect.DeepEqual(orig, copyObj) {
        t.Errorf("DeepCopy did not produce an equivalent object")
    }

    // Modify the copy and check that the original is not affected.
    copyObj.Spec.SecretRef.Name = "changed-secret"
    if orig.Spec.SecretRef.Name == "changed-secret" {
        t.Errorf("DeepCopy is shallow; modification in copy affected the original")
    }
}

// TestDeepCopyObjectAlertProvider tests the DeepCopyObject method for AlertProvider.
func TestDeepCopyObjectAlertProvider(t *testing.T) {
    orig := &AlertProvider{
        TypeMeta:   metav1.TypeMeta{Kind: "AlertProvider", APIVersion: "v1beta1"},
        ObjectMeta: metav1.ObjectMeta{Name: "alert1", Namespace: "default"},
    }

    obj := orig.DeepCopyObject()
    castObj, ok := obj.(*AlertProvider)
    if !ok {
        t.Errorf("DeepCopyObject did not return an *AlertProvider")
    }

    if !reflect.DeepEqual(orig, castObj) {
        t.Errorf("DeepCopyObject did not produce an equivalent object")
    }
}

// TestDeepCopyCanary tests the DeepCopy and DeepCopyObject functions for Canary.
func TestDeepCopyCanary(t *testing.T) {
    orig := &Canary{
        TypeMeta:   metav1.TypeMeta{Kind: "Canary", APIVersion: "v1beta1"},
        ObjectMeta: metav1.ObjectMeta{Name: "canary1", Namespace: "default"},
        Spec: CanarySpec{
            IngressRef: &LocalObjectReference{Name: "ingress1"},
            RouteRef:   &LocalObjectReference{Name: "route1"},
            Service: CanaryService{
                Gateways: []string{"gateway1", "gateway2"},
                Hosts:    []string{"host1", "host2"},
            },
            ProgressDeadlineSeconds: new(int32),
        },
        Status: CanaryStatus{
            Conditions: []CanaryCondition{
                {
                    LastUpdateTime:     metav1.Time{Time: time.Now()},
                    LastTransitionTime: metav1.Time{Time: time.Now()},
                },
            },
        },
    }
    *orig.Spec.ProgressDeadlineSeconds = 30

    copyObj := orig.DeepCopy()
    if copyObj == orig {
        t.Errorf("DeepCopy returned the same pointer")
    }

    if !reflect.DeepEqual(orig, copyObj) {
        t.Errorf("DeepCopy did not produce an equivalent object")
    }

    // Modify nested fields and verify the original remains unchanged.
    copyObj.ObjectMeta.Name = "changed-canary"
    if orig.ObjectMeta.Name == "changed-canary" {
        t.Errorf("Modifying copy's ObjectMeta affected the original")
    }
}

// TestDeepCopyObjectCanary tests the DeepCopyObject method for Canary.
func TestDeepCopyObjectCanary(t *testing.T) {
    orig := &Canary{
        TypeMeta:   metav1.TypeMeta{Kind: "Canary", APIVersion: "v1beta1"},
        ObjectMeta: metav1.ObjectMeta{Name: "canary1", Namespace: "default"},
    }

    obj := orig.DeepCopyObject()
    castObj, ok := obj.(*Canary)
    if !ok {
        t.Errorf("DeepCopyObject did not return an *Canary")
    }

    if !reflect.DeepEqual(orig, castObj) {
        t.Errorf("DeepCopyObject did not produce an equivalent object")
    }
}

// TestDeepCopyCanarySpec tests the deep copy of CanarySpec including maps.
func TestDeepCopyCanarySpec(t *testing.T) {
    progress := int32(60)
    orig := &CanarySpec{
        IngressRef: &LocalObjectReference{Name: "ingress1"},
        RouteRef:   &LocalObjectReference{Name: "route1"},
        ProgressDeadlineSeconds: &progress,
        AutoscalerRef: &AutoscalerRefernce{
            PrimaryScalerQueries: map[string]string{"query": "value"},
        },
    }

    copySpec := orig.DeepCopy()
    if copySpec == orig {
        t.Errorf("DeepCopy returned the same pointer for CanarySpec")
    }

    if !reflect.DeepEqual(orig, copySpec) {
        t.Errorf("DeepCopy did not produce an equivalent CanarySpec")
    }

    // Modify the copy and ensure the original is not affected.
    copySpec.AutoscalerRef.PrimaryScalerQueries["query"] = "modified"
    if orig.AutoscalerRef.PrimaryScalerQueries["query"] == "modified" {
        t.Errorf("Modifying the deep copy of CanarySpec affected the original")
    }
}

// TestDeepCopyMetricTemplate tests the DeepCopy and DeepCopyObject functions for MetricTemplate.
func TestDeepCopyMetricTemplate(t *testing.T) {
    orig := &MetricTemplate{
        TypeMeta:   metav1.TypeMeta{Kind: "MetricTemplate", APIVersion: "v1beta1"},
        ObjectMeta: metav1.ObjectMeta{Name: "mt1", Namespace: "default"},
        Spec: MetricTemplateSpec{
            Provider: MetricTemplateProvider{
                SecretRef: &v1.LocalObjectReference{Name: "secret-mt"},
            },
        },
        Status: MetricTemplateStatus{
            Conditions: []MetricTemplateCondition{
                {
                    LastUpdateTime:     metav1.Time{Time: time.Now()},
                    LastTransitionTime: metav1.Time{Time: time.Now()},
                },
            },
        },
    }

    copyMT := orig.DeepCopy()
    if copyMT == orig {
        t.Errorf("DeepCopy returned the same pointer for MetricTemplate")
    }

    if !reflect.DeepEqual(orig, copyMT) {
        t.Errorf("DeepCopy did not produce an equivalent MetricTemplate")
    }

    // Modify the copy and verify that the original remains unchanged.
    copyMT.Spec.Provider.SecretRef.Name = "changed-secret-mt"
    if orig.Spec.Provider.SecretRef.Name == "changed-secret-mt" {
        t.Errorf("Modifying the deep copy of MetricTemplate affected the original")
    }
}
// TestDeepCopyNilValues tests that calling DeepCopy on nil pointers returns nil.
func TestDeepCopyNilValues(t *testing.T) {
    var ap *AlertProvider = nil
    if got := ap.DeepCopy(); got != nil {
        t.Errorf("Expected nil deep copy for AlertProvider, got %v", got)
    }

    var c *Canary = nil
    if got := c.DeepCopy(); got != nil {
        t.Errorf("Expected nil deep copy for Canary, got %v", got)
    }

    var mt *MetricTemplate = nil
    if got := mt.DeepCopy(); got != nil {
        t.Errorf("Expected nil deep copy for MetricTemplate, got %v", got)
    }
}

// TestDeepCopyAutoscalerRefernce tests the deep copy for AutoscalerRefernce.
func TestDeepCopyAutoscalerRefernce(t *testing.T) {
    orig := &AutoscalerRefernce{
        PrimaryScalerQueries: map[string]string{"key": "value"},
        PrimaryScalerReplicas: &ScalerReplicas{
            MinReplicas: new(int32),
            MaxReplicas: new(int32),
        },
    }
    *orig.PrimaryScalerReplicas.MinReplicas = 1
    *orig.PrimaryScalerReplicas.MaxReplicas = 5

    copyAR := orig.DeepCopy()
    if copyAR == orig {
        t.Errorf("DeepCopy returned the same pointer for AutoscalerRefernce")
    }
    if !reflect.DeepEqual(orig, copyAR) {
        t.Errorf("DeepCopy did not produce an equivalent AutoscalerRefernce")
    }

    // Modify the copy and check that the original is not affected.
    copyAR.PrimaryScalerQueries["key"] = "changed"
    *copyAR.PrimaryScalerReplicas.MinReplicas = 10
    if orig.PrimaryScalerQueries["key"] == "changed" {
        t.Errorf("Modifying the deep copy map affected the original")
    }
    if *orig.PrimaryScalerReplicas.MinReplicas == 10 {
        t.Errorf("Modifying the deep copy ScalerReplicas affected the original")
    }
}

// TestDeepCopyCanaryAnalysis tests the deep copy function for CanaryAnalysis, including slice and map fields.
func TestDeepCopyCanaryAnalysis(t *testing.T) {
    orig := &CanaryAnalysis{
        StepWeights:           []int{10, 20, 30},
        PrimaryReadyThreshold: new(int),
        CanaryReadyThreshold:  new(int),
        Alerts: []CanaryAlert{
                {ProviderRef: CrossNamespaceObjectReference{Name: "provider1"}},
        },
        Metrics: []CanaryMetric{
            {
                TemplateVariables: map[string]string{"var": "val"},
            },
        },
        Webhooks: []CanaryWebhook{
            {
                Metadata: &map[string]string{"hook": "test"},
            },
        },
        // Use an empty slice for Match (from istiov1beta1) since we cannot easily fabricate a full HTTPMatchRequest.
        Match:           []istiov1beta1.HTTPMatchRequest{},
        SessionAffinity: &SessionAffinity{},
    }
    *orig.PrimaryReadyThreshold = 5
    *orig.CanaryReadyThreshold = 3

    copyAnalysis := orig.DeepCopy()
    if copyAnalysis == orig {
        t.Errorf("DeepCopy returned the same pointer for CanaryAnalysis")
    }
    if !reflect.DeepEqual(orig, copyAnalysis) {
        t.Errorf("DeepCopy did not produce an equivalent CanaryAnalysis")
    }

    // Modify the copy and verify the original remains unchanged.
    copyAnalysis.StepWeights[1] = 99
    copyAnalysis.Metrics[0].TemplateVariables["var"] = "changed"
    if orig.StepWeights[1] == 99 {
        t.Errorf("Modifying copy's StepWeights affected the original")
    }
    if orig.Metrics[0].TemplateVariables["var"] == "changed" {
        t.Errorf("Modifying copy's Metrics map affected the original")
    }
}

// TestDeepCopyCustomBackend tests the deep copy for CustomBackend including its nested fields.
func TestDeepCopyCustomBackend(t *testing.T) {
    // We create a CustomBackend with an empty BackendObjectReference and empty Filters.
    orig := &CustomBackend{
        BackendObjectReference: &gatewayapiv1.BackendObjectReference{},
        Filters:                []gatewayapiv1.HTTPRouteFilter{},
    }
    copyBackend := orig.DeepCopy()
    if copyBackend == orig {
        t.Errorf("DeepCopy returned the same pointer for CustomBackend")
    }
    if !reflect.DeepEqual(orig, copyBackend) {
        t.Errorf("DeepCopy did not produce an equivalent CustomBackend")
    }
}

// TestDeepCopyLocalObjectReference tests the deep copy for LocalObjectReference.
func TestDeepCopyLocalObjectReference(t *testing.T) {
    orig := &LocalObjectReference{Name: "local1"}
    copyLoc := orig.DeepCopy()
    if copyLoc == orig {
        t.Errorf("DeepCopy returned the same pointer for LocalObjectReference")
    }
    if !reflect.DeepEqual(orig, copyLoc) {
        t.Errorf("DeepCopy did not produce an equivalent LocalObjectReference")
    }
    copyLoc.Name = "modified"
    if orig.Name == "modified" {
        t.Errorf("Modification in deep copied LocalObjectReference affected the original")
    }
}

// TestDeepCopySessionAffinity tests the deep copy for SessionAffinity, a trivial copy.
func TestDeepCopySessionAffinity(t *testing.T) {
    orig := &SessionAffinity{}
    copySA := orig.DeepCopy()
    if copySA == orig {
        t.Errorf("DeepCopy returned the same pointer for SessionAffinity")
    }
    if !reflect.DeepEqual(orig, copySA) {
        t.Errorf("DeepCopy did not produce an equivalent SessionAffinity")
    }
}
// TestDeepCopyAlertProviderCondition tests the DeepCopy function for AlertProviderCondition.
func TestDeepCopyAlertProviderCondition(t *testing.T) {
    orig := &AlertProviderCondition{
        LastUpdateTime:     metav1.Time{Time: time.Now()},
        LastTransitionTime: metav1.Time{Time: time.Now()},
    }
    copy := orig.DeepCopy()
    if copy == orig {
        t.Errorf("DeepCopy returned the same pointer for AlertProviderCondition")
    }
    if !reflect.DeepEqual(orig, copy) {
        t.Errorf("DeepCopy did not produce an equivalent AlertProviderCondition")
    }
}

// TestDeepCopyCanaryAlert tests the DeepCopy function for CanaryAlert.
func TestDeepCopyCanaryAlert(t *testing.T) {
    orig := &CanaryAlert{
        ProviderRef: CrossNamespaceObjectReference{Name: "alert-provider"},
    }
    copy := orig.DeepCopy()
    if copy == orig {
        t.Errorf("DeepCopy returned the same pointer for CanaryAlert")
    }
    if !reflect.DeepEqual(orig, copy) {
        t.Errorf("DeepCopy did not produce an equivalent CanaryAlert")
    }
}

// TestDeepCopyCanaryMetric tests the DeepCopy function for CanaryMetric.
func TestDeepCopyCanaryMetric(t *testing.T) {
    orig := &CanaryMetric{
        TemplateVariables: map[string]string{"key": "value"},
    }
    copy := orig.DeepCopy()
    if copy == orig {
        t.Errorf("DeepCopy returned the same pointer for CanaryMetric")
    }
    if !reflect.DeepEqual(orig, copy) {
        t.Errorf("DeepCopy did not produce an equivalent CanaryMetric")
    }
    // Change the copy and ensure the original is unchanged.
    copy.TemplateVariables["key"] = "changed"
    if orig.TemplateVariables["key"] == "changed" {
        t.Errorf("Modifying the copy's TemplateVariables affected the original")
    }
}

// TestDeepCopyCustomMetadata tests the DeepCopy function for CustomMetadata.
func TestDeepCopyCustomMetadata(t *testing.T) {
    orig := &CustomMetadata{
        Labels:      map[string]string{"app": "nginx"},
        Annotations: map[string]string{"version": "1.0"},
    }
    copy := orig.DeepCopy()
    if copy == orig {
        t.Errorf("DeepCopy returned the same pointer for CustomMetadata")
    }
    if !reflect.DeepEqual(orig, copy) {
        t.Errorf("DeepCopy did not produce an equivalent CustomMetadata")
    }
    copy.Labels["app"] = "apache"
    if orig.Labels["app"] == "apache" {
        t.Errorf("Deep copy modification affected original CustomMetadata")
    }
}

// TestDeepCopyMetricTemplateModel tests the DeepCopy function for MetricTemplateModel.
func TestDeepCopyMetricTemplateModel(t *testing.T) {
    orig := &MetricTemplateModel{
        Variables: map[string]string{"var1": "foo"},
    }
    copy := orig.DeepCopy()
    if copy == orig {
        t.Errorf("DeepCopy returned the same pointer for MetricTemplateModel")
    }
    if !reflect.DeepEqual(orig, copy) {
        t.Errorf("DeepCopy did not produce an equivalent MetricTemplateModel")
    }
    copy.Variables["var1"] = "bar"
    if orig.Variables["var1"] == "bar" {
        t.Errorf("Modifying the deep copy affected original MetricTemplateModel")
    }
}

// TestDeepCopyMetricTemplateCondition tests the DeepCopy function for MetricTemplateCondition.
func TestDeepCopyMetricTemplateCondition(t *testing.T) {
    orig := &MetricTemplateCondition{
        LastUpdateTime:     metav1.Time{Time: time.Now()},
        LastTransitionTime: metav1.Time{Time: time.Now()},
    }
    copy := orig.DeepCopy()
    if copy == orig {
        t.Errorf("DeepCopy returned the same pointer for MetricTemplateCondition")
    }
    if !reflect.DeepEqual(orig, copy) {
        t.Errorf("DeepCopy did not produce an equivalent MetricTemplateCondition")
    }
}

// TestDeepCopyCrossNamespaceObjectReference tests the DeepCopy function for CrossNamespaceObjectReference.
func TestDeepCopyCrossNamespaceObjectReference(t *testing.T) {
    orig := &CrossNamespaceObjectReference{Name: "ns-ref"}
    copy := orig.DeepCopy()
    if copy == orig {
        t.Errorf("DeepCopy returned the same pointer for CrossNamespaceObjectReference")
    }
    if !reflect.DeepEqual(orig, copy) {
        t.Errorf("DeepCopy did not produce an equivalent CrossNamespaceObjectReference")
    }
}

// TestDeepCopyHTTPRewrite tests the DeepCopy function for HTTPRewrite.
func TestDeepCopyHTTPRewrite(t *testing.T) {
    orig := &HTTPRewrite{}
    copy := orig.DeepCopy()
    if copy == orig {
        t.Errorf("DeepCopy returned the same pointer for HTTPRewrite")
    }
    if !reflect.DeepEqual(orig, copy) {
        t.Errorf("DeepCopy did not produce an equivalent HTTPRewrite")
    }
}

// TestDeepCopyScalerReplicas tests the DeepCopy function for ScalerReplicas.
func TestDeepCopyScalerReplicas(t *testing.T) {
    min := int32(2)
    max := int32(10)
    orig := &ScalerReplicas{
        MinReplicas: &min,
        MaxReplicas: &max,
    }
    copy := orig.DeepCopy()
    if copy == orig {
        t.Errorf("DeepCopy returned the same pointer for ScalerReplicas")
    }
    if !reflect.DeepEqual(orig, copy) {
        t.Errorf("DeepCopy did not produce an equivalent ScalerReplicas")
    }
    *copy.MinReplicas = 5
    if *orig.MinReplicas == 5 {
        t.Errorf("Modifying deep copy ScalerReplicas affected the original")
    }
}