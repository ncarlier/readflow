import React, { useState } from 'react'

import { Box, Button } from '../../components'
import { isDisplayMode } from '../../helpers'
import { useAddToHomescreenPrompt } from '../../hooks'

interface InstallProps {
  onClick?: (e: any) => void
}

const Install = ({ onClick }: InstallProps) => (
  <>
    <p>
      Benefit from a better integration with your system by installing readflow.
      <br />
      This will also activate extra feature such as sharing.
    </p>
    <Button title="Install on your device" variant="primary" onClick={onClick}>
      Install
    </Button>
  </>
)

const Installed = () => (
  <p>
    Great! readflow is installed on your device!
    <br />
    You benefit from a better integration with your system!
  </p>
)

const Uninstallable = () => (
  <p>
    Oh! readflow can&apos;t be installed on this device.
    <br />
    Maybe your device device configuration does not allow it. Sad.
  </p>
)

const InstallationBox = () => {
  const [prompt, promptToInstall] = useAddToHomescreenPrompt()
  const [isInstalled] = useState(isDisplayMode('standalone'))
  const [isInstallable, setInstallableState] = useState(false)

  React.useEffect(
    () => {
      if (prompt) {
        setInstallableState(true)
      }
    },
    [prompt]
  )

  return (
    <Box title="Installation">
      {(() => {
        switch (true) {
          case !isInstalled && isInstallable:
            return <Install onClick={promptToInstall} />
          case isInstalled:
            return <Installed />
          default:
            return <Uninstallable />
        }
      })()}
    </Box>
  )
}

export default InstallationBox
