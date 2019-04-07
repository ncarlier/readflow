import React from 'react'

import { useQuery } from 'react-apollo-hooks'
import { RouteComponentProps } from 'react-router'

import Panel from '../../common/Panel'
import { usePageTitle } from '../../hooks'
import ErrorPanel from '../../error/ErrorPanel'
import { matchResponse } from '../../common/helpers'
import Loader from '../../common/Loader'
import EditApiKeyForm from './EditApiKeyForm'
import { GetApiKeyResponse } from './models';
import { GetApiKey } from './queries';

type AllProps = RouteComponentProps<{id: string}>

export default ({ history, match }: AllProps) => {
  const { id } = match.params
  usePageTitle(`Settings - Edit category #${id}`)
  
  const { data, error, loading } = useQuery<GetApiKeyResponse>(GetApiKey, {
    variables: {id}
  })
  
  const render = matchResponse<GetApiKeyResponse>({
    Loading: () => <Loader />,
    Error: (err) => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: ({apiKey}) => {
      if (apiKey) {
        return (
          <EditApiKeyForm data={apiKey} history={history} />
        )
      } else {
        return (
          <ErrorPanel title="Not found">API key #${id} not found.</ErrorPanel>
        )
      }
    },
    Other: () => <ErrorPanel>Unable to fetch API key #${id} details!</ErrorPanel>
  })

  return (
    <Panel>
      {render(data, error, loading)}
    </Panel>
  )
}
