import React from 'react'

type Props = {
  max?: number
}

export default ({ max = 10 }: Props) => <>{[...Array(max)].map((e, i) => {
  const priority = i + 1
  let label = priority + ""
  if (priority == 1) {
    label = priority + " (lowest)"
  }
  if (priority == max) {
    label = priority + " (highest)"
  }
  return <option key={`priority-${priority}`} value={priority}>
    {label}
  </option>
})}
</>
