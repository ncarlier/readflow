import React from 'react'
import { useQuery } from 'react-apollo-hooks'
import { RouteComponentProps } from 'react-router'

import { matchResponse } from '../../common/helpers'
import Loader from '../../common/Loader'
import Panel from '../../common/Panel'
import ErrorPanel from '../../error/ErrorPanel'
import { usePageTitle } from '../../hooks'
import EditRuleForm from './EditRuleForm'
import { GetRuleResponse } from './models'
import { GetRule } from './queries'

type AllProps = RouteComponentProps<{ id: string }>

export default ({ history, match }: AllProps) => {
  const { id } = match.params
  usePageTitle(`Settings - Edit rule #${id}`)

  const { data, error, loading } = useQuery<GetRuleResponse>(GetRule, {
    variables: { id }
  })

  const render = matchResponse<GetRuleResponse>({
    Loading: () => <Loader />,
    Error: err => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: ({ rule }) => {
      if (rule) {
        return <EditRuleForm data={rule} history={history} />
      } else {
        return <ErrorPanel title="Not found">Rule #${id} not found.</ErrorPanel>
      }
    },
    Other: () => <ErrorPanel>Unable to fetch rule #${id} details!</ErrorPanel>
  })

  return <Panel>{render(data, error, loading)}</Panel>
}
