# analytics-service — README

## Overview
The Analytics Service is the data processing and reporting engine for the Titan platform. It ingests events from all services via Kafka, processes them through streaming pipelines (Kafka Streams for real-time, Spark for batch), and serves aggregated metrics, dashboards, and exportable reports for creators, viewers, and internal operations teams.

## Purpose
Transform raw event data into actionable insights. Provides real-time dashboards for live stream metrics (concurrent viewers, chat rate, gift rate), historical trend analysis (growth charts, retention cohorts), audience demographics (geo, device, language), revenue reporting (MRR, ARPU, LTV), and custom funnel analysis for conversion optimization. Built on a lambda architecture with ClickHouse for fast time-series queries.

## Ownership
**Team:** Data Platform (eng-data@titan.dev)
**SLI:** 99.95% uptime, p99 dashboard load < 2s
**Escalation:** #oncall-analytics
