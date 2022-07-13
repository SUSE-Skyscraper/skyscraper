import { Component, OnInit } from '@angular/core';
import { BackendService, CloudAccountsResponse } from '../backend.service';
import { ActivatedRoute } from '@angular/router';
import { FormArray, FormBuilder, FormControl, FormGroup } from '@angular/forms';

@Component({
  selector: 'app-cloud-accounts',
  templateUrl: './cloud-accounts.component.html',
  styleUrls: ['./cloud-accounts.component.scss'],
})
export class CloudAccountsComponent implements OnInit {
  private cloud = '';
  private tenant_id = '';

  public cloudAccounts: Record<string, string>[] = [];
  public cloudTenantTags: string[] = [];
  public displayedColumns: string[] = [];
  public cloudAccountTagsForm = new FormControl();
  private cloudAccountTagsFormSelectedValues: string[] = [];

  public filtersForm = this.fb.group({
    filters: this.fb.array([]),
  });

  constructor(
    private backendService: BackendService,
    private router: ActivatedRoute,
    private fb: FormBuilder,
  ) {}

  ngOnInit(): void {
    this.cloud = String(this.router.snapshot.paramMap.get('cloud'));
    this.tenant_id = String(this.router.snapshot.paramMap.get('tenant_id'));

    this.getTags();
    this.searchAccounts();
  }

  public updateForm() {
    this.cloudAccountTagsFormSelectedValues = this.cloudAccountTagsForm.value;
    this.displayedColumns = ['name', 'id']
      .concat(this.cloudAccountTagsFormSelectedValues)
      .concat(['actions']);
  }

  get filters() {
    return this.filtersForm.controls['filters'] as FormArray;
  }

  public onFilterSubmit() {
    let filterMap: Map<string, string> = new Map();

    this.filters.controls.forEach((filter) => {
      const key = filter.value['key'];
      const value = filter.value['value'];

      filterMap.set(key, value);
    });

    this.searchAccounts(filterMap);
  }

  private searchAccounts(filter?: Map<string, string>) {
    this.backendService
      .getCloudAccounts(this.cloud, this.tenant_id, filter)
      .subscribe((response: CloudAccountsResponse) => {
        if (response.data.length === 0) {
          this.cloudAccounts = [];
        } else {
          let accounts = [];
          for (let account of response.data) {
            let object: Record<string, string> = {
              name: account.attributes.name,
              cloud: account.attributes.cloud_provider,
              tenant_id: account.attributes.tenant_id,
              account_id: account.attributes.account_id,
            };
            Object.entries(account.attributes.tags_desired).forEach(
              ([key, value]) => {
                object[key] = value;
              },
            );

            accounts.push(object);

            this.cloudAccounts = accounts;
          }
        }
      });
  }

  private getTags() {
    this.backendService
      .getCloudTenantTags(this.cloud, this.tenant_id)
      .subscribe((tags) => {
        this.cloudTenantTags = tags.tags;
        this.initializeForm();
      });
  }

  private initializeForm() {
    for (let i = 0; i < this.cloudTenantTags.length; i++) {
      if (i > 4) {
        break;
      }
      this.cloudAccountTagsFormSelectedValues.push(this.cloudTenantTags[i]);
    }

    this.cloudAccountTagsForm.setValue(this.cloudAccountTagsFormSelectedValues);
    this.updateForm();
  }

  public addFilter() {
    this.filters.push(this.newFilter('', ''));
  }

  public removeFilter(i: number) {
    this.filters.removeAt(i);
  }
  private newFilter(key: string, value: string): FormGroup {
    return this.fb.group({
      key: [{ value: key, disabled: false }],
      value: [{ value: value, disabled: false }],
    });
  }
}
