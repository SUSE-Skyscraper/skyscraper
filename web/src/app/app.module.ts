import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule, IsLoggedIn } from './app-routing.module';
import { AppComponent } from './app.component';
import config from './app.config';
import { OKTA_CONFIG, OktaAuthModule } from '@okta/okta-angular';
import { OktaAuth } from '@okta/okta-auth-js';
import { LoginComponent } from './login/login.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatCardModule } from '@angular/material/card';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';
import { AuthInterceptor } from './auth.interceptor';
import { CloudTenantsComponent } from './cloud-tenants/cloud-tenants.component';
import { MatTableModule } from '@angular/material/table';
import { CloudAccountsComponent } from './cloud-accounts/cloud-accounts.component';
import { MatDialogModule } from '@angular/material/dialog';
import { CloudAccountComponent } from './cloud-account/cloud-account.component';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatOptionModule } from '@angular/material/core';
import { MatSelectModule } from '@angular/material/select';
import { TagsComponent } from './tags/tags.component';
import { AuditLogComponent } from './audit-log/audit-log.component';
import { UsersComponent } from './users/users.component';
import { UserComponent } from './user/user.component';
import { TagFormValidator } from './tag.validator';
import { MatTooltipModule } from '@angular/material/tooltip';

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    DashboardComponent,
    CloudTenantsComponent,
    CloudAccountsComponent,
    CloudAccountComponent,
    TagsComponent,
    AuditLogComponent,
    UsersComponent,
    UserComponent,
    TagFormValidator,
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    AppRoutingModule,
    OktaAuthModule,
    BrowserAnimationsModule,
    MatToolbarModule,
    MatButtonModule,
    MatListModule,
    MatIconModule,
    MatMenuModule,
    MatCardModule,
    MatTableModule,
    MatDialogModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatSnackBarModule,
    MatOptionModule,
    MatSelectModule,
    FormsModule,
    MatTooltipModule,
  ],
  providers: [
    {
      provide: OKTA_CONFIG,
      useFactory: () => {
        const oktaAuth = new OktaAuth(config.oidc);
        return { oktaAuth };
      },
    },
    {
      provide: HTTP_INTERCEPTORS,
      useClass: AuthInterceptor,
      multi: true,
    },
    IsLoggedIn,
  ],
  bootstrap: [AppComponent],
})
export class AppModule {}
