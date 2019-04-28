import gql from 'graphql-tag'

export const GetRules = gql`
  query {
    rules {
      id
      alias
      category_id
      rule
      priority
      created_at
      updated_at
    }
  }
`

export const GetRule = gql`
  query rule($id: ID!) {
    rule(id: $id) {
      id
      alias
      category_id
      rule
      priority
      created_at
      updated_at
    }
  }
`

export const DeleteRules = gql`
  mutation deleteRules($ids: [ID!]!) {
    deleteRules(ids: $ids)
  }
`

export const CreateOrUpdateRule = gql`
  mutation createOrUpdateRule($id: ID, $alias: String!, $category_id: Int!, $rule: String!, $priority: Int) {
    createOrUpdateRule(id: $id, alias: $alias, category_id: $category_id, rule: $rule, priority: $priority) {
      id
      alias
      category_id
      rule
      priority
      created_at
      updated_at
    }
  }
`
