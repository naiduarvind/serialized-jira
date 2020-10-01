#!/usr/bin/env node
import * as cdk from '@aws-cdk/core';
import { SerializedJiraStack } from '../lib/serialized-jira-stack';

const app = new cdk.App();
new SerializedJiraStack(app, 'SerializedJiraStack');
