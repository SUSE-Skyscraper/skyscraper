import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TagsComponent } from './tags.component';
import { MockProvider } from 'ng-mocks';
import { BackendService } from '../backend.service';
import { EMPTY } from 'rxjs';
import { FormBuilder } from '@angular/forms';

describe('TagsComponent', () => {
  let component: TagsComponent;
  let fixture: ComponentFixture<TagsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [TagsComponent],
      providers: [
        MockProvider(BackendService, {
          getCloudAccounts: () => EMPTY,
          getTags: () => EMPTY,
        }),
        FormBuilder,
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(TagsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
