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
import { TenantsComponent } from './tenants/tenants.component';
import { ResourcesComponent } from './resources/resources.component';
import { ResourceComponent } from './resource/resource.component';
import { TagsComponent } from './tags/tags.component';
import { AuditLogComponent } from './audit-log/audit-log.component';
import { UsersComponent } from './users/users.component';
import { UserComponent } from './user/user.component';

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
        this.router.navigate(['/tenants']);
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
    path: 'users',
    component: UsersComponent,
    canActivate: [OktaAuthGuard],
  },
  {
    path: 'users/:id',
    component: UserComponent,
    canActivate: [OktaAuthGuard],
  },
  {
    path: 'tenants',
    component: TenantsComponent,
    canActivate: [OktaAuthGuard],
  },
  {
    path: 'groups/:group/tenants/:tenant_id/resources',
    component: ResourcesComponent,
    canActivate: [OktaAuthGuard],
  },
  {
    path: 'groups/:group/tenants/:tenant_id/resources/:id',
    component: ResourceComponent,
    canActivate: [OktaAuthGuard],
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
