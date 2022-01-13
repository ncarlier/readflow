import { History } from 'history'
import React, { FormEvent, useCallback, useState } from 'react'
import { useMutation } from '@apollo/client'
import { Link } from 'react-router-dom'
import { useFormState } from 'react-use-form-state'

import { Category, CreateOrUpdateCategoryResponse } from '../../categories/models'
import { CreateOrUpdateCategory } from '../../categories/queries'
import { useMessage } from '../../contexts'
import { Button, ErrorPanel, FormInputField, FormSelectField, FormTextareaField, HelpLink } from '../../components'
import { getGQLError, isValidForm } from '../../helpers'

interface EditCategoryFormFields {
  title: string
  rule: string
  notification_strategy: 'none' | 'global' | 'individual'
}

interface Props {
  category: Category
  history: History
}

const EditCategoryForm = ({ category, history }: Props) => {
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text, textarea, select }] = useFormState<EditCategoryFormFields>({
    title: category.title,
    rule: category.rule ? category.rule : '',
    notification_strategy: category.notification_strategy,
  })
  const [editCategoryMutation] = useMutation<CreateOrUpdateCategoryResponse, Category>(CreateOrUpdateCategory)
  const { showMessage } = useMessage()

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
            <HelpLink href="https://docs.readflow.app/en/read-flow/categories/#rule">View rule syntax</HelpLink>
          </FormTextareaField>
          <FormSelectField
            label="Notification strategy"
            {...select('notification_strategy')}
            error={formState.errors.notification_strategy}
            required
          >
            <option value="none">Don&apos;t send any notification</option>
            <option value="individual">Send a notification as soon as an article is received</option>
            <option value="global">Use global notification strategy</option>
          </FormSelectField>
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

export default EditCategoryForm
