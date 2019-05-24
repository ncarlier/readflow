import React, { useCallback, useContext, useState } from 'react'
import { useMutation, useQuery } from 'react-apollo-hooks'
import { RouteComponentProps } from 'react-router'

import { Category } from '../categories/models'
import { getGQLError, getURLParam, matchResponse } from '../common/helpers'
import Loader from '../common/Loader'
import Page from '../common/Page'
import Panel from '../common/Panel'
import { MessageContext } from '../context/MessageContext'
import ErrorPanel from '../error/ErrorPanel'
import AddButton from './components/AddButton'
import ArticleList from './components/ArticleList'
import { DisplayMode } from './components/ArticlesDisplayMode'
import ArticlesPageMenu from './components/ArticlesPageMenu'
import NewArticlesAvailable from './components/NewArticlesAvailable'
import { GetArticlesRequest, GetArticlesResponse } from './models'
import { GetArticles, MarkAllArticlesAsRead } from './queries'

interface Props {
  category?: Category
}

const buildArticlesRequest = (mode: DisplayMode, props: AllProps) => {
  const { category, location } = props
  const params = new URLSearchParams(location.search)

  const req: GetArticlesRequest = {
    limit: getURLParam<number>(params, 'limit', 10),
    sortOrder: getURLParam<string>(params, 'sort', 'asc'),
    status: 'unread'
  }
  switch (mode) {
    case DisplayMode.history:
      req.status = 'read'
      req.sortOrder = getURLParam<string>(params, 'sort', 'desc')
      break
    case DisplayMode.category:
      if (category && category.id) {
        req.category = category.id
        req.status = getURLParam<string>(params, 'status', 'unread')
      }
  }

  return req
}

const buildTitle = (mode: DisplayMode, status: string, category?: Category) => {
  let title = status === 'unread' ? 'to read' : 'read'
  if (category) {
    title = title + ' in "' + category.title + '"'
  }
  return title
}

const computeTotalArticles = (data: GetArticlesResponse, status: string) => {
  const delta = data.articles.entries.filter(a => a.status !== status).length
  return data.articles.totalCount - delta
}

const EmptyMessage = ({ mode }: { mode: DisplayMode }) => {
  switch (mode) {
    case DisplayMode.category:
      return 'no article to read in this category'
    case DisplayMode.history:
      return 'history is empty'
    default:
      return 'no article to read'
  }
}

type AllProps = Props & RouteComponentProps

export default (props: AllProps) => {
  const { category, match } = props

  const { showErrorMessage } = useContext(MessageContext)
  const [reloading, setReloading] = useState(false)

  // Get display mode
  let mode = match.url.startsWith('/history') ? DisplayMode.history : DisplayMode.unread
  mode = category ? DisplayMode.category : mode

  // Build GQL request
  const req = buildArticlesRequest(mode, props)

  const { data, error, loading, fetchMore, refetch } = useQuery<GetArticlesResponse>(GetArticles, {
    variables: req
  })

  const fetchMoreArticles = useCallback(async () => {
    if (!data || !data.articles.hasNext) {
      return
    }
    console.log('fetching more articles...')
    await fetchMore({
      variables: { ...req, afterCursor: data.articles.endCursor, category: category ? category.id : null },
      updateQuery: (prev, { fetchMoreResult }) => {
        if (!fetchMoreResult) return prev
        const nbFetchedArticles = fetchMoreResult.articles.entries.length
        console.log(nbFetchedArticles + ' article(s) fetched')
        const articles = {
          ...fetchMoreResult.articles,
          entries: [...prev.articles.entries, ...fetchMoreResult.articles.entries]
        }
        return { articles }
      }
    })
  }, [data])

  const refresh = useCallback(async () => {
    console.log('re-fetching articles...')
    setReloading(true)
    const { errors } = await refetch()
    if (errors) {
      console.error(errors)
    }
    setReloading(false)
  }, [refetch])

  const markAllArticlesAsReadMutation = useMutation<{ category?: number }>(MarkAllArticlesAsRead)

  const markAllArticlesAsRead = async () => {
    try {
      await markAllArticlesAsReadMutation({
        variables: { category: category ? category.id : null }
      })
      await refresh()
    } catch (err) {
      showErrorMessage(getGQLError(err))
    }
  }

  const markAllAsRead = useCallback(() => {
    markAllArticlesAsRead()
  }, [category])

  const render = matchResponse<GetArticlesResponse>({
    Loading: () => <Loader />,
    Error: err => (
      <Panel>
        <ErrorPanel>{err.message}</ErrorPanel>
      </Panel>
    ),
    Data: d => (
      <>
        {mode === DisplayMode.unread && (
          <NewArticlesAvailable current={computeTotalArticles(d, req.status)} category={category} refresh={refresh} />
        )}
        <ArticleList
          articles={d.articles.entries}
          emptyMessage={EmptyMessage({ mode })}
          filter={a => a.status === req.status}
          hasMore={d.articles.hasNext}
          refetch={refetch}
          fetchMoreArticles={fetchMoreArticles}
        />
        {mode !== DisplayMode.history && <AddButton category={category} onSuccess={refresh} />}
      </>
    ),
    Other: () => (
      <Panel>
        <ErrorPanel>Unable to fetch articles!</ErrorPanel>
      </Panel>
    )
  })

  // Build title
  let title = buildTitle(mode, req.status, category)
  if (data && data.articles) {
    const totalCount = computeTotalArticles(data, req.status)
    const plural = totalCount > 1 ? ' articles ' : ' article '
    title = totalCount + plural + title
  } else title = ' '

  return (
    <Page
      title={title}
      actions={
        <ArticlesPageMenu
          refresh={refresh}
          markAllAsRead={req.status == 'unread' ? markAllAsRead : undefined}
          mode={mode}
        />
      }
    >
      {render(data, error, loading || reloading)}
    </Page>
  )
}
