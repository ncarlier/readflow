import React, { createContext, FC, useCallback, useEffect, useState } from 'react'

const key = 'readflow.localConfiguration'

export type Theme = 'light' | 'dark' | 'auto'

export type SortOrder = 'asc' | 'desc'
export type SortBy = 'key' | 'stars'
export type DisplayMode = 'list' | 'grid'

export interface DisplayPreference {
  order: SortOrder
  by: SortBy
  mode: DisplayMode
}

export interface DisplayPreferences {
  unread: DisplayPreference
  offline: DisplayPreference
  history: DisplayPreference
  starred: DisplayPreference
  [key: string]: DisplayPreference
}

export interface LocalConfiguration {
  display: DisplayPreferences
  limit: number
  theme: Theme
}

const defaultLocalConfiguration: LocalConfiguration = {
  display: {
    unread: { order: 'asc', by: 'key', mode: 'list' },
    offline: { order: 'asc', by: 'key', mode: 'list' },
    starred: { order: 'asc', by: 'stars', mode: 'grid' },
    history: { order: 'desc', by: 'key', mode: 'list' },
  },
  limit: 10,
  theme: 'auto',
}

interface LocalConfigurationContextType {
  localConfiguration: LocalConfiguration
  updateLocalConfiguration: (config: LocalConfiguration) => void
}

const LocalConfigurationContext = createContext<LocalConfigurationContextType>({
  localConfiguration: defaultLocalConfiguration,
  updateLocalConfiguration: () => true,
})

const LocalConfigurationProvider: FC = ({ children }) => {
  const [localConfiguration, setLocalConfiguration] = useState<LocalConfiguration>(() => {
    try {
      const config = window.localStorage.getItem(key)
      if (config === null) {
        window.localStorage.setItem(key, JSON.stringify(defaultLocalConfiguration))
        return defaultLocalConfiguration
      } else {
        return { ...defaultLocalConfiguration, ...JSON.parse(config) }
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
