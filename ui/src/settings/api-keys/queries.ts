import gql from "graphql-tag"

export const GetApiKeys = gql`
  query {
    apiKeys {
      id
      alias
      token
      last_usage_at
      created_at
      updated_at
    }
  }
`

export const GetApiKey = gql`
  query apiKey($id: ID!) {
    apiKey(id: $id) {
      id
      alias
      token
      last_usage_at
      created_at
      updated_at
    }
  }
`

export const DeleteApiKeys = gql`
  mutation deleteAPIKeys($ids: [ID!]!) {
    deleteAPIKeys(ids: $ids)
  }
`

export const CreateOrUpdateApiKey = gql `
  mutation createOrUpdateAPIKey($id: ID, $alias: String!) {
    createOrUpdateAPIKey(id: $id, alias: $alias) {
      id
      alias
      token
      last_usage_at
      created_at
      updated_at
    }
  }
`
