/* eslint-disable react/jsx-no-target-blank */
import React from 'react'

import { TimeAgo } from '.'
import styles from './UserInfos.module.css'
import { useCurrentUser } from '../contexts'
import { AUTHORITY, CLIENT_ID } from '../config'
import { getAPIURL } from '../helpers'

const getAccountURL = () => {
  if (AUTHORITY.includes('realms')) {
    // keycloak specifics
    return `${AUTHORITY}/account?referrer=${CLIENT_ID}&referrer_uri=${encodeURI(document.location.href)}`
  }
  return AUTHORITY
}

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
      {
        AUTHORITY !== 'none' ?
          <a href={getAccountURL()} target="_blank" title="Go to my profile page">
            <img src={getAPIURL(`/avatar/${user.hash}`)} alt={user.username} crossOrigin='anonymous' />
          </a>
          :
          <img src={getAPIURL(`/avatar/${user.hash}`)} alt={user.username} crossOrigin='anonymous' />
      }
    </div>
  )
}

export default UserInfos
