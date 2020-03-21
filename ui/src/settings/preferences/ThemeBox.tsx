import React, { useEffect, useContext } from 'react'

import Box from '../../components/Box'
import FormSelectField from '../../components/FormSelectField'
import { useFormState } from 'react-use-form-state'
import { LocalConfigurationContext, Theme } from '../../context/LocalConfigurationContext'

interface SwitchThemeFormFields {
  theme: Theme
}

const ThemeSwitch = () => {
  const { localConfiguration, updateLocalConfiguration } = useContext(LocalConfigurationContext)

  const [formState, { select }] = useFormState<SwitchThemeFormFields>({
    theme: localConfiguration.theme
  })

  useEffect(() => {
    if (formState.values.theme !== localConfiguration.theme) {
      const { theme } = formState.values
      updateLocalConfiguration({ ...localConfiguration, theme })
    }
  }, [localConfiguration, formState])

  return (
    <form>
      <FormSelectField label="Theme" {...select('theme')}>
        <option value="light">light</option>
        <option value="dark">dark</option>
        <option value="auto">auto (your system will decide)</option>
      </FormSelectField>
    </form>
  )
}

export default () => (
  <Box title="Theme">
    <p>Change the colors of the user interface according to your preferences.</p>
    <ThemeSwitch />
  </Box>
)
