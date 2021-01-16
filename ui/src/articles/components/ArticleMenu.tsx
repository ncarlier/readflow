import React, { CSSProperties } from 'react'
import { useQuery } from '@apollo/client'

import DropdownMenu from '../../components/DropdownMenu'
import Loader from '../../components/Loader'
import { matchResponse } from '../../helpers'
import { GetOutgoingWebhooksResponse } from '../../settings/intergrations/outgoing-webhook/models'
import { GetOutgoingWebhooks } from '../../settings/intergrations/outgoing-webhook/queries'
import { Article } from '../models'
import SendLink from './SendLink'
import OfflineLink from './OfflineLink'
import ShareLink from './ShareLink'
import DownloadLink from './DownloadLink'

interface Props {
  article: Article
  keyboard?: boolean
  style?: CSSProperties
}

const OutgoingWebhooksMenuItems = (attrs: Props) => {
  const { data, error, loading } = useQuery<GetOutgoingWebhooksResponse>(GetOutgoingWebhooks)
  const render = matchResponse<GetOutgoingWebhooksResponse>({
    Loading: () => (
      <li>
        <Loader center />
      </li>
    ),
    Error: (err) => <li>{err.message}</li>,
    Data: ({ outgoingWebhooks }) =>
      outgoingWebhooks.map((webhook) => (
        <li key={`wh-${webhook.id}`}>
          <SendLink webhook={webhook} {...attrs} />
        </li>
      )),
  })

  return <>{render(loading, data, error)}</>
}

export default (props: Props) => {
  const { style, ...attrs } = props
  const nvg: any = window.navigator

  return (
    <DropdownMenu style={style}>
      <ul>
        {nvg.share && (
          <li>
            <ShareLink {...attrs} />
          </li>
        )}
        <li>
          <DownloadLink {...attrs} />
        </li>
        <li>
          <OfflineLink {...attrs} />
        </li>
        <OutgoingWebhooksMenuItems {...attrs} />
      </ul>
    </DropdownMenu>
  )
}
