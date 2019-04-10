import React, { useCallback, useState } from 'react'

import { useFormState } from 'react-use-form-state'
import { useMutation } from 'react-apollo-hooks'

import Button from '../../common/Button'
import FormInputField from '../../common/FormInputField'
import { CreateOrUpdateCategory } from '../../categories/queries'
import { Category } from '../../categories/models'
import ErrorPanel from '../../error/ErrorPanel'
import { getGQLError, isValidForm } from '../../common/helpers'
import { History } from 'history'
import { updateCacheAfterUpdate } from '../../categories/cache'
import { IMessageDispatchProps, connectMessageDispatch } from '../../containers/MessageContainer'

interface EditCategoryFormFields {
  title: string
}

type Props = {
  category: Category
  history: History
}

type AllProps = Props & IMessageDispatchProps

export const EditCategoryForm = ({ category, history, showMessage }: AllProps) => {
  const [errorMessage, setErrorMessage] = useState<string | null>(null) 
  const [formState, { text }] = useFormState<EditCategoryFormFields>({title: category.title})
  const editCategoryMutation = useMutation<Category>(CreateOrUpdateCategory)

  const editCategory = async (category: Category) => {
    try{
      const res = await editCategoryMutation({
        variables: category,
        update: updateCacheAfterUpdate,
      })
      showMessage(`Category edited: ${category.title}`)
      history.goBack()
    } catch (err) {
      setErrorMessage(getGQLError(err))
    }
  }

  const handleOnSubmit = useCallback(() => {
    if (!isValidForm(formState)) {
      setErrorMessage("Please fill out correctly the mandatory fields.")
      return
    }
    editCategory({id: category.id, ...formState.values})
  }, [formState])

  return (
    <>
      <header>
        <h1>Edit category #{category.id}</h1>
      </header>
      <section>
        {errorMessage != null &&
          <ErrorPanel title="Unable to edit category">
            {errorMessage}
          </ErrorPanel>
        }
        <form onSubmit={handleOnSubmit}>
          <FormInputField label="Title"
            {...text('title')}
            error={!formState.validity.label}
            required />
        </form>
      </section>
      <footer>
        <Button title="Back to categories" to="/settings/categories">
          Cancel
        </Button>
        <Button
          title="Edit category"
          onClick={handleOnSubmit}
          primary>
          Update
        </Button>
      </footer>
    </>
  )
}

export default connectMessageDispatch(EditCategoryForm)
