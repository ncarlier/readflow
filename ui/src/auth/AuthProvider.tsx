import React, { createContext, FC, useCallback, useContext, useEffect, useMemo, useState } from 'react'

import { Log, SigninRedirectArgs, SignoutRedirectArgs, User, UserManager } from 'oidc-client-ts'

import { config } from './oidc-configuration'
import { useLocation } from 'react-router-dom'

Log.setLogger(console)
Log.setLevel(Log.WARN)

const hasAuthParams = (search: string): boolean => {
  const params = new URLSearchParams(search)
  return params.has('code') && params.has('state')
}

const getErrorParam = (search: string): string | null => {
  const params = new URLSearchParams(search)
  return params.get('error')
}

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  error?: any
  login: (args?: SigninRedirectArgs | undefined) => Promise<void>
  logout: (args?: SignoutRedirectArgs) => Promise<void>
}

const AuthContext = createContext<AuthContextType>({} as AuthContextType)

export const AuthProvider: FC = ({ children }) => {
  const [userManager] = useState(() => new UserManager(config))
  const [user, setUser] = useState<User | null>(null)
  const [error, setError] = useState<any>()
  const [isLoading, setIsLoading] = useState<boolean>(true)
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false)
  const { search } = useLocation()

  useEffect(() => {
    if (!userManager) return
    if (isAuthenticated) return
    ;(async () => {
      setIsLoading(true)
      //console.debug('auth provider logic...')
      try {
        const error = getErrorParam(search)
        if (error) {
          console.error('error callback from Authority server:', error)
          await userManager.removeUser()
          throw error
        }
        if (hasAuthParams(search)) {
          console.info('callback from Authority server: sign in...')
          const user = await userManager.signinCallback()
          if (user) {
            // clear query params
            window.history.replaceState(null, '', window.location.pathname)
            console.debug('logged user:', user.profile?.preferred_username)
          }
          return
        }
        const user = await userManager.getUser()
        if (user) {
          console.debug('current user:', user?.profile.preferred_username)
          setIsAuthenticated(true)
          setUser(user)
          return
        }
      } catch (err) {
        setError(err)
      } finally {
        setIsLoading(false)
      }
    })()
  }, [userManager, isAuthenticated, search])

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
      setIsAuthenticated(!user.expired)
      setUser(user)
    }
    userManager.events.addUserLoaded(handleUserLoaded)
    // event UserUnloaded (e.g. userManager.removeUser)
    const handleUserUnloaded = () => {
      //console.debug('UserUnLoaded')
      setIsAuthenticated(false)
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

  const login = useCallback(userManager.signinRedirect.bind(userManager), [userManager])
  const logout = useCallback(userManager.signoutRedirect.bind(userManager), [userManager])

  const value = useMemo(
    () => ({
      user,
      isLoading,
      isAuthenticated,
      error,
      login,
      logout,
    }),
    [user, isLoading, isAuthenticated, error, login, logout]
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export const useAuth = () => useContext(AuthContext)
