import React from 'react'
import { useQuery } from '@apollo/client'
import { RouteComponentProps } from 'react-router'

import Loader from '../../../components/Loader'
import Panel from '../../../components/Panel'
import ErrorPanel from '../../../error/ErrorPanel'
import { matchResponse } from '../../../helpers'
import { usePageTitle } from '../../../hooks'
import EditOutgoingWebhookForm from './EditOutgoingWebhookForm'
import { GetOutgoingWebhookResponse } from './models'
import { GetOutgoingWebhook } from './queries'

type AllProps = RouteComponentProps<{ id: string }>

export default ({ history, match }: AllProps) => {
  const { id } = match.params
  usePageTitle(`Settings - Edit outgoing webhook #${id}`)

  const { data, error, loading } = useQuery<GetOutgoingWebhookResponse>(GetOutgoingWebhook, {
    variables: { id },
  })

  const render = matchResponse<GetOutgoingWebhookResponse>({
    Loading: () => <Loader center />,
    Error: (err) => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: ({ outgoingWebhook }) => {
      if (outgoingWebhook) {
        return <EditOutgoingWebhookForm data={outgoingWebhook} history={history} />
      } else {
        return <ErrorPanel title="Not found">Outgoing webhook #${id} not found.</ErrorPanel>
      }
    },
  })

  return <Panel>{render(loading, data, error)}</Panel>
}
