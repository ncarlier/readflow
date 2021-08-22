import React, { FC } from 'react'

import { useOnlineStatus } from '../hooks'

interface Props {
  status: 'online' | 'offline'
}

export const NetworkStatus: FC<Props> = ({ status, children }) => {
  const isOnline = useOnlineStatus()
  const display = (isOnline && status === 'online') || (!isOnline && status === 'offline')

  return display ? <>{children}</> : null
}
