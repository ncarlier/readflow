import React, { FormEvent, useCallback, useContext, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { useFormState } from 'react-use-form-state'

import { Category } from '../../categories/models'
import Button from '../../components/Button'
import FormInputField from '../../components/FormInputField'
import { getGQLError, isValidForm } from '../../helpers'
import Loader from '../../components/Loader'
import Panel from '../../components/Panel'
import { MessageContext } from '../../context/MessageContext'
import ErrorPanel from '../../error/ErrorPanel'
import useOnMountInputValidator from '../../hooks/useOnMountInputValidator'
import { AddNewArticleRequest, Article } from '../models'
import { AddNewArticle } from '../queries'

interface AddArticleFormFields {
  url: string
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
  const [formState, { url }] = useFormState<AddArticleFormFields>({ url: value })
  const onMountValidator = useOnMountInputValidator(formState.validity)
  const addArticleMutation = useMutation<AddNewArticleRequest>(AddNewArticle)

  const addArticle = async (form: AddArticleFormFields) => {
    setLoading(true)
    try {
      const categoryID = category ? category.id : undefined
      const variables = { ...form, category: categoryID }
      const res = await addArticleMutation({ variables })
      setLoading(false)
      const article = res.data.addArticle
      onSuccess(article)
      showMessage(`New article: ${article.title}`)
    } catch (err) {
      setLoading(false)
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
          <FormInputField
            label="URL"
            {...url('url')}
            error={!formState.validity.url}
            required
            autoFocus
            ref={onMountValidator.bind}
          />
        </form>
      </section>
      <footer>
        <Button title="Back to API keys" onClick={onCancel}>
          Cancel
        </Button>
        <Button title="Add new article" onClick={handleOnSubmit} primary>
          Add
        </Button>
      </footer>
    </Panel>
  )
}
