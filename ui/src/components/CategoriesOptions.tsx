import React from 'react'
import { useQuery } from '@apollo/client'

import { GetCategoriesResponse } from '../categories/models'
import { GetCategories } from '../categories/queries'
import { matchResponse } from '../helpers'

export const CategoriesOptions = () => {
  const { data, error, loading } = useQuery<GetCategoriesResponse>(GetCategories)

  const render = matchResponse<GetCategoriesResponse>({
    Loading: () => <option>LOADING...</option>,
    Error: (err) => <option>{err.message}</option>,
    Data: (data) =>
      data.categories.entries.map((category) => (
        <option key={`cat-${category.id}`} value={category.id}>
          {category.title}
        </option>
      )),
  })

  return <>{render(loading, data, error)}</>
}
