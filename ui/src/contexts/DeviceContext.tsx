import React, { createContext, FC, PropsWithChildren, useContext, useState } from 'react'
import { isDisplayMode } from '../helpers'

interface IBeforeInstallPromptEvent extends Event {
  readonly platforms: string[]
  readonly userChoice: Promise<{
    outcome: 'accepted' | 'dismissed'
    platform: string
  }>
  prompt(): Promise<void>
}

interface DeviceContextType {
  isInstalled: boolean
  beforeInstallPromptEvent: IBeforeInstallPromptEvent | null
}

const DeviceContext = createContext<DeviceContextType>({
  isInstalled: false,
  beforeInstallPromptEvent: null,
})

const DeviceProvider: FC<PropsWithChildren> = ({ children }) => {
  const [beforeInstallPromptEvent, setBeforeInstallPromptEventState] = useState<IBeforeInstallPromptEvent | null>(null)
  const [isInstalled, setInstalledState] = useState(isDisplayMode('standalone'))

  React.useEffect(() => {
    const ready = (e: IBeforeInstallPromptEvent) => {
      e.preventDefault()
      setBeforeInstallPromptEventState(e)
    }
    const installed = () => setInstalledState(true)
    window.addEventListener('beforeinstallprompt', ready as any)
    window.addEventListener('appinstalled', installed)
    return () => {
      window.removeEventListener('beforeinstallprompt', ready as any)
      window.removeEventListener('appinstalled', installed)
    }
  }, [])

  return <DeviceContext.Provider value={{ isInstalled, beforeInstallPromptEvent }}>{children}</DeviceContext.Provider>
}

export { DeviceProvider }

export const useDevice = () => useContext(DeviceContext)
