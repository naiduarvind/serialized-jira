#!/usr/bin/env node
import * as cdk from '@aws-cdk/core';
import { SerializedJiraCertificateStack, SerializedJiraStack } from '../lib/serialized-jira-stack';

const app = new cdk.App();
new SerializedJiraCertificateStack(app);
new SerializedJiraStack(app, 'SerializedJiraStack', {
  certificateArn:
    "arn:aws:acm:us-east-1:474105016455:certificate/1a517d3f-5461-466b-ad32-eb6c7b75030b",
});
