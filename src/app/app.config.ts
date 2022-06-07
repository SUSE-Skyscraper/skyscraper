import { environment } from '../environments/environment';

export default {
  oidc: {
    clientId: environment.oktaClientId,
    issuer: environment.oktaIssuer,
    redirectUri: '/callback',
    scopes: ['openid', 'profile', 'email'],
    pkce: true,
  },
  resourceServer: {},
};
