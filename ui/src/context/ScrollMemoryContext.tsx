import React, { ReactNode, useEffect, useState } from 'react'

const ScrollMemoryContext = React.createContext(0)

var cache = new Map<string, number>()

export const setScrollPosition = (key: string, pos: number) => cache.set(key, pos)

interface Props {
  children: ReactNode
}

if ('scrollRestoration' in history) {
  history.scrollRestoration = 'manual'
}

const ScrollMemoryProvider = ({ children }: Props) => {
  const [state, setState] = useState(0)

  const onPopState = () => {
    const key = location.pathname
    if (cache.has(key)) {
      const pos = cache.get(key)
      // console.log(`restoring scroll position for ${key}: ${pos}`)
      setState(pos || 0)
    } else {
      setState(0)
    }
  }

  useEffect(() => {
    // console.log('scrollMemoryContext:init')
    window.addEventListener('popstate', onPopState)
    return () => {
      //console.log('scrollMemoryContext:destroy')
      window.removeEventListener('popstate', onPopState)
      setState(0)
    }
  })

  return <ScrollMemoryContext.Provider value={state}>{children}</ScrollMemoryContext.Provider>
}

export { ScrollMemoryContext, ScrollMemoryProvider }
