import React from 'react'

type Props = {
  name: string
}

export default ({ name }: Props) => (
  <i className="material-icons">{ name }</i>
)
