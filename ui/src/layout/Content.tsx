import React, { ReactNode, useRef } from 'react'
import ScrollToTop from '../components/ScrollToTop'

import useScrollMemory from '../hooks/useScrollMemory'
import styles from './Content.module.css'

interface Props {
  children: ReactNode
  scrollToTop?: boolean
}

export default ({ children, scrollToTop = false }: Props) => {
  const ref = useRef<HTMLDivElement>(null)
  useScrollMemory(ref)

  return (
    <section ref={ref} className={styles.content}>
      {children}
      {scrollToTop && <ScrollToTop parent={ref} />}
    </section>
  )
}
