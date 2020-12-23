import React from 'react'
import { useQuery } from '@apollo/client'
import { Link } from 'react-router-dom'

import DataTable, { OnSelectedFn } from '../../../components/DataTable'
import ButtonIcon from '../../../components/ButtonIcon'
import Loader from '../../../components/Loader'
import TimeAgo from '../../../components/TimeAgo'
import ErrorPanel from '../../../error/ErrorPanel'
import { GetIncomingWebhooksResponse, IncomingWebhook } from './models'
import { GetIncomingWebhooks } from './queries'

const definition = [
  {
    title: 'Alias',
    render: (val: IncomingWebhook) => (
      <Link title="Edit incoming webhook" to={`integrations/incoming-webhooks/${val.id}`} data-test-id={val.alias}>
        {val.alias}
      </Link>
    ),
  },
  {
    title: 'Last usage',
    render: (val: IncomingWebhook) => <TimeAgo dateTime={val.last_usage_at} />,
  },
  {
    title: 'Created',
    render: (val: IncomingWebhook) => <TimeAgo dateTime={val.created_at} />,
  },
  {
    title: 'Updated',
    render: (val: IncomingWebhook) => <TimeAgo dateTime={val.updated_at} />,
  },
  {
    title: '',
    render: (val: IncomingWebhook) => (
      <ButtonIcon title="Edit incoming webhook" as={Link} to={`integrations/incoming-webhooks/${val.id}`} icon="edit">
        {val.alias}
      </ButtonIcon>
    ),
  },
]

interface Props {
  onSelected?: OnSelectedFn
}

export default ({ onSelected }: Props) => {
  const { data, error, loading } = useQuery<GetIncomingWebhooksResponse>(GetIncomingWebhooks)

  return (
    <section>
      {loading && <Loader />}
      {error && <ErrorPanel title="Unable to fetch incoming webhooks">{error.message}</ErrorPanel>}
      {data && (
        <DataTable
          definition={definition}
          data={data.incomingWebhooks}
          onSelected={onSelected}
          emptyMessage="No incoming webhook yet"
        />
      )}
    </section>
  )
}
