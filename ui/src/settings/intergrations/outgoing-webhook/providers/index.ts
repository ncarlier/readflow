import { GenericConfigForm } from './GenericConfigForm'
import { KeeperConfigForm } from './KeeperConfigForm'
import { PocketConfigForm } from './PocketConfigForm'
import { ReadflowConfigForm } from './ReadflowConfigForm'
import { S3ConfigForm } from './S3ConfigForm'
import { ShaarliConfigForm } from './ShaarliConfigForm'
import { WallabagConfigForm } from './WallabagConfigForm'

const providers = {
  generic: {
    label: 'Generic webhook',
    config: GenericConfigForm
  },
  keeper: {
    label: 'Keeper',
    config: KeeperConfigForm
  },
  pocket: {
    label: 'Pocket',
    config: PocketConfigForm
  },
  readflow: {
    label: 'Readflow',
    config: ReadflowConfigForm
  },
  s3: {
    label: 'S3',
    config: S3ConfigForm
  },
  shaarli: {
    label: 'Shaarli',
    config: ShaarliConfigForm
  },
  wallabag: {
    label: 'Wallabag',
    config: WallabagConfigForm
  },
}

export default providers
