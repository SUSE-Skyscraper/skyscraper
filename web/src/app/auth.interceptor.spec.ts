import { TestBed } from '@angular/core/testing';

import { AuthInterceptor } from './auth.interceptor';
import { MockProvider } from 'ng-mocks';
import { OKTA_AUTH } from '@okta/okta-angular';

describe('AuthInterceptor', () => {
  beforeEach(() =>
    TestBed.configureTestingModule({
      providers: [
        AuthInterceptor,
        MockProvider(OKTA_AUTH, {
          getAccessToken: function () {
            return '';
          },
        }),
      ],
    }).compileComponents(),
  );

  it('should be created', () => {
    const interceptor: AuthInterceptor = TestBed.inject(AuthInterceptor);
    expect(interceptor).toBeTruthy();
  });
});
