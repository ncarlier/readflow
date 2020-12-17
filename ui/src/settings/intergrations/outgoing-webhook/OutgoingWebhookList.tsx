import React from 'react'
import { useQuery } from 'react-apollo-hooks'
import { Link } from 'react-router-dom'

import DataTable, { OnSelectedFn } from '../../../components/DataTable'
import Loader from '../../../components/Loader'
import TimeAgo from '../../../components/TimeAgo'
import ErrorPanel from '../../../error/ErrorPanel'
import { matchResponse } from '../../../helpers'
import { GetOutboundServicesResponse, OutboundService } from './models'
import { GetOutboundServices } from './queries'

const definition = [
  {
    title: 'Alias',
    render: (val: OutboundService) => (
      <Link title="Edit outbound service" to={`/settings/archive-services/${val.id}`}>
        {val.alias} {val.is_default && '(default)'}
      </Link>
    ),
  },
  {
    title: 'Provider',
    render: (val: OutboundService) => val.provider,
  },
  {
    title: 'Created',
    render: (val: OutboundService) => <TimeAgo dateTime={val.created_at} />,
  },
  {
    title: 'Updated',
    render: (val: OutboundService) => <TimeAgo dateTime={val.updated_at} />,
  },
]

interface Props {
  onSelected?: OnSelectedFn
}

export default ({ onSelected }: Props) => {
  const { data, error, loading } = useQuery<GetOutboundServicesResponse>(GetOutboundServices)

  const render = matchResponse<GetOutboundServicesResponse>({
    Loading: () => <Loader />,
    Error: (err) => <ErrorPanel title="Unable to fetch outbound services">{err.message}</ErrorPanel>,
    Data: (data) => (
      <DataTable
        definition={definition}
        data={data.outboundServices}
        onSelected={onSelected}
        emptyMessage="No outbound service yet"
      />
    ),
  })

  return <>{render(loading, data, error)}</>
}
