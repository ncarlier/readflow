import React, { useCallback, useState } from 'react'

import Box from '../../components/Box'
import Button from '../../components/Button'
import { isInstalled } from '../../helpers'

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
    <Button title="Install on your device" primary onClick={onClick}>
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
    Maybe your device is too old or the current readflow configuration does not allow it. Sad.
  </p>
)

export default () => {
  const [installed, setInstalled] = useState(isInstalled())
  const { deferredPrompt } = window

  const installHandler = useCallback(() => {
    deferredPrompt.prompt()
    deferredPrompt.userChoice.then((choice: any) => {
      if (choice.outcome === 'accepted') {
        console.log('User accepted the A2HS prompt')
        setInstalled(true)
      } else {
        setInstalled(false)
      }
      window.deferredPrompt = null
    })
  }, [installed])

  return (
    <Box title="Installation">
      {(() => {
        switch (true) {
          case !installed && !!deferredPrompt:
            return <Install onClick={installHandler} />
          case installed:
            return <Installed />
          default:
            return <Uninstallable />
        }
      })()}
    </Box>
  )
}
