import React, { createContext, FC, PropsWithChildren, useContext, useState } from 'react'
import { useAddToHomescreenPrompt } from '../hooks'
import { isDisplayMode } from '../helpers'

interface DeviceContextType {
  isInstalled: boolean
  isInstallable: boolean
  promptToInstall: () => void
}

const DeviceContext = createContext<DeviceContextType>({
  isInstalled: false,
  isInstallable: false,
  promptToInstall: () => {return},
})

const DeviceProvider: FC<PropsWithChildren> = ({ children }) => {
  const [prompt, promptToInstall] = useAddToHomescreenPrompt()
  const [isInstalled] = useState(isDisplayMode('standalone'))
  const [isInstallable, setInstallableState] = useState(false)

  React.useEffect(
    () => {
      if (prompt) {
        setInstallableState(true)
      }
    },
    [prompt]
  )

  return <DeviceContext.Provider value={{ isInstalled, isInstallable, promptToInstall }}>{children}</DeviceContext.Provider>
}

export { DeviceProvider }

export const useDevice = () => useContext(DeviceContext)
