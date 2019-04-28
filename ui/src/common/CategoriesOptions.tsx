import React from 'react'
import { useQuery } from 'react-apollo-hooks'

import { GetCategoriesResponse } from '../categories/models'
import { GetCategories } from '../categories/queries'
import { matchResponse } from './helpers'

export default () => {
  const { data, error, loading } = useQuery<GetCategoriesResponse>(GetCategories)

  const render = matchResponse<GetCategoriesResponse>({
    Loading: () => <option>LOADING...</option>,
    Error: err => <option>{err.message}</option>,
    Data: data =>
      data.categories.map(category => (
        <option key={`cat-${category.id}`} value={category.id}>
          {category.title}
        </option>
      )),
    Other: () => <option>Unable to fetch categories!</option>
  })

  return <>{render(data, error, loading)}</>
}
