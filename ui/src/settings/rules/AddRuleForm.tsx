import React, { useCallback, useState } from 'react'

import { useFormState } from 'react-use-form-state'
import { useMutation } from 'react-apollo-hooks'
import { RouteComponentProps } from 'react-router'

import Panel from '../../common/Panel'
import Button from '../../common/Button'
import { usePageTitle } from '../../hooks'
import { CreateOrUpdateRule } from './queries'
import ErrorPanel from '../../error/ErrorPanel'
import { getGQLError, isValidForm } from '../../common/helpers'
import { updateCacheAfterCreate } from './cache'
import FormInputField from '../../common/FormInputField'
import FormSelectField from '../../common/FormSelectField'
import FormTextareaField from '../../common/FormTextareaField'
import { Rule } from './models'
import { connectMessageDispatch, IMessageDispatchProps } from '../../containers/MessageContainer';
import CategoriesOptions from '../../common/CategoriesOptions'
import PriorityOptions from './PriorityOptions'
import HelpLink from '../../common/HelpLink'

interface AddRuleFormFields {
  alias: string
  category_id: number
  rule: string
  priority: number
}

type AllProps = RouteComponentProps<{}> & IMessageDispatchProps

export const AddRuleForm = ({history, showMessage}: AllProps) => {
  usePageTitle('Settings - Add new rule')

  const [errorMessage, setErrorMessage] = useState<string | null>(null) 
  const [formState, { text, textarea, select }] = useFormState<AddRuleFormFields>(
    {rule: '', alias: '', priority: 0}
  )
  const addRuleMutation = useMutation<Rule>(CreateOrUpdateRule)

  const addRule = async (rule: Rule) => {
    try{
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

  const handleOnClick = useCallback(() => {
    if (!isValidForm(formState)) {
      setErrorMessage("Please fill out correctly the mandatory fields.")
      return
    }
    const { alias, rule, priority, category_id } = formState.values
    addRule({alias, rule, priority: parseInt(priority), category_id: parseInt(category_id)})
  }, [formState])

  return (
    <Panel>
      <header>
        <h1>Add new rule</h1>
      </header>
      <section>
        {errorMessage != null &&
          <ErrorPanel title="Unable to add new rule">
            {errorMessage}
          </ErrorPanel>
        }
        <form>
          <FormInputField label="Alias"
            {...text('alias')}
            error={!formState.validity.alias}
            required />
          <FormTextareaField label="Rule"
            {...textarea('rule')}
            error={!formState.validity.rule}
            required>
            <HelpLink href="https://github.com/antonmedv/expr/wiki/The-Expression-Syntax">
              View rule syntax
            </HelpLink>
          </FormTextareaField>
          <FormSelectField label="Priority"
            {...select('priority')}
            error={!formState.validity.priority}
            required>
            <PriorityOptions />
          </FormSelectField>
          <FormSelectField label="Category"
            {...select('category_id')}
            error={!formState.validity.category_id}
            required>
            <option>Please select a category</option>
            <CategoriesOptions />
          </FormSelectField>
        </form>
      </section>
      <footer>
        <Button title="Back to rules" to="/settings/rules">
          Cancel
        </Button>
        <Button
          title="Add rule"
          onClick={handleOnClick}
          primary>
          Add
        </Button>
      </footer>
    </Panel>
  )
}

export default connectMessageDispatch(AddRuleForm)
