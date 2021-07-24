import React from 'react'
import { useQuery } from '@apollo/client'
import { Link } from 'react-router-dom'
import ButtonIcon from '../../../components/ButtonIcon'

import DataTable, { OnSelectedFn } from '../../../components/DataTable'
import Loader from '../../../components/Loader'
import TimeAgo from '../../../components/TimeAgo'
import ErrorPanel from '../../../error/ErrorPanel'
import Logo from '../../../logos/Logo'
import { GetOutgoingWebhooksResponse, OutgoingWebhook } from './models'
import { GetOutgoingWebhooks } from './queries'
import { matchResponse } from '../../../helpers'

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
    title: 'Created',
    render: (val: OutgoingWebhook) => <TimeAgo dateTime={val.created_at} />,
  },
  {
    title: 'Updated',
    render: (val: OutgoingWebhook) => <TimeAgo dateTime={val.updated_at} />,
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

export default ({ onSelected }: Props) => {
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
