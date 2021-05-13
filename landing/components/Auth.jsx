import React from 'react'
import { AuthProvider } from 'oidc-react'

import oidcConfig from '@/config/oidc'

const Auth = ({children}) => (
  <AuthProvider {...oidcConfig}>{children}</AuthProvider>
)

export default Auth
