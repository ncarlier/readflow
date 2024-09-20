import React from 'react'

import { Box, Button } from '../../components'
import { useDevice } from '../../contexts'

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
    Sadness! readflow can&apos;t be installed on this device from here.
    <br />
    Maybe your device&apos;s browser allows you to do this, but only by using the built-in application menu.
  </p>
)

const InstallationBox = () => {
  const { isInstalled, beforeInstallPromptEvent } = useDevice()

  const install = React.useCallback(() => {
    if (beforeInstallPromptEvent) {
      beforeInstallPromptEvent.prompt()
    }
  }, [beforeInstallPromptEvent])

  return (
    <Box title="Installation">
      {(() => {
        switch (true) {
          case !isInstalled && beforeInstallPromptEvent != null:
            return <Install onClick={install} />
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
