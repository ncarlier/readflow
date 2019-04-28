import { ChangeEvent, useState } from 'react'

export default (initial: string) => {
  const [value, setValue] = useState(initial)

  return {
    value,
    setValue,
    reset: () => setValue(initial),
    bind: {
      value,
      onChange: (e: ChangeEvent<HTMLInputElement>) => setValue(e.target.value)
    }
  }
}
