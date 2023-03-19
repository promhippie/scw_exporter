scw_dashboard_images_count{zone}
: Count of used images

scw_dashboard_ips_count{zone}
: Count of used IPs

scw_dashboard_placement_groups_count{zone}
: Count of placement groups

scw_dashboard_private_nics_count{zone}
: Count of private nics

scw_dashboard_running_servers{zone}
: Count of running servers

scw_dashboard_security_groups_count{zone}
: Count of security groups

scw_dashboard_server_types_count{zone, type}
: Count of servers by type

scw_dashboard_servers_count{zone}
: Count of owned servers

scw_dashboard_snapshots_count{zone}
: Count of used snapshots

scw_dashboard_unused_ips_count{zone}
: Count of unused IPs

scw_dashboard_volumes_bssd_count{zone}
: Count of unused IPs

scw_dashboard_volumes_bssd_total_size{zone}
: Count of unused IPs

scw_dashboard_volumes_count{zone}
: Count of used volumes

scw_dashboard_volumes_lssd_count{zone}
: Count of unused IPs

scw_dashboard_volumes_lssd_total_size{zone}
: Count of unused IPs

scw_request_duration_seconds{collector}
: Histogram of latencies for requests to the api per collector

scw_request_failures_total{collector}
: Total number of failed requests to the api per collector

scw_security_group_created_timestamp{id, name, zone, org, project}
: Timestamp when the security group have been created

scw_security_group_defined{id, name, zone, org, project}
: Constant value of 1 that this security group is defined

scw_security_group_enable_default{id, name, zone, org, project}
: 1 if the security group is enabled by default, 0 otherwise

scw_security_group_inbound_default_policy{id, name, zone, org, project}
: 1 if the security group inbound default policy is accept, 0 otherwise

scw_security_group_modified_timestamp{id, name, zone, org, project}
: Timestamp when the security group have been modified

scw_security_group_outbound_default_policy{id, name, zone, org, project}
: 1 if the security group outbound default policy is accept, 0 otherwise

scw_security_group_project_default{id, name, zone, org, project}
: 1 if the security group is an project default, 0 otherwise

scw_security_group_servers_count{id, name, zone, org, project}
: Number of servers attached to the security group

scw_security_group_stateful{id, name, zone, org, project}
: 1 if the security group is stateful by default, 0 otherwise

scw_server_created_timestamp{id, name, zone, org, project, type, private_ip, public_ip, arch}
: Timestamp when the server have been created

scw_server_modified_timestamp{id, name, zone, org, project, type, private_ip, public_ip, arch}
: Timestamp when the server have been modified

scw_server_private_nic_count{id, name, zone, org, project, type, private_ip, public_ip, arch}
: Number of private nics attached

scw_server_state{id, name, zone, org, project, type, private_ip, public_ip, arch}
: If 1 the server is running, depending on the state otherwise

scw_server_volume_count{id, name, zone, org, project, type, private_ip, public_ip, arch}
: Number of volumes attached

scw_snapshot_available{id, name, zone, org, project}
: Constant value of 1 that this snapshot is available

scw_snapshot_created_timestamp{id, name, zone, org, project}
: Timestamp when the snapshot have been created

scw_snapshot_modified_timestamp{id, name, zone, org, project}
: Timestamp when the snapshot have been modified

scw_snapshot_size_bytes{id, name, zone, org, project}
: Size of the snapshot in bytes

scw_snapshot_state{id, name, zone, org, project}
: State of the snapshot

scw_snapshot_type{id, name, zone, org, project}
: Type of the snapshot

scw_volume_available{id, name, zone, org, project}
: Constant value of 1 that this volume is available

scw_volume_created_timestamp{id, name, zone, org, project}
: Timestamp when the volume have been created

scw_volume_modified_timestamp{id, name, zone, org, project}
: Timestamp when the volume have been modified

scw_volume_size_bytes{id, name, zone, org, project}
: Size of the volume in bytes

scw_volume_state{id, name, zone, org, project}
: State of the snapshot

scw_volume_type{id, name, zone, org, project}
: Type of the snapshot
