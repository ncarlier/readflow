import gql from 'graphql-tag'

export const GetOutgoingWebhooks = gql`
  query {
    outgoingWebhooks {
      id
      alias
      provider
      is_default
      created_at
      updated_at
    }
  }
`

export const GetOutgoingWebhook = gql`
  query outgoingWebhook($id: ID!) {
    outgoingWebhook(id: $id) {
      id
      alias
      provider
      config
      is_default
      created_at
      updated_at
    }
  }
`

export const DeleteOutgoingWebhooks = gql`
  mutation deleteOutgoingWebhooks($ids: [ID!]!) {
    deleteOutgoingWebhooks(ids: $ids)
  }
`

export const CreateOrUpdateOutgoingWebhook = gql`
  mutation createOrUpdateOutgoingWebhook(
    $id: ID
    $alias: String!
    $provider: outgoingWebhookProvider!
    $config: String!
    $is_default: Boolean!
  ) {
    createOrUpdateOutgoingWebhook(
      id: $id
      alias: $alias
      provider: $provider
      config: $config
      is_default: $is_default
    ) {
      id
      alias
      provider
      config
      is_default
      created_at
      updated_at
    }
  }
`
