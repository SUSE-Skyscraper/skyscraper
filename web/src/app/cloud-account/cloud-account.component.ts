import { Component, OnInit } from '@angular/core';
import { BackendService, CloudAccount } from '../backend.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-cloud-account',
  templateUrl: './cloud-account.component.html',
  styleUrls: ['./cloud-account.component.scss'],
})
export class CloudAccountComponent implements OnInit {
  public cloudAccount: CloudAccount | undefined;

  constructor(
    private backendService: BackendService,
    private router: ActivatedRoute,
  ) {}

  ngOnInit(): void {
    const cloud = String(this.router.snapshot.paramMap.get('cloud'));
    const tenant_id = String(this.router.snapshot.paramMap.get('tenant_id'));
    const id = String(this.router.snapshot.paramMap.get('id'));
    this.backendService
      .getCloudAccount(cloud, tenant_id, id)
      .subscribe((cloudAccount: CloudAccount) => {
        this.cloudAccount = cloudAccount;
      });
  }
}
