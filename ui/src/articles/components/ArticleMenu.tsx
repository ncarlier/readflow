import React, { CSSProperties } from 'react'
import { useQuery } from 'react-apollo-hooks'

import DropdownMenu from '../../components/DropdownMenu'
import Loader from '../../components/Loader'
import { matchResponse } from '../../helpers'
import { GetArchiveServicesResponse } from '../../settings/archive-services/models'
import { GetArchiveServices } from '../../settings/archive-services/queries'
import { Article } from '../models'
import ArchiveLink from './ArchiveLink'
import OfflineLink from './OfflineLink'
import ShareLink from './ShareLink'

interface Props {
  article: Article
  keyboard?: boolean
  style?: CSSProperties
}

export default (props: Props) => {
  const { style, ...attrs } = props
  const nvg: any = window.navigator

  const { data, error, loading } = useQuery<GetArchiveServicesResponse>(GetArchiveServices)

  const renderArchiveServices = matchResponse<GetArchiveServicesResponse>({
    Loading: () => (
      <li>
        <Loader />
      </li>
    ),
    Error: (err) => <li>{err.message}</li>,
    Data: ({ archivers }) =>
      archivers.map((service) => (
        <li key={`as-${service.id}`}>
          <ArchiveLink service={service} {...attrs} />
        </li>
      )),
    Other: () => <li>Unknown error</li>,
  })

  return (
    <DropdownMenu style={style}>
      <ul>
        {nvg.share && (
          <li>
            <ShareLink {...attrs} />
          </li>
        )}
        <li>
          <OfflineLink {...attrs} />
        </li>
        {renderArchiveServices(data, error, loading)}
      </ul>
    </DropdownMenu>
  )
}
