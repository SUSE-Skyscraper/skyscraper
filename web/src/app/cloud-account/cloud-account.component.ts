import { AfterViewInit, Component, OnInit, ViewChild } from '@angular/core';
import {
  BackendService,
  CloudAccountItem,
  CloudAccountResponse,
  UpdateCloudAccountRequest,
} from '../backend.service';
import { ActivatedRoute } from '@angular/router';
import { FormArray, FormBuilder, FormGroup } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import {
  AuditLogComponent,
  AuditLogTableItem,
} from '../audit-log/audit-log.component';

@Component({
  selector: 'app-cloud-account',
  templateUrl: './cloud-account.component.html',
  styleUrls: ['./cloud-account.component.scss'],
})
export class CloudAccountComponent implements OnInit {
  cloudAccount: CloudAccountItem | undefined;
  cloud: string = '';
  tenant_id: string = '';
  account_id: string = '';
  auditLogs: AuditLogTableItem[] = [];
  auditLogColumns: string[] = ['user_id', 'message', 'created_at'];

  @ViewChild(AuditLogComponent)
  auditLogComponent: AuditLogComponent | undefined;

  form = this.fb.group({
    tags: this.fb.array([]),
  });

  constructor(
    private fb: FormBuilder,
    private backendService: BackendService,
    private router: ActivatedRoute,
    private snackBar: MatSnackBar,
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
    let update: UpdateCloudAccountRequest = {
      data: {
        tags_desired: {},
      },
    };
    this.tags.controls.forEach((tag) => {
      const key = tag.value['key'];
      update.data.tags_desired[key] = tag.value['value'];
    });

    this.backendService
      .updateCloudAccount(this.cloud, this.tenant_id, this.account_id, update)
      .subscribe((response: CloudAccountResponse) => {
        if (response.data !== null) {
          this.cloudAccount = response.data;
          this.refreshPage();
          this.snackBar.open('Tags Updated', 'close', {
            horizontalPosition: 'center',
            verticalPosition: 'top',
            duration: 10000,
          });
          this.auditLogComponent?.ngOnChanges();
        }
      });
  }

  ngOnInit(): void {
    this.cloud = String(this.router.snapshot.paramMap.get('cloud'));
    this.tenant_id = String(this.router.snapshot.paramMap.get('tenant_id'));
    this.account_id = String(this.router.snapshot.paramMap.get('account_id'));

    this.backendService
      .getCloudAccount(this.cloud, this.tenant_id, this.account_id)
      .subscribe((response: CloudAccountResponse) => {
        if (response.data !== null) {
          this.cloudAccount = response.data;
        }
        this.refreshPage();
      });
  }

  private refreshPage() {
    if (this.cloudAccount === undefined) {
      return;
    }
    this.tags.clear();

    Object.entries(this.cloudAccount.attributes.tags_desired).forEach(
      ([key, value]) => {
        this.tags.push(this.newTag(key, value));
      },
    );
  }

  addTag() {
    this.tags.push(this.newTag('', ''));
  }

  removeTag(i: number) {
    this.tags.removeAt(i);
  }
}
