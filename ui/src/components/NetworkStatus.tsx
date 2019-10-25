import React, { ReactNode } from 'react'

import useOnlineStatus from '../hooks/useOnlineStatus'

interface Props {
  status: 'online' | 'offline'
  children: ReactNode
}

function NetworkStatus({ status, children }: Props) {
  const isOnline = useOnlineStatus()
  const display = (isOnline && status == 'online') || (!isOnline && status == 'offline')

  if (display) {
    return <>{children}</>
  } else {
    return null
  }
}

export default NetworkStatus
