import * as cdk from '@aws-cdk/core';
import * as lambda from '@aws-cdk/aws-lambda';
import * as apigw from '@aws-cdk/aws-apigateway';

export class SerializedJiraStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const lambdaFn = new lambda.Function(this, "SerializedJiraLambdaFn", {
      code: lambda.Code.fromAsset('./../packages/serialized-jira', {exclude: ['*.go', 'BUILD.bazel', 'static/']}),
      runtime: lambda.Runtime.GO_1_X,
      handler: "main",
    })

    // API Gateway
    new apigw.LambdaRestApi(this, 'SerializedJiraAPIEndpoint', {
      handler: lambdaFn
    });
  }
}
