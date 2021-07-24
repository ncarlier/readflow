import React, { CSSProperties } from 'react'

import keeper from './KeeperLogo.svg'
import pocket from './PocketLogo.svg'
import readflow from './ReadflowLogo.svg'
import wallabag from './WallabagLogo.svg'
import webhook from './WebhookLogo.svg'

const logos = {
  keeper,
  pocket,
  readflow,
  wallabag,
  webhook,
  generic: webhook,
}

interface Props {
  name: 'keeper' | 'pocket' | 'readflow' | 'wallabag' | 'webhook' | 'generic'
  title?: string
  style?: CSSProperties
}

export default ({ name, style, title }: Props) => <img src={logos[name]} alt={name} title={title} style={style} />
