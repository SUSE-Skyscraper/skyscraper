import { Component, Inject, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { OKTA_AUTH, OktaAuthStateService } from '@okta/okta-angular';
import { OktaAuth } from '@okta/okta-auth-js';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  title = 'skyscraper-web';
  isAuthenticated: boolean = false;
  email: string = '';

  constructor(
    private authStateService: OktaAuthStateService,
    @Inject(OKTA_AUTH) private oktaAuth: OktaAuth,
    public router: Router,
  ) {
    // Subscribe to authentication state changes
    this.authStateService.authState$.subscribe((authState) => {
      this.isAuthenticated = authState.isAuthenticated === true;

      if (this.isAuthenticated) {
        const userClaim = this.oktaAuth.getUser();
        userClaim.then((resp) => {
          this.email = resp.email === undefined ? '' : resp.email;
        });
      } else {
        this.email = '';
      }
    });
  }

  async ngOnInit() {
    // Get the authentication state for immediate use
    this.isAuthenticated = await this.oktaAuth.isAuthenticated();
  }

  async login() {
    await this.oktaAuth.signInWithRedirect();
  }

  async logout() {
    await this.oktaAuth.signOut();
  }
}
