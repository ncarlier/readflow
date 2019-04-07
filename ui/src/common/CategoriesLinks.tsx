import React from 'react'
import { useQuery } from 'react-apollo-hooks'

import { matchResponse } from './helpers'
import Loader from './Loader'
import LinkIcon from './LinkIcon'
import { GetCategoriesResponse } from '../categories/models'
import { GetCategories } from '../categories/queries'

export default () => {
  const { data, error, loading } = useQuery<GetCategoriesResponse>(GetCategories)
  
  const render = matchResponse<GetCategoriesResponse>({
    Loading: () => <Loader />,
    Error: (err) => <span>{err.message}</span>,
    Data: ({categories}) => <ul>
      {categories.map( category => (
        <li key={`cat-${category.id}`}>  
          <LinkIcon to={`/categories/${category.id}`} icon="bookmark">
            {category.title}
          </LinkIcon>
        </li>
      ))}
    </ul>,
    Other: () => <span>Unable to fetch categories!</span>
  })

  return (<>{render(data, error, loading)}</>)
}
