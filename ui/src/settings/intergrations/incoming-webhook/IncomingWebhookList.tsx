import React from 'react'
import { useQuery } from '@apollo/client'
import { Link } from 'react-router-dom'

import { ButtonIcon, DataTable, ErrorPanel, Loader, TimeAgo, OnSelectedFn } from '../../../components'
import { GetIncomingWebhooksResponse, IncomingWebhook } from './models'
import { GetIncomingWebhooks } from './queries'
import { matchResponse } from '../../../helpers'

const IncomingWebhookDates = ({ val }: { val: IncomingWebhook }) => (
  <small>
    {val.last_usage_at && (
      <>
        Last usage <TimeAgo dateTime={val.last_usage_at} />
        <br />
      </>
    )}
    Created <TimeAgo dateTime={val.created_at} />
    {val.updated_at && (
      <>
        <br />
        Updated <TimeAgo dateTime={val.updated_at} />
      </>
    )}
  </small>
)

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
    title: 'Date(s)',
    render: (val: IncomingWebhook) => <IncomingWebhookDates val={val} />,
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

const IncomingWebhookList = ({ onSelected }: Props) => {
  const { data, error, loading } = useQuery<GetIncomingWebhooksResponse>(GetIncomingWebhooks)

  const render = matchResponse<GetIncomingWebhooksResponse>({
    Loading: () => <Loader center />,
    Error: (err) => <ErrorPanel title="Unable to fetch incoming webhooks">{err.message}</ErrorPanel>,
    Data: (d) => (
      <DataTable
        definition={definition}
        data={d.incomingWebhooks}
        onSelected={onSelected}
        emptyMessage="No incoming webhook yet"
      />
    ),
  })
  return <section>{render(loading, data, error)}</section>
}

export default IncomingWebhookList
