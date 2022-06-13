import { Inject, Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor,
} from '@angular/common/http';
import { Observable } from 'rxjs';
import { OKTA_AUTH } from '@okta/okta-angular';
import { OktaAuth } from '@okta/okta-auth-js';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  private readonly accessToken: string | undefined;

  constructor(@Inject(OKTA_AUTH) private oktaAuth: OktaAuth) {
    this.accessToken = this.oktaAuth.getAccessToken();
  }

  intercept(
    request: HttpRequest<unknown>,
    next: HttpHandler,
  ): Observable<HttpEvent<unknown>> {
    const authRequest = request.clone({
      headers: request.headers.set(
        'Authorization',
        `Bearer ${this.accessToken}`,
      ),
    });
    return next.handle(authRequest);
  }
}
