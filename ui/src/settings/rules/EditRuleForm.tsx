import { History } from 'history'
/* eslint-disable @typescript-eslint/camelcase */
import React, { useCallback, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { useFormState } from 'react-use-form-state'

import Button from '../../common/Button'
import CategoriesOptions from '../../common/CategoriesOptions'
import FormInputField from '../../common/FormInputField'
import FormSelectField from '../../common/FormSelectField'
import FormTextareaField from '../../common/FormTextareaField'
import { getGQLError, isValidForm } from '../../common/helpers'
import HelpLink from '../../common/HelpLink'
import { connectMessageDispatch, IMessageDispatchProps } from '../../containers/MessageContainer'
import ErrorPanel from '../../error/ErrorPanel'
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

type AllProps = Props & IMessageDispatchProps

export const EditRuleForm = ({ data, history, showMessage }: AllProps) => {
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text, select, textarea }] = useFormState<EditRuleFormFields>({
    alias: data.alias,
    rule: data.rule,
    priority: data.priority,
    category_id: data.category_id
  })
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

  const handleOnSubmit = useCallback(() => {
    if (!isValidForm(formState)) {
      setErrorMessage('Please fill out correctly the mandatory fields.')
      return
    }
    const { alias, rule, priority, category_id } = formState.values
    editRule({ id: data.id, alias, rule, priority: parseInt(priority), category_id: parseInt(category_id) })
  }, [formState])

  return (
    <>
      <header>
        <h1>Edit rule #{data.id}</h1>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to edit rule">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField label="Alias" {...text('alias')} error={!formState.validity.alias} required />
          <FormTextareaField label="Rule" {...textarea('rule')} error={!formState.validity.rule} required>
            <HelpLink href="https://github.com/antonmedv/expr/wiki/The-Expression-Syntax">View rule syntax</HelpLink>
          </FormTextareaField>
          <FormSelectField label="Priority" {...select('priority')} error={!formState.validity.priority} required>
            <PriorityOptions />
          </FormSelectField>
          <FormSelectField label="Category" {...select('category_id')} error={!formState.validity.category_id} required>
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

export default connectMessageDispatch(EditRuleForm)
