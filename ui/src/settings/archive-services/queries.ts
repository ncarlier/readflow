import gql from 'graphql-tag'

export const GetArchiveServices = gql`
  query {
    archivers {
      id
      alias
      provider
      is_default
      created_at
      updated_at
    }
  }
`

export const GetArchiveService = gql`
  query archiver($id: ID!) {
    archiver(id: $id) {
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

export const DeleteArchiveServices = gql`
  mutation deleteArchivers($ids: [ID!]!) {
    deleteArchivers(ids: $ids)
  }
`

export const CreateOrUpdateArchiveService = gql`
  mutation createOrUpdateArchiver(
    $id: ID
    $alias: String!
    $provider: provider!
    $config: String!
    $is_default: Boolean!
  ) {
    createOrUpdateArchiver(id: $id, alias: $alias, provider: $provider, config: $config, is_default: $is_default) {
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
