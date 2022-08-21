import React from 'react'

import history from './bg/clean_sweep.svg'
import inbox from './bg/empty_mailbox.svg'
import offline from './bg/dyno.svg'
import starred from './bg/stars.svg'
import to_read from './bg/empty_couch.svg'

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

export const NoArticleBg = ({ name, title }: Props) => (
  <img src={backgrounds[name]} alt={name} className={styles.bg} title={title} />
)
