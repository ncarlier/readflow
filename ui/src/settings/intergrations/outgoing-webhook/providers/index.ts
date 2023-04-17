import React from 'react'
import generic from './generic'
import keeper from './keeper'
import pocket from './pocket'
import readflow from './readflow'
import s3 from './s3'
import shaarli from './shaarli'
import wallabag from './wallabag'

interface ConfigProps {
  onChange(config: any): void
  config: any
  locked?: boolean
}

interface Provider {
  label: string
  form: React.FC<ConfigProps>
  marshal: (config: any) => string[]
}

const providers: Record<string, Provider> = {
  generic,
  keeper,
  pocket,
  readflow,
  s3,
  shaarli,
  wallabag
}

export default providers
