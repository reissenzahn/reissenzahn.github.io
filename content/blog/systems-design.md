---
title: "Systems Design"
date: "2024-08-03"
tags: ["Programming"]
---

## Requirements

### Functional requirements

### Non-functional requirements

- Scalability: The ability of a system to continue to behave as expected in the face of significant upward or downward changes in demand.
  - Vertical scaling
  - Horizontal scaling
- Loose coupling: The components of a system have minimal knowledge of each other; changes to one component generally do not require changes to other components.
  - Service contracts, protocols
  - Distributed monoliths
- Resilience: A measure of how well a system withstands and recovers from errors and faults; the degree to which it can continue to operate correctly in the face of errors and faults.
  - Redundancy
  - Partial failures
  - Circuit breakers and retries
- Reliability: The ability of a system to behave as expected for a given time interval.
- Manageability: The ease (or lack thereof) with which changes can be made to the behavior of a running system, up to and including deploying system components.
  - Configuration changes
  - Feature flags
  - Credential rotation
  - Certificate renewal
  - Deployments
  - Patching
- Maintainability: The ease with which changes can be made to the functionality of a system, most often its code.
- Observability: A measure of how well the internal states of a system can be inferred from its external outputs.
  - Metrics, logging, tracing
  - "Data is not information"

## Fallacies of Distributed Computing

In 1991, L Peter Deutsh formulated the following:

1. The network is reliable: switches fail, routers get misconfigured
2. Latency is zero: it takes time to move data across a network
3. Bandwidth is infinite: a network can only handle so much data at a time
4. The network is secure: don't share secrets in plain text; encrypt everything
5. Topology doesn't change: servers and services come and go
6. There is one administrator: multiple admins lead to heterogeneous solutions
7. Transport cost is zero: moving data around costs time and money
8. The network is homogeneous: every network is different


## Queuing

Queue depth allows for trading off latency for fault tolerance during large load spikes.

Increasing queue depth due to a continuous rate of rejections is almost always bad.


## Noisy Neighbors
The basic problem this  term  refers  to  is  that  other  applications  running  on  the  same  physical  system  as yours can have a noticeable impact on your performance and resource availability.


## Perfect Accounting vs. Janitors


## Sweepers


## State Machines
