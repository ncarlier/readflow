import React from 'react'
import { useQuery } from 'react-apollo-hooks'

import { GetCategoriesResponse } from '../categories/models'
import { GetCategories } from '../categories/queries'
import { matchResponse } from './helpers'
import LinkIcon from './LinkIcon'
import Loader from './Loader'

export default () => {
  const { data, error, loading } = useQuery<GetCategoriesResponse>(GetCategories)

  const render = matchResponse<GetCategoriesResponse>({
    Loading: () => <Loader />,
    Error: err => <span>{err.message}</span>,
    Data: data => (
      <ul>
        {data.categories.map(category => (
          <li key={`cat-${category.id}`}>
            <LinkIcon to={`/categories/${category.id}`} icon="bookmark">
              {category.title}
            </LinkIcon>
          </li>
        ))}
      </ul>
    ),
    Other: () => <span>Unable to fetch categories!</span>
  })

  return <>{render(data, error, loading)}</>
}
