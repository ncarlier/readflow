import React, { CSSProperties } from 'react'

import feedpushr from './FeedpushrLogo.svg'
import keeper from './KeeperLogo.svg'
import pocket from './PocketLogo.svg'
import readflow from './ReadflowLogo.svg'
import s3 from './S3Logo.svg'
import wallabag from './WallabagLogo.svg'
import webhook from './WebhookLogo.svg'

const logos = {
  feedpushr,
  keeper,
  pocket,
  readflow,
  s3,
  wallabag,
  webhook,
  generic: webhook,
}

interface Props {
  name: 'feedpushr' | 'keeper' | 'pocket' | 'readflow' | 's3' | 'wallabag' | 'webhook' | 'generic'
  title?: string
  style?: CSSProperties
}

export const Logo = ({ name, style, title }: Props) => <img src={logos[name]} alt={name} title={title} style={style} />
