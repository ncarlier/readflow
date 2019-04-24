import gql from "graphql-tag"

export const GetDevice = gql`
  query device($id: ID!) {
    device(id: $id) {
      id
      key
      created_at
    }
  }
`

export const DeletePushSubscription = gql`
  mutation deletePushSubscription($id: ID!) {
    deletePushSubscription(id: $id) {
      id
    }
  }
`

export const CreatePushSubscription = gql `
  mutation createPushSubscription($sub: String!) {
    createPushSubscription(sub: $sub) {
      id
      key
      created_at
    }
  }
`
