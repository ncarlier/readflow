import React, { useCallback, useEffect, useRef, useState } from 'react'
import { useApolloClient } from '@apollo/client'

import { Category } from '../../categories/models'
import { Button } from '../../components'
import { GetArticlesResponse } from '../models'
import { GetNbNewArticles } from '../queries'
import { usePageVisibility } from '../../hooks'

const NewArticlesLabel = ({ nb }: { nb: number }) => {
  switch (true) {
    case nb > 1:
      return <>View {nb} new articles</>
    case nb === 1:
      return <>View new article</>
    case nb < 0:
      return <>Refresh</>
    default:
      return null
  }
}

interface Props {
  current: number
  category?: Category
  refresh: () => Promise<any>
}

export const NewArticlesAvailable = ({ current, category, refresh }: Props) => {
  const [nbItems, setNbItems] = useState(0)
  const visibility = usePageVisibility()
  const visibilityDetectionRef = useRef(false)

  const client = useApolloClient()

  const reload = useCallback(async () => {
    try {
      await refresh()
    } finally {
      setNbItems(0)
    }
  }, [refresh])

  const getNbArticlesToRead = useCallback(
    async (lastCount: number) => {
      try {
        const { errors, data } = await client.query<GetArticlesResponse>({
          query: GetNbNewArticles,
          fetchPolicy: 'network-only',
          variables: { category: category ? category.id : undefined },
        })
        if (data) {
          const delta = data.articles.totalCount - lastCount
          setNbItems(delta)
        }
        if (errors) {
          throw new Error(errors[0].message)
        }
      } catch (err) {
        console.error(err)
      }
    },
    [category, client]
  )

  useEffect(() => {
    const timer = setInterval(() => getNbArticlesToRead(current), 60000)
    return () => {
      clearInterval(timer)
    }
  }, [current, getNbArticlesToRead])

  useEffect(() => {
    if (visibilityDetectionRef.current && visibility) {
      reload()
    } else {
      visibilityDetectionRef.current = true
    }
  }, [visibility, reload])

  if (nbItems !== 0) {
    return (
      <span style={{textAlign: 'center'}}>
        <Button variant='flat' onClick={reload}>
          <NewArticlesLabel nb={nbItems} />
        </Button>
      </span>
    )
  }
  return null
}
