import React from 'react'

import styles from './Loader.module.css'

export default () => (
  <section className={styles.overlay}>
    <div className={styles.spinner} />
  </section>
)
