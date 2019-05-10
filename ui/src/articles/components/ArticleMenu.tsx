import React from 'react'
import { useQuery } from 'react-apollo-hooks'

import DropdownMenu from '../../common/DropdownMenu'
import { matchResponse } from '../../common/helpers'
import Loader from '../../common/Loader'
import { GetArchiveServicesResponse } from '../../settings/archive-services/models'
import { GetArchiveServices } from '../../settings/archive-services/queries'
import { Article } from '../models'
import ArchiveLink from './ArchiveLink'
import OfflineLink from './OfflineLink'
import ShareLink from './ShareLink'

interface Props {
  article: Article
  keyboard?: boolean
}

type AllProps = Props

export default (props: AllProps) => {
  const nvg: any = window.navigator

  const { data, error, loading } = useQuery<GetArchiveServicesResponse>(GetArchiveServices)

  const renderArchiveServices = matchResponse<GetArchiveServicesResponse>({
    Loading: () => (
      <li>
        <Loader />
      </li>
    ),
    Error: err => <li>{err.message}</li>,
    Data: ({ archivers }) =>
      archivers.map(service => (
        <li key={`as-${service.id}`}>
          <ArchiveLink service={service} {...props} />
        </li>
      )),
    Other: () => <li>Unknown error</li>
  })

  return (
    <DropdownMenu>
      <ul>
        {nvg.share && (
          <li>
            <ShareLink {...props} />
          </li>
        )}
        <li>
          <OfflineLink {...props} />
        </li>
        {renderArchiveServices(data, error, loading)}
      </ul>
    </DropdownMenu>
  )
}
