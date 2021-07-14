import React from 'react'

import logo from './feedpushr.svg'

interface Props {
  maxWidth?: string
}

export default ({ maxWidth = '2em' }: Props) => (
  <img src={logo} alt="feedpushr" style={{ maxWidth, verticalAlign: 'middle' }} />
)
