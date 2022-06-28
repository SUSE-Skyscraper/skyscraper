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

  getProfile(): Observable<Profile> {
    const url = new URL('/api/v1/profile', this.host);
    return this.http.get<Profile>(url.href);
  }

  getCloudTenants(): Observable<CloudTenant[]> {
    const url = new URL('/api/v1/cloud_tenants', this.host);
    return this.http.get<CloudTenant[]>(url.href);
  }

  getCloudAccount(
    cloud: string,
    tenantId: string,
    accountId: string,
  ): Observable<CloudAccount> {
    const url = new URL(
      `/api/v1/cloud_tenants/cloud/${cloud}/tenant/${tenantId}/accounts/${accountId}`,
      this.host,
    );
    return this.http.get<CloudAccount>(url.href);
  }

  getCloudAccounts(
    cloud: string,
    tenantId: string,
  ): Observable<CloudAccount[]> {
    const url = new URL(
      `/api/v1/cloud_tenants/cloud/${cloud}/tenant/${tenantId}/accounts`,
      this.host,
    );
    return this.http.get<CloudAccount[]>(url.href);
  }
}

export interface Profile {
  email: string;
}

export interface CloudTenant {
  cloud_provider: string;
  tenant_id: string;
  name: string;
  active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CloudAccount {
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
