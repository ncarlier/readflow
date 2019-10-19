import { useCallback, useEffect, useState } from 'react'

const key = 'LocalConfiguration'

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
}

const initialConfiguration: LocalConfiguration = {
  sortOrders: {
    unread: 'asc',
    offline: 'asc',
    history: 'desc'
  },
  limit: 10
}

export default (): [LocalConfiguration, (value: LocalConfiguration) => void] => {
  const [state, updateState] = useState<LocalConfiguration>(() => {
    try {
      const localConfiguration = window.localStorage.getItem(key)
      if (localConfiguration === null) {
        window.localStorage.setItem(key, JSON.stringify(initialConfiguration))
        return initialConfiguration
      } else {
        return JSON.parse(localConfiguration)
      }
    } catch {
      return initialConfiguration
    }
  })
  const localStorageChanged = (e: StorageEvent) => {
    if (e.key === key) {
      updateState(JSON.parse(e.newValue as string))
    }
  }
  const setState = useCallback(
    (value: LocalConfiguration) => {
      window.localStorage.setItem(key, JSON.stringify(value))
      updateState(value)
    },
    [key, updateState]
  )
  useEffect(() => {
    window.addEventListener('storage', localStorageChanged)
    return () => {
      window.removeEventListener('storage', localStorageChanged)
    }
  })
  return [state, setState]
}
