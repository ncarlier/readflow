import gql from 'graphql-tag'

export const GetIncomingWebhooks = gql`
  query {
    incomingWebhooks {
      id
      alias
      token
      script
      last_usage_at
      created_at
      updated_at
    }
  }
`

export const GetIncomingWebhook = gql`
  query incomingWebhook($id: ID!) {
    incomingWebhook(id: $id) {
      id
      alias
      token
      email
      script
      last_usage_at
      created_at
      updated_at
    }
  }
`

export const DeleteIncomingWebhooks = gql`
  mutation deleteIncomingWebhooks($ids: [ID!]!) {
    deleteIncomingWebhooks(ids: $ids)
  }
`

export const CreateOrUpdateIncomingWebhook = gql`
  mutation createOrUpdateIncomingWebhook($id: ID, $alias: String!, $script: String!) {
    createOrUpdateIncomingWebhook(id: $id, alias: $alias, script: $script) {
      id
      alias
      token
      email
      script
      last_usage_at
      created_at
      updated_at
    }
  }
`
