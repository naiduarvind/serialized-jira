import path = require("path")
import * as cdk from '@aws-cdk/core';
import * as lambda from '@aws-cdk/aws-lambda';
import * as apigw from '@aws-cdk/aws-apigateway';
import * as assets from '@aws-cdk/aws-s3-assets';

export class SerializedJiraStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const myLambdaAsset = new assets.Asset(
      // @ts-ignore - this expects Construct not cdk.Construct :thinking:
      this,
      "HelloGoServerLambdaFnZip",
      {
        path: path.join(__dirname, "lambda"),
      }
    )

    const lambdaFn = new lambda.Function(this, "HelloGoServerLambdaFn", {
      code: lambda.Code.fromBucket(
        myLambdaAsset.bucket,
        myLambdaAsset.s3ObjectKey
      ),
      runtime: lambda.Runtime.GO_1_X,
      handler: "main",
    })

    // API Gateway
    new apigw.LambdaRestApi(
      // @ts-ignore - this expects Construct not cdk.Construct :thinking:
      this,
      "HelloGoServerLambdaFnEndpoint",
      {
        handler: lambdaFn,
      }
    )
  }
}
