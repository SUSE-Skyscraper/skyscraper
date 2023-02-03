import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TenantsComponent } from './tenants.component';
import { MockProvider } from 'ng-mocks';
import { HttpClient } from '@angular/common/http';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { BackendService } from '../backend.service';
import { EMPTY } from 'rxjs';

describe('CloudTenantsComponent', () => {
  let component: TenantsComponent;
  let fixture: ComponentFixture<TenantsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [TenantsComponent],
      providers: [
        MockProvider(BackendService, {
          getCloudTenants: () => EMPTY,
        }),
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(TenantsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
