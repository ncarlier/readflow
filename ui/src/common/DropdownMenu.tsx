import React, { ReactNode, useEffect, createRef } from 'react'

import styles from './DropdownMenu.module.css'
import ButtonIcon from './ButtonIcon'

type Props = {
  children: ReactNode
}

export default ({children}: Props) => {
  const ref = createRef<HTMLDetailsElement>()
  useEffect(() => {
    document.addEventListener('click', handleClickOutside, {capture: true})
    return () => {
      document.removeEventListener('click', handleClickOutside, {capture: true})
    }
  }, [ref])
  
  const handleClickOutside = (e: MouseEvent) => {
    const $el = e.target
    if (!($el instanceof Element)) return
    if (ref.current) {
      const $details = $el.closest("details")
      if ($details === ref.current && $details.hasAttribute('open')) {
        e.preventDefault()
      }
      ref.current.removeAttribute('open')
    }
  }

  return (
    <details ref={ref} className={styles.menu}>
      <summary><ButtonIcon icon="more_vert" /></summary>
      <nav className={styles.nav} tabIndex={-1}>{children}</nav>
    </details>
  )
}
