/* eslint-disable react/jsx-no-target-blank */
import React from 'react'

import authService from '../auth'
import { TimeAgo } from '.'
import { getRoboHash } from '../helpers'
import styles from './UserInfos.module.css'
import { useCurrentUser } from '../contexts'

export const UserInfos = () => {
  const user = useCurrentUser()
  if (!user) {
    return null
  }
  return (
    <div className={styles.userInfos}>
      <span>
        <strong title={user.username}>{user.username}</strong>
        <small>
          Member <TimeAgo dateTime={user.created_at} />
        </small>
      </span>
      <a href={authService.getAccountUrl()} target="_blank" title="Go to my profile page">
        <img src={getRoboHash(user.hash, '48')} alt={user.username} />
      </a>
    </div>
  )
}

export default UserInfos
