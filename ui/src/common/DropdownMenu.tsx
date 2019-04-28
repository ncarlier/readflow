import React, { createRef, ReactNode, useEffect } from 'react'

import ButtonIcon from './ButtonIcon'
import styles from './DropdownMenu.module.css'

interface Props {
  children: ReactNode
}

export default ({ children }: Props) => {
  const ref = createRef<HTMLDetailsElement>()

  const handleClickOutside = (e: MouseEvent) => {
    const $el = e.target
    if (!($el instanceof Element)) return
    if (ref.current) {
      const isButton = $el.parentElement && $el.parentElement.tagName === 'BUTTON'
      const $details = $el.closest('details')
      if (isButton && $details === ref.current && $details.hasAttribute('open')) {
        e.preventDefault()
      }
      ref.current.removeAttribute('open')
    }
  }

  useEffect(() => {
    document.addEventListener('click', handleClickOutside, { capture: true })
    return () => {
      document.removeEventListener('click', handleClickOutside, { capture: true })
    }
  }, [ref])

  return (
    <details ref={ref} className={styles.menu}>
      <summary>
        <ButtonIcon icon="more_vert" />
      </summary>
      <nav className={styles.nav} tabIndex={-1}>
        {children}
      </nav>
    </details>
  )
}
