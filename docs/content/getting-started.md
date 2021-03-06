---
title: "Getting Started"
date: 2018-05-02T00:00:00+00:00
anchor: "getting-started"
weight: 10
---

## Installation

We won't cover further details how to properly setup [Prometheus](https://prometheus.io) itself, we will only cover some basic setup based on [docker-compose](https://docs.docker.com/compose/). But if you want to run this exporter without [docker-compose](https://docs.docker.com/compose/) you should be able to adopt that to your needs.

First of all we need to prepare a configuration for [Prometheus](https://prometheus.io) that includes the exporter as a target based on a static host mapping which is just the [docker-compose](https://docs.docker.com/compose/) container name, e.g. `scw-exporter`.

{{< gist tboerger ae28347ce6c6aeebc5c74e1a00ff6d8c "prometheus.yml" >}}

After preparing the configuration we need to create the `docker-compose.yml` within the same folder, this `docker-compose.yml` starts a simple [Prometheus](https://prometheus.io) instance together with the exporter. Don't forget to update the exporter envrionment variables with the required credentials.

{{< gist tboerger ae28347ce6c6aeebc5c74e1a00ff6d8c "docker-compose.yml" >}}

Since our `latest` Docker tag always refers to the `master` branch of the Git repository you should always use some fixed version. You can see all available tags at our [DockerHub repository](https://hub.docker.com/r/promhippie/scw-exporter/tags/), there you will see that we also provide a manifest, you can easily start the exporter on various architectures without any change to the image name. You should apply a change like this to the `docker-compose.yml`:

{{< gist tboerger ae28347ce6c6aeebc5c74e1a00ff6d8c "tag.diff" >}}

If you want to access the exporter directly you should bind it to a local port, otherwise only [Prometheus](https://prometheus.io) will have access to the exporter. For debugging purpose or just to discover all available metrics directly you can apply this change to your `docker-compose.yml`, after that you can access it directly at [http://localhost:9503/metrics](http://localhost:9503/metrics):

{{< gist tboerger ae28347ce6c6aeebc5c74e1a00ff6d8c "port.diff" >}}

Finally the exporter should be configured fine, let's start this stack with [docker-compose](https://docs.docker.com/compose/), you just need to execute `docker-compose up` within the directory where you have stored the `prometheus.yml` and `docker-compose.yml`.

{{< gist tboerger ae28347ce6c6aeebc5c74e1a00ff6d8c "output.log" >}}

That's all, the exporter should be up and running. Have fun with it and hopefully you will gather interesting metrics and never run into issues. You can access the exporter at [http://localhost:9503/metrics](http://localhost:9503/metrics) and [Prometheus](https://prometheus.io) at [http://localhost:9090](http://localhost:9090). Below you can see the [Prometheus](https://prometheus.io) targets and some simple query:

{{< figure src="targets.png" title="Prometheus targets including Scaleway exporter" >}}

{{< figure src="query.png" title="Prometheus query based on Scaleway exporter" >}}

## Kubernetes

Currently we have not prepared a deployment for Kubernetes, but this is something we will provide for sure. Most interesting will be the integration into the [Prometheus Operator](https://coreos.com/operators/prometheus/docs/latest/), so stay tuned.

## Configuration

SCW_EXPORTER_TOKEN
: Access token for the Scaleway API, required for authentication

SCW_EXPORTER_ORG
: Organization for the Scaleway API, required for authentication

SCW_EXPORTER_REGION
: Region for the Scaleway API, required for authentication

SCW_EXPORTER_LOG_LEVEL
: Only log messages with given severity, defaults to `info`

SCW_EXPORTER_LOG_PRETTY
: Enable pretty messages for logging, defaults to `false`

SCW_EXPORTER_WEB_ADDRESS
: Address to bind the metrics server, defaults to `0.0.0.0:9503`

SCW_EXPORTER_WEB_PATH
: Path to bind the metrics server, defaults to `/metrics`

SCW_EXPORTER_REQUEST_TIMEOUT
: Request timeout as duration, defaults to `5s`

SCW_EXPORTER_COLLECTOR_DASHBOARD
: Enable collector for dashboard, defaults to `true`

SCW_EXPORTER_COLLECTOR_SECURITY_GROUPS
: Enable collector for security groups, defaults to `true`

SCW_EXPORTER_COLLECTOR_SERVERS
: Enable collector for servers, defaults to `true`

SCW_EXPORTER_COLLECTOR_SNAPSHOTS
: Enable collector for snapshots, defaults to `true`

SCW_EXPORTER_COLLECTOR_VOLUMES
: Enable collector for volumes, defaults to `true`

## Metrics

scw_request_duration_seconds
: Histogram of latencies for requests to the Scaleway API per collector

scw_request_failures_total
: Total number of failed requests to the Scaleway API per collector

scw_dashboard_running_servers
: Count of running servers

scw_dashboard_servers_count
: Count of owned servers

scw_dashboard_volumes_count
: Count of used volumes

scw_dashboard_images_count
: Count of used images

scw_dashboard_snapshots_count
: Count of used snapshots

scw_dashboard_ips_count
: Count of used IPs

scw_security_group_defined
: Constant value of 1 that this security group is defined

scw_security_group_enable_default
: 1 if the security group is enabled by default, 0 otherwise

scw_security_group_organization_default
: 1 if the security group is an organization default, 0 otherwise

scw_server_running
: If 1 the server is running, 0 otherwise

scw_server_created_timestamp
: Timestamp when the server have been created

scw_server_modified_timestamp
: Timestamp when the server have been modified

scw_snapshot_available
: Constant value of 1 that this snapshot is available

scw_snapshot_size_bytes
: Size of the snapshot in bytes

scw_snapshot_created_timestamp
: Timestamp when the snapshot have been created

scw_snapshot_modified_timestamp
: Timestamp when the snapshot have been modified

scw_volume_available
: Constant value of 1 that this volume is available

scw_volume_size_bytes
: Size of the volume in bytes

scw_volume_created_timestamp
: Timestamp when the volume have been created

scw_volume_modified_timestamp
: Timestamp when the volume have been modified
