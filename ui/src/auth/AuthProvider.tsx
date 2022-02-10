import React, { createContext, FC, useCallback, useContext, useEffect, useMemo, useState } from 'react'

import { Log, SigninRedirectArgs, User, UserManager } from 'oidc-client-ts'

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
  isLoading: boolean
  error?: any
  login: (args?: SigninRedirectArgs | undefined) => Promise<void>
}

const AuthContext = createContext<AuthContextType>({} as AuthContextType)

export const AuthProvider: FC = ({ children }) => {
  const [userManager] = useState(() => new UserManager(config))
  const [user, setUser] = useState<User | null>(null)
  const [error, setError] = useState<any>()
  const [isLoading, setIsLoading] = useState<boolean>(true)
  const { search } = useLocation()

  useEffect(() => {
    if (!userManager) return
    if (user) return
    ;(async () => {
      setIsLoading(true)
      //console.debug('auth porvider logic...')
      try {
        const error = getErrorParam(search)
        if (error) {
          console.error('error callback from Authority server:', error)
          await userManager.removeUser()
          throw error
        }
        if (hasAuthParams(search)) {
          console.info('callback from Authority server: sign in...')
          const _user = await userManager.signinCallback()
          if (_user) {
            // clear query params
            window.history.replaceState(null, '', window.location.pathname)
            console.debug('logged user:', _user.profile?.preferred_username)
          }
          return
        }
        const _user = await userManager.getUser()
        if (_user) {
          console.debug('current user:', _user?.profile.preferred_username)
          setUser(_user)
          return
        }
      } catch (err) {
        setError(err)
      } finally {
        setIsLoading(false)
      }
    })()
  }, [userManager, user, search])

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
      //console.debug('UserLoaded', user)
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

  const login = useCallback(userManager.signinRedirect.bind(userManager), [userManager])

  const memoedValue = useMemo(
    () => ({
      user,
      isLoading,
      error,
      login,
    }),
    [user, isLoading, error, login]
  )

  return <AuthContext.Provider value={memoedValue}>{children}</AuthContext.Provider>
}

export const useAuth = () => useContext(AuthContext)
