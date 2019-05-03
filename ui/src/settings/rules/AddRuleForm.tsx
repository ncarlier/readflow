/* eslint-disable @typescript-eslint/camelcase */
import React, { FormEvent, useCallback, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { RouteComponentProps } from 'react-router'
import { useFormState } from 'react-use-form-state'

import Button from '../../common/Button'
import CategoriesOptions from '../../common/CategoriesOptions'
import FormInputField from '../../common/FormInputField'
import FormSelectField from '../../common/FormSelectField'
import FormTextareaField from '../../common/FormTextareaField'
import { getGQLError, isValidForm } from '../../common/helpers'
import HelpLink from '../../common/HelpLink'
import Panel from '../../common/Panel'
import { connectMessageDispatch, IMessageDispatchProps } from '../../containers/MessageContainer'
import ErrorPanel from '../../error/ErrorPanel'
import { usePageTitle } from '../../hooks'
import useOnMountInputValidator from '../../hooks/useOnMountInputValidator'
import { updateCacheAfterCreate } from './cache'
import { Rule } from './models'
import PriorityOptions from './PriorityOptions'
import { CreateOrUpdateRule } from './queries'

interface AddRuleFormFields {
  alias: string
  category_id: number
  rule: string
  priority: number
}

type AllProps = RouteComponentProps<{}> & IMessageDispatchProps

export const AddRuleForm = ({ history, showMessage }: AllProps) => {
  usePageTitle('Settings - Add new rule')

  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text, textarea, select }] = useFormState<AddRuleFormFields>({ rule: '', alias: '', priority: 0 })
  const onMountValidator = useOnMountInputValidator(formState.validity)
  const addRuleMutation = useMutation<Rule>(CreateOrUpdateRule)

  const addRule = async (rule: Rule) => {
    try {
      const res = await addRuleMutation({
        variables: rule,
        update: updateCacheAfterCreate
      })
      showMessage(`New rule: ${res.data.createOrUpdateRule.id}`)
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
      addRule({ alias, rule, priority: parseInt(priority), category_id: parseInt(category_id) })
    },
    [formState]
  )

  return (
    <Panel>
      <header>
        <h1>Add new rule</h1>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to add new rule">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField
            label="Alias"
            {...text('alias')}
            error={!formState.validity.alias}
            required
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
            <option>Please select a category</option>
            <CategoriesOptions />
          </FormSelectField>
        </form>
      </section>
      <footer>
        <Button title="Back to rules" to="/settings/rules">
          Cancel
        </Button>
        <Button title="Add rule" onClick={handleOnSubmit} primary>
          Add
        </Button>
      </footer>
    </Panel>
  )
}

export default connectMessageDispatch(AddRuleForm)
