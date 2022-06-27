--------------------------------------------------------------------------------------------------------------------
-- Cloud Tenants
--------------------------------------------------------------------------------------------------------------------

-- name: CreateCloudTenant :exec
insert into cloud_tenants (cloud, tenant_id, name)
values ($1, $2, $3)
on conflict (cloud, tenant_id) do update set name = $3, updated_at = now();

-- name: GetCloudTenants :many
select * from cloud_tenants
    order by cloud, tenant_id;

-- name: GetCloudTenant :one
select * from cloud_tenants
    where cloud = $1 and tenant_id = $2;

--------------------------------------------------------------------------------------------------------------------
-- Cloud Account Metadata
--------------------------------------------------------------------------------------------------------------------

-- name: CreateCloudAccountMetadata :exec
insert into cloud_account_metadata (cloud, tenant_id, account_id, name, tags_current, tags_desired)
    VALUES ($1, $2, $3, $4, $5, $6)
    ON CONFLICT (cloud, tenant_id, account_id) DO UPDATE SET
        tags_current = $5,
        tags_desired = $6,
        name = $4,
        updated_at = now();

-- name: GetCloudAllAccountMetadata :many
select * from cloud_account_metadata
    order by cloud, tenant_id, account_id;

-- name: GetCloudAllAccountMetadataForCloud :many
select * from cloud_account_metadata
    where cloud = $1
    order by tenant_id, account_id;

-- name: GetCloudAllAccountMetadataForCloudAndTenant :many
select * from cloud_account_metadata
    where cloud = $1 and tenant_id = $2
    order by account_id;

-- name: GetCloudAccountMetadata :one
select * from cloud_account_metadata
    where cloud = $1 and tenant_id = $2 and account_id = $3;