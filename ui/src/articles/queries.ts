import gql from 'graphql-tag'

export const GetArticles = gql`
  query articles(
    $limit: Int
    $sortOrder: sortOrder
    $status: status
    $starred: Boolean
    $category: Int
    $afterCursor: Int
    $query: String
  ) {
    articles(
      limit: $limit
      sortOrder: $sortOrder
      status: $status
      starred: $starred
      category: $category
      afterCursor: $afterCursor
      query: $query
    ) {
      totalCount
      endCursor
      hasNext
      entries {
        id
        title
        text
        url
        image
        status
        starred
        category {
          id
          title
        }
        created_at
      }
    }
  }
`

export const GetArticle = gql`
  query article($id: ID!) {
    article(id: $id) {
      id
      title
      text
      html
      url
      status
      starred
      category {
        id
        title
      }
      created_at
    }
  }
`

export const GetFullArticle = gql`
  query article($id: ID!) {
    article(id: $id) {
      id
      title
      text
      html
      url
      image
      status
      starred
      category {
        id
        title
      }
      created_at
    }
  }
`

export const UpdateArticle = gql`
  mutation updateArticle($id: ID!, $status: status, $starred: Boolean) {
    updateArticle(id: $id, status: $status, starred: $starred) {
      article {
        id
        status
        starred
        category {
          id
          unread
        }
        updated_at
      }
      _all
      _starred
    }
  }
`

export const SendArticleToOutgoingWebhook = gql`
  mutation sendArticleToOutgoingWebhook($id: ID!, $alias: String!) {
    sendArticleToOutgoingWebhook(id: $id, alias: $alias)
  }
`

export const MarkAllArticlesAsRead = gql`
  mutation markAllArticlesAsRead($category: ID) {
    markAllArticlesAsRead(category: $category) {
      _all
      entries {
        id
        unread
      }
    }
  }
`
export const AddNewArticle = gql`
  mutation addArticle($url: String!, $category: ID) {
    addArticle(url: $url, category: $category) {
      id
      title
      text
      html
      url
      image
      status
      starred
      category {
        id
        unread
      }
      created_at
    }
  }
`

export const GetNbNewArticles = gql`
  query articles($category: Int) {
    articles(limit: 1, status: unread, category: $category) {
      totalCount
    }
  }
`
