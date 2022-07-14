import { Component, OnInit, ViewChild } from '@angular/core';
import {
  BackendService,
  CloudAccountItem,
  CloudAccountResponse,
  UpdateCloudAccountRequest,
} from '../backend.service';
import { ActivatedRoute } from '@angular/router';
import { FormArray, FormBuilder, FormGroup } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { AuditLogComponent } from '../audit-log/audit-log.component';

@Component({
  selector: 'app-cloud-account',
  templateUrl: './cloud-account.component.html',
  styleUrls: ['./cloud-account.component.scss'],
})
export class CloudAccountComponent implements OnInit {
  cloudAccount: CloudAccountItem | undefined;
  id: string = '';

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

  public ngOnInit(): void {
    this.id = String(this.router.snapshot.paramMap.get('id'));

    this.backendService
      .getCloudAccount(this.id)
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

  //--------------------------------------------------------------------------------------------------------------------
  // Update Account Form
  //--------------------------------------------------------------------------------------------------------------------

  public onSubmit() {
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
      .updateCloudAccount(this.id, update)
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

  public get tags() {
    return this.form.controls['tags'] as FormArray;
  }

  public addTag() {
    this.tags.push(this.newTag('', ''));
  }

  public removeTag(i: number) {
    this.tags.removeAt(i);
  }

  private newTag(key: string, value: string): FormGroup {
    return this.fb.group({
      key: [{ value: key, disabled: false }],
      value: [{ value: value, disabled: false }],
    });
  }
}
