/* eslint-disable react/jsx-no-target-blank */
import React from 'react'
import gql from 'graphql-tag'
import { useQuery } from '@apollo/client'

import authService from '../auth'
import { ErrorPanel, Loader, TimeAgo } from '.'
import { matchResponse, getRoboHash } from '../helpers'
import styles from './UserInfos.module.css'

export const GetCurrentUser = gql`
  query {
    me {
      username
      hash
      hashid
      plan
      customer_id
      last_login_at
      created_at
    }
  }
`

export interface User {
  username: string
  hash: string
  hashid: string
  plan: string
  customer_id: string
  created_at: string
  last_login_at: string
}

export interface GetCurrentUserResponse {
  me: User
}

export const UserInfos = () => {
  const { data, error, loading } = useQuery<GetCurrentUserResponse>(GetCurrentUser)

  const render = matchResponse<GetCurrentUserResponse>({
    Loading: () => <Loader />,
    Error: (err) => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: (data) => (
      <>
        <span>
          <strong title={data.me.username}>{data.me.username}</strong>
          <small>
            Member <TimeAgo dateTime={data.me.created_at} />
          </small>
        </span>
        <a href={authService.getAccountUrl()} target="_blank" title="Go to my profile page">
          <img src={getRoboHash(data.me.hash, '48')} alt={data.me.username} />
        </a>
      </>
    ),
  })

  return <div className={styles.userInfos}>{render(loading, data, error)}</div>
}

export default UserInfos
