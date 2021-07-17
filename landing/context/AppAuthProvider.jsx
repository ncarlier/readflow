import React from 'react'
import { AuthProvider, AuthContext } from 'oidc-react'

import oidcConfig from '@/config/oidc'

const AppAuthProvider = ({children}) => (
  <AuthProvider {...oidcConfig}>{children}</AuthProvider>
)

export default AppAuthProvider
