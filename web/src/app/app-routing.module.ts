import { Inject, Injectable, NgModule } from '@angular/core';
import { Router, RouterModule, Routes } from '@angular/router';
import {
  OKTA_AUTH,
  OktaAuthGuard,
  OktaCallbackComponent,
} from '@okta/okta-angular';
import { LoginComponent } from './login/login.component';
import { OktaAuth } from '@okta/okta-auth-js';
import { DashboardComponent } from './dashboard/dashboard.component';
import { CloudTenantsComponent } from './cloud-tenants/cloud-tenants.component';
import { CloudAccountsComponent } from './cloud-accounts/cloud-accounts.component';
import { CloudAccountComponent } from './cloud-account/cloud-account.component';
import { TagsComponent } from './tags/tags.component';
import { AuditLogComponent } from './audit-log/audit-log.component';

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
        this.router.navigate(['/cloud_tenants']);
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
    canActivate: [OktaAuthGuard],
  },
  {
    path: 'cloud_accounts',
    component: CloudAccountsComponent,
    canActivate: [OktaAuthGuard],
  },
  {
    path: 'audit_log',
    component: AuditLogComponent,
    canActivate: [OktaAuthGuard],
  },
  {
    path: 'tags',
    component: TagsComponent,
    canActivate: [OktaAuthGuard],
  },
  {
    path: 'cloud_accounts/:id',
    component: CloudAccountComponent,
    canActivate: [OktaAuthGuard],
  },
  {
    path: 'cloud_tenants',
    component: CloudTenantsComponent,
    canActivate: [OktaAuthGuard],
  },
  {
    path: 'cloud_tenants/cloud/:cloud/tenant/:tenant_id/accounts',
    component: CloudAccountsComponent,
    canActivate: [OktaAuthGuard],
  },
  {
    path: 'cloud_tenants/cloud/:cloud/tenant/:tenant_id/accounts/:account_id',
    component: CloudAccountComponent,
    canActivate: [OktaAuthGuard],
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
