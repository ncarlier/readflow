/* eslint-disable @typescript-eslint/camelcase */
import React, { FormEvent, useCallback, useContext, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { RouteComponentProps } from 'react-router'
import { useFormState } from 'react-use-form-state'

import Button from '../../components/Button'
import CategoriesOptions from '../../components/CategoriesOptions'
import FormInputField from '../../components/FormInputField'
import FormSelectField from '../../components/FormSelectField'
import FormTextareaField from '../../components/FormTextareaField'
import { getGQLError, isValidForm } from '../../helpers'
import HelpLink from '../../components/HelpLink'
import Panel from '../../components/Panel'
import { MessageContext } from '../../context/MessageContext'
import ErrorPanel from '../../error/ErrorPanel'
import { usePageTitle } from '../../hooks'
import useOnMountInputValidator from '../../hooks/useOnMountInputValidator'
import { updateCacheAfterCreate } from './cache'
import { Rule, CreateOrUpdateRuleResponse } from './models'
import PriorityOptions from './PriorityOptions'
import { CreateOrUpdateRule } from './queries'

interface AddRuleFormFields {
  alias: string
  category_id: number
  rule: string
  priority: number
}

type AllProps = RouteComponentProps<{}>

export default ({ history }: AllProps) => {
  usePageTitle('Settings - Add new rule')

  const { showMessage } = useContext(MessageContext)
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text, textarea, select }] = useFormState<AddRuleFormFields>({ rule: '', alias: '', priority: 0 })
  const onMountValidator = useOnMountInputValidator(formState.validity)
  const addRuleMutation = useMutation<CreateOrUpdateRuleResponse, Rule>(CreateOrUpdateRule)

  const addRule = async (rule: Rule) => {
    try {
      const res = await addRuleMutation({
        variables: rule,
        update: updateCacheAfterCreate
      })
      showMessage(`New rule: ${res.data!.createOrUpdateRule.alias}`)
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
