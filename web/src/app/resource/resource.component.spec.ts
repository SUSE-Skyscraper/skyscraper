import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ResourceComponent } from './resource.component';
import { MockProvider } from 'ng-mocks';
import { BackendService } from '../backend.service';
import { EMPTY } from 'rxjs';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { FormBuilder } from '@angular/forms';
import { MatSnackBarModule } from '@angular/material/snack-bar';

describe('CloudAccountComponent', () => {
  let component: ResourceComponent;
  let fixture: ComponentFixture<ResourceComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MatSnackBarModule],
      declarations: [ResourceComponent],
      providers: [
        MockProvider(BackendService, {
          getCloudAccount: () => EMPTY,
        }),
        FormBuilder,
        {
          provide: ActivatedRoute,
          useValue: {
            snapshot: {
              paramMap: convertToParamMap({
                cloud: 'AWS',
                tenant_id: '12345',
                account_id: '12345',
              }),
            },
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(ResourceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
