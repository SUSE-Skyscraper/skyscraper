import { Component, OnInit } from '@angular/core';
import { BackendService, CloudAccount } from '../backend.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-cloud-accounts',
  templateUrl: './cloud-accounts.component.html',
  styleUrls: ['./cloud-accounts.component.scss'],
})
export class CloudAccountsComponent implements OnInit {
  public cloudAccounts: CloudAccount[] = [];
  public displayedColumns: string[] = [
    'cloud_provider',
    'tenant_id',
    'account_id',
    'name',
    'created_at',
    'updated_at',
    'active',
    'actions',
  ];

  constructor(
    private backendService: BackendService,
    private router: ActivatedRoute,
  ) {}

  ngOnInit(): void {
    const cloud = String(this.router.snapshot.paramMap.get('cloud'));
    const tenant_id = String(this.router.snapshot.paramMap.get('tenant_id'));
    this.backendService
      .getCloudAccounts(cloud, tenant_id)
      .subscribe((cloudAccounts: CloudAccount[]) => {
        this.cloudAccounts = cloudAccounts;
      });
  }
}
