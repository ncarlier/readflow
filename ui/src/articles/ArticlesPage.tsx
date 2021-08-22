import React, { useCallback, useContext, useState } from 'react'
import { NetworkStatus, useQuery } from '@apollo/client'
import { RouteComponentProps } from 'react-router'

import { Category } from '../categories/models'
import { LocalConfiguration, LocalConfigurationContext } from '../contexts/LocalConfigurationContext'
import { Center, ErrorPanel, Loader, Panel } from '../components'
import { getURLParam, matchResponse } from '../helpers'
import { Appbar, Page } from '../layout'
import { ArticleList, ArticlesPageMenu, NewArticlesAvailable, Search } from './components'
import { ArticleStatus, GetArticlesRequest, GetArticlesResponse } from './models'
import { GetArticles } from './queries'
import { useMedia } from '../hooks'

type Variant = 'unread' | 'history' | 'starred'

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
    sortOrder: getURLParam(params, 'order', localConfig.display.unread.order),
    status: 'unread',
    starred: null,
    category: null,
    afterCursor: null,
    query: getURLParam(params, 'query', ''),
  }
  switch (variant) {
    case 'history':
      req.status = 'read'
      req.sortOrder = getURLParam(params, 'order', localConfig.display.history.order)
      break
    case 'starred':
      req.status = null
      req.starred = true
      req.sortOrder = getURLParam(params, 'order', localConfig.display.starred.order)
      req.sortBy = getURLParam(params, 'by', localConfig.display.starred.by)
      break
    case 'unread':
      if (category && category.id) {
        req.category = category.id
        req.status = getURLParam<ArticleStatus>(params, 'status', 'unread')
        const sortKey = `cat_${category.id}`
        if (Object.prototype.hasOwnProperty.call(localConfig.display, sortKey)) {
          req.sortOrder = getURLParam(params, 'order', localConfig.display[sortKey].order)
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

const buildTitle = (status: string | null, category?: Category) => {
  let title = ''
  if (status) {
    title = status === 'unread' ? 'to read' : 'read'
  }
  if (category) {
    title = title + ' in "' + category.title + '"'
  }
  return title
}

const computeTotalArticles = (data: GetArticlesResponse, status: string | null) => {
  let delta = 0
  if (status) {
    delta = data.articles.entries.filter((a) => a.status !== status).length
  }
  return data.articles.totalCount - delta
}

const EmptyMessage = ({ variant }: { variant: Variant }) => {
  switch (variant) {
    case 'starred':
      return 'no starred article'
    case 'history':
      return 'history is empty'
    default:
      return 'no article to read'
  }
}

type AllProps = Props & RouteComponentProps

export default (props: AllProps) => {
  const { variant, category } = props

  const { localConfiguration } = useContext(LocalConfigurationContext)
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
    console.log('re-fetching articles...')
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
          {variant === 'unread' && (
            <NewArticlesAvailable current={computeTotalArticles(d, req.status)} category={category} refresh={refresh} />
          )}
          <ArticleList
            articles={entries}
            emptyMessage={EmptyMessage({ variant })}
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
  let title = buildTitle(req.status, category)
  if (data && data.articles) {
    const totalCount = computeTotalArticles(data, req.status)
    const plural = totalCount > 1 ? ' articles ' : ' article '
    title = totalCount + plural + title
  } else title = ' '

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
      {refetching && <Loader center />}
      {render(loading, data, error)}
    </Page>
  )
}
