import React, { CSSProperties, FormEvent, MouseEvent, useCallback, useState } from 'react'
import { useMutation } from '@apollo/client'
import { useFormState } from 'react-use-form-state'

import { Category } from '../../categories/models'
import { useMessage } from '../../contexts'
import { Button, CategoriesOptions, ErrorPanel, FormInputField, FormSelectField, Loader, Panel } from '../../components'
import { getGQLError, isValidForm } from '../../helpers'
import { AddNewArticleRequest, AddNewArticleResponse, Article } from '../models'
import { AddNewArticle } from '../queries'
import { updateCacheAfterCreate } from '../cache'

interface AddArticleFormFields {
  url: string
  category?: number
}

interface Props {
  value?: string
  category?: Category
  style?: CSSProperties
  onSuccess: (article: Article) => void
  onCancel: (e: any) => void
}

export const AddArticleForm = ({ value, category, style, onSuccess, onCancel }: Props) => {
  const [loading, setLoading] = useState(false)
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const { showMessage } = useMessage()
  const [formState, { url, select }] = useFormState<AddArticleFormFields>({
    url: value,
    category: category ? category.id : undefined,
  })
  const [addArticleMutation] = useMutation<AddNewArticleResponse, AddNewArticleRequest>(AddNewArticle)

  const addArticle = useCallback(
    async (form: AddArticleFormFields) => {
      setLoading(true)
      try {
        const variables = { ...form }
        const res = await addArticleMutation({
          variables,
          update: updateCacheAfterCreate,
        })
        setLoading(false)
        if (res.data) {
          const article = res.data.addArticle
          showMessage(`New article: ${article.title}`)
          onSuccess(article)
        }
      } catch (err) {
        setLoading(false)
        setErrorMessage(getGQLError(err))
      }
    },
    [addArticleMutation, onSuccess, showMessage]
  )

  const handleOnSubmit = useCallback(
    (e: FormEvent | MouseEvent) => {
      e.preventDefault()
      if (!isValidForm(formState)) {
        setErrorMessage('Please fill out correctly the mandatory fields.')
        return
      }
      addArticle(formState.values)
    },
    [formState, addArticle]
  )

  return (
    <Panel>
      {loading && <Loader blur />}
      <header>
        <h1>Add new article</h1>
      </header>
      <section style={style}>
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
