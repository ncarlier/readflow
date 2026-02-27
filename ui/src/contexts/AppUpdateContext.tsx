import React, { createContext, FC, PropsWithChildren, useContext, useState, useEffect } from 'react'

interface AppUpdateContextType {
  updateAvailable: boolean
  registration: ServiceWorkerRegistration | null
}

const AppUpdateContext = createContext<AppUpdateContextType>({
  updateAvailable: false,
  registration: null,
})

const AppUpdateProvider: FC<PropsWithChildren> = ({ children }) => {
  const [updateAvailable, setUpdateAvailableState] = useState(false)
  const [registration, setRegistration] = useState<ServiceWorkerRegistration | null>(null)

  useEffect(() => {
    const onUpdate = (event: CustomEvent<ServiceWorkerRegistration>) => {
      setRegistration(event.detail)
      setUpdateAvailableState(true)
    }
    window.addEventListener('sw-update', onUpdate as EventListener)
    return () => window.removeEventListener('sw-update', onUpdate as EventListener)
  }, [])

  return <AppUpdateContext.Provider value={{ updateAvailable, registration }}>{children}</AppUpdateContext.Provider>
}

export { AppUpdateProvider }
export const useAppUpdate = () => useContext(AppUpdateContext)
