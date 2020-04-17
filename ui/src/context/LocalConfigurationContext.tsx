import React, { createContext, ReactNode, useCallback, useEffect, useState } from 'react'

const key = 'LocalConfiguration'

export type Theme = 'light' | 'dark' | 'auto'

export type SortOrder = 'asc' | 'desc'

interface SortOrders {
  unread: SortOrder
  offline: SortOrder
  history: SortOrder
  [key: string]: SortOrder
}

export interface LocalConfiguration {
  sortOrders: SortOrders
  limit: number
  theme: Theme
}

const defaultLocalConfiguration: LocalConfiguration = {
  sortOrders: {
    unread: 'asc',
    offline: 'asc',
    history: 'desc'
  },
  limit: 10,
  theme: 'auto'
}

interface LocalConfigurationContextType {
  localConfiguration: LocalConfiguration
  updateLocalConfiguration: (config: LocalConfiguration) => void
}

const LocalConfigurationContext = createContext<LocalConfigurationContextType>({
  localConfiguration: defaultLocalConfiguration,
  updateLocalConfiguration: () => true
})

interface Props {
  children: ReactNode
}

const LocalConfigurationProvider = ({ children }: Props) => {
  const [localConfiguration, setLocalConfiguration] = useState<LocalConfiguration>(() => {
    try {
      const config = window.localStorage.getItem(key)
      if (config === null) {
        window.localStorage.setItem(key, JSON.stringify(defaultLocalConfiguration))
        return defaultLocalConfiguration
      } else {
        return JSON.parse(config)
      }
    } catch {
      return defaultLocalConfiguration
    }
  })

  const localStorageChanged = (e: StorageEvent) => {
    if (e.key === key) {
      setLocalConfiguration(JSON.parse(e.newValue as string))
    }
  }

  const updateLocalConfiguration = useCallback(
    (value: LocalConfiguration) => {
      window.localStorage.setItem(key, JSON.stringify(value))
      setLocalConfiguration(value)
    },
    [setLocalConfiguration]
  )

  useEffect(() => {
    window.addEventListener('storage', localStorageChanged)
    return () => {
      window.removeEventListener('storage', localStorageChanged)
    }
  })

  return (
    <LocalConfigurationContext.Provider value={{ localConfiguration, updateLocalConfiguration }}>
      {children}
    </LocalConfigurationContext.Provider>
  )
}

export { LocalConfigurationContext, LocalConfigurationProvider }
