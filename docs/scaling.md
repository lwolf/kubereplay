# Scaling

Lets say we have the following objects:
* original deployment - web-api
* shadow deployment - web-api-gor
* harvester - example-harvester

There are several possible scenarios regarding the scaling.

* Changing the number of replicas of the real deployment (web-api). 
During reconciliation, the number of replicas in this deployment will be used to calculate the new split.
 Behavior is the same as when new deployment is created.

| Harvester segment  | Deployment Replicas Change   | Deployment Replicas | Shadow Deployment Replicas | 
| ------------------ | --------------------- | -------------------------- | ----------------------- |
| 50 | 5->20 | 5(10) | 5(10) |
| 50 | 10->6 | 10(3)  | 10(3) |
| 100 | 0->5 | 0(0)  | ?(5) |
| 10 | 1->10 | 1(9)  | 1(1) |
  
* Changing the number of replicas of the shadow deployment (web-api-gor). 
Number of replicas is fully controlled by Kubereplay. Any changes to replica count will be reverted during the next reconciliation cycle.

* Changing `segment` value in the Harvester spec.
During reconciliation, both deployments will be scaled. New replica counts will be calculated using new segment value and sum of replicas of deployments.

| Harvester segment change  | Deployment Replicas | Shadow Deployment Replicas | 
| ------------------ | --------------------- | -------------------------- |
| 50->20 | 10(18) | 10(2) |
| 50->10 | 3(5)  | 3(1) |
| 20->100 | 1(0)  | 5(6) |

