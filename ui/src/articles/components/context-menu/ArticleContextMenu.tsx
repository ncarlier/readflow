import React from 'react'

import DropdownMenu, { DropDownOrigin } from '../../../components/DropdownMenu'
import { Article } from '../../models'
import OfflineLink from './OfflineLink'
import ShareLink from './ShareLink'
import DownloadLink from './DownloadLink'
import OutgoingWebhooksMenuItems from './OutgoingWebhooksMenuItems'

interface Props {
  article: Article
  keyboard?: boolean
  origin?: DropDownOrigin
}

export default (props: Props) => {
  const { origin, ...attrs } = props
  const nvg: any = window.navigator

  return (
    <DropdownMenu origin={origin} title="More actions...">
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
