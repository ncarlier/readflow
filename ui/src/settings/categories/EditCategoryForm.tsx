import { History } from 'history'
import React, { FormEvent, useCallback, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { useFormState } from 'react-use-form-state'

import { updateCacheAfterUpdate } from '../../categories/cache'
import { Category } from '../../categories/models'
import { CreateOrUpdateCategory } from '../../categories/queries'
import Button from '../../common/Button'
import FormInputField from '../../common/FormInputField'
import { getGQLError, isValidForm } from '../../common/helpers'
import { connectMessageDispatch, IMessageDispatchProps } from '../../containers/MessageContainer'
import ErrorPanel from '../../error/ErrorPanel'
import useOnMountInputValidator from '../../hooks/useOnMountInputValidator'

interface EditCategoryFormFields {
  title: string
}

interface Props {
  category: Category
  history: History
}

type AllProps = Props & IMessageDispatchProps

export const EditCategoryForm = ({ category, history, showMessage }: AllProps) => {
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text }] = useFormState<EditCategoryFormFields>({ title: category.title })
  const onMountValidator = useOnMountInputValidator(formState.validity)
  const editCategoryMutation = useMutation<Category>(CreateOrUpdateCategory)

  const editCategory = async (category: Category) => {
    try {
      await editCategoryMutation({
        variables: category,
        update: updateCacheAfterUpdate
      })
      showMessage(`Category edited: ${category.title}`)
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
      editCategory({ id: category.id, ...formState.values })
    },
    [formState]
  )

  return (
    <>
      <header>
        <h1>Edit category #{category.id}</h1>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to edit category">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField
            label="Title"
            {...text('title')}
            error={!formState.validity.title}
            required
            ref={onMountValidator.bind}
          />
        </form>
      </section>
      <footer>
        <Button title="Back to categories" to="/settings/categories">
          Cancel
        </Button>
        <Button title="Edit category" onClick={handleOnSubmit} primary>
          Update
        </Button>
      </footer>
    </>
  )
}

export default connectMessageDispatch(EditCategoryForm)
