import { useCallback, useEffect, useState } from 'react'

const key = 'LocalConfiguration'

export interface LocalConfiguration {
  defaultSortOrder: string
  sortOrders: Map<string, string>
}

const initialConfiguration: LocalConfiguration = {
  defaultSortOrder: 'asc',
  sortOrders: new Map([['unread', 'asc'], ['offline', 'asc'], ['history', 'desc']])
}

export default (): [LocalConfiguration, (value: LocalConfiguration) => void] => {
  const [state, updateState] = useState(() => {
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
