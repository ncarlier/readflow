import React from 'react'
import { useQuery } from 'react-apollo-hooks'
import { RouteComponentProps } from 'react-router'

import { GetCategoryResponse } from '../../categories/models'
import { GetCategory } from '../../categories/queries'
import Loader from '../../components/Loader'
import Panel from '../../components/Panel'
import ErrorPanel from '../../error/ErrorPanel'
import { matchResponse } from '../../helpers'
import { usePageTitle } from '../../hooks'
import EditCategoryForm from './EditCategoryForm'

type AllProps = RouteComponentProps<{ id: string }>

export default ({ history, match }: AllProps) => {
  const { id } = match.params
  usePageTitle(`Settings - Edit category #${id}`)

  const { data, error, loading } = useQuery<GetCategoryResponse>(GetCategory, {
    variables: { id },
  })

  const render = matchResponse<GetCategoryResponse>({
    Loading: () => <Loader />,
    Error: (err) => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: ({ category }) => {
      if (category) {
        return <EditCategoryForm category={category} history={history} />
      } else {
        return <ErrorPanel title="Not found">Category #${id} not found.</ErrorPanel>
      }
    },
    Other: () => <ErrorPanel>Unable to fetch category #${id} details!</ErrorPanel>,
  })

  return <Panel>{render(data, error, loading)}</Panel>
}
