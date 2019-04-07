import gql from "graphql-tag"

export const GetCategories = gql`
  query {
    categories {
      id
      title
      created_at
      updated_at
    }
  }
`

export const GetCategory = gql`
  query category($id: ID!) {
    category(id: $id) {
      id
      title
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

export const CreateOrUpdateCategory = gql `
  mutation createOrUpdateCategory($id: ID, $title: String!) {
    createOrUpdateCategory(id: $id, title: $title) {
      id
      title
      created_at
      updated_at
    }
  }
`
