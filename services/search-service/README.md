# search-service — README

## Overview
The Search Service provides full-text and semantic search across all content on the Titan platform. It indexes creators, channels, streams, recordings, clips, and community content using Elasticsearch. Supports typo-tolerant search, faceted filtering, autocomplete suggestions, and personalized result ranking based on user behavior.

## Purpose
Enable users to quickly find relevant content across the platform. Provides a unified search API that queries a multi-index Elasticsearch cluster with custom analyzers for each content type. Supports boolean queries, field boosting (title > description > tags), geo-filtering for events, date range filters, and personalized ranking where user history boosts preferred categories. Autocomplete returns real-time suggestions as users type.

## Ownership
**Team:** Discovery (eng-discovery@titan.dev)
**SLI:** 99.95% uptime, p99 search < 200ms
**Escalation:** #oncall-search
