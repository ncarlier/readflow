import React, { useCallback, useState } from 'react'

import { useFormState } from 'react-use-form-state'
import { useMutation } from 'react-apollo-hooks'

import Panel from '../../common/Panel'
import Button from '../../common/Button'
import FormInputField from '../../common/FormInputField'
import ErrorPanel from '../../error/ErrorPanel'
import { getGQLError, isValidForm } from '../../common/helpers'
import { AddNewArticle } from '../queries'
import { Article, AddNewArticleRequest } from '../models'
import { Category } from '../../categories/models'
import Loader from '../../common/Loader'

interface AddArticleFormFields {
  url: string
}

type Props = {
  value?: string
  category?: Category
  onSuccess: (article: Article) => void
  onCancel: (e: any) => void
}

type AllProps = Props

export default ({ value, category, onSuccess, onCancel }: AllProps) => {
  const [loading, setLoading] = useState(false)
  const [errorMessage, setErrorMessage] = useState<string | null>(null) 
  const [formState, { url }] = useFormState<AddArticleFormFields>({url: value})
  const addArticleMutation = useMutation<AddNewArticleRequest>(AddNewArticle)

  const addArticle = async (form: AddArticleFormFields) => {
    setLoading(true)
    try {
      const categoryID = category ? category.id : undefined
      const variables = {...form, category: categoryID}
      const res = await addArticleMutation({variables})
      setLoading(false)
      onSuccess(res.data.addArticle)
    } catch (err) {
      setLoading(false)
      setErrorMessage(getGQLError(err))
    }
  }

  const handleOnClick = useCallback(() => {
    if (!isValidForm(formState)) {
      setErrorMessage("Please fill out correctly the mandatory fields.")
      return
    }
    addArticle(formState.values)
  }, [formState])

  return (
    <Panel>
      { loading && <Loader blur /> }
      <header>
        <h1>Add new article</h1>
      </header>
      <section>
        {errorMessage != null &&
          <ErrorPanel title="Unable to add new article">
            {errorMessage}
          </ErrorPanel>
        }
        <form>
          <FormInputField label="URL"
            {...url('url')}
            error={!formState.validity.url}
            required />
        </form>
      </section>
      <footer>
        <Button title="Back to API keys" onClick={onCancel} >
          Cancel
        </Button>
        <Button
          title="Add new article"
          onClick={handleOnClick}
          primary>
          Add
        </Button>
      </footer>
    </Panel>
  )
}
