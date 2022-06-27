import { Component, OnInit } from '@angular/core';
import { BackendService, CloudTenant } from '../backend.service';

@Component({
  selector: 'app-cloud-tenants',
  templateUrl: './cloud-tenants.component.html',
  styleUrls: ['./cloud-tenants.component.scss'],
})
export class CloudTenantsComponent implements OnInit {
  public cloudTenants: CloudTenant[] = [];
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
      this.cloudTenants = tenants;
    });
  }
}
