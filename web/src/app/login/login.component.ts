import { Component, Inject, OnInit } from '@angular/core';
import { OKTA_AUTH } from '@okta/okta-angular';
import { OktaAuth } from '@okta/okta-auth-js';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
})
export class LoginComponent {
  constructor(@Inject(OKTA_AUTH) private oktaAuth: OktaAuth) {}

  async login() {
    await this.oktaAuth.signInWithRedirect();
  }
}
