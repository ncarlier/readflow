import React from 'react'

import { ReactComponent as history } from './bg/min/clean_sweep.svg'
import { ReactComponent as inbox } from './bg/min/empty_mailbox.svg'
import { ReactComponent as offline } from './bg/min/dyno.svg'
import { ReactComponent as starred } from './bg/min/stars.svg'
import { ReactComponent as to_read } from './bg/min/empty_couch.svg'

import styles from './NoArticleBg.module.css'

const backgrounds = {
  history,
  inbox,
  offline,
  starred,
  to_read,
}

interface Props {
  name: 'history' | 'inbox' | 'offline' | 'starred' | 'to_read'
  title?: string
}

export const NoArticleBg = ({ name, title }: Props) => {
  const Bg = backgrounds[name]
  return <Bg className={styles.bg} title={title}/>
}
