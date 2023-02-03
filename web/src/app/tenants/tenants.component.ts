import { Component, OnInit } from '@angular/core';
import { BackendService, CloudTenantItem } from '../backend.service';

@Component({
  selector: 'app-cloud-tenants',
  templateUrl: './tenants.component.html',
  styleUrls: ['./tenants.component.scss'],
})
export class TenantsComponent implements OnInit {
  public cloudTenants: CloudTenantItem[] = [];
  public displayedColumns: string[] = [
    'cloud_provider',
    'tenant_id',
    'accounts',
    'name',
    'created_at',
    'updated_at',
    'active',
  ];

  constructor(private backendService: BackendService) {}

  ngOnInit(): void {
    this.backendService.getCloudTenants().subscribe((tenants) => {
      this.cloudTenants = tenants.data;
    });
  }
}
