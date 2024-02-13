package warden

import (
	"context"
	"fmt"
	"testing"

	"github.com/kyma-incubator/reconciler/pkg/reconciler"
	"github.com/kyma-incubator/reconciler/pkg/reconciler/kubernetes/mocks"
	"github.com/kyma-incubator/reconciler/pkg/reconciler/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"sigs.k8s.io/kustomize/kyaml/yaml"

	chartmocks "github.com/kyma-incubator/reconciler/pkg/reconciler/chart/mocks"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCleanupWardenAdmissionCertVolumeMounts_Run(t *testing.T) {

	t.Run("no admission image override present", func(t *testing.T) {
		chartProvider := &chartmocks.Provider{}
		chartValuesYAML := getWardenValuesYAML(false, "")
		chartValues, err := unmarshalTestValues(chartValuesYAML)
		require.NoError(t, err)
		chartProvider.On("Configuration", mock.Anything).Return(chartValues, nil)
		ctx := context.Background()
		mockClient := &mocks.Client{}
		context := &service.ActionContext{
			Context:       ctx,
			Logger:        zap.NewNop().Sugar(),
			KubeClient:    mockClient,
			Task:          &reconciler.Task{Version: "testversion"},
			ChartProvider: chartProvider,
		}
		action := &CleanupWardenAdmissionCertVolumeMounts{}
		err = action.Run(context)
		require.NoError(t, err)
	})

	t.Run("admission image override present - doesnt qualify for cleanup", func(t *testing.T) {
		chartProvider := &chartmocks.Provider{}
		chartValuesYAML := getWardenValuesYAML(true, "europe-docker.pkg.dev/kyma-project/prod/warden/admission:v20230529-2de25d82")
		chartValues, err := unmarshalTestValues(chartValuesYAML)
		require.NoError(t, err)
		chartProvider.On("Configuration", mock.Anything).Return(chartValues, nil)
		ctx := context.Background()
		mockClient := &mocks.Client{}
		context := &service.ActionContext{
			Context:       ctx,
			Logger:        zap.NewNop().Sugar(),
			KubeClient:    mockClient,
			Task:          &reconciler.Task{Version: "testversion"},
			ChartProvider: chartProvider,
		}
		action := &CleanupWardenAdmissionCertVolumeMounts{}
		err = action.Run(context)
		require.NoError(t, err)
	})

	t.Run("admission image override present but doesnt have the : delimiter", func(t *testing.T) {

		chartProvider := &chartmocks.Provider{}
		chartValuesYAML := getWardenValuesYAML(true, "europe-docker.pkg.dev/kyma-project/prod/warden/admission")
		chartValues, err := unmarshalTestValues(chartValuesYAML)
		require.NoError(t, err)
		chartProvider.On("Configuration", mock.Anything).Return(chartValues, nil)
		ctx := context.Background()
		mockClient := &mocks.Client{}
		context := &service.ActionContext{
			Context:       ctx,
			Logger:        zap.NewNop().Sugar(),
			KubeClient:    mockClient,
			Task:          &reconciler.Task{Version: "testversion"},
			ChartProvider: chartProvider,
		}
		action := &CleanupWardenAdmissionCertVolumeMounts{}
		err = action.Run(context)
		require.NoError(t, err)
	})

	t.Run("admission image override present - qualifies for cleanup but no warden admission deployment found", func(t *testing.T) {
		chartProvider := &chartmocks.Provider{}
		chartValuesYAML := getWardenValuesYAML(true, "europe-docker.pkg.dev/kyma-project/prod/warden/admission:0.10.0")
		chartValues, err := unmarshalTestValues(chartValuesYAML)
		require.NoError(t, err)
		chartProvider.On("Configuration", mock.Anything).Return(chartValues, nil)
		ctx := context.Background()
		mockClient := &mocks.Client{}
		mockClient.On("GetDeployment", ctx, wardenAdmissionDeploymentName, wardenAdmissionDeploymentNamespace).Return(nil, errors.NewNotFound(schema.GroupResource{}, "test error"))
		context := &service.ActionContext{
			Context:       ctx,
			Logger:        zap.NewNop().Sugar(),
			KubeClient:    mockClient,
			Task:          &reconciler.Task{Version: "testversion"},
			ChartProvider: chartProvider,
		}
		action := &CleanupWardenAdmissionCertVolumeMounts{}
		err = action.Run(context)
		require.NoError(t, err)
	})

	t.Run("error while reading warden admission deployment", func(t *testing.T) {
		chartProvider := &chartmocks.Provider{}
		chartValuesYAML := getWardenValuesYAML(true, "europe-docker.pkg.dev/kyma-project/prod/warden/admission:0.10.1")
		chartValues, err := unmarshalTestValues(chartValuesYAML)
		require.NoError(t, err)
		chartProvider.On("Configuration", mock.Anything).Return(chartValues, nil)
		ctx := context.Background()
		mockClient := &mocks.Client{}
		mockClient.On("GetDeployment", ctx, wardenAdmissionDeploymentName, wardenAdmissionDeploymentNamespace).Return(nil, errors.NewBadRequest("test error"))
		context := &service.ActionContext{
			Context:       ctx,
			Logger:        zap.NewNop().Sugar(),
			KubeClient:    mockClient,
			Task:          &reconciler.Task{Version: "testversion"},
			ChartProvider: chartProvider,
		}
		action := &CleanupWardenAdmissionCertVolumeMounts{}
		err = action.Run(context)
		require.Error(t, err)
	})

	t.Run("warden admission deployment found in v 0.10.1 but no admission container", func(t *testing.T) {
		chartProvider := &chartmocks.Provider{}
		chartValuesYAML := getWardenValuesYAML(true, "europe-docker.pkg.dev/kyma-project/prod/warden/admission:0.10.1")
		chartValues, err := unmarshalTestValues(chartValuesYAML)
		require.NoError(t, err)
		chartProvider.On("Configuration", mock.Anything).Return(chartValues, nil)
		ctx := context.Background()
		mockClient := &mocks.Client{}
		podSpec := corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "foo"}, {Name: "bar"},
			},
		}
		mockClient.On("GetDeployment", ctx, wardenAdmissionDeploymentName, wardenAdmissionDeploymentNamespace).Return(fixedDeploymentWithPodSpec(podSpec), nil)
		context := &service.ActionContext{
			Context:    ctx,
			Logger:     zap.NewNop().Sugar(),
			KubeClient: mockClient,
			Task: &reconciler.Task{
				Version: "testversion",
			},
			ChartProvider: chartProvider,
		}
		action := &CleanupWardenAdmissionCertVolumeMounts{}
		err = action.Run(context)
		require.NoError(t, err)
	})

	t.Run("warden admission deployment found in v 0.10.1 but no volumemounts", func(t *testing.T) {
		chartProvider := &chartmocks.Provider{}
		chartValuesYAML := getWardenValuesYAML(true, "europe-docker.pkg.dev/kyma-project/prod/warden/admission:0.10.1")
		chartValues, err := unmarshalTestValues(chartValuesYAML)
		require.NoError(t, err)
		chartProvider.On("Configuration", mock.Anything).Return(chartValues, nil)
		ctx := context.Background()
		mockClient := &mocks.Client{}
		mockClient.On("GetDeployment", ctx, wardenAdmissionDeploymentName, wardenAdmissionDeploymentNamespace).Return(fixedDeploymentWith(nil, nil), nil)
		context := &service.ActionContext{
			Context:    ctx,
			Logger:     zap.NewNop().Sugar(),
			KubeClient: mockClient,
			Task: &reconciler.Task{
				Version: "testversion",
			},
			ChartProvider: chartProvider,
		}
		action := &CleanupWardenAdmissionCertVolumeMounts{}
		err = action.Run(context)
		require.NoError(t, err)
	})

	t.Run("warden admission deployment found in v 0.10.1 but only volumes found w/o volumemounts", func(t *testing.T) {
		vm := []corev1.Volume{
			{Name: "foo"}, {Name: volumeName},
		}

		chartProvider := &chartmocks.Provider{}
		chartValuesYAML := getWardenValuesYAML(true, "europe-docker.pkg.dev/kyma-project/prod/warden/admission:0.10.1")
		chartValues, err := unmarshalTestValues(chartValuesYAML)
		require.NoError(t, err)
		chartProvider.On("Configuration", mock.Anything).Return(chartValues, nil)
		ctx := context.Background()
		mockClient := &mocks.Client{}
		mockClient.On("GetDeployment", ctx, wardenAdmissionDeploymentName, wardenAdmissionDeploymentNamespace).Return(fixedDeploymentWith(nil, vm), nil)
		context := &service.ActionContext{
			Context:    ctx,
			Logger:     zap.NewNop().Sugar(),
			KubeClient: mockClient,
			Task: &reconciler.Task{
				Version: "testversion",
			},
			ChartProvider: chartProvider,
		}
		action := &CleanupWardenAdmissionCertVolumeMounts{}
		err = action.Run(context)
		require.NoError(t, err)
	})

	t.Run("warden admission deployment found in v 0.10.2 for cleanup with volumemounts", func(t *testing.T) {

		vms := []corev1.VolumeMount{
			{Name: "foo"}, {Name: volumeName},
		}
		vm := []corev1.Volume{
			{Name: "foo"}, {Name: volumeName},
		}

		chartProvider := &chartmocks.Provider{}
		chartValuesYAML := getWardenValuesYAML(true, "europe-docker.pkg.dev/kyma-project/prod/warden/admission:0.10.2")
		chartValues, err := unmarshalTestValues(chartValuesYAML)
		require.NoError(t, err)
		chartProvider.On("Configuration", mock.Anything).Return(chartValues, nil)

		ctx := context.Background()
		mockClient := &mocks.Client{}
		mockClient.On("GetDeployment", ctx, wardenAdmissionDeploymentName, wardenAdmissionDeploymentNamespace).Return(fixedDeploymentWith(vms, vm), nil)
		mockClient.On("PatchUsingStrategy", ctx, "Deployment", wardenAdmissionDeploymentName, wardenAdmissionDeploymentNamespace, []byte(`[{"op": "remove", "path": "/spec/template/spec/containers/0/volumeMounts/1"},{"op": "remove", "path": "/spec/template/spec/volumes/1"}]`), mock.Anything).Return(nil)
		context := &service.ActionContext{
			Context:    ctx,
			Logger:     zap.NewNop().Sugar(),
			KubeClient: mockClient,
			Task: &reconciler.Task{
				Version: "testversion",
			},
			ChartProvider: chartProvider,
		}
		action := &CleanupWardenAdmissionCertVolumeMounts{}
		err = action.Run(context)
		require.NoError(t, err)
	})

	t.Run("warden admission deployment found in v 0.10.2 for cleanup with volumemounts and multiple containers", func(t *testing.T) {

		vms := []corev1.VolumeMount{
			{Name: "foo"}, {Name: volumeName},
		}
		vm := []corev1.Volume{
			{Name: "foo"}, {Name: volumeName},
		}

		containers := []corev1.Container{
			{Name: "foo"}, {Name: "bar"}, {Name: containerName, VolumeMounts: vms},
		}

		podSpec := corev1.PodSpec{
			Containers: containers,
			Volumes:    vm,
		}

		chartProvider := &chartmocks.Provider{}
		chartValuesYAML := getWardenValuesYAML(true, "europe-docker.pkg.dev/kyma-project/prod/warden/admission:0.10.2")
		chartValues, err := unmarshalTestValues(chartValuesYAML)
		require.NoError(t, err)
		chartProvider.On("Configuration", mock.Anything).Return(chartValues, nil)
		ctx := context.Background()
		mockClient := &mocks.Client{}
		mockClient.On("GetDeployment", ctx, wardenAdmissionDeploymentName, wardenAdmissionDeploymentNamespace).Return(fixedDeploymentWithPodSpec(podSpec), nil)
		mockClient.On("PatchUsingStrategy", ctx, "Deployment", wardenAdmissionDeploymentName, wardenAdmissionDeploymentNamespace, []byte(`[{"op": "remove", "path": "/spec/template/spec/containers/2/volumeMounts/1"},{"op": "remove", "path": "/spec/template/spec/volumes/1"}]`), mock.Anything).Return(nil)
		context := &service.ActionContext{
			Context:    ctx,
			Logger:     zap.NewNop().Sugar(),
			KubeClient: mockClient,
			Task: &reconciler.Task{
				Version: "testversion",
			},
			ChartProvider: chartProvider,
		}
		action := &CleanupWardenAdmissionCertVolumeMounts{}
		err = action.Run(context)
		require.NoError(t, err)
	})

	t.Run("warden admission deployment found in v 0.10.2 for cleanup with volumemounts - handle error", func(t *testing.T) {

		vms := []corev1.VolumeMount{
			{Name: volumeName}, {Name: "foo"},
		}
		vm := []corev1.Volume{
			{Name: "foo"}, {Name: volumeName},
		}

		chartProvider := &chartmocks.Provider{}
		chartValuesYAML := getWardenValuesYAML(true, "europe-docker.pkg.dev/kyma-project/prod/warden/admission:0.10.2")
		chartValues, err := unmarshalTestValues(chartValuesYAML)
		require.NoError(t, err)
		chartProvider.On("Configuration", mock.Anything).Return(chartValues, nil)
		ctx := context.Background()
		mockClient := &mocks.Client{}
		mockClient.On("GetDeployment", ctx, wardenAdmissionDeploymentName, wardenAdmissionDeploymentNamespace).Return(fixedDeploymentWith(vms, vm), nil)
		mockClient.On("PatchUsingStrategy", ctx, "Deployment", wardenAdmissionDeploymentName, wardenAdmissionDeploymentNamespace, []byte(`[{"op": "remove", "path": "/spec/template/spec/containers/0/volumeMounts/0"},{"op": "remove", "path": "/spec/template/spec/volumes/1"}]`), mock.Anything).Return(errors.NewBadRequest(""))
		context := &service.ActionContext{
			Context:    ctx,
			Logger:     zap.NewNop().Sugar(),
			KubeClient: mockClient,
			Task: &reconciler.Task{
				Version: "testversion",
			},
			ChartProvider: chartProvider,
		}
		action := &CleanupWardenAdmissionCertVolumeMounts{}
		err = action.Run(context)
		require.Error(t, err)
	})
}

func fixedDeploymentWith(volumeMounts []corev1.VolumeMount, volumes []corev1.Volume) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      wardenAdmissionDeploymentName,
			Namespace: wardenAdmissionDeploymentNamespace,
		},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:         containerName,
							VolumeMounts: volumeMounts,
						},
					},
					Volumes: volumes,
				},
			},
		},
	}
}

func fixedDeploymentWithPodSpec(podSpec corev1.PodSpec) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      wardenAdmissionDeploymentName,
			Namespace: wardenAdmissionDeploymentNamespace,
		},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: podSpec,
			},
		},
	}
}

func getWardenValuesYAML(withImage bool, image string) string {
	if !withImage {
		return `
    global:
      admission:`
	}

	return fmt.Sprintf(`
    global:
      admission:
        image: %s`,
		image,
	)
}

func unmarshalTestValues(yamlValues string) (map[string]interface{}, error) {
	var values map[string]interface{}
	err := yaml.Unmarshal([]byte(yamlValues), &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}
