import React, { FunctionComponent } from 'react'

import { useMedia } from '../hooks'

interface Props {
  mobile?: boolean
  desktop?: boolean
}

const DeviceSpecific: FunctionComponent<Props> = ({ mobile, desktop, children }) => {
  // const isMobile = useMedia('(max-width: 400px)')
  const isMobile = useMedia('(max-width: 767px)')
  const isDesktop = useMedia('(min-width: 767px)')
  if ((isMobile && mobile) || (isDesktop && desktop)) {
    return <>{children}</>
  }
  return null
}

export default DeviceSpecific
