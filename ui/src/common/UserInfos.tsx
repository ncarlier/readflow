import React, { useState, useEffect } from 'react'

import styles from './UserInfos.module.css'
import TimeAgo from './TimeAgo'
import gql from 'graphql-tag'
import { useQuery } from 'react-apollo-hooks'
import { matchResponse } from './helpers'
import Loader from './Loader'
import ErrorPanel from '../error/ErrorPanel'
import authService from '../auth/AuthService'

const GetCurrentUser = gql`
  query {
    me {
      username
      hash
      last_login_at
      created_at
    }
  }
`

type User = {
  username: string
  hash: string
  created_at: string
  last_login_at: string
}

export type GetCurrentUserResponse = {
  me: User
}

export default () => {
  const { data, error, loading } = useQuery<GetCurrentUserResponse>(GetCurrentUser)
  
  const render = matchResponse<GetCurrentUserResponse>({
    Loading: () => <Loader />,
    Error: (err) => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: ({me}) => (
      <>
        <span>
          <strong title={me.username}>{me.username}</strong>
          <small>Member <TimeAgo dateTime={me.created_at} /></small>
        </span>
        <a href={authService.getAccountUrl()} target="_blank" title="Go to my profile page">
          <img src={`https://www.gravatar.com/avatar/${me.hash}?d=mp&s=42"`} />
        </a>
      </>
    ),
    Other: () => <ErrorPanel>Unable to fetch current user infos!</ErrorPanel>
  })

  return (
    <div className={styles.userInfos}>
      {render(data, error, loading)}
    </div>
  )
}
