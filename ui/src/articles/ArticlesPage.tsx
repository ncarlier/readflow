import React, { useCallback, useState } from 'react'
import { NetworkStatus, useQuery } from '@apollo/client'
import { RouteComponentProps } from 'react-router'

import { Category } from '../categories/models'
import { LocalConfiguration, useLocalConfiguration } from '../contexts/LocalConfigurationContext'
import { Center, ErrorPanel, Loader, Panel } from '../components'
import { getURLParam, matchResponse } from '../helpers'
import { Appbar, Page } from '../layout'
import { ArticleList, ArticlesPageMenu, NewArticlesAvailable, NoArticleBg, Search } from './components'
import { ArticleStatus, GetArticlesRequest, GetArticlesResponse } from './models'
import { GetArticles } from './queries'
import { useMedia } from '../hooks'

type Variant = 'inbox' | 'to_read' | 'history' | 'starred'

const variants = {
  inbox: {
    emptyLabel: 'your inbox is full of air',
    getTitle: (nb: number) => `${nb} article${nb > 1 ? 's' : ''} received`,
  },
  to_read: {
    emptyLabel: 'nothing more to read, time to relax',
    getTitle: (nb: number) => `${nb} article${nb > 1 ? 's' : ''} to read`,
  },
  history: {
    emptyLabel: 'your reading history is all clean',
    getTitle: (nb: number) => `${nb} read article${nb > 1 ? 's' : ''}`,
  },
  starred: {
    emptyLabel: 'only the sky is full of stars',
    getTitle: (nb: number) => `${nb} stared article${nb > 1 ? 's' : ''}`,
  },
}

interface Props {
  variant: Variant
  category?: Category
}

const buildArticlesRequest = (variant: Variant, props: AllProps, localConfig: LocalConfiguration) => {
  const { category, location } = props
  const params = new URLSearchParams(location.search)

  const req: GetArticlesRequest = {
    limit: getURLParam(params, 'limit', localConfig.limit),
    sortBy: null,
    sortOrder: getURLParam(params, 'order', localConfig.display[variant].order),
    status: null,
    starred: null,
    category: null,
    afterCursor: null,
    query: getURLParam(params, 'query', ''),
  }
  switch (variant) {
    case 'history':
      req.status = 'read'
      req.starred = false
      break
    case 'starred':
      req.starred = true
      req.sortBy = getURLParam(params, 'by', localConfig.display.starred.by)
      break
    default:
      req.status = getURLParam<ArticleStatus>(params, 'status', variant)
      if (category && category.id) {
        req.category = category.id
        const configKey = `cat_${category.id}`
        if (Object.prototype.hasOwnProperty.call(localConfig.display, configKey)) {
          req.sortOrder = getURLParam(params, 'order', localConfig.display[configKey].order)
        }
        if (getURLParam(params, 'starred', '' as string) === 'true') {
          req.starred = true
          req.status = null
        }
      }
  }

  return req
}

const getDisplayMode = (variant: Variant, localConfig: LocalConfiguration, category?: Category) => {
  if (category && category.id) {
    const key = `cat_${category.id}`
    if (Object.prototype.hasOwnProperty.call(localConfig.display, key)) {
      return localConfig.display[key].mode
    }
  }
  return localConfig.display[variant].mode
}

const computeTotalArticles = (data: GetArticlesResponse, status: string | null) => {
  let delta = 0
  if (status) {
    delta = data.articles.entries.filter((a) => a.status !== status).length
  }
  return data.articles.totalCount - delta
}

type AllProps = Props & RouteComponentProps

export default (props: AllProps) => {
  const { variant, category } = props

  const { localConfiguration } = useLocalConfiguration()
  const [req] = useState<GetArticlesRequest>(buildArticlesRequest(variant, props, localConfiguration))
  const { data, error, fetchMore, refetch, networkStatus } = useQuery<GetArticlesResponse>(GetArticles, {
    variables: req,
    notifyOnNetworkStatusChange: true,
  })
  const isMobileDisplay = useMedia('(max-width: 767px)')

  const fetchMoreArticles = useCallback(async () => {
    if (!data || !data.articles.hasNext) {
      return
    }
    console.log('fetching more articles...')
    await fetchMore({
      variables: { ...req, afterCursor: data.articles.endCursor },
      updateQuery: (prev, { fetchMoreResult }) => {
        if (!fetchMoreResult) return prev
        const nbFetchedArticles = fetchMoreResult.articles.entries.length
        console.log(nbFetchedArticles + ' article(s) fetched')
        let { entries } = prev.articles
        if (req.status) {
          entries = entries.filter((a) => a.status === req.status)
        }
        const ids = new Set(entries.map((a) => a.id))
        const articles = {
          ...fetchMoreResult.articles,
          entries: [...entries, ...fetchMoreResult.articles.entries.filter((a) => !ids.has(a.id))],
        }
        return { articles }
      },
    })
  }, [data, fetchMore, req])

  const refresh = useCallback(async () => {
    console.log('re-fetching articles ...')
    const { errors } = await refetch()
    if (errors) {
      console.error(errors)
    }
  }, [refetch])

  const render = matchResponse<GetArticlesResponse>({
    Loading: () => (
      <Center>
        <Loader />,
      </Center>
    ),
    Error: (err) => (
      <Panel>
        <ErrorPanel>{err.message}</ErrorPanel>
      </Panel>
    ),
    Data: (d) => {
      const mode = getDisplayMode(variant, localConfiguration, category)
      let { entries } = d.articles
      if (req.status) {
        entries = entries.filter((a) => a.status === req.status)
      }
      return (
        <>
          {variant === 'inbox' && (
            <NewArticlesAvailable current={computeTotalArticles(d, req.status)} category={category} refresh={refresh} />
          )}
          <ArticleList
            articles={entries}
            empty={<NoArticleBg name={variant} title={variants[variant].emptyLabel} />}
            hasMore={d.articles.hasNext}
            refetch={refetch}
            swipeable={isMobileDisplay && variant !== 'starred'}
            fetchMoreArticles={fetchMoreArticles}
            variant={mode}
          />
        </>
      )
    },
  })

  // Build title
  const nbArticles = data && data.articles ? computeTotalArticles(data, req.status) : 0
  let variantKey = variant
  if (category && req.status === 'read') variantKey = 'history'
  if (category && req.status === 'to_read') variantKey = 'to_read'
  if (category && req.starred) variantKey = 'starred'
  const title = variants[variantKey].getTitle(nbArticles) + (category ? ` in "${category?.title}"` : '')

  const $header = (
    <Appbar title={title}>
      <Search req={req} />
      <ArticlesPageMenu refresh={refresh} req={req} variant={variant} />
    </Appbar>
  )

  const loading = networkStatus === NetworkStatus.loading
  const refetching = networkStatus === NetworkStatus.refetch

  return (
    <Page title={title} header={$header} scrollToTop>
      {refetching && <Loader style={{position: 'fixed', right: '2em', top: '3em'}} />}
      {render(loading, data, error)}
    </Page>
  )
}
