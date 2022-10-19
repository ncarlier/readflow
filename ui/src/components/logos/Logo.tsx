import React, { CSSProperties } from 'react'

import { ReactComponent as feedpushr } from './min/FeedpushrLogo.svg'
import { ReactComponent as keeper } from './min/KeeperLogo.svg'
import { ReactComponent as pocket } from './min/PocketLogo.svg'
import { ReactComponent as readflow } from './min/ReadflowLogo.svg'
import { ReactComponent as s3 } from './min/S3Logo.svg'
import { ReactComponent as wallabag } from './min/WallabagLogo.svg'
import { ReactComponent as webhook } from './min/WebhookLogo.svg'


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

export const Logo = ({ name, style, title }: Props) => {
  const Logo = logos[name]
  return <Logo title={title} style={style} />
}
