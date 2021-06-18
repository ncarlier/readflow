/* eslint-disable react/jsx-no-target-blank */
import gql from 'graphql-tag'
import React from 'react'
import { useQuery } from '@apollo/client'

import authService from '../auth'
import ErrorPanel from '../error/ErrorPanel'
import { matchResponse } from '../helpers'
import Loader from './Loader'
import TimeAgo from './TimeAgo'
import styles from './UserInfos.module.css'
import { getRoboHash } from '../helpers/avatar'

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

export default () => {
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
