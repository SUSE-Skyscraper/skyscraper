import { Component, OnInit } from '@angular/core';
import {
  BackendService,
  CloudAccount,
  UpdateCloudAccount,
} from '../backend.service';
import { ActivatedRoute } from '@angular/router';
import { FormArray, FormBuilder, FormGroup } from '@angular/forms';

@Component({
  selector: 'app-cloud-account',
  templateUrl: './cloud-account.component.html',
  styleUrls: ['./cloud-account.component.scss'],
})
export class CloudAccountComponent implements OnInit {
  cloudAccount: CloudAccount | undefined;
  cloud: string = '';
  tenant_id: string = '';
  id: string = '';

  form = this.fb.group({
    tags: this.fb.array([]),
  });

  constructor(
    private fb: FormBuilder,
    private backendService: BackendService,
    private router: ActivatedRoute,
  ) {}

  get tags() {
    return this.form.controls['tags'] as FormArray;
  }

  newTag(key: string, value: string): FormGroup {
    return this.fb.group({
      key: [{ value: key, disabled: false }],
      value: [{ value: value, disabled: false }],
    });
  }

  onSubmit() {
    let update: UpdateCloudAccount = {
      tags_desired: {},
    };
    this.tags.controls.forEach((tag) => {
      const key = tag.value['key'];
      update.tags_desired[key] = tag.value['value'];
    });

    this.backendService
      .updateCloudAccount(this.cloud, this.tenant_id, this.id, update)
      .subscribe((cloudAccount: CloudAccount) => {
        this.cloudAccount = cloudAccount;
        this.refreshForm();
      });
  }

  ngOnInit(): void {
    this.cloud = String(this.router.snapshot.paramMap.get('cloud'));
    this.tenant_id = String(this.router.snapshot.paramMap.get('tenant_id'));
    this.id = String(this.router.snapshot.paramMap.get('id'));

    this.backendService
      .getCloudAccount(this.cloud, this.tenant_id, this.id)
      .subscribe((cloudAccount: CloudAccount) => {
        this.cloudAccount = cloudAccount;
        this.refreshForm();
      });
  }

  private refreshForm() {
    if (this.cloudAccount === undefined) {
      return;
    }
    this.tags.clear();

    Object.entries(this.cloudAccount.tags_desired).forEach(([key, value]) => {
      this.tags.push(this.newTag(key, value));
    });
  }
}
