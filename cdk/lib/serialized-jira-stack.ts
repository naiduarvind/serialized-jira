import * as cdk from '@aws-cdk/core';
import * as kms from '@aws-cdk/aws-kms';
import * as lambda from '@aws-cdk/aws-lambda';
import * as apigw from '@aws-cdk/aws-apigateway';
import { Bucket, HttpMethods } from '@aws-cdk/aws-s3';
import * as acm from "@aws-cdk/aws-certificatemanager";
import * as s3deploy from '@aws-cdk/aws-s3-deployment';
import { EndpointType } from '@aws-cdk/aws-apigateway';
import { CfnOutput, RemovalPolicy } from "@aws-cdk/core";
import { Certificate } from '@aws-cdk/aws-certificatemanager';

const domainName = "jira.thebility.engineer";

export class SerializedJiraCertificateStack extends cdk.Stack {
  constructor(scope: cdk.Construct) {
    super(scope, 'SerializedJiraCertificateStack', {
      env: { region: "us-east-1" },
    });

    const acmCert = new acm.Certificate(this, "SerializedJiraCustomDomainCertificate", {
      domainName: domainName,
      validation: acm.CertificateValidation.fromDns(),
    });

    new CfnOutput(this, "SerializedJiraCertificateArn", {
      value: acmCert.certificateArn,
    });
  }
}

interface SerializedJiraProps {
  certificateArn: string;
}

export class SerializedJiraStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props: SerializedJiraProps) {
    super(scope, id);

    const s3Bucket = new Bucket(this, "SerializedJiraAssetBucket", {
      bucketName: "serialized-jira-assets.thebility.engineer",
      cors: [
        {
          allowedOrigins: ['*'],
          allowedMethods: [ HttpMethods.GET ],
          maxAge: 3000,
          allowedHeaders: ['*']
        }
      ],
      publicReadAccess: true,
      versioned: true,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });

    new s3deploy.BucketDeployment(this, 'SerializedJiraDeployAsset', {
      sources: [s3deploy.Source.asset('./../packages/serialized-jira/static')],
      destinationBucket: s3Bucket,
      destinationKeyPrefix: 'static'
    });

    const lambdaFn = new lambda.Function(this, "SerializedJiraLambdaFn", {
      code: lambda.Code.fromAsset('./../packages/serialized-jira', { exclude: ['*.go', '*.bazel', 'static/**'] }),
      runtime: lambda.Runtime.GO_1_X,
      handler: "main",
      environment: {
        SECRETHUB_IDENTITY_PROVIDER: "aws",
      }
    });

    const kmsKey = new kms.Key(this, "SerializedJiraKMSKey", {
      description: "KMS Key used by Secret Hub by Serialized Jira Lambda",
      removalPolicy: RemovalPolicy.DESTROY,
      trustAccountIdentities: true,
    });
    kmsKey.addAlias("serialized-jira-service-key");
    kmsKey.grantEncryptDecrypt(lambdaFn);

    const customDomain = new apigw.DomainName(this, "SerializedJiraCustomDomain", {
      domainName: domainName,
      certificate: Certificate.fromCertificateArn(this, 'Certificate', props.certificateArn),
      endpointType: EndpointType.EDGE,
      securityPolicy: apigw.SecurityPolicy.TLS_1_2,
    });

    const apiGw = new apigw.LambdaRestApi(this, "SerializedJiraAPIEndpoint", {
      handler: lambdaFn,
    });

    customDomain.addBasePathMapping(apiGw);
  }
}