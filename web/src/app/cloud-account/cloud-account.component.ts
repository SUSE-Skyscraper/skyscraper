import { Component, OnInit } from '@angular/core';
import { BackendService, CloudAccount } from '../backend.service';
import { ActivatedRoute } from '@angular/router';
import { FormArray, FormBuilder, FormGroup } from '@angular/forms';

@Component({
  selector: 'app-cloud-account',
  templateUrl: './cloud-account.component.html',
  styleUrls: ['./cloud-account.component.scss'],
})
export class CloudAccountComponent implements OnInit {
  cloudAccount: CloudAccount | undefined;
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
      key: [{ value: key, disabled: true }],
      value: [{ value: value, disabled: false }],
    });
  }

  onSubmit() {
    console.log(JSON.stringify(this.form.value));
  }

  ngOnInit(): void {
    const cloud = String(this.router.snapshot.paramMap.get('cloud'));
    const tenant_id = String(this.router.snapshot.paramMap.get('tenant_id'));
    const id = String(this.router.snapshot.paramMap.get('id'));
    this.backendService
      .getCloudAccount(cloud, tenant_id, id)
      .subscribe((cloudAccount: CloudAccount) => {
        this.cloudAccount = cloudAccount;
        Object.keys(cloudAccount.tags_current).forEach((key) => {
          this.tags.push(this.newTag(key, cloudAccount.tags_current[key]));
        });
      });
  }
}
