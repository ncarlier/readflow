import React from 'react'

import logo from './webhook.svg'

interface Props {
  maxWidth?: string
}

export default ({ maxWidth = '2em' }: Props) => (
  <img src={logo} alt="webhook" style={{ maxWidth, verticalAlign: 'middle' }} />
)
