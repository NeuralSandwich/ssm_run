package awsCmds

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/pkg/errors"
)

// SSMSendCmd sends a command to an EC2 windows instance
func SSMSendCmd(windows bool, instanceID, cmdString string) (string, error) {
	svc, err := initSSMClient()
	if err != nil {
		return "", errors.Wrap(err, "Failed to Initialise SSM Client")
	}

	// Set the DocumentName based on the platform type of the instance
	docName := "AWS-RunShellScript"
	if windows {
		docName = "AWS-RunPowerShellScript"
	}

	cmdMap := make(map[string][]string)
	cmdMap["commands"] = []string{cmdString}

	sendCmdInput := ssm.SendCommandInput{
		DocumentName: aws.String(docName),
		InstanceIds:  []string{instanceID},
		Parameters:   cmdMap,
	}

	sendCmdRequest := svc.SendCommandRequest(&sendCmdInput)

	sendCmdResponse, err := sendCmdRequest.Send()
	if err != nil {
		return "", errors.Wrap(err, "Problem making SSM Command Request")
	}
	return *sendCmdResponse.Command.CommandId, nil
}

// SSMCommandInfo describes an SSM Command Invocation
type SSMCommandInfo struct {
	InstanceID    string
	CommandID     string
	CommandStatus string
	StdOut        string
	StdErr        string
}

// SSMGetCmd queries SSM for a command invocation
func SSMGetCmd(instanceID, cmdID string) (SSMCommandInfo, error) {
	var cmdInfo SSMCommandInfo

	svc, err := initSSMClient()
	if err != nil {
		return cmdInfo, errors.Wrap(err, "Failed to Initialise SSM Client")
	}

	getCmdInfo := ssm.GetCommandInvocationInput{
		CommandId:  aws.String(cmdID),
		InstanceId: aws.String(instanceID),
	}

	getCmdRequest := svc.GetCommandInvocationRequest(&getCmdInfo)

	invocation, err := getCmdRequest.Send()
	if err != nil {
		return cmdInfo, errors.Wrap(err, "Problem initially getting SSM Command Invocation")
	}

	cmdInfo.InstanceID = *invocation.InstanceId
	cmdInfo.CommandID = *invocation.CommandId
	cmdInfo.CommandStatus = string(invocation.Status)
	cmdInfo.StdOut = *invocation.StandardOutputContent
	cmdInfo.StdErr = *invocation.StandardErrorContent

	return cmdInfo, nil
}

// EC2InstanceStatus defines the Status of an EC2 instance
type EC2InstanceStatus struct {
	InstanceID     string
	SystemStatus   string
	InstanceStatus string
}

// Ec2DescribeInstanceStatus is a wrapper that returns the healthcheck status of an EC2 instance
func Ec2DescribeInstanceStatus(instanceID string) (EC2InstanceStatus, error) {
	var instStatus EC2InstanceStatus

	svc, err := initEC2Client()
	if err != nil {
		return instStatus, errors.Wrap(err, "Failed to Initialise SSM Client")
	}

	statusInput := ec2.DescribeInstanceStatusInput{
		InstanceIds: []string{instanceID},
	}

	statusRequest := svc.DescribeInstanceStatusRequest(&statusInput)

	status, err := statusRequest.Send()
	if err != nil {
		return instStatus, errors.Wrap(err, "Problem initially getting SSM Command Invocation")
	}

	instStatus.InstanceID = instanceID
	instStatus.SystemStatus = string(status.InstanceStatuses[0].InstanceStatus.Status)
	instStatus.InstanceStatus = string(status.InstanceStatuses[0].SystemStatus.Status)

	return instStatus, nil
}

// initialise the SSM Client
func initSSMClient() (*ssm.SSM, error) {
	// initialse with default config
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to laod default AWS SDK Config")
	}

	// Set the AWS Region that the service clients should use
	cfg.Region = endpoints.EuWest1RegionID

	return ssm.New(cfg), nil
}

// Initialise the EC2 Client
func initEC2Client() (*ec2.EC2, error) {
	// initialse with default config
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to laod default AWS SDK Config")
	}

	// Set the AWS Region that the service clients should use
	cfg.Region = endpoints.EuWest1RegionID

	return ec2.New(cfg), nil
}
