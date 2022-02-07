/* eslint-disable react/jsx-no-target-blank */
import React from 'react'

import { TimeAgo } from '.'
import { getRoboHash } from '../helpers'
import styles from './UserInfos.module.css'
import { useCurrentUser } from '../contexts'
import { AUTHORITY, CLIENT_ID } from '../constants'

const getAccountURL = () =>
  AUTHORITY !== 'mock'
    ? AUTHORITY + '/account?referrer=' + CLIENT_ID + '&referrer_uri=' + encodeURI(document.location.href)
    : ''

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
      <a href={getAccountURL()} target="_blank" title="Go to my profile page">
        <img src={getRoboHash(user.hash, '48')} alt={user.username} />
      </a>
    </div>
  )
}

export default UserInfos
