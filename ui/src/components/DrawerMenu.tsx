import React, { createRef, CSSProperties, FC, PropsWithChildren, useCallback, useEffect, useState } from 'react'
import Drawer from 'react-bottom-drawer'

import { ButtonIcon } from '.'
import { useKeyboard } from '../hooks'

interface Props extends PropsWithChildren {
  style?: CSSProperties
  title?: string
  icon?: string
  kbs?: string
}

export const DrawerMenu: FC<Props> = (props) => {
  const { children, icon = 'more_vert', kbs, ...attrs } = props
  const [isVisible, setIsVisible] = useState(false)
  const ref = createRef<HTMLDivElement>()

  const toggle = React.useCallback(() => {
    setIsVisible(!isVisible)
  }, [isVisible])
  const onClose = React.useCallback(() => {
    setIsVisible(false)
  }, [])
  const handleClickMenu = () => setIsVisible(true)
  const handleClickOutside = useCallback(
    (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setIsVisible(false)
      }
    },
    [ref]
  )
  useEffect(() => {
    document.addEventListener('click', handleClickOutside, { capture: true })
    return () => {
      document.removeEventListener('click', handleClickOutside, { capture: true })
    }
  }, [handleClickOutside])

  useKeyboard(kbs, toggle, !!kbs)
  attrs.title += kbs ? ` [${kbs}]` : ''
  
  return (
    <>
      <ButtonIcon icon={icon} {...attrs} onClick={handleClickMenu} />
      <div ref={ref}>
        <Drawer isVisible={isVisible} onClose={onClose} className="drawer">
          {children}
        </Drawer>
      </div>
    </>
  )
}
