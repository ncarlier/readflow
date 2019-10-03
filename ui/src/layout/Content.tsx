import React, { ReactNode, useRef } from 'react'

import { ScrollMemoryProvider } from '../context/ScrollMemoryContext'
import useScrollMemory from '../hooks/useScrollMemory'
import styles from './Content.module.css'

interface Props {
  children: ReactNode
}

export default ({ children }: Props) => {
  const ref = useRef<HTMLDivElement>(null)
  useScrollMemory(ref)

  return (
    <ScrollMemoryProvider>
      <section ref={ref} className={styles.content}>
        {children}
      </section>
    </ScrollMemoryProvider>
  )
}
