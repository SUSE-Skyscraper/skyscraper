import { environment } from '../environments/environment';

export default {
  backend: {
    host: environment.backendServer,
    validator: environment.validatorServer,
  },
  oidc: {
    clientId: environment.oktaClientId,
    issuer: environment.oktaIssuer,
    redirectUri: '/callback',
    scopes: ['openid', 'profile', 'email'],
    pkce: true,
  },
  resourceServer: {},
};
