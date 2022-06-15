import { Component, Inject, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { OKTA_AUTH, OktaAuthStateService } from '@okta/okta-angular';
import { OktaAuth } from '@okta/okta-auth-js';
import { BackendService, Profile } from './backend.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  title = 'skyscraper-web';
  isAuthenticated: boolean = false;
  profile: Profile | undefined;

  constructor(
    private authStateService: OktaAuthStateService,
    @Inject(OKTA_AUTH) private oktaAuth: OktaAuth,
    private backendService: BackendService,
    public router: Router,
  ) {
    // Subscribe to authentication state changes
    this.authStateService.authState$.subscribe((authState) => {
      this.isAuthenticated = authState.isAuthenticated === true;

      if (this.isAuthenticated) {
        this.backendService.getProfile().subscribe((profile) => {
          this.profile = profile;
        });
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
