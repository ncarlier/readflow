import React, { FormEvent, MouseEvent, useCallback, useState } from 'react'
import { useMutation } from '@apollo/client'
import { useFormState } from 'react-use-form-state'

import { useMessage } from '../../contexts'
import { Button, CategoriesOptions, ErrorPanel, FormCheckboxField, FormInputField, FormSelectField, FormTextareaField, Loader, Panel } from '../../components'
import { getGQLError, isValidForm } from '../../helpers'
import { Article, UpdateArticleRequest, UpdateArticleResponse } from '../models'
import { UpdateFullArticle } from '../queries'
import { updateCacheAfterUpdate } from '../cache'

interface EditArticleFormFields {
  title: string
  text: string
  category_id?: number
  refresh: boolean
}

interface Props {
  article: Article
  onSuccess: (article: Article) => void
  onCancel: (e: any) => void
}

export const EditArticleForm = ({ article, onSuccess, onCancel }: Props) => {
  const [loading, setLoading] = useState(false)
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const { showMessage } = useMessage()
  const [formState, { text, textarea, select, checkbox }] = useFormState<EditArticleFormFields>({
    title: article.title,
    text: article.text,
    category_id: article.category?.id,
    refresh: false,
  })
  const [editArticleMutation] = useMutation<UpdateArticleResponse, UpdateArticleRequest>(UpdateFullArticle)

  const editArticle = useCallback(
    async (form: EditArticleFormFields) => {
      setLoading(true)
      try {
        const variables = { id: article.id, ...form }
        const res = await editArticleMutation({
          variables,
          update: updateCacheAfterUpdate,
        })
        if (res.data) {
          const {article} = res.data.updateArticle
          showMessage(`Article "${article.title}" (#${article.id}) edited`)
          onSuccess(article)
        }
      } catch (err) {
        setErrorMessage(getGQLError(err))
      } finally {
        setLoading(false)
      }
    },
    [article, editArticleMutation, onSuccess, showMessage]
  )

  const handleOnSubmit = useCallback(
    (e: FormEvent | MouseEvent) => {
      e.preventDefault()
      if (!isValidForm(formState)) {
        setErrorMessage('Please fill out correctly the mandatory fields.')
        return
      }
      editArticle(formState.values)
    },
    [formState, editArticle]
  )

  return (
    <Panel>
      {loading && <Loader blur />}
      <header>
        <h1>Edit article</h1>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to edit article">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField label="Title" {...text('title')} error={formState.errors.title} pattern=".*\S+.*" maxLength={256} autoFocus />
          <FormTextareaField label="Text" {...textarea('text')} error={formState.errors.text} maxLength={512} />
          <FormSelectField label="Category" {...select('category_id')} error={formState.errors.category_id}>
            <option>Optional category</option>
            <CategoriesOptions />
          </FormSelectField>
          <FormCheckboxField label="Refresh content" {...checkbox('refresh')} />
        </form>
      </section>
      <footer>
        <Button title="Cancel" onClick={onCancel}>
          Cancel
        </Button>
        <Button title="Update article" onClick={handleOnSubmit} variant="primary">
          Update
        </Button>
      </footer>
    </Panel>
  )
}
