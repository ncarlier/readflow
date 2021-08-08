import React, { createRef, CSSProperties, FunctionComponent, useCallback, useEffect, useState } from 'react'
import Drawer from 'react-bottom-drawer'

import ButtonIcon from './ButtonIcon'

interface Props {
  style?: CSSProperties
  title?: string
  icon?: string
}

const DrawerMenu: FunctionComponent<Props> = (props) => {
  const { children, icon = 'more_vert', ...attrs } = props
  const [isVisible, setIsVisible] = useState(false)
  const ref = createRef<HTMLDivElement>()

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

export default DrawerMenu
