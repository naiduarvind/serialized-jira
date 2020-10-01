import * as cdk from '@aws-cdk/core';
import * as lambda from '@aws-cdk/aws-lambda';
import * as apigw from '@aws-cdk/aws-apigateway';

export class SerializedJiraStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // INFO: defines AWS Lambda resource for Serialized Jira
    const serializedJira = new lambda.Function(this, 'SerializedJira', {
      runtime: lambda.Runtime.GO_1_X, // INFO: execution environment
      code: lambda.Code.fromAsset('../packages/serialized-jira'), // INFO: code assets from lambda directory
      handler: 'index.handler' // INFO: file is "serialized-jira", function is "handler"
    });

    new apigw.LambdaRestApi(this, 'Endpoint', {
      handler: serializedJira
    });
  }
}
