import React, { FormEvent, useCallback, useContext, useState } from 'react'
import { useMutation } from '@apollo/client'
import { RouteComponentProps } from 'react-router'
import { Link } from 'react-router-dom'
import { useFormState } from 'react-use-form-state'

import { updateCacheAfterCreate } from '../../categories/cache'
import { Category, CreateOrUpdateCategoryResponse } from '../../categories/models'
import { CreateOrUpdateCategory } from '../../categories/queries'
import Button from '../../components/Button'
import FormInputField from '../../components/FormInputField'
import FormTextareaField from '../../components/FormTextareaField'
import HelpLink from '../../components/HelpLink'
import Panel from '../../components/Panel'
import { MessageContext } from '../../context/MessageContext'
import ErrorPanel from '../../error/ErrorPanel'
import { getGQLError, isValidForm } from '../../helpers'
import { usePageTitle } from '../../hooks'

interface AddCategoryFormFields {
  title: string
  rule: string
}

type AllProps = RouteComponentProps<{}>

export default ({ history }: AllProps) => {
  usePageTitle('Settings - Add new category')

  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text, textarea }] = useFormState<AddCategoryFormFields>({
    title: '',
    rule: '',
  })
  const [addCategoryMutation] = useMutation<CreateOrUpdateCategoryResponse, Category>(CreateOrUpdateCategory)
  const { showMessage } = useContext(MessageContext)

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
