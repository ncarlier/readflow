import gql from 'graphql-tag'

export const GetCategories = gql`
  query {
    categories {
      _inbox
      _to_read
      _starred
      entries {
        id
        title
        rule
        notification_strategy
        inbox
        created_at
        updated_at
      }
    }
  }
`

export const GetCategory = gql`
  query category($id: ID!) {
    category(id: $id) {
      id
      title
      rule
      notification_strategy
      created_at
      updated_at
    }
  }
`

export const DeleteCategories = gql`
  mutation deleteCategories($ids: [ID!]!) {
    deleteCategories(ids: $ids)
  }
`

export const CreateOrUpdateCategory = gql`
  mutation createOrUpdateCategory(
    $id: ID
    $title: String
    $rule: String
    $notification_strategy: notificationStrategy
  ) {
    createOrUpdateCategory(id: $id, title: $title, rule: $rule, notification_strategy: $notification_strategy) {
      id
      title
      rule
      notification_strategy
      created_at
      updated_at
    }
  }
`
