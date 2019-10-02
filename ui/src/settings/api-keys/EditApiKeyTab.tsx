import React from 'react'
import { useQuery } from 'react-apollo-hooks'
import { RouteComponentProps } from 'react-router'

import { matchResponse } from '../../helpers'
import Loader from '../../components/Loader'
import Panel from '../../components/Panel'
import ErrorPanel from '../../error/ErrorPanel'
import { usePageTitle } from '../../hooks'
import EditApiKeyForm from './EditApiKeyForm'
import { GetApiKeyResponse } from './models'
import { GetApiKey } from './queries'

type AllProps = RouteComponentProps<{ id: string }>

export default ({ history, match }: AllProps) => {
  const { id } = match.params
  usePageTitle(`Settings - Edit API key #${id}`)

  const { data, error, loading } = useQuery<GetApiKeyResponse>(GetApiKey, {
    variables: { id }
  })

  const render = matchResponse<GetApiKeyResponse>({
    Loading: () => <Loader />,
    Error: err => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: ({ apiKey }) => {
      if (apiKey) {
        return <EditApiKeyForm data={apiKey} history={history} />
      } else {
        return <ErrorPanel title="Not found">API key #${id} not found.</ErrorPanel>
      }
    },
    Other: () => <ErrorPanel>Unable to fetch API key #${id} details!</ErrorPanel>
  })

  return <Panel>{render(data, error, loading)}</Panel>
}
