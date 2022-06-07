import { Inject, Injectable, NgModule } from '@angular/core';
import { Router, RouterModule, Routes } from '@angular/router';
import { OKTA_AUTH, OktaCallbackComponent } from '@okta/okta-angular';
import { LoginComponent } from './login/login.component';
import { OktaAuth } from '@okta/okta-auth-js';
import { DashboardComponent } from './dashboard/dashboard.component';

@Injectable()
export class IsLoggedIn {
  constructor(
    private router: Router,
    @Inject(OKTA_AUTH) private authService: OktaAuth,
  ) {}

  resolve(): void {
    const isAuthenticatedClaim = this.authService.isAuthenticated();
    isAuthenticatedClaim.then((resp) => {
      if (resp.valueOf()) {
        this.router.navigate(['/dashboard']);
      }
    });
  }
}

const routes: Routes = [
  {
    path: '',
    redirectTo: '/login',
    pathMatch: 'full',
  },
  {
    path: 'login',
    component: LoginComponent,
    resolve: [IsLoggedIn],
  },
  {
    path: 'callback',
    component: OktaCallbackComponent,
  },
  {
    path: 'dashboard',
    component: DashboardComponent,
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
