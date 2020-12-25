import { History } from 'history'
import React, { FormEvent, useCallback, useContext, useState } from 'react'
import { useMutation } from '@apollo/client'
import { Link } from 'react-router-dom'
import { useFormState } from 'react-use-form-state'

import { Category, CreateOrUpdateCategoryResponse } from '../../categories/models'
import { CreateOrUpdateCategory } from '../../categories/queries'
import Button from '../../components/Button'
import FormInputField from '../../components/FormInputField'
import FormTextareaField from '../../components/FormTextareaField'
import HelpLink from '../../components/HelpLink'
import { MessageContext } from '../../context/MessageContext'
import ErrorPanel from '../../error/ErrorPanel'
import { getGQLError, isValidForm } from '../../helpers'

interface EditCategoryFormFields {
  title: string
  rule: string
}

interface Props {
  category: Category
  history: History
}

export default ({ category, history }: Props) => {
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text, textarea }] = useFormState<EditCategoryFormFields>({
    title: category.title,
    rule: category.rule ? category.rule : '',
  })
  const [editCategoryMutation] = useMutation<CreateOrUpdateCategoryResponse, Category>(CreateOrUpdateCategory)
  const { showMessage } = useContext(MessageContext)

  const editCategory = useCallback(
    async (category: Category) => {
      try {
        await editCategoryMutation({
          variables: category,
          // update: updateCacheAfterUpdate
        })
        showMessage(`Category edited: ${category.title}`)
        history.goBack()
      } catch (err) {
        setErrorMessage(getGQLError(err))
      }
    },
    [editCategoryMutation, showMessage, history]
  )

  const handleOnSubmit = useCallback(
    (e: FormEvent | MouseEvent) => {
      e.preventDefault()
      if (!isValidForm(formState)) {
        setErrorMessage('Please fill out correctly the mandatory fields.')
        return
      }
      editCategory({ id: category.id, ...formState.values })
    },
    [category, formState, editCategory]
  )

  return (
    <>
      <header>
        <h1>Edit category #{category.id}</h1>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to edit category">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField label="Title" {...text('title')} error={formState.errors.title} required autoFocus />
          <FormTextareaField label="Rule" {...textarea('rule')}>
            <HelpLink href="https://about.readflow.app/docs/en/read-flow/categories/#rule">View rule syntax</HelpLink>
          </FormTextareaField>
        </form>
      </section>
      <footer>
        <Button title="Back to categories" as={Link} to="/settings/categories">
          Cancel
        </Button>
        <Button title="Edit category" onClick={handleOnSubmit} variant="primary">
          Update
        </Button>
      </footer>
    </>
  )
}
