import * as cdk from '@aws-cdk/core';
import { CfnOutput } from "@aws-cdk/core";
import * as lambda from '@aws-cdk/aws-lambda';
import * as apigw from '@aws-cdk/aws-apigateway';
import * as acm from "@aws-cdk/aws-certificatemanager";
import { EndpointType } from '@aws-cdk/aws-apigateway';
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

    const certArn = acmCert.certificateArn;

    new CfnOutput(this, "SerializedJiraCertificateArn", {
      value: certArn,
    });
  }
}

interface SerializedJiraProps {
  certificateArn: string;
}

export class SerializedJiraStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props: SerializedJiraProps) {
    super(scope, id);

    // TODO: Upload static files into the serialized-jira.thebility.engineer S3 bucket

    const lambdaFn = new lambda.Function(this, "SerializedJiraLambdaFn", {
      // TODO: Exclude certain artifacts, files, and directories
      code: lambda.Code.fromAsset('./../packages/serialized-jira'),
      runtime: lambda.Runtime.GO_1_X,
      handler: "main",
    })

    const customDomain = new apigw.DomainName(this, 'SerializedJiraCustomDomain', {
      domainName: domainName,
      certificate: Certificate.fromCertificateArn(this, 'SerializedJiraCertificate', props.certificateArn),
      endpointType: EndpointType.EDGE,
    });

    new apigw.LambdaRestApi(this, 'SerializedJiraAPIEndpoint', {
      handler: lambdaFn,
    });
  }
}
