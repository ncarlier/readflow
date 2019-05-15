/* eslint-disable @typescript-eslint/camelcase */
import { AuthService } from './'

const fakeUser = {
  expired: false,
  access_token: null
}

export class MockAuthService implements AuthService {
  public getUser() {
    return Promise.resolve(fakeUser)
  }

  public getAccountUrl() {
    return '/account'
  }

  public login() {
    return Promise.resolve({})
  }

  public renewToken() {
    return Promise.resolve(fakeUser)
  }

  public logout() {
    return Promise.resolve({})
  }
}
