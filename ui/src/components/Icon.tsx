import React from 'react'

interface Props {
  name: string
}

export const Icon = ({ name }: Props) => <i className="material-icons">{name}</i>
