import React, { FormEvent, useCallback, useState } from 'react'
import { useMutation } from '@apollo/client'
import { RouteComponentProps } from 'react-router'
import { Link } from 'react-router-dom'
import { useFormState } from 'react-use-form-state'

import { updateCacheAfterCreate } from '../../categories/cache'
import { Category, CreateOrUpdateCategoryResponse } from '../../categories/models'
import { CreateOrUpdateCategory } from '../../categories/queries'
import { useMessage } from '../../contexts'
import {
  Button,
  ErrorPanel,
  FormInputField,
  FormSelectField,
  FormTextareaField,
  HelpLink,
  Panel,
} from '../../components'
import { getGQLError, isValidForm } from '../../helpers'
import { usePageTitle } from '../../hooks'

interface AddCategoryFormFields {
  title: string
  rule: string
  notification_strategy: 'none' | 'global' | 'individual'
}

const AddCategoryForm = ({ history }: RouteComponentProps) => {
  usePageTitle('Settings - Add new category')

  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text, textarea, select }] = useFormState<AddCategoryFormFields>({
    title: '',
    rule: '',
    notification_strategy: 'none',
  })
  const [addCategoryMutation] = useMutation<CreateOrUpdateCategoryResponse, Category>(CreateOrUpdateCategory)
  const { showMessage } = useMessage()

  const addNewCategory = useCallback(
    async (category: Category) => {
      try {
        await addCategoryMutation({
          variables: category,
          update: updateCacheAfterCreate,
        })
        showMessage(`New category: ${category.title}`)
        // console.log('New category', res)
        history.goBack()
      } catch (err) {
        setErrorMessage(getGQLError(err))
      }
    },
    [addCategoryMutation, showMessage, history]
  )

  const handleOnSubmit = useCallback(
    (e: FormEvent | MouseEvent) => {
      e.preventDefault()
      if (!isValidForm(formState)) {
        setErrorMessage('Please fill out correctly the mandatory fields.')
        return
      }
      addNewCategory(formState.values)
    },
    [formState, addNewCategory]
  )

  return (
    <Panel>
      <header>
        <h1>Add new category</h1>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to add new category">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField label="Title" {...text('title')} error={formState.errors.title} required autoFocus />
          <FormTextareaField label="Rule" {...textarea('rule')} error={formState.errors.rule}>
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
        <Button title="Add category" onClick={handleOnSubmit} variant="primary">
          Add
        </Button>
      </footer>
    </Panel>
  )
}

export default AddCategoryForm
