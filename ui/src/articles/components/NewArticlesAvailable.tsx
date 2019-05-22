import React, { useCallback, useEffect, useState } from 'react'
import { useApolloClient } from 'react-apollo-hooks'

import { Category } from '../../categories/models'
import Button from '../../common/Button'
import Loader from '../../common/Loader'
import Panel from '../../common/Panel'
import { GetArticlesResponse } from '../models'
import { GetNbNewArticles } from '../queries'

const renderLabel = (nb: number) => {
  switch (true) {
    case nb > 1:
      return `View ${nb} new articles`
    case nb === 1:
      return 'View new article'
    case nb < 0:
      return 'View new articles'
    default:
      return ''
  }
}

interface Props {
  current: number
  category?: Category
  refresh: () => Promise<any>
}

export default ({ current, category, refresh }: Props) => {
  const [loading, setLoading] = useState(false)
  const [nbItems, setNbItems] = useState(0)

  const client = useApolloClient()

  const reload = useCallback(async () => {
    setLoading(true)
    await refresh()
  }, [refresh])

  const getNbArticlesToRead = async () => {
    // console.log('getNbArticlesToRead...')
    try {
      const { errors, data } = await client.query<GetArticlesResponse>({
        query: GetNbNewArticles,
        fetchPolicy: 'network-only',
        variables: { category: category ? category.id : undefined }
      })
      if (data) {
        const delta = data.articles.totalCount - current
        setNbItems(delta)
      }
      if (errors) {
        throw new Error(errors[0])
      }
    } catch (err) {
      console.error(err)
    }
  }

  useEffect(() => {
    const timer = setInterval(() => getNbArticlesToRead(), 60000)
    return () => {
      clearInterval(timer)
    }
  }, [])

  switch (true) {
    case loading:
      return <Loader />
    case nbItems != 0:
      return (
        <Panel>
          <Button onClick={reload}>{renderLabel(nbItems)}</Button>
        </Panel>
      )
    default:
      return null
  }
}
