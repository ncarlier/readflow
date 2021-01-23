import React, { useCallback } from 'react'
import Icon from './Icon'

import styles from './Stars.module.css'

interface StarProps {
  icon?: string
  value: number
  onClick: (value: number) => void
}

const Star = ({ value, icon = 'star', onClick }: StarProps) => {
  /*
  const handleMouseEnter = useCallback(() => (name !== 'star' ? setName('star') : null), [name])
  const handleMouseLeave = useCallback(() => (active ? setName('star') : setName('star_outline')), [active])
  */

  return (
    <button onClick={() => onClick(value)} /*onMouseEnter={handleMouseEnter} onMouseLeave={handleMouseLeave}*/>
      <Icon name={icon} />
    </button>
  )
}

interface Props {
  value: number
  size?: number
  onChange: (value: number) => void
}

export default (props: Props) => {
  const { value, size = 5, onChange } = props
  const stars = Array.from(Array(size).keys())

  const handleClick = useCallback(
    (val: number) => {
      if (val !== value) {
        onChange(val)
      }
    },
    [value, onChange]
  )

  return (
    <ol className={styles.stars}>
      <li>
        <Star value={0} onClick={handleClick} icon="star_outline" />
      </li>
      {stars.map((i) => (
        <li key={i} className={value > i ? styles.active : ''}>
          <Star value={i + 1} onClick={handleClick} />
        </li>
      ))}
    </ol>
  )
}
