import React from 'react'
import { useQuery } from '@apollo/client'
import { RouteComponentProps } from 'react-router'

import { GetCategoryResponse } from '../../categories/models'
import { GetCategory } from '../../categories/queries'
import { ErrorPanel, Loader, Panel } from '../../components'
import { matchResponse } from '../../helpers'
import { usePageTitle } from '../../hooks'
import EditCategoryForm from './EditCategoryForm'

const EditCategoryTab = ({ history, match }: RouteComponentProps<{ id: string }>) => {
  const { id } = match.params
  usePageTitle(`Settings - Edit category #${id}`)

  const { data, error, loading } = useQuery<GetCategoryResponse>(GetCategory, {
    variables: { id },
  })

  const render = matchResponse<GetCategoryResponse>({
    Loading: () => <Loader center />,
    Error: (err) => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: ({ category }) => {
      if (category) {
        return <EditCategoryForm category={category} history={history} />
      } else {
        return <ErrorPanel title="Not found">Category #${id} not found.</ErrorPanel>
      }
    },
  })

  return <Panel>{render(loading, data, error)}</Panel>
}

export default EditCategoryTab
