import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CloudAccountComponent } from './cloud-account.component';
import { MockProvider } from 'ng-mocks';
import { BackendService } from '../backend.service';
import { EMPTY } from 'rxjs';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { FormBuilder } from '@angular/forms';

describe('CloudAccountComponent', () => {
  let component: CloudAccountComponent;
  let fixture: ComponentFixture<CloudAccountComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [CloudAccountComponent],
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

    fixture = TestBed.createComponent(CloudAccountComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
