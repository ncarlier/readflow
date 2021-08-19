import React from 'react'

interface Props {
  total: number
  value: number
}

export default ({ total, value }: Props) => (
  <ul>
    <li>
      Downloading: &nbsp;
      <progress value={value} max={total}>
        {value}
      </progress>
      &nbsp;
      {value} / {total} bytes
    </li>
  </ul>
)
