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
  private cloud: string | null = null;
  private tenant_id: string | null = null;

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
    let cloud = this.router.snapshot.paramMap.get('cloud');
    let tenant_id = this.router.snapshot.paramMap.get('tenant_id');

    if (cloud !== null && tenant_id !== null) {
      this.cloud = cloud;
      this.tenant_id = tenant_id;
    }

    this.getTags();
    this.searchAccounts();
  }

  public updateForm() {
    this.cloudAccountTagsFormSelectedValues = this.cloudAccountTagsForm.value;
    this.displayedColumns = ['name', 'id']
      .concat(this.cloudAccountTagsFormSelectedValues)
      .concat(['actions']);

    if (this.cloud === null || this.tenant_id === null) {
      this.displayedColumns.unshift('cloud', 'tenant_id');
    }
  }

  get filters() {
    return this.filtersForm.controls['filters'] as FormArray;
  }

  public onFilterSubmit() {
    this.searchAccounts();
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

  private getTags() {
    this.backendService.getTags().subscribe((response) => {
      if (response.data !== null && response.data.length !== 0) {
        let tags: string[] = [];
        response.data.forEach((tag) => {
          tags.push(tag.attributes.key);
        });
        this.cloudTenantTags = tags;
      }

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

  private searchAccounts() {
    let filterMap: Map<string, string> = new Map();

    if (this.cloud !== null && this.tenant_id !== null) {
      filterMap.set('cloud', this.cloud);
      filterMap.set('tenant_id', this.tenant_id);
    }

    this.filters.controls.forEach((filter) => {
      const key = filter.value['key'];
      const value = filter.value['value'];

      filterMap.set(key, value);
    });

    this.backendService
      .getCloudAccounts(filterMap)
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
}
