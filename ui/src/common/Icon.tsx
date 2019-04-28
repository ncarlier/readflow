import React from 'react'

interface Props {
  name: string
}

export default ({ name }: Props) => <i className="material-icons">{name}</i>
