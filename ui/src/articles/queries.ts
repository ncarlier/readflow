import gql from 'graphql-tag'

export const GetArticles = gql`
  query articles($limit: Int!, $sortOrder: sortOrder!, $status: status!, $category: Int, $afterCursor: Int) {
    articles(limit: $limit, sortOrder: $sortOrder, status: $status, category: $category, afterCursor: $afterCursor) {
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
      category {
        id
        title
      }
      created_at
    }
  }
`

export const UpdateArticleStatus = gql`
  mutation updateArticleStatus($id: ID!, $status: status!) {
    updateArticleStatus(id: $id, status: $status) {
      article {
        id
        status
        category {
          id
          unread
        }
        updated_at
      }
      _all {
        id
        unread
      }
    }
  }
`

export const ArchiveArticle = gql`
  mutation archiveArticle($id: ID!, $archiver: String!) {
    archiveArticle(id: $id, archiver: $archiver)
  }
`

export const MarkAllArticlesAsRead = gql`
  mutation markAllArticlesAsRead($category: ID) {
    markAllArticlesAsRead(category: $category) {
      id
      unread
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
      category {
        id
        unread
      }
      created_at
    }
  }
`
