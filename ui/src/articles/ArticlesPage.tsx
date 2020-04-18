import React, { useCallback, useContext, useState } from 'react'
import { useQuery } from 'react-apollo-hooks'
import { RouteComponentProps } from 'react-router'

import { Category } from '../categories/models'
import Loader from '../components/Loader'
import Panel from '../components/Panel'
import { LocalConfiguration, LocalConfigurationContext } from '../context/LocalConfigurationContext'
import ErrorPanel from '../error/ErrorPanel'
import { getURLParam, matchResponse } from '../helpers'
import Appbar from '../layout/Appbar'
import Page from '../layout/Page'
import AddButton from './components/AddButton'
import ArticleList from './components/ArticleList'
import ArticlesPageMenu from './components/ArticlesPageMenu'
import NewArticlesAvailable from './components/NewArticlesAvailable'
import { ArticleStatus, GetArticlesRequest, GetArticlesResponse } from './models'
import { GetArticles } from './queries'

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
    sortOrder: getURLParam(params, 'sort', localConfig.sortOrders.unread),
    status: 'unread',
    starred: null,
    category: null,
    afterCursor: null
  }
  switch (variant) {
    case 'history':
      req.status = 'read'
      req.sortOrder = getURLParam(params, 'sort', localConfig.sortOrders.history)
      break
    case 'starred':
      req.status = null
      req.starred = true
      req.sortOrder = getURLParam(params, 'sort', localConfig.sortOrders.starred)
      break
    case 'unread':
      if (category && category.id) {
        req.category = category.id
        req.status = getURLParam<ArticleStatus>(params, 'status', 'unread')
        const sortKey = `cat_${category.id}`
        if (Object.prototype.hasOwnProperty.call(localConfig.sortOrders, sortKey)) {
          req.sortOrder = getURLParam(params, 'sort', localConfig.sortOrders[sortKey])
        }
      }
  }

  return req
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
    delta = data.articles.entries.filter(a => a.status !== status).length
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

  const [reloading, setReloading] = useState(false)
  const { localConfiguration } = useContext(LocalConfigurationContext)
  const [req] = useState<GetArticlesRequest>(buildArticlesRequest(variant, props, localConfiguration))
  const { data, error, loading, fetchMore, refetch } = useQuery<GetArticlesResponse>(GetArticles, {
    variables: req,
  })

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
          entries = entries.filter(a => a.status === req.status)
        }
        const articles = {
          ...fetchMoreResult.articles,
          entries: [...entries, ...fetchMoreResult.articles.entries]
        }
        return { articles }
      }
    })
  }, [data, fetchMore, req])

  const refresh = useCallback(async () => {
    console.log('re-fetching articles...')
    setReloading(true)
    const { errors } = await refetch()
    if (errors) {
      console.error(errors)
    }
    setReloading(false)
  }, [refetch])

  const render = matchResponse<GetArticlesResponse>({
    Loading: () => <Loader />,
    Error: err => (
      <Panel>
        <ErrorPanel>{err.message}</ErrorPanel>
      </Panel>
    ),
    Data: d => {
      let { entries } = d.articles
      if (req.status) {
        entries = entries.filter(a => a.status === req.status)
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
            fetchMoreArticles={fetchMoreArticles}
          />
          {variant === 'unread' && <AddButton category={category} onSuccess={refresh} />}
        </>
      )
    },
    Other: () => (
      <Panel>
        <ErrorPanel>Unable to fetch articles!</ErrorPanel>
      </Panel>
    )
  })

  // Build title
  let title = buildTitle(req.status, category)
  if (data && data.articles) {
    const totalCount = computeTotalArticles(data, req.status)
    const plural = totalCount > 1 ? ' articles ' : ' article '
    title = totalCount + plural + title
  } else title = ' '

  const $actions = <ArticlesPageMenu refresh={refresh} req={req} variant={variant} />

  return (
    <Page title={title} header={<Appbar title={title} actions={$actions} />}>
      {render(data, error, loading || reloading)}
    </Page>
  )
}
