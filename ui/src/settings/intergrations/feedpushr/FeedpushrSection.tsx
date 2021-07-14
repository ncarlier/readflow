import React from 'react'
import { useQuery } from '@apollo/client'

import Button from '../../../components/Button'
import Loader from '../../../components/Loader'
import { GetCurrentUser, GetCurrentUserResponse } from '../../../components/UserInfos'
import ErrorPanel from '../../../error/ErrorPanel'
import { matchResponse } from '../../../helpers'
import FeedpushrLogo from './FeedpushrLogo'

export default () => {
  const { data, error, loading } = useQuery<GetCurrentUserResponse>(GetCurrentUser)

  const render = matchResponse<GetCurrentUserResponse>({
    Loading: () => <Loader center />,
    Error: (err) => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: ({ me }) => {
      if (me.plan !== 'premium') {
        return
      }
      return (
        <section>
          <header>
            <h2>
              <FeedpushrLogo />
              Feedpushr
            </h2>
            <Button as={'a'} href={`https://feedpushr.nunux.org/${me.hashid}`} target="_blank" title="Manage my feeds">
              Manage
            </Button>
          </header>
          <p>Thanks to your user plan you can manage RSS feeds with Feedpushr and send articles to your readflow.</p>
        </section>
      )
    },
  })

  return <>{render(loading, data, error)}</>
}
