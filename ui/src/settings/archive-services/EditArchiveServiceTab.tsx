import React from 'react'

import { useQuery } from 'react-apollo-hooks'
import { RouteComponentProps } from 'react-router'

import Panel from '../../common/Panel'
import { usePageTitle } from '../../hooks'
import ErrorPanel from '../../error/ErrorPanel'
import { matchResponse } from '../../common/helpers'
import Loader from '../../common/Loader'
import EditArchiveServiceForm from './EditArchiveServiceForm'
import { GetArchiveService } from './queries'
import { GetArchiveServiceResponse } from './models'

type AllProps = RouteComponentProps<{id: string}>

export default ({ history, match }: AllProps) => {
  const { id } = match.params
  usePageTitle(`Settings - Edit archive service #${id}`)
  
  const { data, error, loading } = useQuery<GetArchiveServiceResponse>(GetArchiveService, {
    variables: {id}
  })
  
  const render = matchResponse<GetArchiveServiceResponse>({
    Loading: () => <Loader />,
    Error: (err) => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: ({archiver}) => {
      if (archiver) {
        return (
          <EditArchiveServiceForm data={archiver} history={history} />
        )
      } else {
        return (
          <ErrorPanel title="Not found">Archive service #${id} not found.</ErrorPanel>
        )
      }
    },
    Other: () => <ErrorPanel>Unable to fetch archive service #${id} details!</ErrorPanel>
  })

  return (
    <Panel>
      {render(data, error, loading)}
    </Panel>
  )
}
