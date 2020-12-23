import React from 'react'
import { useQuery } from '@apollo/client'
import { RouteComponentProps } from 'react-router'

import Loader from '../../../components/Loader'
import Panel from '../../../components/Panel'
import ErrorPanel from '../../../error/ErrorPanel'
import { matchResponse } from '../../../helpers'
import { usePageTitle } from '../../../hooks'
import EditIncomingWebhookForm from './EditIncomingWebhookForm'
import { GetIncomingWebhookResponse } from './models'
import { GetIncomingWebhook } from './queries'

type AllProps = RouteComponentProps<{ id: string }>

export default ({ history, match }: AllProps) => {
  const { id } = match.params
  usePageTitle(`Settings - Edit incoming Webhook #${id}`)

  const { data, error, loading } = useQuery<GetIncomingWebhookResponse>(GetIncomingWebhook, {
    variables: { id },
  })

  const render = matchResponse<GetIncomingWebhookResponse>({
    Loading: () => <Loader />,
    Error: (err) => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: ({ incomingWebhook }) => {
      if (incomingWebhook) {
        return <EditIncomingWebhookForm data={incomingWebhook} history={history} />
      } else {
        return <ErrorPanel title="Not found">Incoming webhook #${id} not found.</ErrorPanel>
      }
    },
  })

  return <Panel>{render(loading, data, error)}</Panel>
}
