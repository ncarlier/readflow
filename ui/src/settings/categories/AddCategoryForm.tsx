import React, { FormEvent, useCallback, useContext, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { RouteComponentProps } from 'react-router'
import { useFormState } from 'react-use-form-state'

import { updateCacheAfterCreate } from '../../categories/cache'
import { Category } from '../../categories/models'
import { CreateOrUpdateCategory } from '../../categories/queries'
import Button from '../../components/Button'
import FormInputField from '../../components/FormInputField'
import { getGQLError, isValidForm } from '../../helpers'
import Panel from '../../components/Panel'
import { MessageContext } from '../../context/MessageContext'
import ErrorPanel from '../../error/ErrorPanel'
import { usePageTitle } from '../../hooks'
import useOnMountInputValidator from '../../hooks/useOnMountInputValidator'

interface AddCategoryFormFields {
  title: string
}

type AllProps = RouteComponentProps<{}>

export default ({ history }: AllProps) => {
  usePageTitle('Settings - Add new category')

  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text }] = useFormState<AddCategoryFormFields>()
  const onMountValidator = useOnMountInputValidator(formState.validity)
  const addCategoryMutation = useMutation<Category>(CreateOrUpdateCategory)
  const { showMessage } = useContext(MessageContext)

  const addNewCategory = async (category: Category) => {
    try {
      await addCategoryMutation({
        variables: category,
        update: updateCacheAfterCreate
      })
      showMessage(`New category: ${category.title}`)
      // console.log('New category', res)
      history.goBack()
    } catch (err) {
      setErrorMessage(getGQLError(err))
    }
  }

  const handleOnSubmit = useCallback(
    (e: FormEvent<HTMLFormElement>) => {
      e.preventDefault()
      if (!isValidForm(formState, onMountValidator)) {
        setErrorMessage('Please fill out correctly the mandatory fields.')
        return
      }
      addNewCategory(formState.values)
    },
    [formState]
  )

  return (
    <Panel>
      <header>
        <h1>Add new category</h1>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to add new category">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField
            label="Title"
            {...text('title')}
            error={!formState.validity.title}
            required
            autoFocus
            ref={onMountValidator.bind}
          />
        </form>
      </section>
      <footer>
        <Button title="Back to categories" to="/settings/categories">
          Cancel
        </Button>
        <Button title="Add category" onClick={handleOnSubmit} primary>
          Add
        </Button>
      </footer>
    </Panel>
  )
}
