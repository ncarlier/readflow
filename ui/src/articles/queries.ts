import gql from 'graphql-tag'

export const GetArticles = gql`
  query articles(
    $limit: Int
    $sortBy: sortBy
    $sortOrder: sortOrder
    $status: status
    $starred: Boolean
    $category: Int
    $afterCursor: Int
    $query: String
  ) {
    articles(
      limit: $limit
      sortBy: $sortBy
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
        stars
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
      stars
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
      stars
      category {
        id
        title
      }
      created_at
    }
  }
`

export const UpdateArticle = gql`
  mutation updateArticle($id: ID!, $status: status, $stars: Int) {
    updateArticle(id: $id, status: $status, stars: $stars) {
      article {
        id
        status
        stars
        category {
          id
          inbox
        }
        updated_at
      }
      _inbox
      _to_read
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
  mutation markAllArticlesAsRead($status: status!, $category: ID) {
    markAllArticlesAsRead(status: $status, category: $category) {
      _inbox
      entries {
        id
        inbox
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
      stars
      category {
        id
        inbox
      }
      created_at
    }
  }
`

export const GetNbNewArticles = gql`
  query articles($category: Int) {
    articles(limit: 1, status: inbox, category: $category) {
      totalCount
    }
  }
`
