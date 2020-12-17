import gql from 'graphql-tag'

export const GetOutboundServices = gql`
  query {
    outboundServices {
      id
      alias
      provider
      is_default
      created_at
      updated_at
    }
  }
`

export const GetOutboundService = gql`
  query outboundService($id: ID!) {
    outboundService(id: $id) {
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

export const DeleteOutboundServices = gql`
  mutation deleteOutboundServices($ids: [ID!]!) {
    deleteOutboundServices(ids: $ids)
  }
`

export const CreateOrUpdateOutboundService = gql`
  mutation createOrUpdateOutboundService(
    $id: ID
    $alias: String!
    $provider: provider!
    $config: String!
    $is_default: Boolean!
  ) {
    createOrUpdateOutboundService(
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
