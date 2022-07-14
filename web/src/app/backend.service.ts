import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import config from './app.config';

@Injectable({
  providedIn: 'root',
})
export class BackendService {
  private readonly host: string;

  constructor(private http: HttpClient) {
    this.host = config.backend.host;
  }

  getAuditLogs(
    resource_id?: string,
    resource_type?: string,
  ): Observable<AuditLogsResponse> {
    const url = new URL('/api/v1/audit_logs', this.host);
    if (resource_id !== undefined) {
      url.searchParams.append('resource_id', resource_id);
    }
    if (resource_type !== undefined) {
      url.searchParams.append('resource_type', resource_type);
    }

    return this.http.get<AuditLogsResponse>(url.href);
  }

  getProfile(): Observable<UserResponse> {
    const url = new URL('/api/v1/profile', this.host);

    return this.http.get<UserResponse>(url.href);
  }

  getTags(): Observable<TagsResponse> {
    const url = new URL(`/api/v1/tags`, this.host);

    return this.http.get<TagsResponse>(url.href);
  }

  updateTag(id: string, update: UpdateTagRequest): Observable<TagsResponse> {
    const url = new URL(`/api/v1/tags/${id}`, this.host);

    return this.http.put<TagResponse>(url.href, update);
  }

  createTag(update: CreateTagRequest): Observable<TagsResponse> {
    const url = new URL(`/api/v1/tags`, this.host);

    return this.http.post<TagResponse>(url.href, update);
  }

  getCloudTenants(): Observable<CloudTenantsResponse> {
    const url = new URL('/api/v1/cloud_tenants', this.host);

    return this.http.get<CloudTenantsResponse>(url.href);
  }

  getCloudAccount(
    cloud: string,
    tenantId: string,
    accountId: string,
  ): Observable<CloudAccountResponse> {
    const url = new URL(
      `/api/v1/cloud_tenants/cloud/${cloud}/tenant/${tenantId}/accounts/${accountId}`,
      this.host,
    );
    return this.http.get<CloudAccountResponse>(url.href);
  }

  updateCloudAccount(
    cloud: string,
    tenantId: string,
    accountId: string,
    update: UpdateCloudAccountRequest,
  ): Observable<CloudAccountResponse> {
    const url = new URL(
      `/api/v1/cloud_tenants/cloud/${cloud}/tenant/${tenantId}/accounts/${accountId}`,
      this.host,
    );
    return this.http.put<CloudAccountResponse>(url.href, update);
  }

  getCloudAccounts(
    filter: Map<string, string>,
  ): Observable<CloudAccountsResponse> {
    const url = new URL(`/api/v1/cloud_accounts`, this.host);
    if (filter !== undefined) {
      filter.forEach((value, key) => {
        url.searchParams.append(key, value);
      });
    }
    return this.http.get<CloudAccountsResponse>(url.href);
  }
}

export interface UserAttributes {
  username: string;
}

export interface UserItem {
  id: string;
  type: string;
  attributes: UserAttributes;
}

export interface UserResponse {
  data: UserItem;
}

export interface CloudTenantItem {
  cloud_provider: string;
  tenant_id: string;
  name: string;
  active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CloudTenantsResponse {
  data: CloudTenantItem[];
}

export interface CloudAccountResponse {
  data: CloudAccountItem | null;
}

export interface CloudAccountsResponse {
  data: CloudAccountItem[];
}

export interface CloudAccountItem {
  id: string;
  type: string;
  attributes: CloudAccountAttributes;
}

export interface CloudAccountAttributes {
  cloud_provider: string;
  tenant_id: string;
  account_id: string;
  name: string;
  active: boolean;
  tags_desired: { [key: string]: string };
  tags_current: { [key: string]: string };
  tags_drift_detected: boolean;
  created_at: string;
  updated_at: string;
}

export interface UpdateCloudAccountRequestData {
  tags_desired: { [key: string]: string };
}

export interface UpdateCloudAccountRequest {
  data: UpdateCloudAccountRequestData;
}

export interface TagItemAttributes {
  display_name: string;
  required: boolean;
  description: string;
  key: string;
}

export interface TagItem {
  id: string;
  type: string;
  attributes: TagItemAttributes;
}

export interface TagsResponse {
  data: TagItem[] | null;
}

export interface TagResponse {
  data: TagItem[] | null;
}

export interface UpdateTagRequestData {
  display_name: string;
  key: string;
  required: boolean;
  description: string;
}

export interface UpdateTagRequest {
  data: UpdateTagRequestData;
}

export interface CreateTagRequestData {
  display_name: string;
  key: string;
  required: boolean;
  description: string;
}

export interface CreateTagRequest {
  data: CreateTagRequestData;
}

export interface AuditLogAttributes {
  message: string;
  user_id: string;
  resource_type: string;
  resource_id: string;
  created_at: string;
}

export interface RelationshipData {
  id: string;
  type: string;
}

export interface Relationship {
  data: RelationshipData;
}

export interface AuditLogItem {
  id: string;
  type: string;
  attributes: AuditLogAttributes;
  relationships: { [key: string]: Relationship };
}

export interface AuditLogsResponse {
  data: AuditLogItem[] | null;
  included: UserItem[] | null;
}
