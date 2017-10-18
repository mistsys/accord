package aws_params

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type MockedSSMAPI struct {
	ssmiface.SSMAPI
	mSecureStrings map[string]*ssm.Parameter
	mStrings       map[string]*ssm.Parameter
}

func NewMockedSSMAPI(secureStrings map[string]string, insecureStrings map[string]string) *MockedSSMAPI {
	secStrings := make(map[string]*ssm.Parameter)
	insecStrings := make(map[string]*ssm.Parameter)
	for name, str := range secureStrings {
		secStrings[name] = &ssm.Parameter{
			Name:  aws.String(name),
			Type:  aws.String("Secure String"),
			Value: aws.String(str),
		}
	}

	for name, str := range insecureStrings {
		insecStrings[name] = &ssm.Parameter{
			Name:  aws.String(name),
			Type:  aws.String("String"),
			Value: aws.String(str),
		}
	}

	return &MockedSSMAPI{
		mSecureStrings: secStrings,
		mStrings:       insecStrings,
	}
}

// we only need GetParameter but ssmiface requires all of these
// Generated with this command: impl 'c *MockedSSMAPI' github.com/aws/aws-sdk-go/service/ssm/ssmiface.SSMAPI

func (c *MockedSSMAPI) AddTagsToResource(*ssm.AddTagsToResourceInput) (*ssm.AddTagsToResourceOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) AddTagsToResourceWithContext(aws.Context, *ssm.AddTagsToResourceInput, ...request.Option) (*ssm.AddTagsToResourceOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) AddTagsToResourceRequest(*ssm.AddTagsToResourceInput) (*request.Request, *ssm.AddTagsToResourceOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CancelCommand(*ssm.CancelCommandInput) (*ssm.CancelCommandOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CancelCommandWithContext(aws.Context, *ssm.CancelCommandInput, ...request.Option) (*ssm.CancelCommandOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CancelCommandRequest(*ssm.CancelCommandInput) (*request.Request, *ssm.CancelCommandOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreateActivation(*ssm.CreateActivationInput) (*ssm.CreateActivationOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreateActivationWithContext(aws.Context, *ssm.CreateActivationInput, ...request.Option) (*ssm.CreateActivationOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreateActivationRequest(*ssm.CreateActivationInput) (*request.Request, *ssm.CreateActivationOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreateAssociation(*ssm.CreateAssociationInput) (*ssm.CreateAssociationOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreateAssociationWithContext(aws.Context, *ssm.CreateAssociationInput, ...request.Option) (*ssm.CreateAssociationOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreateAssociationRequest(*ssm.CreateAssociationInput) (*request.Request, *ssm.CreateAssociationOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreateAssociationBatch(*ssm.CreateAssociationBatchInput) (*ssm.CreateAssociationBatchOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreateAssociationBatchWithContext(aws.Context, *ssm.CreateAssociationBatchInput, ...request.Option) (*ssm.CreateAssociationBatchOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreateAssociationBatchRequest(*ssm.CreateAssociationBatchInput) (*request.Request, *ssm.CreateAssociationBatchOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreateDocument(*ssm.CreateDocumentInput) (*ssm.CreateDocumentOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreateDocumentWithContext(aws.Context, *ssm.CreateDocumentInput, ...request.Option) (*ssm.CreateDocumentOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreateDocumentRequest(*ssm.CreateDocumentInput) (*request.Request, *ssm.CreateDocumentOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreateMaintenanceWindow(*ssm.CreateMaintenanceWindowInput) (*ssm.CreateMaintenanceWindowOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreateMaintenanceWindowWithContext(aws.Context, *ssm.CreateMaintenanceWindowInput, ...request.Option) (*ssm.CreateMaintenanceWindowOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreateMaintenanceWindowRequest(*ssm.CreateMaintenanceWindowInput) (*request.Request, *ssm.CreateMaintenanceWindowOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreatePatchBaseline(*ssm.CreatePatchBaselineInput) (*ssm.CreatePatchBaselineOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreatePatchBaselineWithContext(aws.Context, *ssm.CreatePatchBaselineInput, ...request.Option) (*ssm.CreatePatchBaselineOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreatePatchBaselineRequest(*ssm.CreatePatchBaselineInput) (*request.Request, *ssm.CreatePatchBaselineOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreateResourceDataSync(*ssm.CreateResourceDataSyncInput) (*ssm.CreateResourceDataSyncOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreateResourceDataSyncWithContext(aws.Context, *ssm.CreateResourceDataSyncInput, ...request.Option) (*ssm.CreateResourceDataSyncOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) CreateResourceDataSyncRequest(*ssm.CreateResourceDataSyncInput) (*request.Request, *ssm.CreateResourceDataSyncOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteActivation(*ssm.DeleteActivationInput) (*ssm.DeleteActivationOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteActivationWithContext(aws.Context, *ssm.DeleteActivationInput, ...request.Option) (*ssm.DeleteActivationOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteActivationRequest(*ssm.DeleteActivationInput) (*request.Request, *ssm.DeleteActivationOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteAssociation(*ssm.DeleteAssociationInput) (*ssm.DeleteAssociationOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteAssociationWithContext(aws.Context, *ssm.DeleteAssociationInput, ...request.Option) (*ssm.DeleteAssociationOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteAssociationRequest(*ssm.DeleteAssociationInput) (*request.Request, *ssm.DeleteAssociationOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteDocument(*ssm.DeleteDocumentInput) (*ssm.DeleteDocumentOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteDocumentWithContext(aws.Context, *ssm.DeleteDocumentInput, ...request.Option) (*ssm.DeleteDocumentOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteDocumentRequest(*ssm.DeleteDocumentInput) (*request.Request, *ssm.DeleteDocumentOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteMaintenanceWindow(*ssm.DeleteMaintenanceWindowInput) (*ssm.DeleteMaintenanceWindowOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteMaintenanceWindowWithContext(aws.Context, *ssm.DeleteMaintenanceWindowInput, ...request.Option) (*ssm.DeleteMaintenanceWindowOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteMaintenanceWindowRequest(*ssm.DeleteMaintenanceWindowInput) (*request.Request, *ssm.DeleteMaintenanceWindowOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteParameter(*ssm.DeleteParameterInput) (*ssm.DeleteParameterOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteParameterWithContext(aws.Context, *ssm.DeleteParameterInput, ...request.Option) (*ssm.DeleteParameterOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteParameterRequest(*ssm.DeleteParameterInput) (*request.Request, *ssm.DeleteParameterOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteParameters(*ssm.DeleteParametersInput) (*ssm.DeleteParametersOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteParametersWithContext(aws.Context, *ssm.DeleteParametersInput, ...request.Option) (*ssm.DeleteParametersOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteParametersRequest(*ssm.DeleteParametersInput) (*request.Request, *ssm.DeleteParametersOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeletePatchBaseline(*ssm.DeletePatchBaselineInput) (*ssm.DeletePatchBaselineOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeletePatchBaselineWithContext(aws.Context, *ssm.DeletePatchBaselineInput, ...request.Option) (*ssm.DeletePatchBaselineOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeletePatchBaselineRequest(*ssm.DeletePatchBaselineInput) (*request.Request, *ssm.DeletePatchBaselineOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteResourceDataSync(*ssm.DeleteResourceDataSyncInput) (*ssm.DeleteResourceDataSyncOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteResourceDataSyncWithContext(aws.Context, *ssm.DeleteResourceDataSyncInput, ...request.Option) (*ssm.DeleteResourceDataSyncOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeleteResourceDataSyncRequest(*ssm.DeleteResourceDataSyncInput) (*request.Request, *ssm.DeleteResourceDataSyncOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeregisterManagedInstance(*ssm.DeregisterManagedInstanceInput) (*ssm.DeregisterManagedInstanceOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeregisterManagedInstanceWithContext(aws.Context, *ssm.DeregisterManagedInstanceInput, ...request.Option) (*ssm.DeregisterManagedInstanceOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeregisterManagedInstanceRequest(*ssm.DeregisterManagedInstanceInput) (*request.Request, *ssm.DeregisterManagedInstanceOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeregisterPatchBaselineForPatchGroup(*ssm.DeregisterPatchBaselineForPatchGroupInput) (*ssm.DeregisterPatchBaselineForPatchGroupOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeregisterPatchBaselineForPatchGroupWithContext(aws.Context, *ssm.DeregisterPatchBaselineForPatchGroupInput, ...request.Option) (*ssm.DeregisterPatchBaselineForPatchGroupOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeregisterPatchBaselineForPatchGroupRequest(*ssm.DeregisterPatchBaselineForPatchGroupInput) (*request.Request, *ssm.DeregisterPatchBaselineForPatchGroupOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeregisterTargetFromMaintenanceWindow(*ssm.DeregisterTargetFromMaintenanceWindowInput) (*ssm.DeregisterTargetFromMaintenanceWindowOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeregisterTargetFromMaintenanceWindowWithContext(aws.Context, *ssm.DeregisterTargetFromMaintenanceWindowInput, ...request.Option) (*ssm.DeregisterTargetFromMaintenanceWindowOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeregisterTargetFromMaintenanceWindowRequest(*ssm.DeregisterTargetFromMaintenanceWindowInput) (*request.Request, *ssm.DeregisterTargetFromMaintenanceWindowOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeregisterTaskFromMaintenanceWindow(*ssm.DeregisterTaskFromMaintenanceWindowInput) (*ssm.DeregisterTaskFromMaintenanceWindowOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeregisterTaskFromMaintenanceWindowWithContext(aws.Context, *ssm.DeregisterTaskFromMaintenanceWindowInput, ...request.Option) (*ssm.DeregisterTaskFromMaintenanceWindowOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DeregisterTaskFromMaintenanceWindowRequest(*ssm.DeregisterTaskFromMaintenanceWindowInput) (*request.Request, *ssm.DeregisterTaskFromMaintenanceWindowOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeActivations(*ssm.DescribeActivationsInput) (*ssm.DescribeActivationsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeActivationsWithContext(aws.Context, *ssm.DescribeActivationsInput, ...request.Option) (*ssm.DescribeActivationsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeActivationsRequest(*ssm.DescribeActivationsInput) (*request.Request, *ssm.DescribeActivationsOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeActivationsPages(*ssm.DescribeActivationsInput, func(*ssm.DescribeActivationsOutput, bool) bool) error {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeActivationsPagesWithContext(aws.Context, *ssm.DescribeActivationsInput, func(*ssm.DescribeActivationsOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeAssociation(*ssm.DescribeAssociationInput) (*ssm.DescribeAssociationOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeAssociationWithContext(aws.Context, *ssm.DescribeAssociationInput, ...request.Option) (*ssm.DescribeAssociationOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeAssociationRequest(*ssm.DescribeAssociationInput) (*request.Request, *ssm.DescribeAssociationOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeAutomationExecutions(*ssm.DescribeAutomationExecutionsInput) (*ssm.DescribeAutomationExecutionsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeAutomationExecutionsWithContext(aws.Context, *ssm.DescribeAutomationExecutionsInput, ...request.Option) (*ssm.DescribeAutomationExecutionsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeAutomationExecutionsRequest(*ssm.DescribeAutomationExecutionsInput) (*request.Request, *ssm.DescribeAutomationExecutionsOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeAvailablePatches(*ssm.DescribeAvailablePatchesInput) (*ssm.DescribeAvailablePatchesOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeAvailablePatchesWithContext(aws.Context, *ssm.DescribeAvailablePatchesInput, ...request.Option) (*ssm.DescribeAvailablePatchesOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeAvailablePatchesRequest(*ssm.DescribeAvailablePatchesInput) (*request.Request, *ssm.DescribeAvailablePatchesOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeDocument(*ssm.DescribeDocumentInput) (*ssm.DescribeDocumentOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeDocumentWithContext(aws.Context, *ssm.DescribeDocumentInput, ...request.Option) (*ssm.DescribeDocumentOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeDocumentRequest(*ssm.DescribeDocumentInput) (*request.Request, *ssm.DescribeDocumentOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeDocumentPermission(*ssm.DescribeDocumentPermissionInput) (*ssm.DescribeDocumentPermissionOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeDocumentPermissionWithContext(aws.Context, *ssm.DescribeDocumentPermissionInput, ...request.Option) (*ssm.DescribeDocumentPermissionOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeDocumentPermissionRequest(*ssm.DescribeDocumentPermissionInput) (*request.Request, *ssm.DescribeDocumentPermissionOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeEffectiveInstanceAssociations(*ssm.DescribeEffectiveInstanceAssociationsInput) (*ssm.DescribeEffectiveInstanceAssociationsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeEffectiveInstanceAssociationsWithContext(aws.Context, *ssm.DescribeEffectiveInstanceAssociationsInput, ...request.Option) (*ssm.DescribeEffectiveInstanceAssociationsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeEffectiveInstanceAssociationsRequest(*ssm.DescribeEffectiveInstanceAssociationsInput) (*request.Request, *ssm.DescribeEffectiveInstanceAssociationsOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeEffectivePatchesForPatchBaseline(*ssm.DescribeEffectivePatchesForPatchBaselineInput) (*ssm.DescribeEffectivePatchesForPatchBaselineOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeEffectivePatchesForPatchBaselineWithContext(aws.Context, *ssm.DescribeEffectivePatchesForPatchBaselineInput, ...request.Option) (*ssm.DescribeEffectivePatchesForPatchBaselineOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeEffectivePatchesForPatchBaselineRequest(*ssm.DescribeEffectivePatchesForPatchBaselineInput) (*request.Request, *ssm.DescribeEffectivePatchesForPatchBaselineOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeInstanceAssociationsStatus(*ssm.DescribeInstanceAssociationsStatusInput) (*ssm.DescribeInstanceAssociationsStatusOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeInstanceAssociationsStatusWithContext(aws.Context, *ssm.DescribeInstanceAssociationsStatusInput, ...request.Option) (*ssm.DescribeInstanceAssociationsStatusOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeInstanceAssociationsStatusRequest(*ssm.DescribeInstanceAssociationsStatusInput) (*request.Request, *ssm.DescribeInstanceAssociationsStatusOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeInstanceInformation(*ssm.DescribeInstanceInformationInput) (*ssm.DescribeInstanceInformationOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeInstanceInformationWithContext(aws.Context, *ssm.DescribeInstanceInformationInput, ...request.Option) (*ssm.DescribeInstanceInformationOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeInstanceInformationRequest(*ssm.DescribeInstanceInformationInput) (*request.Request, *ssm.DescribeInstanceInformationOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeInstanceInformationPages(*ssm.DescribeInstanceInformationInput, func(*ssm.DescribeInstanceInformationOutput, bool) bool) error {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeInstanceInformationPagesWithContext(aws.Context, *ssm.DescribeInstanceInformationInput, func(*ssm.DescribeInstanceInformationOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeInstancePatchStates(*ssm.DescribeInstancePatchStatesInput) (*ssm.DescribeInstancePatchStatesOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeInstancePatchStatesWithContext(aws.Context, *ssm.DescribeInstancePatchStatesInput, ...request.Option) (*ssm.DescribeInstancePatchStatesOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeInstancePatchStatesRequest(*ssm.DescribeInstancePatchStatesInput) (*request.Request, *ssm.DescribeInstancePatchStatesOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeInstancePatchStatesForPatchGroup(*ssm.DescribeInstancePatchStatesForPatchGroupInput) (*ssm.DescribeInstancePatchStatesForPatchGroupOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeInstancePatchStatesForPatchGroupWithContext(aws.Context, *ssm.DescribeInstancePatchStatesForPatchGroupInput, ...request.Option) (*ssm.DescribeInstancePatchStatesForPatchGroupOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeInstancePatchStatesForPatchGroupRequest(*ssm.DescribeInstancePatchStatesForPatchGroupInput) (*request.Request, *ssm.DescribeInstancePatchStatesForPatchGroupOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeInstancePatches(*ssm.DescribeInstancePatchesInput) (*ssm.DescribeInstancePatchesOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeInstancePatchesWithContext(aws.Context, *ssm.DescribeInstancePatchesInput, ...request.Option) (*ssm.DescribeInstancePatchesOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeInstancePatchesRequest(*ssm.DescribeInstancePatchesInput) (*request.Request, *ssm.DescribeInstancePatchesOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeMaintenanceWindowExecutionTaskInvocations(*ssm.DescribeMaintenanceWindowExecutionTaskInvocationsInput) (*ssm.DescribeMaintenanceWindowExecutionTaskInvocationsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeMaintenanceWindowExecutionTaskInvocationsWithContext(aws.Context, *ssm.DescribeMaintenanceWindowExecutionTaskInvocationsInput, ...request.Option) (*ssm.DescribeMaintenanceWindowExecutionTaskInvocationsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeMaintenanceWindowExecutionTaskInvocationsRequest(*ssm.DescribeMaintenanceWindowExecutionTaskInvocationsInput) (*request.Request, *ssm.DescribeMaintenanceWindowExecutionTaskInvocationsOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeMaintenanceWindowExecutionTasks(*ssm.DescribeMaintenanceWindowExecutionTasksInput) (*ssm.DescribeMaintenanceWindowExecutionTasksOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeMaintenanceWindowExecutionTasksWithContext(aws.Context, *ssm.DescribeMaintenanceWindowExecutionTasksInput, ...request.Option) (*ssm.DescribeMaintenanceWindowExecutionTasksOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeMaintenanceWindowExecutionTasksRequest(*ssm.DescribeMaintenanceWindowExecutionTasksInput) (*request.Request, *ssm.DescribeMaintenanceWindowExecutionTasksOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeMaintenanceWindowExecutions(*ssm.DescribeMaintenanceWindowExecutionsInput) (*ssm.DescribeMaintenanceWindowExecutionsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeMaintenanceWindowExecutionsWithContext(aws.Context, *ssm.DescribeMaintenanceWindowExecutionsInput, ...request.Option) (*ssm.DescribeMaintenanceWindowExecutionsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeMaintenanceWindowExecutionsRequest(*ssm.DescribeMaintenanceWindowExecutionsInput) (*request.Request, *ssm.DescribeMaintenanceWindowExecutionsOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeMaintenanceWindowTargets(*ssm.DescribeMaintenanceWindowTargetsInput) (*ssm.DescribeMaintenanceWindowTargetsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeMaintenanceWindowTargetsWithContext(aws.Context, *ssm.DescribeMaintenanceWindowTargetsInput, ...request.Option) (*ssm.DescribeMaintenanceWindowTargetsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeMaintenanceWindowTargetsRequest(*ssm.DescribeMaintenanceWindowTargetsInput) (*request.Request, *ssm.DescribeMaintenanceWindowTargetsOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeMaintenanceWindowTasks(*ssm.DescribeMaintenanceWindowTasksInput) (*ssm.DescribeMaintenanceWindowTasksOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeMaintenanceWindowTasksWithContext(aws.Context, *ssm.DescribeMaintenanceWindowTasksInput, ...request.Option) (*ssm.DescribeMaintenanceWindowTasksOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeMaintenanceWindowTasksRequest(*ssm.DescribeMaintenanceWindowTasksInput) (*request.Request, *ssm.DescribeMaintenanceWindowTasksOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeMaintenanceWindows(*ssm.DescribeMaintenanceWindowsInput) (*ssm.DescribeMaintenanceWindowsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeMaintenanceWindowsWithContext(aws.Context, *ssm.DescribeMaintenanceWindowsInput, ...request.Option) (*ssm.DescribeMaintenanceWindowsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeMaintenanceWindowsRequest(*ssm.DescribeMaintenanceWindowsInput) (*request.Request, *ssm.DescribeMaintenanceWindowsOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeParameters(*ssm.DescribeParametersInput) (*ssm.DescribeParametersOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeParametersWithContext(aws.Context, *ssm.DescribeParametersInput, ...request.Option) (*ssm.DescribeParametersOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeParametersRequest(*ssm.DescribeParametersInput) (*request.Request, *ssm.DescribeParametersOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeParametersPages(*ssm.DescribeParametersInput, func(*ssm.DescribeParametersOutput, bool) bool) error {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribeParametersPagesWithContext(aws.Context, *ssm.DescribeParametersInput, func(*ssm.DescribeParametersOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribePatchBaselines(*ssm.DescribePatchBaselinesInput) (*ssm.DescribePatchBaselinesOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribePatchBaselinesWithContext(aws.Context, *ssm.DescribePatchBaselinesInput, ...request.Option) (*ssm.DescribePatchBaselinesOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribePatchBaselinesRequest(*ssm.DescribePatchBaselinesInput) (*request.Request, *ssm.DescribePatchBaselinesOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribePatchGroupState(*ssm.DescribePatchGroupStateInput) (*ssm.DescribePatchGroupStateOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribePatchGroupStateWithContext(aws.Context, *ssm.DescribePatchGroupStateInput, ...request.Option) (*ssm.DescribePatchGroupStateOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribePatchGroupStateRequest(*ssm.DescribePatchGroupStateInput) (*request.Request, *ssm.DescribePatchGroupStateOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribePatchGroups(*ssm.DescribePatchGroupsInput) (*ssm.DescribePatchGroupsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribePatchGroupsWithContext(aws.Context, *ssm.DescribePatchGroupsInput, ...request.Option) (*ssm.DescribePatchGroupsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) DescribePatchGroupsRequest(*ssm.DescribePatchGroupsInput) (*request.Request, *ssm.DescribePatchGroupsOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetAutomationExecution(*ssm.GetAutomationExecutionInput) (*ssm.GetAutomationExecutionOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetAutomationExecutionWithContext(aws.Context, *ssm.GetAutomationExecutionInput, ...request.Option) (*ssm.GetAutomationExecutionOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetAutomationExecutionRequest(*ssm.GetAutomationExecutionInput) (*request.Request, *ssm.GetAutomationExecutionOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetCommandInvocation(*ssm.GetCommandInvocationInput) (*ssm.GetCommandInvocationOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetCommandInvocationWithContext(aws.Context, *ssm.GetCommandInvocationInput, ...request.Option) (*ssm.GetCommandInvocationOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetCommandInvocationRequest(*ssm.GetCommandInvocationInput) (*request.Request, *ssm.GetCommandInvocationOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetDefaultPatchBaseline(*ssm.GetDefaultPatchBaselineInput) (*ssm.GetDefaultPatchBaselineOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetDefaultPatchBaselineWithContext(aws.Context, *ssm.GetDefaultPatchBaselineInput, ...request.Option) (*ssm.GetDefaultPatchBaselineOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetDefaultPatchBaselineRequest(*ssm.GetDefaultPatchBaselineInput) (*request.Request, *ssm.GetDefaultPatchBaselineOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetDeployablePatchSnapshotForInstance(*ssm.GetDeployablePatchSnapshotForInstanceInput) (*ssm.GetDeployablePatchSnapshotForInstanceOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetDeployablePatchSnapshotForInstanceWithContext(aws.Context, *ssm.GetDeployablePatchSnapshotForInstanceInput, ...request.Option) (*ssm.GetDeployablePatchSnapshotForInstanceOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetDeployablePatchSnapshotForInstanceRequest(*ssm.GetDeployablePatchSnapshotForInstanceInput) (*request.Request, *ssm.GetDeployablePatchSnapshotForInstanceOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetDocument(*ssm.GetDocumentInput) (*ssm.GetDocumentOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetDocumentWithContext(aws.Context, *ssm.GetDocumentInput, ...request.Option) (*ssm.GetDocumentOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetDocumentRequest(*ssm.GetDocumentInput) (*request.Request, *ssm.GetDocumentOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetInventory(*ssm.GetInventoryInput) (*ssm.GetInventoryOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetInventoryWithContext(aws.Context, *ssm.GetInventoryInput, ...request.Option) (*ssm.GetInventoryOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetInventoryRequest(*ssm.GetInventoryInput) (*request.Request, *ssm.GetInventoryOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetInventorySchema(*ssm.GetInventorySchemaInput) (*ssm.GetInventorySchemaOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetInventorySchemaWithContext(aws.Context, *ssm.GetInventorySchemaInput, ...request.Option) (*ssm.GetInventorySchemaOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetInventorySchemaRequest(*ssm.GetInventorySchemaInput) (*request.Request, *ssm.GetInventorySchemaOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetMaintenanceWindow(*ssm.GetMaintenanceWindowInput) (*ssm.GetMaintenanceWindowOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetMaintenanceWindowWithContext(aws.Context, *ssm.GetMaintenanceWindowInput, ...request.Option) (*ssm.GetMaintenanceWindowOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetMaintenanceWindowRequest(*ssm.GetMaintenanceWindowInput) (*request.Request, *ssm.GetMaintenanceWindowOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetMaintenanceWindowExecution(*ssm.GetMaintenanceWindowExecutionInput) (*ssm.GetMaintenanceWindowExecutionOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetMaintenanceWindowExecutionWithContext(aws.Context, *ssm.GetMaintenanceWindowExecutionInput, ...request.Option) (*ssm.GetMaintenanceWindowExecutionOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetMaintenanceWindowExecutionRequest(*ssm.GetMaintenanceWindowExecutionInput) (*request.Request, *ssm.GetMaintenanceWindowExecutionOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetMaintenanceWindowExecutionTask(*ssm.GetMaintenanceWindowExecutionTaskInput) (*ssm.GetMaintenanceWindowExecutionTaskOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetMaintenanceWindowExecutionTaskWithContext(aws.Context, *ssm.GetMaintenanceWindowExecutionTaskInput, ...request.Option) (*ssm.GetMaintenanceWindowExecutionTaskOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetMaintenanceWindowExecutionTaskRequest(*ssm.GetMaintenanceWindowExecutionTaskInput) (*request.Request, *ssm.GetMaintenanceWindowExecutionTaskOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetMaintenanceWindowExecutionTaskInvocation(*ssm.GetMaintenanceWindowExecutionTaskInvocationInput) (*ssm.GetMaintenanceWindowExecutionTaskInvocationOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetMaintenanceWindowExecutionTaskInvocationWithContext(aws.Context, *ssm.GetMaintenanceWindowExecutionTaskInvocationInput, ...request.Option) (*ssm.GetMaintenanceWindowExecutionTaskInvocationOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetMaintenanceWindowExecutionTaskInvocationRequest(*ssm.GetMaintenanceWindowExecutionTaskInvocationInput) (*request.Request, *ssm.GetMaintenanceWindowExecutionTaskInvocationOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetMaintenanceWindowTask(*ssm.GetMaintenanceWindowTaskInput) (*ssm.GetMaintenanceWindowTaskOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetMaintenanceWindowTaskWithContext(aws.Context, *ssm.GetMaintenanceWindowTaskInput, ...request.Option) (*ssm.GetMaintenanceWindowTaskOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetMaintenanceWindowTaskRequest(*ssm.GetMaintenanceWindowTaskInput) (*request.Request, *ssm.GetMaintenanceWindowTaskOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetParameter(*ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetParameterWithContext(aws.Context, *ssm.GetParameterInput, ...request.Option) (*ssm.GetParameterOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetParameterRequest(*ssm.GetParameterInput) (*request.Request, *ssm.GetParameterOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetParameterHistory(*ssm.GetParameterHistoryInput) (*ssm.GetParameterHistoryOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetParameterHistoryWithContext(aws.Context, *ssm.GetParameterHistoryInput, ...request.Option) (*ssm.GetParameterHistoryOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetParameterHistoryRequest(*ssm.GetParameterHistoryInput) (*request.Request, *ssm.GetParameterHistoryOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetParameterHistoryPages(*ssm.GetParameterHistoryInput, func(*ssm.GetParameterHistoryOutput, bool) bool) error {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetParameterHistoryPagesWithContext(aws.Context, *ssm.GetParameterHistoryInput, func(*ssm.GetParameterHistoryOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetParameters(i *ssm.GetParametersInput) (*ssm.GetParametersOutput, error) {
	params := []*ssm.Parameter{}
	if *i.WithDecryption == true {
		for _, n := range i.Names {
			param, ok := c.mSecureStrings[*n]
			if !ok {
				return nil, fmt.Errorf("No parameter %s", n)
			}
			params = append(params, param)
		}
	} else {
		for _, n := range i.Names {
			param, ok := c.mSecureStrings[*n]
			if !ok {
				return nil, fmt.Errorf("No parameter %s", n)
			}
			params = append(params, param)
		}
	}
	return &ssm.GetParametersOutput{Parameters: params}, nil
}

func (c *MockedSSMAPI) GetParametersWithContext(aws.Context, *ssm.GetParametersInput, ...request.Option) (*ssm.GetParametersOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetParametersRequest(*ssm.GetParametersInput) (*request.Request, *ssm.GetParametersOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetParametersByPath(*ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetParametersByPathWithContext(aws.Context, *ssm.GetParametersByPathInput, ...request.Option) (*ssm.GetParametersByPathOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetParametersByPathRequest(*ssm.GetParametersByPathInput) (*request.Request, *ssm.GetParametersByPathOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetParametersByPathPages(*ssm.GetParametersByPathInput, func(*ssm.GetParametersByPathOutput, bool) bool) error {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetParametersByPathPagesWithContext(aws.Context, *ssm.GetParametersByPathInput, func(*ssm.GetParametersByPathOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetPatchBaseline(*ssm.GetPatchBaselineInput) (*ssm.GetPatchBaselineOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetPatchBaselineWithContext(aws.Context, *ssm.GetPatchBaselineInput, ...request.Option) (*ssm.GetPatchBaselineOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetPatchBaselineRequest(*ssm.GetPatchBaselineInput) (*request.Request, *ssm.GetPatchBaselineOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetPatchBaselineForPatchGroup(*ssm.GetPatchBaselineForPatchGroupInput) (*ssm.GetPatchBaselineForPatchGroupOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetPatchBaselineForPatchGroupWithContext(aws.Context, *ssm.GetPatchBaselineForPatchGroupInput, ...request.Option) (*ssm.GetPatchBaselineForPatchGroupOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) GetPatchBaselineForPatchGroupRequest(*ssm.GetPatchBaselineForPatchGroupInput) (*request.Request, *ssm.GetPatchBaselineForPatchGroupOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListAssociationVersions(*ssm.ListAssociationVersionsInput) (*ssm.ListAssociationVersionsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListAssociationVersionsWithContext(aws.Context, *ssm.ListAssociationVersionsInput, ...request.Option) (*ssm.ListAssociationVersionsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListAssociationVersionsRequest(*ssm.ListAssociationVersionsInput) (*request.Request, *ssm.ListAssociationVersionsOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListAssociations(*ssm.ListAssociationsInput) (*ssm.ListAssociationsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListAssociationsWithContext(aws.Context, *ssm.ListAssociationsInput, ...request.Option) (*ssm.ListAssociationsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListAssociationsRequest(*ssm.ListAssociationsInput) (*request.Request, *ssm.ListAssociationsOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListAssociationsPages(*ssm.ListAssociationsInput, func(*ssm.ListAssociationsOutput, bool) bool) error {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListAssociationsPagesWithContext(aws.Context, *ssm.ListAssociationsInput, func(*ssm.ListAssociationsOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListCommandInvocations(*ssm.ListCommandInvocationsInput) (*ssm.ListCommandInvocationsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListCommandInvocationsWithContext(aws.Context, *ssm.ListCommandInvocationsInput, ...request.Option) (*ssm.ListCommandInvocationsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListCommandInvocationsRequest(*ssm.ListCommandInvocationsInput) (*request.Request, *ssm.ListCommandInvocationsOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListCommandInvocationsPages(*ssm.ListCommandInvocationsInput, func(*ssm.ListCommandInvocationsOutput, bool) bool) error {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListCommandInvocationsPagesWithContext(aws.Context, *ssm.ListCommandInvocationsInput, func(*ssm.ListCommandInvocationsOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListCommands(*ssm.ListCommandsInput) (*ssm.ListCommandsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListCommandsWithContext(aws.Context, *ssm.ListCommandsInput, ...request.Option) (*ssm.ListCommandsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListCommandsRequest(*ssm.ListCommandsInput) (*request.Request, *ssm.ListCommandsOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListCommandsPages(*ssm.ListCommandsInput, func(*ssm.ListCommandsOutput, bool) bool) error {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListCommandsPagesWithContext(aws.Context, *ssm.ListCommandsInput, func(*ssm.ListCommandsOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListComplianceItems(*ssm.ListComplianceItemsInput) (*ssm.ListComplianceItemsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListComplianceItemsWithContext(aws.Context, *ssm.ListComplianceItemsInput, ...request.Option) (*ssm.ListComplianceItemsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListComplianceItemsRequest(*ssm.ListComplianceItemsInput) (*request.Request, *ssm.ListComplianceItemsOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListComplianceSummaries(*ssm.ListComplianceSummariesInput) (*ssm.ListComplianceSummariesOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListComplianceSummariesWithContext(aws.Context, *ssm.ListComplianceSummariesInput, ...request.Option) (*ssm.ListComplianceSummariesOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListComplianceSummariesRequest(*ssm.ListComplianceSummariesInput) (*request.Request, *ssm.ListComplianceSummariesOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListDocumentVersions(*ssm.ListDocumentVersionsInput) (*ssm.ListDocumentVersionsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListDocumentVersionsWithContext(aws.Context, *ssm.ListDocumentVersionsInput, ...request.Option) (*ssm.ListDocumentVersionsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListDocumentVersionsRequest(*ssm.ListDocumentVersionsInput) (*request.Request, *ssm.ListDocumentVersionsOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListDocuments(*ssm.ListDocumentsInput) (*ssm.ListDocumentsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListDocumentsWithContext(aws.Context, *ssm.ListDocumentsInput, ...request.Option) (*ssm.ListDocumentsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListDocumentsRequest(*ssm.ListDocumentsInput) (*request.Request, *ssm.ListDocumentsOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListDocumentsPages(*ssm.ListDocumentsInput, func(*ssm.ListDocumentsOutput, bool) bool) error {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListDocumentsPagesWithContext(aws.Context, *ssm.ListDocumentsInput, func(*ssm.ListDocumentsOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListInventoryEntries(*ssm.ListInventoryEntriesInput) (*ssm.ListInventoryEntriesOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListInventoryEntriesWithContext(aws.Context, *ssm.ListInventoryEntriesInput, ...request.Option) (*ssm.ListInventoryEntriesOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListInventoryEntriesRequest(*ssm.ListInventoryEntriesInput) (*request.Request, *ssm.ListInventoryEntriesOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListResourceComplianceSummaries(*ssm.ListResourceComplianceSummariesInput) (*ssm.ListResourceComplianceSummariesOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListResourceComplianceSummariesWithContext(aws.Context, *ssm.ListResourceComplianceSummariesInput, ...request.Option) (*ssm.ListResourceComplianceSummariesOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListResourceComplianceSummariesRequest(*ssm.ListResourceComplianceSummariesInput) (*request.Request, *ssm.ListResourceComplianceSummariesOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListResourceDataSync(*ssm.ListResourceDataSyncInput) (*ssm.ListResourceDataSyncOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListResourceDataSyncWithContext(aws.Context, *ssm.ListResourceDataSyncInput, ...request.Option) (*ssm.ListResourceDataSyncOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListResourceDataSyncRequest(*ssm.ListResourceDataSyncInput) (*request.Request, *ssm.ListResourceDataSyncOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListTagsForResource(*ssm.ListTagsForResourceInput) (*ssm.ListTagsForResourceOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListTagsForResourceWithContext(aws.Context, *ssm.ListTagsForResourceInput, ...request.Option) (*ssm.ListTagsForResourceOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ListTagsForResourceRequest(*ssm.ListTagsForResourceInput) (*request.Request, *ssm.ListTagsForResourceOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ModifyDocumentPermission(*ssm.ModifyDocumentPermissionInput) (*ssm.ModifyDocumentPermissionOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ModifyDocumentPermissionWithContext(aws.Context, *ssm.ModifyDocumentPermissionInput, ...request.Option) (*ssm.ModifyDocumentPermissionOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) ModifyDocumentPermissionRequest(*ssm.ModifyDocumentPermissionInput) (*request.Request, *ssm.ModifyDocumentPermissionOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) PutComplianceItems(*ssm.PutComplianceItemsInput) (*ssm.PutComplianceItemsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) PutComplianceItemsWithContext(aws.Context, *ssm.PutComplianceItemsInput, ...request.Option) (*ssm.PutComplianceItemsOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) PutComplianceItemsRequest(*ssm.PutComplianceItemsInput) (*request.Request, *ssm.PutComplianceItemsOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) PutInventory(*ssm.PutInventoryInput) (*ssm.PutInventoryOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) PutInventoryWithContext(aws.Context, *ssm.PutInventoryInput, ...request.Option) (*ssm.PutInventoryOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) PutInventoryRequest(*ssm.PutInventoryInput) (*request.Request, *ssm.PutInventoryOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) PutParameter(*ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) PutParameterWithContext(aws.Context, *ssm.PutParameterInput, ...request.Option) (*ssm.PutParameterOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) PutParameterRequest(*ssm.PutParameterInput) (*request.Request, *ssm.PutParameterOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) RegisterDefaultPatchBaseline(*ssm.RegisterDefaultPatchBaselineInput) (*ssm.RegisterDefaultPatchBaselineOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) RegisterDefaultPatchBaselineWithContext(aws.Context, *ssm.RegisterDefaultPatchBaselineInput, ...request.Option) (*ssm.RegisterDefaultPatchBaselineOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) RegisterDefaultPatchBaselineRequest(*ssm.RegisterDefaultPatchBaselineInput) (*request.Request, *ssm.RegisterDefaultPatchBaselineOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) RegisterPatchBaselineForPatchGroup(*ssm.RegisterPatchBaselineForPatchGroupInput) (*ssm.RegisterPatchBaselineForPatchGroupOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) RegisterPatchBaselineForPatchGroupWithContext(aws.Context, *ssm.RegisterPatchBaselineForPatchGroupInput, ...request.Option) (*ssm.RegisterPatchBaselineForPatchGroupOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) RegisterPatchBaselineForPatchGroupRequest(*ssm.RegisterPatchBaselineForPatchGroupInput) (*request.Request, *ssm.RegisterPatchBaselineForPatchGroupOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) RegisterTargetWithMaintenanceWindow(*ssm.RegisterTargetWithMaintenanceWindowInput) (*ssm.RegisterTargetWithMaintenanceWindowOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) RegisterTargetWithMaintenanceWindowWithContext(aws.Context, *ssm.RegisterTargetWithMaintenanceWindowInput, ...request.Option) (*ssm.RegisterTargetWithMaintenanceWindowOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) RegisterTargetWithMaintenanceWindowRequest(*ssm.RegisterTargetWithMaintenanceWindowInput) (*request.Request, *ssm.RegisterTargetWithMaintenanceWindowOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) RegisterTaskWithMaintenanceWindow(*ssm.RegisterTaskWithMaintenanceWindowInput) (*ssm.RegisterTaskWithMaintenanceWindowOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) RegisterTaskWithMaintenanceWindowWithContext(aws.Context, *ssm.RegisterTaskWithMaintenanceWindowInput, ...request.Option) (*ssm.RegisterTaskWithMaintenanceWindowOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) RegisterTaskWithMaintenanceWindowRequest(*ssm.RegisterTaskWithMaintenanceWindowInput) (*request.Request, *ssm.RegisterTaskWithMaintenanceWindowOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) RemoveTagsFromResource(*ssm.RemoveTagsFromResourceInput) (*ssm.RemoveTagsFromResourceOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) RemoveTagsFromResourceWithContext(aws.Context, *ssm.RemoveTagsFromResourceInput, ...request.Option) (*ssm.RemoveTagsFromResourceOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) RemoveTagsFromResourceRequest(*ssm.RemoveTagsFromResourceInput) (*request.Request, *ssm.RemoveTagsFromResourceOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) SendAutomationSignal(*ssm.SendAutomationSignalInput) (*ssm.SendAutomationSignalOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) SendAutomationSignalWithContext(aws.Context, *ssm.SendAutomationSignalInput, ...request.Option) (*ssm.SendAutomationSignalOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) SendAutomationSignalRequest(*ssm.SendAutomationSignalInput) (*request.Request, *ssm.SendAutomationSignalOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) SendCommand(*ssm.SendCommandInput) (*ssm.SendCommandOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) SendCommandWithContext(aws.Context, *ssm.SendCommandInput, ...request.Option) (*ssm.SendCommandOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) SendCommandRequest(*ssm.SendCommandInput) (*request.Request, *ssm.SendCommandOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) StartAutomationExecution(*ssm.StartAutomationExecutionInput) (*ssm.StartAutomationExecutionOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) StartAutomationExecutionWithContext(aws.Context, *ssm.StartAutomationExecutionInput, ...request.Option) (*ssm.StartAutomationExecutionOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) StartAutomationExecutionRequest(*ssm.StartAutomationExecutionInput) (*request.Request, *ssm.StartAutomationExecutionOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) StopAutomationExecution(*ssm.StopAutomationExecutionInput) (*ssm.StopAutomationExecutionOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) StopAutomationExecutionWithContext(aws.Context, *ssm.StopAutomationExecutionInput, ...request.Option) (*ssm.StopAutomationExecutionOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) StopAutomationExecutionRequest(*ssm.StopAutomationExecutionInput) (*request.Request, *ssm.StopAutomationExecutionOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateAssociation(*ssm.UpdateAssociationInput) (*ssm.UpdateAssociationOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateAssociationWithContext(aws.Context, *ssm.UpdateAssociationInput, ...request.Option) (*ssm.UpdateAssociationOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateAssociationRequest(*ssm.UpdateAssociationInput) (*request.Request, *ssm.UpdateAssociationOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateAssociationStatus(*ssm.UpdateAssociationStatusInput) (*ssm.UpdateAssociationStatusOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateAssociationStatusWithContext(aws.Context, *ssm.UpdateAssociationStatusInput, ...request.Option) (*ssm.UpdateAssociationStatusOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateAssociationStatusRequest(*ssm.UpdateAssociationStatusInput) (*request.Request, *ssm.UpdateAssociationStatusOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateDocument(*ssm.UpdateDocumentInput) (*ssm.UpdateDocumentOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateDocumentWithContext(aws.Context, *ssm.UpdateDocumentInput, ...request.Option) (*ssm.UpdateDocumentOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateDocumentRequest(*ssm.UpdateDocumentInput) (*request.Request, *ssm.UpdateDocumentOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateDocumentDefaultVersion(*ssm.UpdateDocumentDefaultVersionInput) (*ssm.UpdateDocumentDefaultVersionOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateDocumentDefaultVersionWithContext(aws.Context, *ssm.UpdateDocumentDefaultVersionInput, ...request.Option) (*ssm.UpdateDocumentDefaultVersionOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateDocumentDefaultVersionRequest(*ssm.UpdateDocumentDefaultVersionInput) (*request.Request, *ssm.UpdateDocumentDefaultVersionOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateMaintenanceWindow(*ssm.UpdateMaintenanceWindowInput) (*ssm.UpdateMaintenanceWindowOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateMaintenanceWindowWithContext(aws.Context, *ssm.UpdateMaintenanceWindowInput, ...request.Option) (*ssm.UpdateMaintenanceWindowOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateMaintenanceWindowRequest(*ssm.UpdateMaintenanceWindowInput) (*request.Request, *ssm.UpdateMaintenanceWindowOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateMaintenanceWindowTarget(*ssm.UpdateMaintenanceWindowTargetInput) (*ssm.UpdateMaintenanceWindowTargetOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateMaintenanceWindowTargetWithContext(aws.Context, *ssm.UpdateMaintenanceWindowTargetInput, ...request.Option) (*ssm.UpdateMaintenanceWindowTargetOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateMaintenanceWindowTargetRequest(*ssm.UpdateMaintenanceWindowTargetInput) (*request.Request, *ssm.UpdateMaintenanceWindowTargetOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateMaintenanceWindowTask(*ssm.UpdateMaintenanceWindowTaskInput) (*ssm.UpdateMaintenanceWindowTaskOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateMaintenanceWindowTaskWithContext(aws.Context, *ssm.UpdateMaintenanceWindowTaskInput, ...request.Option) (*ssm.UpdateMaintenanceWindowTaskOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateMaintenanceWindowTaskRequest(*ssm.UpdateMaintenanceWindowTaskInput) (*request.Request, *ssm.UpdateMaintenanceWindowTaskOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateManagedInstanceRole(*ssm.UpdateManagedInstanceRoleInput) (*ssm.UpdateManagedInstanceRoleOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateManagedInstanceRoleWithContext(aws.Context, *ssm.UpdateManagedInstanceRoleInput, ...request.Option) (*ssm.UpdateManagedInstanceRoleOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdateManagedInstanceRoleRequest(*ssm.UpdateManagedInstanceRoleInput) (*request.Request, *ssm.UpdateManagedInstanceRoleOutput) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdatePatchBaseline(*ssm.UpdatePatchBaselineInput) (*ssm.UpdatePatchBaselineOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdatePatchBaselineWithContext(aws.Context, *ssm.UpdatePatchBaselineInput, ...request.Option) (*ssm.UpdatePatchBaselineOutput, error) {
	panic("not implemented")
}

func (c *MockedSSMAPI) UpdatePatchBaselineRequest(*ssm.UpdatePatchBaselineInput) (*request.Request, *ssm.UpdatePatchBaselineOutput) {
	panic("not implemented")
}
