import React, { FC, PropsWithChildren } from 'react'

import { useMedia } from '../hooks'

interface Props extends PropsWithChildren {
  mobile?: boolean
  desktop?: boolean
}

export const DeviceSpecific: FC<Props> = ({ mobile, desktop, children }) => {
  // const isMobile = useMedia('(max-width: 400px)')
  const isMobile = useMedia('(max-width: 767px)')
  const isDesktop = useMedia('(min-width: 767px)')
  return (isMobile && mobile) || (isDesktop && desktop) ? <>{children}</> : null
}

export default DeviceSpecific
