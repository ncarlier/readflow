import React from 'react'
import { useQuery } from '@apollo/client'
import { Link } from 'react-router-dom'

import { ButtonIcon, ErrorPanel, Loader, TimeAgo, Logo, DataTable, OnSelectedFn } from '../../../components'
import { GetOutgoingWebhooksResponse, OutgoingWebhook } from './models'
import { GetOutgoingWebhooks } from './queries'
import { matchResponse } from '../../../helpers'

const OutgoingWebhookDates = ({ val }: { val: OutgoingWebhook }) => (
  <small>
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
    render: (val: OutgoingWebhook) => (
      <Link title="Edit outgoing webhook" to={`integrations/outgoing-webhooks/${val.id}`}>
        {val.alias} {val.is_default && '(default)'}
      </Link>
    ),
  },
  {
    title: 'Provider',
    render: (val: OutgoingWebhook) => (
      <Logo name={val.provider} style={{ maxWidth: '2em' }} title={`${val.provider} provider`} />
    ),
  },
  {
    title: 'Date(s)',
    render: (val: OutgoingWebhook) => <OutgoingWebhookDates val={val} />,
  },
  {
    title: '',
    render: (val: OutgoingWebhook) => (
      <ButtonIcon title="Edit outgoing webhook" as={Link} to={`integrations/outgoing-webhooks/${val.id}`} icon="edit">
        {val.alias}
      </ButtonIcon>
    ),
  },
]

interface Props {
  onSelected?: OnSelectedFn
}

const OutgoingWebhookList = ({ onSelected }: Props) => {
  const { data, error, loading } = useQuery<GetOutgoingWebhooksResponse>(GetOutgoingWebhooks)

  const render = matchResponse<GetOutgoingWebhooksResponse>({
    Loading: () => <Loader center />,
    Error: (err) => <ErrorPanel title="Unable to fetch outgoing webhooks">{err.message}</ErrorPanel>,
    Data: (d) => (
      <DataTable
        definition={definition}
        data={d.outgoingWebhooks}
        onSelected={onSelected}
        emptyMessage="No outgoing webhook yet"
      />
    ),
  })
  return <section>{render(loading, data, error)}</section>
}

export default OutgoingWebhookList
