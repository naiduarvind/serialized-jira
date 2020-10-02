import { expect as expectCDK, matchTemplate, MatchStyle, haveResource } from '@aws-cdk/assert';
import * as cdk from '@aws-cdk/core';
import * as lambda from '@aws-cdk/aws-lambda';
import * as apigw from '@aws-cdk/aws-apigateway';
import * as SerializedJira from '../lib/serialized-jira-stack';

test('Lambda Created', () => {
    const app = new cdk.App();
    // INFO: WHEN
    const stack = new SerializedJira.SerializedJiraStack(app, 'SerializedJiraTestStack');
    // INFO: THEN
    expectCDK(stack).to(haveResource("AWS::Lambda::Function"));
});

test('API Gateway Created', () => {
  const app = new cdk.App();
  // INFO: WHEN
  const stack = new SerializedJira.SerializedJiraStack(app, 'SerializedJiraTestStack');
  // INFO: THEN
  expectCDK(stack).to(haveResource("AWS::SNS::Topic"));
});