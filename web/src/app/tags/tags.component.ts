import { Component, OnInit } from '@angular/core';
import {
  BackendService,
  CreateTagRequest,
  TagItem,
  UpdateTagRequest,
  UpdateTagRequestData,
} from '../backend.service';
import { FormBuilder, FormControl } from '@angular/forms';

@Component({
  selector: 'app-tags',
  templateUrl: './tags.component.html',
  styleUrls: ['./tags.component.scss'],
})
export class TagsComponent implements OnInit {
  public tags: TagItem[] | null = null;
  public newTagForm = this.fb.group({
    display_name: new FormControl(),
    key: new FormControl(),
    description: new FormControl(),
  });

  constructor(
    private backendService: BackendService,
    private fb: FormBuilder,
  ) {}

  ngOnInit(): void {
    this.refreshData();
  }

  public onTagsCreateSubmit() {
    let displayName: string = this.newTagForm.get('display_name')?.value;
    let key: string = this.newTagForm.get('key')?.value;
    let description: string = this.newTagForm.get('description')?.value;

    let create: CreateTagRequest = {
      data: {
        display_name: displayName,
        key: key,
        required: false,
        description: description,
      },
    };

    this.backendService.createTag(create).subscribe((_) => {
      this.refreshData();
      this.newTagForm.reset();
    });
  }

  public onTagsUpdateSubmit(tag: TagItem) {
    let update: UpdateTagRequest = {
      data: {
        display_name: tag.attributes.display_name,
        key: tag.attributes.key,
        description: tag.attributes.description,
      },
    };
    this.backendService.updateTag(tag.id, update).subscribe((_) => {
      this.refreshData();
    });
  }

  private refreshData() {
    this.backendService.getTags().subscribe((tags) => {
      this.tags = tags.data;
    });
  }
}
