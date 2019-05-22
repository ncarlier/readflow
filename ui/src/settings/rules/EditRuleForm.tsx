import { History } from 'history'
/* eslint-disable @typescript-eslint/camelcase */
import React, { FormEvent, useCallback, useContext, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { useFormState } from 'react-use-form-state'

import Button from '../../common/Button'
import CategoriesOptions from '../../common/CategoriesOptions'
import FormInputField from '../../common/FormInputField'
import FormSelectField from '../../common/FormSelectField'
import FormTextareaField from '../../common/FormTextareaField'
import { getGQLError, isValidForm } from '../../common/helpers'
import HelpLink from '../../common/HelpLink'
import { MessageContext } from '../../context/MessageContext'
import ErrorPanel from '../../error/ErrorPanel'
import useOnMountInputValidator from '../../hooks/useOnMountInputValidator'
import { updateCacheAfterUpdate } from './cache'
import { Rule } from './models'
import PriorityOptions from './PriorityOptions'
import { CreateOrUpdateRule } from './queries'

interface EditRuleFormFields {
  alias: string
  rule: string
  priority: number
  category_id: number
}

interface Props {
  data: Rule
  history: History
}

export default ({ data, history }: Props) => {
  const { showMessage } = useContext(MessageContext)
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text, select, textarea }] = useFormState<EditRuleFormFields>({
    alias: data.alias,
    rule: data.rule,
    priority: data.priority,
    category_id: data.category_id
  })
  const onMountValidator = useOnMountInputValidator(formState.validity)
  const editRuleMutation = useMutation<Rule>(CreateOrUpdateRule)

  const editRule = async (rule: Rule) => {
    try {
      await editRuleMutation({
        variables: rule,
        update: updateCacheAfterUpdate
      })
      showMessage(`Rule edited: ${rule.alias}`)
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
      const { alias, rule, priority, category_id } = formState.values
      editRule({ id: data.id, alias, rule, priority: parseInt(priority), category_id: parseInt(category_id) })
    },
    [formState]
  )

  return (
    <>
      <header>
        <h1>Edit rule #{data.id}</h1>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to edit rule">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField
            label="Alias"
            {...text('alias')}
            error={!formState.validity.alias}
            required
            autoFocus
            ref={onMountValidator.bind}
          />
          <FormTextareaField
            label="Rule"
            {...textarea('rule')}
            error={!formState.validity.rule}
            required
            ref={onMountValidator.bind}
          >
            <HelpLink href="https://about.readflow.app/docs/en/read-flow/organize/rules/#syntax">
              View rule syntax
            </HelpLink>
          </FormTextareaField>
          <FormSelectField
            label="Priority"
            {...select('priority')}
            error={!formState.validity.priority}
            required
            ref={onMountValidator.bind}
          >
            <PriorityOptions />
          </FormSelectField>
          <FormSelectField
            label="Category"
            {...select('category_id')}
            error={!formState.validity.category_id}
            required
            ref={onMountValidator.bind}
          >
            <CategoriesOptions />
          </FormSelectField>
        </form>
      </section>
      <footer>
        <Button title="Back to rules" to="/settings/rules">
          Cancel
        </Button>
        <Button title="Edit rule" onClick={handleOnSubmit} primary>
          Update
        </Button>
      </footer>
    </>
  )
}
