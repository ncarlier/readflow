import React, { forwardRef, ReactNode, Ref } from 'react'

import styles from './Content.module.css'

interface Props {
  children: ReactNode
}

export default forwardRef(({ children }: Props, ref: Ref<any>) => (
  <section ref={ref} className={styles.content}>
    {children}
  </section>
))
