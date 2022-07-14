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

  public tableDisplayedColumns: string[] = [];
  public tableData: Record<string, string>[] = [];

  // Display Tags dynamic form
  public displayTagsForm = new FormControl();
  public tags: string[] = [];

  // Search filter dynamic form
  public searchFiltersForm = this.fb.group({
    filters: this.fb.array([]),
  });
  public searchFilters: string[] = [];

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

  private getTags() {
    let filters: string[] = ['cloud', 'tenant_id'];
    let tags: string[] = [];

    this.backendService.getTags().subscribe((response) => {
      if (response.data !== null && response.data.length !== 0) {
        response.data.forEach((tag) => {
          tags.push(tag.attributes.key);
          filters.push(tag.attributes.key);
        });
      }
      this.tags = tags;
      this.searchFilters = filters;

      this.initializeTagForm();
    });
  }

  private searchAccounts() {
    let filterMap: Map<string, string> = new Map();

    if (this.cloud !== null && this.tenant_id !== null) {
      filterMap.set('cloud', this.cloud);
      filterMap.set('tenant_id', this.tenant_id);
    }

    this.filterFormArray.controls.forEach((filter) => {
      const key = filter.value['key'];
      const value = filter.value['value'];

      filterMap.set(key, value);
    });

    this.backendService
      .getCloudAccounts(filterMap)
      .subscribe((response: CloudAccountsResponse) => {
        let accounts = [];

        if (response.data.length !== 0) {
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
          }
        }

        this.tableData = accounts;
      });
  }

  //--------------------------------------------------------------------------------------------------------------------
  // Tags Dynamic Form
  //--------------------------------------------------------------------------------------------------------------------

  public onDisplayTagsChange() {
    this.refreshTagForm();
  }

  private initializeTagForm() {
    let initialTags: string[] = [];

    // Display the first four tags by default
    for (let i = 0; i < this.tags.length; i++) {
      if (i > 4) {
        break;
      }
      initialTags.push(this.tags[i]);
    }

    this.displayTagsForm.setValue(initialTags);

    this.refreshTagForm();
  }

  private refreshTagForm() {
    this.tableDisplayedColumns = ['name', 'id']
      .concat(this.displayTagsForm.value)
      .concat(['actions']);

    if (this.cloud === null || this.tenant_id === null) {
      this.tableDisplayedColumns.unshift('cloud', 'tenant_id');
    }
  }

  //--------------------------------------------------------------------------------------------------------------------
  // Search Filters Dynamic Form
  //--------------------------------------------------------------------------------------------------------------------

  get filterFormArray() {
    return this.searchFiltersForm.controls['filters'] as FormArray;
  }

  public onFilterSubmit() {
    this.searchAccounts();
  }

  public addFilter() {
    this.filterFormArray.push(this.newFilter('', ''));
  }

  public removeFilter(i: number) {
    this.filterFormArray.removeAt(i);
  }

  private newFilter(key: string, value: string): FormGroup {
    return this.fb.group({
      key: [{ value: key, disabled: false }],
      value: [{ value: value, disabled: false }],
    });
  }
}
