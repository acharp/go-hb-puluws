package main

import (
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an IAM role.
		role, err := iam.NewRole(ctx, "go-hb-puluws", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [{
					"Sid": "",
					"Effect": "Allow",
					"Principal": {
						"Service": "lambda.amazonaws.com"
					},
					"Action": "sts:AssumeRole"
				}]
			}`),
		})
		if err != nil {
			return err
		}

		// Attach a policy to allow writing logs to CloudWatch
		logPolicy, err := iam.NewRolePolicy(ctx, "lambda-log-policy", &iam.RolePolicyArgs{
			Role: role.Name,
			Policy: pulumi.String(`{
		        "Version": "2012-10-17",
		        "Statement": [{
		            "Effect": "Allow",
		            "Action": [
		                "logs:CreateLogGroup",
		                "logs:CreateLogStream",
		                "logs:PutLogEvents"
		            ],
		            "Resource": "arn:aws:logs:*:*:*"
		        }]
		    }`),
		})

		// Set arguments for constructing the function resource.
		args := &lambda.FunctionArgs{
			Handler: pulumi.String("handler"),
			Role:    role.Arn,
			Runtime: pulumi.String("go1.x"),
			Code:    pulumi.NewFileArchive("./handler/handler.zip"),
			Environment: &lambda.FunctionEnvironmentArgs{
				Variables: pulumi.StringMap{
					"PHONE_NUMBER": pulumi.String(os.Getenv("PHONE_NUMBER")),
					"MBIRD_TEST":   pulumi.String(os.Getenv("MBIRD_TEST")),
					"MBIRD_LIVE":   pulumi.String(os.Getenv("MBIRD_LIVE")),
				},
			},
		}
		if err != nil {
			return err
		}

		// Create the lambda using the args.
		function, err := lambda.NewFunction(
			ctx,
			"go-hb",
			args,
			pulumi.DependsOn([]pulumi.Resource{logPolicy}),
		)
		if err != nil {
			return err
		}

		// Create the cloudwatch event rule
		eventRule, err := cloudwatch.NewEventRule(ctx, "go-hb", &cloudwatch.EventRuleArgs{
			ScheduleExpression: pulumi.String("cron(0 13 * * ? *)"),
		})
		if err != nil {
			return err
		}

		_, err = lambda.NewPermission(ctx, "go-hb", &lambda.PermissionArgs{
			Action:    pulumi.String("lambda:InvokeFunction"),
			Principal: pulumi.String("events.amazonaws.com"),
			SourceArn: eventRule.Arn,
			Function:  function.Name,
		}, pulumi.Parent(eventRule))
		if err != nil {
			return err
		}

		// Add the lambda as a target for the event rule
		_, err = cloudwatch.NewEventTarget(ctx, "go-hb", &cloudwatch.EventTargetArgs{
			Rule: eventRule.Name,
			Arn:  function.Arn,
		}, pulumi.Parent(eventRule))
		if err != nil {
			return err
		}

		// Export the lambda ARN.
		ctx.Export("lambda", function.Arn)

		return nil
	})
}
