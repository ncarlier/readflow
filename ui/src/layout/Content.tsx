import React, { FC, useRef } from 'react'

import { ScrollToTop } from '../components'
import { useScrollMemory } from '../hooks'
import styles from './Content.module.css'

interface Props {
  scrollToTop?: boolean
}

export const Content: FC<Props> = ({ children, scrollToTop = false }) => {
  const ref = useRef<HTMLDivElement>(null)
  useScrollMemory(ref)

  return (
    <section ref={ref} className={styles.content}>
      {children}
      {scrollToTop && <ScrollToTop parent={ref} />}
    </section>
  )
}
