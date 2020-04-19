/* eslint-disable react/jsx-no-target-blank */
import gql from 'graphql-tag'
import React from 'react'
import { useQuery } from 'react-apollo-hooks'

import authService from '../auth'
import ErrorPanel from '../error/ErrorPanel'
import { matchResponse } from '../helpers'
import Loader from './Loader'
import TimeAgo from './TimeAgo'
import styles from './UserInfos.module.css'

export const GetCurrentUser = gql`
  query {
    me {
      username
      hash
      plan
      last_login_at
      created_at
    }
  }
`

interface User {
  username: string
  hash: string
  plan: string
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
          <img src={`https://www.gravatar.com/avatar/${data.me.hash}?d=mp&s=42"`} alt={data.me.username} />
        </a>
      </>
    ),
    Other: () => <ErrorPanel>Unable to fetch current user infos!</ErrorPanel>,
  })

  return <div className={styles.userInfos}>{render(data, error, loading)}</div>
}
