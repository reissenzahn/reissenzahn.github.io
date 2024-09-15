---
title: "Monitoring"
date: "2024-08-03"
tags: ["Programming"]
subtitle: "Some monitoring notes"
---

## Overview

- More visibility is required as system become more complicated.
- Observability is a system property that reflects how well the internal states of a system can be inferred from knowledge of its external outputs. A system can be considered observable when it is possible to quickly and consistently ask novel questions about it with minimal prior knowledge and without having to re-instrument or build new code. An observable system lets you ask it questions that you have not thought of yet.
  - Traditionally, monitoring focuses on asking questions in the hope of identifying or predicting some expected or previously observed failure modes.
  - At a certain level of complexity, the number of "unknown unknowns" starts to overwhelm the number of "known unknowns". Failures are more often unpredicted and monitoring for every possible failure mode becomes effectively impossible.
  - Understanding all possible failure (or non-failure) states in a complex system is pretty much impossible. The first step to achieving observability is to stop looking for specific, expected failure modes.
  - The three pillars of observability:
    - Tracing: Follows a request as it propagates through a system, allowing the entire end-to-end request flow to be reconstructed as a DAG called a trace.
    - Metrics: The collection of numerical data points representing the state of various aspects of a system at specific points in time.  Collections  of  data  points, representing  observations  of  the  same  subject  at  various  times,  are  particularly useful  for  visualization  and  mathematical  analysis,  and  can  be  used  to  highlight trends, identify anomalies, and predict future behavior
    - Logging: Logging is the process of appending records of noteworthy events to an immuta‐ ble  record—the  log—for  later  review  or  analysis.  A  log  can  take  a  variety  of forms, from a continuously appended file on disk to a full-text search engine like Elasticsearch. Logs provides valuable, context-rich insight into application- specific events emitted by processes. However, it’s important that log entries are properly structured; not doing so can sharply limit their utility
  - A truly observable system will interweave tracing, metrics and logging so that each can reference the others. For instance, metrics might be used to track down a subset of misbehaving traces and those traces might highlight logs that could help to find the underlying cause of the behavior.
  -

## Logs

Logging  is  the  act  of  recording  events  that  occur  during  the  running  of  a  program.
It is often an undervalued activity in programming because it is additional work that
has little immediate payback for the programmer.
During  the  normal  operations  of  a  program,  logging  is  an  overhead,  taking  up
processing  cycles  to  write  to  a  file,  database,  or  even  to  the  screen.  In  addition,
unmanaged  logs  can  cause  problems.  The  classic  case  of  logfiles  getting  so  big  that
they take up all the available disk space and crash the server is too real and happens
too often.
However, when something happens, and you want to find out the sequence of events
that led to it, logs become an invaluable diagnostic resource. Logs can also be moni‐
tored in real time, and alerts can be sent out when needed.

### Log Rotation

### Log Retention

### Sensitive User Data

- Logging complex types

> My analogy for user-data is that its like a pool of toxic sludge. You might be willing to poke it with a stick to see how deep it is, but you definitely wouldn't dive into it.


### Compliance



## Metrics

### Faults
Zero-ed metrics.

### Errors

### Latency

### Saturation

#### Starvation

Note our technique here for identifying the starvation: a metric. Starvation makes for
a good argument for recording and sampling metrics. One of the ways you can detect
and solve starvation is by logging when work is accomplished, and then determining
if your rate of work is as high as you expect it.

### Cache Metrics

https://guava.dev/releases/snapshot/api/docs/com/google/common/cache/CacheStats.html

### Queue Metrics


### Database Query Metrics


## Alarms

### Severity levels

### Aggregate hierarchies

### Alarm fatigue

### Anomaly detection
Very skeptical.


## Tracing

By tracking requests as they propagate through the system (even across process, network and security boundaries) tracing can help you to pinpoint component failures, identify performance bottlenecks and analyze service dependencies.

There are two fundamental concepts:

- Spans: A span describes a unit of work performed by a request, such as a fork in the execution flow or hop across the network, as it propagates through a system. Each span has an associated name, a start time and a duration. They can be and typically are nested and ordered to model casual relationships.
- Traces: 


## OpenTelemetry
- OTel is an effort to standardize how telemetry data--traces, metrics and (eventually) logs--are expressed, collected and transferred.
- This seeks to unify the instrumentation space around a single vendor-neutral specification that standardizes how telemetry data is collected and sent to backend platforms.
- It has the following core components:
  - Specifications: describe the requirements and expectations for all OTel APIs, SDKs and protocols.
  - API: language-specific interfaces and implementations based on the specifications that can be used to add OTel to an application.
  - SDK: The concrete OTel implementations that sit between the APIs and the Exporters, providing functionality like state tracking and batching data for transmission.
  - Exporters: In-process SDK plugins that are capable of sending data to a specific destination. This decouples the instrumentation from the backend.
  - Collector: An optional vendor-agnostic service that can receive and process telemetry data before forwarding it to one or more destinations.
