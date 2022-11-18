// This file is required by karma.conf.js and loads recursively all the .spec and framework files

import 'zone.js/testing';
import { getTestBed } from '@angular/core/testing';
import {
  BrowserDynamicTestingModule,
  platformBrowserDynamicTesting,
} from '@angular/platform-browser-dynamic/testing';
import { ngMocks } from 'ng-mocks';
import { OKTA_AUTH, OktaAuthStateService } from '@okta/okta-angular';
import { EMPTY } from 'rxjs';

ngMocks.autoSpy('jasmine');
ngMocks.defaultMock(OktaAuthStateService, () => ({
  authState$: EMPTY,
}));

// First, initialize the Angular testing environment.
getTestBed().initTestEnvironment(
  BrowserDynamicTestingModule,
  platformBrowserDynamicTesting(),
);
