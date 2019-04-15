import { Log, UserManager } from 'oidc-client'
import { AUTHORITY, CLIENT_ID } from '../constants'

export class AuthService {
  public userManager: UserManager

  constructor() {
    const settings = {
      authority: AUTHORITY,
      client_id: CLIENT_ID,
      redirect_uri: `${document.location.origin}/signin-callback.html`,
      silent_redirect_uri: `${document.location.origin}/silent-renew.html`,
      post_logout_redirect_uri: document.location.origin,
      response_type: 'id_token token',
      scope: 'openid profile email'
    }
    this.userManager = new UserManager(settings)

    Log.logger = console
    Log.level = Log.WARN
  }

  public getUser() {
    return this.userManager.getUser()
  }

  public getAccountUrl() {
    return AUTHORITY + 
      '/account?referrer=' +
      CLIENT_ID +
      '&referrer_uri=' +
      encodeURI(document.location.href)
  }

  public login() {
    return this.userManager.signinRedirect()
  }

  public renewToken() {
    return this.userManager.signinSilent()
  }

  public logout() {
    return this.userManager.signoutRedirect()
  }
}

export default new AuthService()
