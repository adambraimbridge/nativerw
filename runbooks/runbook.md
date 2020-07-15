# UPP - Native read/write

The purpose of this service is to read/write content to and from Mongo native store.

## Code

nativestorereaderwriter

## Primary URL

<https://upp-prod-publish-glb.upp.ft.com/__nativerw/>

## Service Tier

Platinum

## Lifecycle Stage

Production

## Delivered By

content

## Supported By

content

## Known About By

- dimitar.terziev
- hristo.georgiev
- elitsa.pavlova
- elina.kaneva
- kalin.arsov
- ivan.nikolov
- miroslav.gatsanoga
- mihail.mihaylov
- tsvetan.dimitrov
- georgi.ivanov
- robert.marinov

## Host Platform

AWS

## Architecture

Writes any raw content/data from native CMS in mongoDB without transformation. The same data can then be read from here 
just like from the original CMS.

## Contains Personal Data

No

## Contains Sensitive Data

No

## Dependencies

- upp-mongodb

## Failover Architecture Type

ActivePassive

## Failover Process Type

FullyAutomated

## Failback Process Type

FullyAutomated

## Failover Details

The service is Publish cluster.
The failover guide for the cluster is located here:
<https://github.com/Financial-Times/upp-docs/blob/master/failover-guides/publishing-cluster>

## Data Recovery Process Type

NotApplicable

## Data Recovery Details

The service does not store data, so it does not require any data recovery steps.

## Release Process Type

PartiallyAutomated

## Rollback Process Type

Manual

## Release Details

Manual failover is needed when a new version of
the service is deployed to production.
Otherwise, an automated failover is going to take place when releasing.
For more details about the failover process please see: <https://github.com/Financial-Times/upp-docs/blob/master/failover-guides/publishing-cluster>

## Key Management Process Type

Manual

## Key Management Details

To access the service clients need to provide basic auth credentials.
To rotate credentials you need to login to a particular cluster and update varnish-auth secrets.

## Monitoring

Service in UPP K8S Publish clusters:

- Publish-Prod-EU health: <https://upp-prod-publish-eu.ft.com/__health/__pods-health?service-name=nativerw>
- Publish-Prod-US health: <https://upp-prod-publish-eu.ft.com/__health/__pods-health?service-name=nativerw>

## First Line Troubleshooting

<https://github.com/Financial-Times/upp-docs/tree/master/guides/ops/first-line-troubleshooting>

## Second Line Troubleshooting

Please refer to the GitHub repository README for troubleshooting information.
