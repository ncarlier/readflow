import React, { FormEvent, MouseEvent, useCallback, useContext, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { useFormState } from 'react-use-form-state'

import { Category } from '../../categories/models'
import Button from '../../components/Button'
import CategoriesOptions from '../../components/CategoriesOptions'
import FormInputField from '../../components/FormInputField'
import FormSelectField from '../../components/FormSelectField'
import Loader from '../../components/Loader'
import Panel from '../../components/Panel'
import { MessageContext } from '../../context/MessageContext'
import ErrorPanel from '../../error/ErrorPanel'
import { getGQLError, isValidForm } from '../../helpers'
import { AddNewArticleRequest, AddNewArticleResponse, Article } from '../models'
import { AddNewArticle } from '../queries'

interface AddArticleFormFields {
  url: string
  category?: number
}

interface Props {
  value?: string
  category?: Category
  onSuccess: (article: Article) => void
  onCancel: (e: any) => void
}

export default ({ value, category, onSuccess, onCancel }: Props) => {
  const [loading, setLoading] = useState(false)
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const { showMessage } = useContext(MessageContext)
  const [formState, { url, select }] = useFormState<AddArticleFormFields>({
    url: value,
    category: category ? category.id : undefined
  })
  const [addArticleMutation] = useMutation<AddNewArticleResponse, AddNewArticleRequest>(AddNewArticle)

  const addArticle = async (form: AddArticleFormFields) => {
    setLoading(true)
    try {
      const variables = { ...form }
      const res = await addArticleMutation({ variables })
      setLoading(false)
      if (res.data) {
        const article = res.data.addArticle
        onSuccess(article)
        showMessage(`New article: ${article.title}`)
      }
    } catch (err) {
      setLoading(false)
      setErrorMessage(getGQLError(err))
    }
  }

  const handleOnSubmit = useCallback(
    (e: FormEvent | MouseEvent) => {
      e.preventDefault()
      if (!isValidForm(formState)) {
        setErrorMessage('Please fill out correctly the mandatory fields.')
        return
      }
      addArticle(formState.values)
    },
    [formState]
  )

  return (
    <Panel>
      {loading && <Loader blur />}
      <header>
        <h1>Add new article</h1>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to add new article">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField label="URL" {...url('url')} error={formState.errors.url} required autoFocus />
          <FormSelectField label="Category" {...select('category')} error={formState.errors.category}>
            <option>Optional category</option>
            <CategoriesOptions />
          </FormSelectField>
        </form>
      </section>
      <footer>
        <Button title="Back to API keys" onClick={onCancel}>
          Cancel
        </Button>
        <Button title="Add new article" onClick={handleOnSubmit} variant="primary">
          Add
        </Button>
      </footer>
    </Panel>
  )
}
