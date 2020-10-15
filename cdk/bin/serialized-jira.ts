#!/usr/bin/env node
import * as cdk from '@aws-cdk/core';
import { SerializedJiraCertificateStack, SerializedJiraStack } from '../lib/serialized-jira-stack';

const app = new cdk.App();
new SerializedJiraCertificateStack(app);
new SerializedJiraStack(app, 'SerializedJiraStack', {
  certificateArn:
    "arn:aws:acm:us-east-1:815311713915:certificate/6faa7582-2c69-4165-86d3-b5727e15770e",
});
