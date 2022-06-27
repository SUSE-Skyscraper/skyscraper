import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CloudTenantsComponent } from './cloud-tenants.component';
import { MockProvider } from 'ng-mocks';
import { HttpClient } from '@angular/common/http';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { BackendService } from '../backend.service';
import { EMPTY } from 'rxjs';

describe('CloudTenantsComponent', () => {
  let component: CloudTenantsComponent;
  let fixture: ComponentFixture<CloudTenantsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [CloudTenantsComponent],
      providers: [
        MockProvider(BackendService, {
          getCloudTenants: () => EMPTY,
        }),
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(CloudTenantsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
