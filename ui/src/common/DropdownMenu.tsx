import React, { createRef, ReactNode, useEffect, MouseEventHandler, CSSProperties } from 'react'

import ButtonIcon from './ButtonIcon'
import styles from './DropdownMenu.module.css'

interface Props {
  children: ReactNode
  style?: CSSProperties
}

export default ({ children, style }: Props) => {
  const ref = createRef<HTMLDetailsElement>()

  const handleClickOutside = (e: MouseEvent) => {
    const $el = e.target
    if (!($el instanceof Element)) return
    if (ref.current) {
      const isButton = $el.parentElement && $el.parentElement.tagName === 'BUTTON'
      const $details = $el.closest('details')
      if (!isButton || $details !== ref.current) {
        ref.current.removeAttribute('open')
      }
    }
  }

  const handleClickMenu: MouseEventHandler = (e) => {
    e.preventDefault()
    const $el = e.currentTarget
    const $summary: any = $el.parentNode
    $summary.click()
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
        <ButtonIcon icon="more_vert" onClick={handleClickMenu}/>
      </summary>
      <nav className={styles.nav} style={style} tabIndex={-1}>
        {children}
      </nav>
    </details>
  )
}
