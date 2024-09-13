import React, { createRef, CSSProperties, FC, MouseEventHandler, PropsWithChildren, useCallback, useEffect } from 'react'

import { ButtonIcon } from '.'
import styles from './DropDownMenu.module.css'

export interface DropDownOrigin {
  vertical: 'top' | 'bottom' | number
  horizontal: 'left' | 'right' | number
}

const getDropdownOriginStyle = (origin: DropDownOrigin = { horizontal: 'left', vertical: 'bottom' }): CSSProperties => {
  const result: CSSProperties = {}
  if (typeof origin.horizontal === 'number') {
    result.left = origin.horizontal + 'px'
  } else if (origin.horizontal === 'left') {
    result.right = 0
  } else {
    result.left = '100%'
  }
  if (typeof origin.vertical === 'number') {
    result.bottom = origin.vertical + 'px'
  } else if (origin.vertical === 'bottom') {
    result.top = '100%'
  } else {
    result.bottom = '100%'
  }

  return result
}

interface Props extends PropsWithChildren {
  style?: CSSProperties
  title?: string
  icon?: string
  origin?: DropDownOrigin
}

export const DropDownMenu: FC<Props> = (props) => {
  const { children, icon = 'more_vert', origin, ...attrs } = props
  const ref = createRef<HTMLDetailsElement>()

  const handleClickOutside = useCallback(
    (e: MouseEvent) => {
      const $el = e.target
      if (!($el instanceof Element)) return
      if (ref.current && !ref.current.contains($el)) {
        ref.current.removeAttribute('open')
      }
    },
    [ref]
  )

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
  }, [handleClickOutside])
  
  return (
    <details ref={ref} className={styles.menu}>
      <summary>
        <ButtonIcon icon={icon} {...attrs} onClick={handleClickMenu} />
      </summary>
      <nav className={styles.nav} style={getDropdownOriginStyle(origin)} tabIndex={-1}>
        {children}
      </nav>
    </details>
  )
}
