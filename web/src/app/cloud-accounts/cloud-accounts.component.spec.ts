import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CloudAccountsComponent } from './cloud-accounts.component';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { MockProvider } from 'ng-mocks';
import { BackendService } from '../backend.service';
import { EMPTY } from 'rxjs';

describe('CloudAccountsComponent', () => {
  let component: CloudAccountsComponent;
  let fixture: ComponentFixture<CloudAccountsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [CloudAccountsComponent],
      providers: [
        MockProvider(BackendService, {
          getCloudAccounts: () => EMPTY,
        }),
        {
          provide: ActivatedRoute,
          useValue: {
            snapshot: {
              paramMap: convertToParamMap({
                cloud: 'AWS',
                tenant_id: '12345',
              }),
            },
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(CloudAccountsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
