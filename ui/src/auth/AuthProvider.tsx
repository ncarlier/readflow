import React, { createContext, FC, PropsWithChildren, useCallback, useContext, useEffect, useMemo, useRef, useState } from 'react'

import { Log, SigninRedirectArgs, SignoutRedirectArgs, User, UserManager } from 'oidc-client-ts'

import { config } from './oidc-configuration'
import { useHistory, useLocation } from 'react-router-dom'
import { clearAuthParams, hasAuthParams } from './helper'

Log.setLogger(console)
Log.setLevel(Log.WARN)

interface AuthContextType {
  user: User | null
  isLoading: boolean
  error?: any
  login: (args?: SigninRedirectArgs | undefined) => Promise<void>
  logout: (args?: SignoutRedirectArgs) => Promise<void>
}

const redirectKey = 'readflow.redirect'

const AuthContext = createContext<AuthContextType>({} as AuthContextType)

export const AuthProvider: FC<PropsWithChildren> = ({ children }) => {
  const [userManager] = useState(() => new UserManager(config))
  const [user, setUser] = useState<User | null>(null)
  const [error, setError] = useState<any>()
  const [isLoading, setIsLoading] = useState<boolean>(false)
  const { search } = useLocation()
  const history = useHistory()
  
  const login = useCallback(async () => {
    try {
      localStorage.setItem(redirectKey, JSON.stringify(history.location))
      console.debug('saving location', JSON.stringify(history.location))
      await userManager.signinRedirect()
    } catch (err) {
      setError(err)
    }
  }, [userManager])
  const logout = useCallback(userManager.signoutRedirect.bind(userManager), [userManager])
  const handleLoginFlow = useCallback(async () => {
    setIsLoading(true)
    try {
      const params = new URLSearchParams(search)
      // handle error callback:
      if (params.has('error')) {
        const error = params.get('error')
        console.error('error callback from Authority server:', error)
        await userManager.removeUser()
        throw error
      }
      // handle login callback:
      if (hasAuthParams(params)) {
        console.info('callback from Authority server: sign in...')
        const user = await userManager.signinCallback()
        if (user) {
          console.debug('logged user:', user.profile?.sub)
          setUser(user)
          const redirect = localStorage.getItem(redirectKey)
          if (redirect) {
            localStorage.removeItem(redirectKey)
            console.debug('restoring location', redirect)
            return history.replace(JSON.parse(redirect))
          }
          return history.replace({
            search: clearAuthParams(params),
          })
        }
      }
      // otherwise handle user state:
      const user = await userManager.getUser()
      if (user) {
        console.debug('authenticated user:', user?.profile.sub)
        setUser(user)
      } else {
        console.info('user not authenticated, redirecting to sign-in page...')
        login()
      }
    } catch (err) {
      setError(err)
    } finally {
      setIsLoading(false)
    }
  }, [userManager, search, login])

  // main login flow
  const didInitialize = useRef<boolean>(false)
  useEffect(() => {
    if (didInitialize.current) return
    didInitialize.current = true
    console.info('exectuting login flow')
    handleLoginFlow()
  }, [handleLoginFlow])

  // userManager events handlers:
  useEffect(() => {
    if (!userManager) return
    // event UserSignedOut (e.g. external sign out)
    const handleUserSignedOut = () => {
      console.log('user signed out from Authority server: sign out...')
      userManager.signoutRedirect()
    }
    userManager.events.addUserSignedOut(handleUserSignedOut)
    // event UserLoaded (e.g. initial load, silent renew success)
    const handleUserLoaded = (user: User) => {
      //console.debug('UserLoaded', user, user.expired)
      setUser(user)
    }
    userManager.events.addUserLoaded(handleUserLoaded)
    // event UserUnloaded (e.g. userManager.removeUser)
    const handleUserUnloaded = () => {
      //console.debug('UserUnLoaded')
      setUser(null)
    }
    userManager.events.addUserUnloaded(handleUserUnloaded)
    // event SilentRenewError (silent renew error)
    const handleSilentRenewError = (err: Error) => {
      //console.debug('SilentRenewError', err)
      setError(err)
    }
    userManager.events.addSilentRenewError(handleSilentRenewError)

    // clear stale state...
    console.debug('clearing authentication stale state...')
    userManager.clearStaleState()
    return () => {
      userManager.events.removeUserSignedOut(handleUserSignedOut)
      userManager.events.removeUserLoaded(handleUserLoaded)
      userManager.events.removeUserUnloaded(handleUserUnloaded)
      userManager.events.removeSilentRenewError(handleSilentRenewError)
    }
  }, [userManager])

  const value = useMemo(
    () => ({
      user,
      isLoading,
      error,
      login,
      logout,
    }),
    [user, isLoading, error, login, logout]
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export const useAuth = () => useContext(AuthContext)
